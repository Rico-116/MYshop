package Service

import (
	"MYshop/dao"
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"encoding/json"
	"errors"

	"go.uber.org/zap"
)

// 如果 MQ 处理失败，要恢复 Redis 库存和用户标记
const seckillRollbackScript = `
local stockKey = KEYS[1]
local userKey = KEYS[2]

if redis.call('EXISTS', userKey) == 1 then
    redis.call('INCR', stockKey)
    redis.call('DEL', userKey)
    return 1
end

return 0
`

// 启动秒杀订单消费者
func StartSeckillOrderConsumer() error {
	if err := util.MQChannel.Qos(1, 0, false); err != nil {
		return err
	}

	msgs, err := util.MQChannel.Consume(
		util.SeckillOrderQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	logger.Log.Info("秒杀订单消费者启动成功",
		zap.String("queue", util.SeckillOrderQueue),
	)

	go func() {
		for d := range msgs {
			var msg models.SeckillOrderMessage

			if err := json.Unmarshal(d.Body, &msg); err != nil {
				logger.Log.Error("秒杀消息反序列化失败",
					zap.Error(err),
					zap.ByteString("body", d.Body),
				)
				_ = d.Ack(false)
				continue
			}

			if err := HandleSeckillOrderMessage(msg); err != nil {
				logger.Log.Error("处理秒杀订单失败",
					zap.Error(err),
					zap.Uint("activity_id", msg.ActivityId),
					zap.Uint("user_id", msg.UserId),
					zap.Uint("address_id", msg.AddressId),
				)
				_ = d.Ack(false)
				continue
			}

			_ = d.Ack(false)
		}
	}()

	return nil
}

// 处理秒杀消息
func HandleSeckillOrderMessage(msg models.SeckillOrderMessage) error {
	if err := createSeckillOrderByMessage(msg); err != nil {
		// 订单创建失败，恢复 Redis 库存
		rollbackRedisSeckill(msg.ActivityId, msg.UserId)

		// 写入失败结果，前端轮询时能看到失败原因
		_ = saveSeckillResult(msg.ActivityId, msg.UserId, models.SeckillResult{
			ActivityId: msg.ActivityId,
			Status:     models.SeckillResultFail,
			StatusText: err.Error(),
		})

		return err
	}

	return nil
}

// 真正创建秒杀订单
func createSeckillOrderByMessage(msg models.SeckillOrderMessage) error {
	if msg.ActivityId == 0 || msg.UserId == 0 || msg.AddressId == 0 {
		return errors.New("秒杀消息参数错误")
	}

	// MQ 可能重复投递，所以消费者这里也要查一次是否已经创建过
	exists, existOrderNo, err := dao.CheckSeckillOrderExists(msg.ActivityId, msg.UserId)
	if err != nil {
		return err
	}
	if exists {
		_ = saveSeckillResult(msg.ActivityId, msg.UserId, models.SeckillResult{
			ActivityId: msg.ActivityId,
			Status:     models.SeckillResultSuccess,
			StatusText: "秒杀成功",
			OrderNo:    existOrderNo,
		})
		return nil
	}

	// 查询活动
	activity, err := dao.GetSeckillActivityById(msg.ActivityId)
	if err != nil {
		return err
	}
	if activity == nil || activity.Id == 0 {
		return errors.New("秒杀活动不存在")
	}
	if activity.Status != models.SeckillActivityStatusOn {
		return errors.New("秒杀活动未开启")
	}

	// 查询秒杀商品信息
	item, err := dao.GetSeckillOrderBuyItem(msg.ActivityId)
	if err != nil {
		return err
	}
	if item == nil || item.ActivityId == 0 {
		return errors.New("秒杀商品不存在")
	}
	if item.ProductStatus != 1 {
		return errors.New("商品已下架")
	}
	if item.SkuStatus != 1 {
		return errors.New("商品规格不可购买")
	}

	// 查询收货地址
	addr, err := dao.GetAddressByIdAndUserId(msg.AddressId, msg.UserId)
	if err != nil {
		return err
	}
	if addr == nil || addr.Id == 0 {
		return errors.New("收货地址不存在")
	}

	orderNo := generateOrderNo(msg.UserId)

	tx := util.Db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	// 1. 写 seckill_order，靠唯一索引防止重复秒杀
	if err := dao.CreateSeckillOrderTx(tx, msg.ActivityId, msg.UserId, orderNo); err != nil {
		if errors.Is(err, dao.ErrSeckillOrderDuplicated) {
			exists, existOrderNo, e := dao.CheckSeckillOrderExists(msg.ActivityId, msg.UserId)
			if e != nil {
				return e
			}
			if exists {
				committed = true
				tx.Rollback()
				_ = saveSeckillResult(msg.ActivityId, msg.UserId, models.SeckillResult{
					ActivityId: msg.ActivityId,
					Status:     models.SeckillResultSuccess,
					StatusText: "秒杀成功",
					OrderNo:    existOrderNo,
				})
				return nil
			}
		}
		return err
	}

	// 2. 扣秒杀活动库存：sold_count + 1
	if err := dao.DeductSeckillActivityStockTx(tx, msg.ActivityId); err != nil {
		return err
	}

	// 3. 扣普通 SKU 库存
	if err := dao.DeductSkuStockTx(tx, item.SkuId, 1); err != nil {
		return err
	}

	// 4. 创建普通订单主表
	order := &models.Order{
		OrderNo:               orderNo,
		UserId:                msg.UserId,
		Status:                models.OrderStatusUnpaid,
		TotalAmount:           item.Price,
		PayAmount:             item.Price,
		CouponAmount:          0,
		FreightAmount:         0,
		ReceiverName:          addr.ReceiverName,
		ReceiverPhone:         addr.ReceiverPhone,
		ReceiverProvince:      addr.Province,
		ReceiverCity:          addr.City,
		ReceiverDistrict:      addr.District,
		ReceiverDetailAddress: addr.DetailAddress,
		Remark:                "秒杀订单",
	}

	orderId, err := dao.CreateOrderTx(tx, order)
	if err != nil {
		return err
	}

	// 5. 创建订单明细
	orderItem := &models.OrderItem{
		OrderId:      orderId,
		OrderNo:      orderNo,
		UserId:       msg.UserId,
		ProductId:    item.ProductId,
		SkuId:        item.SkuId,
		ProductName:  item.ProductName,
		ProductImage: item.ProductImage,
		SkuName:      item.SkuName,
		Price:        item.Price,
		Quantity:     1,
		TotalAmount:  item.Price,
	}

	if err := dao.CreateOrderItemTx(tx, orderItem); err != nil {
		return err
	}

	// 6. 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}
	committed = true

	// 7. 发送普通订单的延迟关单消息
	if err := PublishOrderCloseDelay(orderNo, msg.UserId); err != nil {
		logger.Log.Error("秒杀订单发送延迟关单消息失败",
			zap.Error(err),
			zap.String("order_no", orderNo),
		)
	}

	// 8. 写入成功结果
	_ = saveSeckillResult(msg.ActivityId, msg.UserId, models.SeckillResult{
		ActivityId: msg.ActivityId,
		Status:     models.SeckillResultSuccess,
		StatusText: "秒杀成功",
		OrderNo:    orderNo,
	})

	return nil
}

// 恢复 Redis 秒杀库存
func rollbackRedisSeckill(activityId uint, userId uint) {
	stockKey := util.SeckillStockKey(activityId)
	userKey := util.SeckillUserKey(activityId, userId)

	_ = util.RDB.Eval(
		util.Ctx,
		seckillRollbackScript,
		[]string{stockKey, userKey},
	).Err()
}
