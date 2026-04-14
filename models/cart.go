package models

import "time"

type Cart struct {
	Id        uint      `json:"id"`
	UserId    uint      `json:"user_id"`
	ProductId uint      `json:"product_id"`
	SkuId     uint      `json:"sku_id"`
	Quantity  uint      `json:"quantity"`
	Checked   int       `json:"checked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type AddCartRequest struct {
	SkuId    uint `json:"sku_id" binding:"required"`
	Quantity uint `json:"quantity" binding:"required"`
}

type CartItem struct {
	CartId        uint    `json:"cart_id" `
	ProductId     uint    `json:"product_id" `
	SkuId         uint    `json:"sku_id" `
	ProductName   string  `json:"product_name" `
	MainImage     string  `json:"main_image" `
	SkuName       string  `json:"sku_name" `
	Price         float64 `json:"price" `
	Stock         int     `json:"stock" `
	Quantity      int     `json:"quantity" `
	Checked       int     `json:"checked" `
	Subtotal      float64 `json:"subtotal" `
	Invalid       int     `json:"invalid" `
	ProductStatus int     `json:"-"`
	SkuStatus     int     `json:"-"`
}
type UpdateCartQuantityRequest struct {
	CartId   uint `json:"cart_id" binding:"required"`
	Quantity uint `json:"quantity" binding:"required"`
}
type UpdateCartCheckRequest struct {
	CartId  uint `json:"cart_id" binding:"required"`
	Checked int  `json:"checked" binding:"required"`
}
type DeleteCartRequest struct {
	CartId uint `json:"cart_id" binding:"required"`
}
type CartDisplayItem struct {
	CartId      uint    `json:"cart_id" `
	SkuId       uint    `json:"sku_id" `
	ProductId   uint    `json:"product_id" `
	Title       string  `json:"title" `
	Image       string  `json:"image" `
	Price       float64 `json:"price" `
	Stock       int     `json:"stock" `
	Quantity    int     `json:"quantity" `
	Checked     int     `json:"checked" `
	TotalAmount float64 `json:"total_amount" `
}
