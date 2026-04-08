package models

import "time"

type Product struct {
	Id          uint      `json:"id"`
	CategoryId  uint      `json:"category_id"`
	Name        string    `json:"name"`
	Subtitle    string    `json:"subtitle"`
	MainImage   string    `json:"main_image"`
	Status      int       `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Rating      float64   `json:"rating"`
	RatingCount int       `json:"rating_count"`
	ClickCount  int       `json:"click_count"`
	Price       float64   `json:"price"`
}
