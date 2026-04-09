package dao

import (
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"go.uber.org/zap"
)

func GetProductList() ([]models.Product, error) {
	var list []models.Product
	sql := "select id,category_id,name,subtitle,main_image,status,description,created_at,updated_at,rating,rating_count,price from product where status=1 order  by id desc"
	err := util.Db.Raw(sql).Scan(&list).Error
	if err != nil {
		logger.Log.Error("查询商品列表失败", zap.Error(err))
		return nil, err
	}
	return list, nil
}
func GetProductById(id int) (*models.Product, error) {
	var product models.Product
	sql := "select id,category_id,name,subtitle,main_image,status,description,created_at,updated_at,rating,rating_count,click_count,price from product where id=? AND status=1 LIMIT 1"
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
			rating_count,
			click_count,
			price
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
func GetProductByIDs(id []int) ([]models.Product, error) {
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
			rating_count,
			click_count,
			price
		FROM product
		WHERE id IN ? AND status = 1
		ORDER BY id DESC
	`
	err := util.Db.Raw(sql, id).Scan(&list).Error
	if err != nil {
		logger.Log.Error("查询分类商品失败", zap.Error(err))
		return nil, err
	}
	return list, nil
}

func GetProductSkuListByProductId(productId int) ([]models.ProductSku, error) {
	var list []models.ProductSku
	sql := `
		SELECT 
			id,
			product_id,
			sku_code,
			sku_name,
			price,
			stock,
			image,
			status,
			created_at,
			updated_at
		FROM product_sku
		WHERE product_id = ? AND status = 1
		ORDER BY id ASC
	`
	err := util.Db.Raw(sql, productId).Scan(&list).Error
	if err != nil {
		logger.Log.Error("查询商品SKU失败", zap.Error(err))
		return nil, err
	}
	return list, nil
}
func GetMinSkuPriceByProductId(productId int) (float64, error) {
	var minPrice float64

	sql := `
		SELECT MIN(price)
		FROM product_sku
		WHERE product_id = ? AND status = 1 AND stock > 0
	`
	err := util.Db.Raw(sql, productId).Scan(&minPrice).Error
	if err != nil {
		return 0, err
	}

	if minPrice > 0 {
		return minPrice, nil
	}

	sql = `
		SELECT MIN(price)
		FROM product_sku
		WHERE product_id = ? AND status = 1
	`
	err = util.Db.Raw(sql, productId).Scan(&minPrice).Error
	if err != nil {
		return 0, err
	}

	return minPrice, nil
}
func SyncProductMinPrice(productId int) error {
	minPrice, err := GetMinSkuPriceByProductId(productId)
	if err != nil {
		return err
	}

	sql := `UPDATE product SET price = ? WHERE id = ?`
	return util.Db.Exec(sql, minPrice, productId).Error
}
