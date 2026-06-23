package models

import "time"

const (
	SeckillActivityStatusOff         = 0
	SeckillActivityStatusOn          = 1
	AdminSeckillListStatusClosed     = 0
	AdminSeckillListStatusNotStarted = 1
	AdminSeckillListStatusRunning    = 2
	AdminSeckillListStatusEnded      = 3
	SeckillResultQueuing             = 0 // 排队中
	SeckillResultSuccess             = 1 // 秒杀成功
	SeckillResultFail                = 2 // 秒杀失败
)

type SeckillSubmitRequest struct {
	ActivityId string `json:"activity_id" binding:"required"`
	AddressId  string `json:"address_id" binding:"required"`
}
type SeckillSubmitResult struct {
	ActivityId uint   `json:"activity_id"`
	Status     int    `json:"status"`
	StatusText string `json:"status_text"`
}
type SeckillResult struct {
	ActivityId uint   `json:"activity_id"`
	Status     int    `json:"status"`
	StatusText string `json:"status_text"`
	OrderNo    string `json:"order_no,omitempty"`
}
type UserSeckillActivityVO struct {
	Id            uint    `json:"id"`
	ProductId     uint    `json:"product_id"`
	SkuId         uint    `json:"sku_id"`
	ProductName   string  `json:"product_name"`
	ProductImage  string  `json:"product_image"`
	SkuName       string  `json:"sku_name"`
	OriginalPrice float64 `json:"original_price"`
	SeckillPrice  float64 `json:"seckill_price"`
	Stock         int     `json:"stock"`
	SoldCount     int     `json:"sold_count"`
	RemainStock   int     `json:"remain_stock"`
	Status        int     `json:"status"`
	StatusText    string  `json:"status_text"`
	StartTime     string  `json:"start_time"`
	EndTime       string  `json:"end_time"`
}
type UserSeckillListResult struct {
	List     []UserSeckillActivityVO `json:"list"`
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
}
type SeckillOrder struct {
	Id         uint      `json:"id"`
	ActivityId uint      `json:"activity_id"`
	UserId     uint      `json:"user_id"`
	OrderNo    string    `json:"order_no"`
	Status     int       `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
type SeckillOrderBuyItem struct {
	ActivityId    uint    `json:"activity_id"`
	ProductId     uint    `json:"product_id"`
	SkuId         uint    `json:"sku_id"`
	ProductName   string  `json:"product_name"`
	ProductImage  string  `json:"product_image"`
	ProductStatus int     `json:"product_status"`
	SkuName       string  `json:"sku_name"`
	SkuStatus     int     `json:"sku_status"`
	Price         float64 `json:"price"`
	Stock         int     `json:"stock"`
}
type SeckillOrderMessage struct {
	ActivityId uint `json:"activity_id"`
	UserId     uint `json:"user_id"`
	AddressId  uint `json:"address_id"`
}
type SeckillActivity struct {
	Id           uint      `json:"id" gorm:"column:id;primaryKey"`
	ProductId    uint      `json:"product_id"`
	SkuId        uint      `json:"sku_id"`
	SeckillPrice float64   `json:"seckill_price"`
	Stock        int       `json:"stock"`
	SoldCount    int       `json:"sold_count"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"create_time"`
	UpdatedAt    time.Time `json:"update_time"`
}
type AdminCreateSeckillActivityRequest struct {
	ProductId    uint    `json:"product_id" binding:"required"`
	SkuId        uint    `json:"sku_id" binding:"required"`
	SeckillPrice float64 `json:"seckill_price" binding:"required"`
	Stock        int     `json:"stock" binding:"required"`
	StartTime    string  `json:"start_time" binding:"required"`
	EndTime      string  `json:"end_time" binding:"required"`
}
type AdminUpdateSeckillActivityRequest struct {
	Id           uint    `json:"id" binding:"required"`
	SeckillPrice float64 `json:"seckill_price" binding:"required"`
	Stock        int     `json:"stock" binding:"required"`
	StartTime    string  `json:"start_time" binding:"required"`
	EndTime      string  `json:"end_time" binding:"required"`
	Status       int     `json:"status" binding:"required"`
}
type AdminSeckillActivityVO struct {
	Id            uint    `json:"id"`
	ProductId     uint    `json:"product_id"`
	SkuId         uint    `json:"sku_id"`
	ProductName   string  `json:"product_name"`
	ProductImage  string  `json:"product_image"`
	SkuName       string  `json:"sku_name"`
	OriginalPrice float64 `json:"original_price"`
	SeckillPrice  float64 `json:"seckill_price"`
	Stock         int     `json:"stock"`
	SoldCount     int     `json:"sold_count"`
	RemainStock   int     `json:"remain_stock"`
	Status        int     `json:"status"`
	StatusText    string  `json:"status_text"`
	StartTime     string  `json:"start_time"`
	EndTime       string  `json:"end_time"`
	CreateTime    string  `json:"create_time"`
}
type AdminSeckillListResult struct {
	List     []AdminSeckillActivityVO `json:"list"`
	Total    int64                    `json:"total"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"page_size"`
}
