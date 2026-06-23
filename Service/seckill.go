package Service

import (
	"MYshop/dao"
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// lua脚本：三件事，一是判断用户是否抢过，二是内存是否足够，三是扣redis库存并标记用户已抢
const seckillCheckScript = `local stockKey = KEYS[1]
local userKey = KEYS[2]

if redis.call('EXISTS', userKey) == 1 then
    return 2
end

local stock = tonumber(redis.call('GET', stockKey))
if stock == nil then
    return 3
end

if stock <= 0 then
    return 1
end

redis.call('DECR', stockKey)
redis.call('SET', userKey, '1', 'EX', ARGV[1])
return 0`

func GetUserSeckillActivityList(page int, pageSize int) (*models.UserSeckillListResult, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	key := util.SeckillListKey(page, pageSize)
	val, err := util.RDB.Get(util.Ctx, key).Result()
	if err == nil {
		var result models.UserSeckillListResult
		if e := json.Unmarshal([]byte(val), &result); e == nil {
			refreshUserSeckillListStatusAndStock(&result)
			return &result, nil
		} else {
			logger.Log.Warn("秒杀列表缓存反序列化失败", zap.Error(e))
		}
	}
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.Log.Warn("读取秒杀列表缓存失败", zap.Error(err))
	}

	result, err := rebuildUserSeckillActivityListCache(page, pageSize)
	if err != nil {
		return nil, err
	}
	refreshUserSeckillListStatusAndStock(result)
	return result, nil
}

func rebuildUserSeckillActivityListCache(page int, pageSize int) (*models.UserSeckillListResult, error) {
	key := util.SeckillListKey(page, pageSize)
	lockKey := util.SeckillListLockKey(page, pageSize)
	locked, err := util.RDB.SetNX(util.Ctx, lockKey, "1", 5*time.Second).Result()
	if err != nil {
		logger.Log.Warn("秒杀列表缓存重建加锁失败", zap.Error(err))
	}

	if !locked {
		time.Sleep(80 * time.Millisecond)
		val, err := util.RDB.Get(util.Ctx, key).Result()
		if err == nil {
			var result models.UserSeckillListResult
			if e := json.Unmarshal([]byte(val), &result); e == nil {
				return &result, nil
			}
		}
		return queryUserSeckillActivityList(page, pageSize)
	}

	defer func() {
		if err := util.RDB.Del(util.Ctx, lockKey).Err(); err != nil {
			logger.Log.Warn("秒杀列表缓存重建解锁失败", zap.Error(err))
		}
	}()

	val, err := util.RDB.Get(util.Ctx, key).Result()
	if err == nil {
		var result models.UserSeckillListResult
		if e := json.Unmarshal([]byte(val), &result); e == nil {
			return &result, nil
		}
	}

	result, err := queryUserSeckillActivityList(page, pageSize)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(result)
	if err != nil {
		return result, nil
	}
	if err := util.RDB.Set(util.Ctx, key, data, withJitter(20*time.Second, 20)).Err(); err != nil {
		logger.Log.Warn("写入秒杀列表缓存失败", zap.Error(err))
	}
	return result, nil
}

func queryUserSeckillActivityList(page int, pageSize int) (*models.UserSeckillListResult, error) {
	list, total, err := dao.GetUserSeckillActivityList(page, pageSize)
	if err != nil {
		return nil, err
	}

	for i := range list {
		list[i].Status, list[i].StatusText = getUserSeckillStatus(list[i])
	}

	return &models.UserSeckillListResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func refreshUserSeckillListStatusAndStock(result *models.UserSeckillListResult) {
	if result == nil {
		return
	}
	for i := range result.List {
		result.List[i].Status, result.List[i].StatusText = getUserSeckillStatus(result.List[i])
		stock, err := util.RDB.Get(util.Ctx, util.SeckillStockKey(result.List[i].Id)).Int()
		if err == nil {
			result.List[i].RemainStock = stock
			result.List[i].SoldCount = result.List[i].Stock - stock
			if result.List[i].SoldCount < 0 {
				result.List[i].SoldCount = 0
			}
		}
	}
}

func DeleteUserSeckillListCache() {
	var cursor uint64
	pattern := util.SeckillListKeyPrefix + ":*"
	for {
		keys, nextCursor, err := util.RDB.Scan(util.Ctx, cursor, pattern, 100).Result()
		if err != nil {
			logger.Log.Warn("扫描秒杀列表缓存失败", zap.Error(err))
			return
		}
		if len(keys) > 0 {
			if err := util.RDB.Del(util.Ctx, keys...).Err(); err != nil {
				logger.Log.Warn("删除秒杀列表缓存失败", zap.Error(err))
				return
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			return
		}
	}
}

func getUserSeckillStatus(activity models.UserSeckillActivityVO) (int, string) {
	start, _ := parseSeckillTime(activity.StartTime)
	end, _ := parseSeckillTime(activity.EndTime)
	now := time.Now()

	if now.Before(start) {
		return models.AdminSeckillListStatusNotStarted, "未开始"
	}

	if now.Before(end) {
		return models.AdminSeckillListStatusRunning, "进行中"
	}

	return models.AdminSeckillListStatusEnded, "已结束"
}
func SubmitSeckill(userId uint, req models.SeckillSubmitRequest) (*models.SeckillSubmitResult, error) {
	aid, err := strconv.Atoi(req.ActivityId)
	if err != nil {
		return nil, err
	}
	as, err := strconv.Atoi(req.AddressId)
	if err != nil {
		return nil, err
	}
	if userId == 0 {
		return nil, errors.New("用户未登录")
	}
	if req.ActivityId == "" {
		return nil, errors.New("秒杀活动ID不能为空")
	}
	if req.AddressId == "" {
		return nil, errors.New("请选择收货地址")
	}

	// 1. 查活动
	activity, err := dao.GetSeckillActivityById(uint(aid))
	if err != nil {
		return nil, err
	}
	if activity == nil || activity.Id == 0 {
		return nil, errors.New("秒杀活动不存在")
	}
	if activity.Status != models.SeckillActivityStatusOn {
		return nil, errors.New("秒杀活动未开启")
	}

	// 2. 判断活动时间
	now := time.Now()
	if now.Before(activity.StartTime) {
		return nil, errors.New("秒杀活动未开始")
	}
	if !now.Before(activity.EndTime) {
		return nil, errors.New("秒杀活动已结束")
	}

	// 3. 校验收货地址必须属于当前用户
	addr, err := dao.GetAddressByIdAndUserId(uint(as), userId)
	if err != nil {
		return nil, err
	}
	if addr == nil || addr.Id == 0 {
		return nil, errors.New("收货地址不存在")
	}

	// 4. 数据库兜底：判断是否已经成功秒杀过
	exists, _, err := dao.CheckSeckillOrderExists(uint(aid), userId)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("你已经参与过该秒杀活动")
	}

	// 5. Redis Lua 原子扣库存 + 防重复
	if err := checkAndDeductSeckillStock(uint(aid), userId); err != nil {
		return nil, err
	}

	// 6. 先写一个排队中的结果，方便前端轮询
	_ = saveSeckillResult(uint(aid), userId, models.SeckillResult{
		ActivityId: uint(aid),
		Status:     models.SeckillResultQueuing,
		StatusText: "排队中",
	})

	// 7. 发送 MQ，让消费者异步创建订单
	msg := models.SeckillOrderMessage{
		ActivityId: uint(aid),
		UserId:     userId,
		AddressId:  uint(as),
	}

	if err := PublishSeckillOrderMessage(msg); err != nil {
		// 如果 MQ 发送失败，要把 Redis 预扣库存还回去
		rollbackRedisSeckill(uint(aid), userId)
		_ = saveSeckillResult(uint(aid), userId, models.SeckillResult{
			ActivityId: uint(aid),
			Status:     models.SeckillResultFail,
			StatusText: "秒杀失败，请稍后重试",
		})

		return nil, errors.New("秒杀请求提交失败")
	}
	aid, err = strconv.Atoi(req.ActivityId)
	if err != nil {
		return nil, err
	}
	// 8. 这里不返回订单号，只返回排队中
	return &models.SeckillSubmitResult{
		ActivityId: uint(aid),
		Status:     models.SeckillResultQueuing,
		StatusText: "排队中",
	}, nil
}

func checkAndDeductSeckillStock(activityId uint, userId uint) error {
	stockKey := util.SeckillStockKey(activityId)
	userKey := util.SeckillUserKey(activityId, userId)
	res, err := evalSeckillCheckScript(stockKey, userKey)
	if err != nil {
		return errors.New("redis秒杀校验失败")
	}
	switch res {
	case 0:
		return nil
	case 1:
		return errors.New("秒杀商品已抢光")
	case 2:
		exists, _, err := dao.CheckSeckillOrderExists(activityId, userId)
		if err != nil {
			return err
		}
		if !exists {
			if err := util.RDB.Del(util.Ctx, userKey).Err(); err != nil {
				return errors.New("redis秒杀校验失败")
			}
			res, err = evalSeckillCheckScript(stockKey, userKey)
			if err != nil {
				return errors.New("redis秒杀校验失败")
			}
			if res == 0 {
				return nil
			}
			if res == 1 {
				return errors.New("秒杀商品已抢光")
			}
			if res == 3 {
				return errors.New("秒杀库存未初始化")
			}
		}
		return errors.New("你已经参与过该秒杀活动")
	case 3:
		return errors.New("秒杀库存未初始化")
	default:
		return errors.New("秒杀失败")
	}
}

func evalSeckillCheckScript(stockKey string, userKey string) (int, error) {
	return util.RDB.Eval(
		util.Ctx,
		seckillCheckScript,
		[]string{stockKey, userKey},
		24*60*60, // 用户抢购标记保存 24 小时
	).Int()
}

func GetSeckillResult(userId uint, activityId uint) (*models.SeckillResult, error) {
	if userId == 0 {
		return nil, errors.New("用户未登录")
	}
	if activityId == 0 {
		return nil, errors.New("秒杀活动ID不能为空")
	}

	key := util.SeckillResultKey(activityId, userId)

	val, err := util.RDB.Get(util.Ctx, key).Result()
	if err != nil {
		return nil, errors.New("暂未查询到秒杀结果")
	}

	var result models.SeckillResult
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, errors.New("秒杀结果解析失败")
	}

	return &result, nil
}
func saveSeckillResult(activityId uint, userId uint, result models.SeckillResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return util.RDB.Set(
		util.Ctx,
		util.SeckillResultKey(activityId, userId),
		data,
		24*time.Hour,
	).Err()
}
