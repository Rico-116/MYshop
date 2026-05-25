package models

import "time"

const (
	OrderStatusAll      = -1
	OrderStatusCanceled = 0
	OrderStatusUnpaid   = 1
	OrderStatusPaid     = 2
	OrderStatusShipped  = 3
	OrderStatusFinished = 4

	OderStatusShipped  = OrderStatusShipped
	OderStatusFinished = OrderStatusFinished
)

func GetOrderStatusText(status int) string {
	switch status {
	case OrderStatusCanceled:
		return "已取消"
	case OrderStatusUnpaid:
		return "待支付"
	case OrderStatusPaid:
		return "待发货"
	case OrderStatusShipped:
		return "待收货"
	case OrderStatusFinished:
		return "已完成"
	default:
		return "未知状态"
	}
}

type OrderTab struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
}

func GetOrderTabs() []OrderTab {
	return []OrderTab{
		{Name: "全部", Status: OrderStatusAll},
		{Name: "待付款", Status: OrderStatusUnpaid},
		{Name: "待发货", Status: OrderStatusPaid},
		{Name: "待收货", Status: OrderStatusShipped},
		{Name: "已完成", Status: OrderStatusFinished},
		{Name: "已取消", Status: OrderStatusCanceled},
	}
}

type OrderListRequest struct {
	Status   int `json:"status"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}
type OrderListVO struct {
	Id                    uint          `json:"id"`
	OrderNo               string        `json:"order_no"`
	UserId                uint          `json:"user_id"`
	Status                int           `json:"status"`
	StatusText            string        `json:"status_text"`
	TotalAmount           float64       `json:"total_amount"`
	PayAmount             float64       `json:"pay_amount"`
	FreightAmount         float64       `json:"freight_amount"`
	CouponAmount          float64       `json:"coupon_amount"`
	ReceiverName          string        `json:"receiver_name"`
	ReceiverPhone         string        `json:"receiver_phone"`
	ReceiverProvince      string        `json:"receiver_province"`
	ReceiverCity          string        `json:"receiver_city"`
	ReceiverDistrict      string        `json:"receiver_district"`
	ReceiverDetailAddress string        `json:"receiver_detail_address"`
	Remark                string        `json:"remark"`
	PayTime               *time.Time    `json:"pay_time"`
	DeliveryTime          *time.Time    `json:"delivery_time"`
	FinishTime            *time.Time    `json:"finish_time"`
	CloseTime             *time.Time    `json:"close_time"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
	Items                 []OrderItemVO `json:"items" gorm:"-"`
}
type OrderItemVO struct {
	OrderNo      string  `json:"order_no"`
	ProductId    uint    `json:"product_id"`
	SkuId        uint    `json:"sku_id"`
	ProductName  string  `json:"product_name"`
	ProductImage string  `json:"product_image"`
	SkuName      string  `json:"sku_name"`
	Price        float64 `json:"price"`
	Quantity     int     `json:"quantity"`
	TotalAmount  float64 `json:"total_amount"`
}
type OrderListResult struct {
	List     []OrderListVO `json:"list"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Tabs     []OrderTab    `json:"tabs"`
}
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
	UserDeleted           int        `json:"user_deleted"`
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
	PayAmount   float64 `json:"pay_amount"`
	ItemCount   int     `json:"item_count"`
	Status      int     `json:"status"`
	StatusText  string  `json:"status_text"`
}
type PayOrderRequest struct {
	OrderNO string `json:"order_no"`
}
type PayOrderResult struct {
	OrderNo    string  `json:"order_no"`
	PayAmount  float64 `json:"pay_amount"`
	Status     int     `json:"status"`
	StatusText string  `json:"status_text"`
	PayTime    string  `json:"pay_time"`
}
type PayPageResult struct {
	OrderId               uint          `json:"order_id"`
	OrderNo               string        `json:"order_no"`
	PayAmount             float64       `json:"pay_amount"`
	Status                int           `json:"status"`
	StatusText            string        `json:"status_text"`
	CreatedAt             string        `json:"created_at"`
	ReceiverName          string        `json:"receiver_name"`
	ReceiverPhone         string        `json:"receiver_phone"`
	ReceiverProvince      string        `json:"receiver_province"`
	ReceiverCity          string        `json:"receiver_city"`
	ReceiverDistrict      string        `json:"receiver_district"`
	ReceiverDetailAddress string        `json:"receiver_detail_address"`
	Items                 []OrderItemVO `json:"items" gorm:"-"`
}
type OrderPreviewRequest struct {
	SourceType string `json:"source_type"`
	CartIds    []uint `json:"cart_ids"`
	SkuId      uint   `json:"sku_id"`
	Quantity   int    `json:"quantity"`
}
type OrderPreviewItemVO struct {
	ProductId    uint    `json:"product_id"`
	SkuId        uint    `json:"sku_id"`
	ProductName  string  `json:"product_name"`
	ProductImage string  `json:"product_image"`
	SkuName      string  `json:"sku_name"`
	Price        float64 `json:"price"`
	Quantity     int     `json:"quantity"`
	TotalAmount  float64 `json:"total_amount"`
	Stock        int     `json:"stock"`
}
type OrderPreviewResult struct {
	Items          []OrderPreviewItemVO `json:"items"`
	AddressList    []UserAddress        `json:"address_list"`
	DefaultAddress *UserAddress         `json:"default_address"`
	TotalAmount    float64              `json:"total_amount"`
	FreightAmount  float64              `json:"freight_amount"`
	CouponAmount   float64              `json:"coupon_amount"`
	PayAmount      float64              `json:"pay_amount"`
	ItemCount      int                  `json:"item_count"`
}

type OrderDetailResult struct {
	OrderId               uint          `json:"order_id"`
	OrderNo               string        `json:"order_no"`
	UserId                uint          `json:"user_id"`
	Status                int           `json:"status"`
	StatusText            string        `json:"status_text"`
	TotalAmount           float64       `json:"total_amount"`
	PayAmount             float64       `json:"pay_amount"`
	FreightAmount         float64       `json:"freight_amount"`
	CouponAmount          float64       `json:"coupon_amount"`
	ReceiverName          string        `json:"receiver_name"`
	ReceiverPhone         string        `json:"receiver_phone"`
	ReceiveProvince       string        `json:"receive_province"`
	ReceiverCity          string        `json:"receiver_city"`
	ReceiverDistrict      string        `json:"receiver_district"`
	ReceiverDetailAddress string        `json:"receiver_detail_address"`
	Remark                string        `json:"remark"`
	PayTime               *time.Time    `json:"pay_time"`
	DeliverTime           *time.Time    `json:"deliver_time"`
	FinishTime            *time.Time    `json:"finish_time"`
	CloseTime             *time.Time    `json:"close_time"`
	CreateTime            time.Time     `json:"create_time"`
	UpdateTime            time.Time     `json:"update_time"`
	Items                 []OrderItemVO `json:"items" gorm:"-"`
}

type CancelOrderRequest struct {
	OrderNo string `json:"order_no" binding:"required"`
}

type CancelOrderResult struct {
	OrderNo    string `json:"order_no"`
	Status     int    `json:"status"`
	StatusText string `json:"status_text"`
}

type AdminShipOrderRequest struct {
	OrderNo string `json:"order_no" binding:"required"`
}

type AdminCancelOrderRequest struct {
	OrderNo string `json:"order_no" binding:"required"`
}
