package models

import "time"

type Banner struct {
	Id        uint      `json:"id"`
	Title     string    `json:"title"`
	ImageURL  string    `json:"image"`
	Link      string    `json:"link"`
	Sort      int       `json:"sort"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time `json:"updated"`
}
