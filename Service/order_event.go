package Service

import (
	"MYshop/models"
	"MYshop/util"
	"encoding/json"
)

// PublishOrderCloseDelay 生产者发送延迟关单消息
func PublishOrderCloseDelay(orderNo string, userId uint) error {
	msg := models.OrderCloseMessage{
		OrderNo: orderNo,
		UserId:  userId,
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return util.PublishWithConfirm(
		util.OrderDelayExchange,
		util.OrderCloseDelayRoutingKey,
		body,
	)
}
