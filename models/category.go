package models

type Category struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	ParentId uint   `json:"parent_id"`
	Sort     uint   `json:"sort"`
	Status   uint   `json:"status"`
	Icon     string `json:"icon"`
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
}
type CategoryTree struct {
	Id       uint           `json:"id"`
	Name     string         `json:"name"`
	Icon     string         `json:"icon,omitempty"`
	Children []CategoryTree `json:"parent_id"`
}
