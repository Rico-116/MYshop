package models

import "time"

type User struct {
	Username string    `json:"username"`
	Password string    `json:"password"`
	Nickname string    `json:"nickname"`
	Phone    string    `json:"phone"`
	Email    string    `json:"email"`
	Avatar   string    `json:"avatar"`
	Status   int       `json:"status"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt string    `json:"update_at"`
}

func (User) TableName() string {
	return "user"
}

type SendRegisterCodeRequest struct {
	Email string `json:"email"`
}

type RegisterRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm"`
	Email           string `json:"email"`
	Code            string `json:"code"`
	Nickname        string `json:"nickname"`
	Phone           string `json:"phone"`
	Avatar          string `json:"avatar"`
}
