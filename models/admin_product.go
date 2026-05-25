package models

type AdminCreateProductRequest struct {
	CategoryId  uint   `json:"category_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Subtitle    string `json:"subtitle"`
	MainImage   string `json:"main_image"`
	Description string `json:"description"`
	Status      *int   `json:"status"`
}

type AdminCreateSkuRequest struct {
	ProductId uint    `json:"product_id" binding:"required"`
	SkuCode   string  `json:"sku_code" binding:"required"`
	SkuName   string  `json:"sku_name" binding:"required"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
	Image     string  `json:"image"`
	Status    *int    `json:"status"`
}

type AdminUpdateProductRequest struct {
	CategoryId  *uint   `json:"category_id"`
	Name        *string `json:"name"`
	Subtitle    *string `json:"subtitle"`
	MainImage   *string `json:"main_image"`
	Description *string `json:"description"`
	Status      *int    `json:"status"`
}

type AdminUpdateSkuRequest struct {
	SkuCode *string  `json:"sku_code"`
	SkuName *string  `json:"sku_name"`
	Price   *float64 `json:"price"`
	Stock   *int     `json:"stock"`
	Image   *string  `json:"image"`
	Status  *int     `json:"status"`
}
