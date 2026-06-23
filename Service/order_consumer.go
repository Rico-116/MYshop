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

// StartOrderCloseConsumer 监听真正的关单队列
func StartOrderCloseConsumer() error {
	if err := util.MQChannel.Qos(1, 0, false); err != nil {
		return err
	}
	msgs, err := util.MQChannel.Consume(
		util.OrderCloseExecuteQueue,
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
	logger.Log.Info("订单关单消费者启动成功",
		zap.String("queue", util.OrderCloseExecuteQueue),
	)

	go func() {
		for d := range msgs {
			logger.Log.Info("收到延迟关单消息",
				zap.ByteString("body", d.Body),
			)

			var msg models.OrderCloseMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				logger.Log.Error("关单消息反序列化失败",
					zap.ByteString("body", d.Body),
					zap.Error(err),
				)
				_ = d.Ack(false)
				continue
			}

			if err := HandleOrderClose(msg.OrderNo); err != nil {
				logger.Log.Error("处理超时关单失败",
					zap.String("order_no", msg.OrderNo),
					zap.Error(err),
				)

				// 先不要无限重回队列，否则可能疯狂刷日志
				// 调试阶段建议先 Ack，定位业务错误
				_ = d.Ack(false)
				continue
			}

			logger.Log.Info("处理超时关单成功",
				zap.String("order_no", msg.OrderNo),
			)

			_ = d.Ack(false)
		}
	}()

	return nil
}
func HandleOrderClose(orderNo string) error {
	if orderNo == "" {
		return errors.New("OrderNo为空")
	}
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
	rows, err := dao.UpdateOrderToCanceledIfUnpaidTx(tx, orderNo)
	if err != nil {
		return err
	}
	if rows == 0 {
		if err := tx.Commit().Error; err != nil {
			return err
		}
		committed = true //更新失败，说明已支付不需要恢复库存，此时committed=false,不需要回滚
		return nil
	}
	items, err := dao.GetOrderItemsByOrderNo(orderNo)
	if err != nil {
		return err
	}
	for _, item := range items {
		if err := dao.RestoreSkuStockTx(tx, item.SkuId, item.Quantity); err != nil {
			return err
		}
	}

	seckillOrder, err := dao.CancelSeckillOrderAndRestoreActivityStockTx(tx, orderNo)
	if err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	committed = true
	if seckillOrder != nil {
		rollbackRedisSeckill(seckillOrder.ActivityId, seckillOrder.UserId)
	}
	return nil
}
