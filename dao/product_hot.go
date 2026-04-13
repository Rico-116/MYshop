package dao

import (
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"go.uber.org/zap"
)

func BatchAddProductClickCount(click map[int]int64) error {
	if len(click) == 0 {
		return nil
	}
	for productID, delta := range click {
		sql := "UPDATE product SET click_count = click_count + ? WHERE id = ?"
		if err := util.Db.Exec(sql, delta, productID).Error; err != nil {
			logger.Log.Error("回写商品点击量失败", zap.Int("product_id", productID), zap.Int64("delta", delta), zap.Error(err))
			return err
		}
	}
	return nil
}
func GetDefaultHotProducts(limit int) ([]models.Product, error) {
	if limit <= 0 {
		limit = 8
	}

	sqlStr := `
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
		WHERE status = 1
		ORDER BY click_count DESC, rating DESC, rating_count DESC, id ASC
		LIMIT ?
	`

	var products []models.Product
	err := util.Db.Raw(sqlStr, limit).Scan(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}
