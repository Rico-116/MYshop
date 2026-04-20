package models

import "time"

type AdminUser struct {
	Id        uint      `json:"id" gorm:"column:id;primaryKey"`
	Username  string    `json:"username" gorm:"column:username"`
	Password  string    `json:"password" gorm:"column:password"`
	Nickname  string    `json:"nickname" gorm:"column:nickname"`
	Status    int       `json:"status" gorm:"column:status"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}
