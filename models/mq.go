package models

type OrderCloseMessage struct {
	OrderNo string `json:"order_no"`
	UserId  uint   `json:"user_id"`
}
