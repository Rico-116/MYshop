package models

import "time"

const (
	OrderStatusCanceled = 0
	OrderStatusUnpaid   = 1
	OrderStatusPaid     = 2
	OrderStatusShipped  = 3
	OrderStatusFinished = 4

	// 兼容你原来拼写
	OderStatusShipped  = OrderStatusShipped
	OderStatusFinished = OrderStatusFinished
)

type Order struct {
	Id                    uint       `json:"id"`
	OrderNo               string     `json:"order_no"`
	UserId                uint       `json:"user_id"`
	Status                int        `json:"status"`
	TotalAmount           float64    `json:"total_amount"`   // 总金额
	PayAmount             float64    `json:"pay_amount"`     // 实际支付金额
	CouponAmount          float64    `json:"coupon_amount"`  // 优惠金额
	FreightAmount         float64    `json:"freight_amount"` // 运费
	ReceiverName          string     `json:"receiver_name"`
	ReceiverPhone         string     `json:"receiver_phone"`
	ReceiverProvince      string     `json:"receiver_province"`
	ReceiverCity          string     `json:"receiver_city"`
	ReceiverDistrict      string     `json:"receiver_district"`
	ReceiverDetailAddress string     `json:"receiver_detail_address"`
	Remark                string     `json:"remark"`
	PayTime               *time.Time `json:"pay_time"`
	DeliveryTime          *time.Time `json:"delivery_time"`
	FinishTime            *time.Time `json:"finish_time"`
	CloseTime             *time.Time `json:"close_time"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type OrderItem struct {
	Id           uint      `json:"id"`
	OrderId      uint      `json:"order_id"`
	OrderNo      string    `json:"order_no"`
	UserId       uint      `json:"user_id"`
	ProductId    uint      `json:"product_id"`
	SkuId        uint      `json:"sku_id"`
	ProductName  string    `json:"product_name"`
	ProductImage string    `json:"product_image"`
	SkuName      string    `json:"sku_name"`
	Price        float64   `json:"price"`
	Quantity     int       `json:"quantity"`
	TotalAmount  float64   `json:"total_amount"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// 创建订单请求
type CreateOrderRequest struct {
	SourceType string `json:"source_type"` // cart / direct
	AddressId  uint   `json:"address_id"`
	Remark     string `json:"remark"`

	// 购物车下单
	CartIds []uint `json:"cart_ids"`

	// 直接购买
	SkuId    uint `json:"sku_id"`
	Quantity int  `json:"quantity"`
}

// 下单时统一组装的商品结构
type OrderBuyItem struct {
	CartId        uint    `json:"cart_id"`
	ProductId     uint    `json:"product_id"`
	SkuId         uint    `json:"sku_id"`
	Quantity      int     `json:"quantity"`
	ProductName   string  `json:"product_name"`
	ProductImage  string  `json:"product_image"`
	ProductStatus int     `json:"product_status"`
	SkuName       string  `json:"sku_name"`
	Price         float64 `json:"price"`
	Stock         int     `json:"stock"`
	SkuStatus     int     `json:"sku_status"`
	TotalAmount   float64 `json:"total_amount"`
}

type CreateOrderResult struct {
	OrderId     uint    `json:"order_id"`
	OrderNo     string  `json:"order_no"`
	TotalAmount float64 `json:"total_amount"`
	ItemCount   int     `json:"item_count"`
	Status      int     `json:"status"`
}
