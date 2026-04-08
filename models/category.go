package models

type Category struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	ParentId uint   `json:"parent_id"`
	Sort     uint   `json:"sort"`
	Status   uint   `json:"status"`
	Icon     string `json:"icon"` //分类图标识别
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
}
type CategoryTree struct {
	Id       uint           `json:"id"`
	Name     string         `json:"name"`
	Icon     string         `json:"icon,omitempty"`
	Children []CategoryTree `json:"children"` //这里更改要告诉前端
}
type CategoryDisplay struct {
	CurrentCategory Category   `json:"current_category"`
	SubCategories   []Category `json:"sub_categories"`
	ProductList     []Product  `json:"product_list"`
	IsLeaf          bool       `json:"is_leaf"` //只有叶子节点才有商品，那么在不是叶子节点是要怎么办
}
