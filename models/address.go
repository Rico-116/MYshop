package models

import "time"

type UserAddress struct {
	Id            uint      `json:"id"`
	UserId        uint      `json:"user_id"`
	ReceiverName  string    `json:"receiver_name"`
	ReceiverPhone string    `json:"receiver_phone"`
	Province      string    `json:"province"`
	City          string    `json:"city"`
	District      string    `json:"district"`
	DetailAddress string    `json:"detail_address"`
	IsDefault     int       `json:"is_default"`
	Status        int       `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AddAddressRequest struct {
	ReceiverName  string `json:"receiver_name"`
	ReceiverPhone string `json:"receiver_phone"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	DetailAddress string `json:"detail_address"`
	IsDefault     int    `json:"is_default"`
}

type UpdateAddressRequest struct {
	Id            uint   `json:"id"`
	ReceiverName  string `json:"receiver_name"`
	ReceiverPhone string `json:"receiver_phone"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	DetailAddress string `json:"detail_address"`
	IsDefault     int    `json:"is_default"`
}
