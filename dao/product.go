package dao

import (
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"go.uber.org/zap"
)

func GetProductList() ([]models.Product, error) {
	var list []models.Product
	sql := "select id,category_id,name,subtitle,main_image,status,description,created_at,updated_at,reting,rating_count from product where status=1 order  by id desc"
	err := util.Db.Raw(sql).Scan(&list).Error
	if err != nil {
		logger.Log.Error("查询商品列表失败", zap.Error(err))
		return nil, err
	}
	return list, nil
}
func GetProductById(id int) (*models.Product, error) {
	var product models.Product
	sql := "select id,category_id,name,subtitle,main_image,status,description,created_at,updated_at,reting,rating_count from product where id=? AND status=1 LIMIT 1"
	err := util.Db.Raw(sql, id).Scan(&product).Error
	if err != nil {
		logger.Log.Error("查询商品详情页失败", zap.Error(err))
		return nil, err
	}
	return &product, nil
}
func GetProductByCategoryId(id int) ([]models.Product, error) {
	var list []models.Product
	sql := `
		SELECT 
			id,
			category_id,
			name,
			subtitle,
			main_image,
			status,
			description,
			created_at,
			updated_at,
			rating,
			rating_count
		FROM product
		WHERE category_id = ? AND status = 1
		ORDER BY id DESC
	`
	err := util.Db.Raw(sql, id).Scan(&list).Error
	if err != nil {
		logger.Log.Error("按分类查询商品失败", zap.Error(err))
		return nil, err
	}
	return list, nil
}
