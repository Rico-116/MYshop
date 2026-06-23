package dao

import (
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"go.uber.org/zap"
)

func SearchProductList(req models.SearchProductRequest) ([]models.Product, int64, error) {
	var list []models.Product
	var total int64

	db := util.Db.Table("product").Where("product.status = 1")
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("product.name LIKE ? OR product.subtitle LIKE ? OR product.description LIKE ?", like, like, like)
	}
	if req.CategoryId != nil && *req.CategoryId > 0 {
		db = db.Where("product.category_id = ?", *req.CategoryId)
	}
	if req.MinPrice != nil {
		db = db.Where("product.price >= ?", *req.MinPrice)
	}
	if req.MaxPrice != nil {
		db = db.Where("product.price <= ?", *req.MaxPrice)
	}

	if err := db.Count(&total).Error; err != nil {
		fields := append(searchRequestLogFields(req), zap.Error(err))
		logger.Log.Error("搜索商品数量统计失败", fields...)
		return nil, 0, err
	}

	orderBy := "product.id desc"
	switch req.Sort {
	case "price_asc":
		orderBy = "product.price asc, product.id desc"
	case "price_desc":
		orderBy = "product.price desc, product.id desc"
	case "newest":
		orderBy = "product.created_at desc, product.id desc"
	case "hot":
		orderBy = "product.click_count desc, product.id desc"
	}

	offset := (req.Page - 1) * req.PageSize
	err := db.Select(
		"product.id",
		"product.category_id",
		"product.name",
		"product.subtitle",
		"product.main_image",
		"product.status",
		"product.description",
		"product.created_at",
		"product.updated_at",
		"product.rating",
		"product.rating_count",
		"product.click_count",
		"product.price",
		"COALESCE(SUM(product_sku.stock), 0) AS stock",
	).
		Joins("LEFT JOIN product_sku ON product_sku.product_id = product.id AND product_sku.status = 1").
		Group("product.id").
		Order(orderBy).
		Offset(offset).
		Limit(req.PageSize).
		Scan(&list).Error
	if err != nil {
		fields := append(searchRequestLogFields(req), zap.Error(err))
		logger.Log.Error("搜索商品列表失败", fields...)
		return nil, 0, err
	}

	return list, total, nil
}

func searchRequestLogFields(req models.SearchProductRequest) []zap.Field {
	fields := []zap.Field{
		zap.String("keyword", req.Keyword),
		zap.String("sort", req.Sort),
		zap.Int("page", req.Page),
		zap.Int("page_size", req.PageSize),
	}
	if req.CategoryId != nil {
		fields = append(fields, zap.Uint("category_id", *req.CategoryId))
	}
	if req.MinPrice != nil {
		fields = append(fields, zap.Float64("min_price", *req.MinPrice))
	}
	if req.MaxPrice != nil {
		fields = append(fields, zap.Float64("max_price", *req.MaxPrice))
	}
	return fields
}
