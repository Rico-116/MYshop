package Service

import (
	"MYshop/dao"
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"errors"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

func parseSeckillTime(timeStr string) (time.Time, error) {
	timeStr = strings.TrimSpace(timeStr)
	if timeStr == "" {
		return time.Time{}, errors.New("时间不能为空")
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
	if err != nil {
		return time.Time{}, errors.New("错误时间格式，请使用：2026-05-28 10:00:00")
	}
	return t, nil
}

func AdminCreateSeckillActivity(req models.AdminCreateSeckillActivityRequest) error {
	if req.ProductId == 0 {
		return errors.New("商品ID不能为空")
	}
	if req.SkuId == 0 {
		return errors.New("SKU ID不能为空")
	}
	if req.SeckillPrice <= 0 {
		return errors.New("秒杀价必须大于0")
	}
	if req.Stock <= 0 {
		return errors.New("秒杀库存必须大于0")
	}
	startTime, err := parseSeckillTime(req.StartTime)
	if err != nil {
		return err
	}
	endTime, err := parseSeckillTime(req.EndTime)
	if err != nil {
		return err
	}
	if !endTime.After(startTime) {
		return errors.New("结束时间必须晚于开始时间")
	}
	if err := dao.CheckProductSkuExist(req.ProductId, req.SkuId); err != nil {
		return err
	}
	activity := &models.SeckillActivity{
		ProductId:    req.ProductId,
		SkuId:        req.SkuId,
		SeckillPrice: req.SeckillPrice,
		Stock:        req.Stock,
		SoldCount:    0,
		StartTime:    startTime,
		EndTime:      endTime,
		Status:       models.SeckillActivityStatusOn,
	}
	if err := dao.CreateSeckillActivity(activity); err != nil {
		return err
	}
	DeleteUserSeckillListCache()
	return syncSeckillStockToRedis(activity)
}

func AdminGetSeckillActivityList(page int, pageSize int, status int) (*models.AdminSeckillListResult, error) {
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	if pageSize > 100 {
		pageSize = 100
	}
	list, total, err := dao.GetAdminSeckillActivityList(page, pageSize, status)
	if err != nil {
		return nil, err
	}
	for i := range list {
		list[i].Status, list[i].StatusText = getAdminSeckillListStatus(list[i])
	}
	return &models.AdminSeckillListResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func getAdminSeckillListStatus(activity models.AdminSeckillActivityVO) (int, string) {
	if activity.Status == models.SeckillActivityStatusOff {
		return models.AdminSeckillListStatusClosed, "已关闭"
	}

	now := time.Now()
	startTime, startErr := parseSeckillTime(activity.StartTime)
	endTime, endErr := parseSeckillTime(activity.EndTime)
	if startErr != nil || endErr != nil {
		return models.AdminSeckillListStatusEnded, "已结束"
	}

	if now.Before(startTime) {
		return models.AdminSeckillListStatusNotStarted, "未开始"
	}
	if now.Before(endTime) {
		return models.AdminSeckillListStatusRunning, "进行中"
	}
	return models.AdminSeckillListStatusEnded, "已结束"
}

func AdminUpdateSeckillActivity(req models.AdminUpdateSeckillActivityRequest) error {
	if req.Id == 0 {
		return errors.New("秒杀活动ID不能为空")
	}
	if req.SeckillPrice <= 0 {
		return errors.New("秒杀价必须大于0")
	}
	if req.Stock <= 0 {
		return errors.New("秒杀库存必须大于0")
	}

	oldActivity, err := dao.GetSeckillActivityById(req.Id)
	if err != nil {
		return err
	}
	if oldActivity == nil || oldActivity.Id == 0 {
		return errors.New("秒杀活动不存在")
	}

	if oldActivity.SoldCount > req.Stock {
		return errors.New("秒杀库存不能小于已售数量")
	}
	startTime, err := parseSeckillTime(req.StartTime)
	if err != nil {
		return err
	}

	endTime, err := parseSeckillTime(req.EndTime)
	if err != nil {
		return err
	}

	if !endTime.After(startTime) {
		return errors.New("结束时间必须晚于开始时间")
	}

	if req.Status != models.SeckillActivityStatusOn &&
		req.Status != models.SeckillActivityStatusOff {
		return errors.New("秒杀活动状态错误")
	}
	activity := &models.SeckillActivity{
		Id:           req.Id,
		SeckillPrice: req.SeckillPrice,
		Stock:        req.Stock,
		SoldCount:    oldActivity.SoldCount,
		StartTime:    startTime,
		EndTime:      endTime,
		Status:       req.Status,
	}

	if err := dao.UpdateSeckillActivity(activity); err != nil {
		return err
	}
	DeleteUserSeckillListCache()

	if oldActivity.StartTime.After(time.Now()) {
		return syncUnstartedSeckillStockToRedis(activity)
	}
	return nil
}

func AdminCloseSeckillActivity(activityId uint) error {
	if activityId == 0 {
		return errors.New("秒杀活动ID不能为空")
	}

	activity, err := dao.GetSeckillActivityById(activityId)
	if err != nil {
		return err
	}
	if activity == nil || activity.Id == 0 {
		return errors.New("秒杀活动不存在")
	}

	if err := dao.CloseSeckillActivity(activityId); err != nil {
		return err
	}
	DeleteUserSeckillListCache()
	return nil
}

func AdminInitSeckillStock(activityId uint) error {
	if activityId == 0 {
		return errors.New("秒杀活动ID不能为空")
	}
	activity, err := dao.GetSeckillActivityById(activityId)
	if err != nil {
		return err
	}
	if activity == nil || activity.Id == 0 {
		return errors.New("秒杀活动不存在")
	}
	if activity.Status != models.SeckillActivityStatusOn {
		return errors.New("秒杀活动未开启")
	}
	DeleteUserSeckillListCache()
	return syncSeckillStockToRedis(activity)
}

func LoadUnfinishedSeckillStockToRedis() error {
	activities, err := dao.GetUnfinishedSeckillActivities()
	if err != nil {
		return err
	}

	for i := range activities {
		if err := syncSeckillStockToRedis(&activities[i]); err != nil {
			return err
		}
	}

	logger.Log.Info("未结束秒杀活动库存已加载到Redis", zap.Int("count", len(activities)))
	return nil
}

func syncUnstartedSeckillStockToRedis(activity *models.SeckillActivity) error {
	if activity.Status != models.SeckillActivityStatusOn || !activity.EndTime.After(time.Now()) {
		key := util.SeckillStockKey(activity.Id)
		if err := util.RDB.Del(util.Ctx, key).Err(); err != nil {
			logger.Log.Error("删除未开始秒杀活动Redis库存失败",
				zap.Error(err),
				zap.Uint("activity_id", activity.Id),
			)
			return err
		}
		return nil
	}
	return syncSeckillStockToRedis(activity)
}

func syncSeckillStockToRedis(activity *models.SeckillActivity) error {
	remainStock := activity.Stock - activity.SoldCount
	if remainStock < 0 {
		remainStock = 0
	}
	key := util.SeckillStockKey(activity.Id)
	err := util.RDB.Set(
		util.Ctx,
		key,
		remainStock,
		0,
	).Err()
	if err != nil {
		logger.Log.Error("同步秒杀库存到Redis失败",
			zap.Error(err),
			zap.Uint("activity_id", activity.Id),
		)
		return err
	}
	return nil
}

func ParseIntQuery(val string, defaultVal int) int {
	val = strings.TrimSpace(val)

	if val == "" {
		return defaultVal
	}

	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}

	return n
}

func ParseUintQuery(val string) uint {
	val = strings.TrimSpace(val)

	if val == "" {
		return 0
	}

	n, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}

	return uint(n)
}
