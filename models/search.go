package models

type SearchProductRequest struct {
	Keyword    string   `json:"keyword"`
	CategoryId *uint    `json:"category_id"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
	Sort       string   `json:"sort"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
}

type SearchProductResult struct {
	List     []Product `json:"list"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
	Keyword  string    `json:"keyword"`
}
