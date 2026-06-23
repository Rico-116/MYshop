package Service

import (
	"MYshop/models"
	"MYshop/util"
	"encoding/json"
)

func PublishSeckillOrderMessage(msg models.SeckillOrderMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return util.PublishWithConfirm(
		util.SeckillOrderExchange,
		util.SeckillOrderRoutingKey,
		body,
	)
}
