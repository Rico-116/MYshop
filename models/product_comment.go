package models

import "time"

type ProductComment struct {
	Id            uint      `json:"id" gorm:"column:id;primaryKey"`
	ProductId     uint      `json:"product_id" gorm:"column:product_id"`
	UserId        uint      `json:"user_id" gorm:"column:user_id"`
	OrderNo       string    `json:"order_no" gorm:"column:order_no"`
	SkuId         uint      `json:"sku_id" gorm:"column:sku_id"`
	Rating        *int      `json:"rating" gorm:"column:rating"`
	Content       string    `json:"content" gorm:"column:content"`
	RootId        uint      `json:"root_id" gorm:"column:root_id"`
	ParentId      uint      `json:"parent_id" gorm:"column:parent_id"`
	ReplyToUserId uint      `json:"reply_to_user_id" gorm:"column:reply_to_user_id"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (ProductComment) TableName() string {
	return "product_comment"
}

type CreateProductCommentRequest struct {
	ProductId uint   `json:"product_id" binding:"required"`
	OrderNo   string `json:"order_no"`
	SkuId     uint   `json:"sku_id"`
	Rating    int    `json:"rating"`
	Content   string `json:"content" binding:"required"`
	ParentId  uint   `json:"parent_id"`
}

type ProductCommentVO struct {
	Id              uint               `json:"id"`
	ProductId       uint               `json:"product_id"`
	UserId          uint               `json:"user_id"`
	Nickname        string             `json:"nickname"`
	Username        string             `json:"username"`
	Avatar          string             `json:"avatar"`
	OrderNo         string             `json:"order_no,omitempty"`
	SkuId           uint               `json:"sku_id,omitempty"`
	Rating          *int               `json:"rating,omitempty"`
	Content         string             `json:"content"`
	RootId          uint               `json:"root_id"`
	ParentId        uint               `json:"parent_id"`
	ReplyToUserId   uint               `json:"reply_to_user_id,omitempty"`
	ReplyToNickname string             `json:"reply_to_nickname,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	CreatedTime     string             `json:"created_time"`
	Children        []ProductCommentVO `json:"children" gorm:"-"`
}

type ProductCommentListResult struct {
	List     []ProductCommentVO `json:"list"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}
