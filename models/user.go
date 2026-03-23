package models

import "time"

type User struct {
	UserId    uint      `json:"id" gorm:"column:id;primaryKey"`
	Username  string    `json:"username" gorm:"column:username"`
	Password  string    `json:"password" gorm:"column:password"`
	Nickname  string    `json:"nickname" gorm:"column:nickname"`
	Phone     string    `json:"phone" gorm:"column:phone"`
	Email     string    `json:"email" gorm:"column:email"`
	Avatar    string    `json:"avatar" gorm:"column:avatar"`
	Status    int       `json:"status" gorm:"column:status"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:update_at"`
}
type SendLoginCodeRequest struct {
	Email string `json:"email" binding:"required"`
}
type EmailLoginRequest struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

func (User) TableName() string {
	return "user"
}

type SendRegisterCodeRequest struct {
	Email string `json:"email"`
}

type RegisterRequest struct {
	UserId          string `json:"id"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm"`
	Email           string `json:"email"`
	Code            string `json:"code"`
	Nickname        string `json:"nickname"`
	Phone           string `json:"phone"`
	Avatar          string `json:"avatar"`
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
