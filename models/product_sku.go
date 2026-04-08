package models

import "time"

type ProductSku struct {
	Id        uint      `json:"id"`
	ProductId uint      `json:"product_id"`
	SkuCode   string    `json:"sku_code"`
	SkuName   string    `json:"sku_name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	Image     string    `json:"image"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
