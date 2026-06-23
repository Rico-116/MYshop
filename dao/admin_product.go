package dao

import (
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func AdminCreateProduct(product *models.Product) error {
	err := util.Db.Table("product").Create(product).Error
	if err != nil {
		logger.Log.Error("管理员新增商品失败",
			zap.Error(err),
			zap.Uint("category_id", product.CategoryId),
			zap.String("name", product.Name),
		)
		return err
	}
	return nil
}

func AdminGetProductById(id int) (*models.Product, error) {
	var product models.Product
	err := util.Db.Table("product").Select("id", "category_id", "name", "subtitle", "main_image", "status", "description", "created_at", "updated_at", "rating", "rating_count", "click_count", "price").Where("id = ?", id).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Log.Error("查询商品失败", zap.Error(err), zap.Uint("id", uint(id)))
		return nil, err
	}
	return &product, nil
}

func AdminGetProductList(keyword string, category *uint, status *int, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64
	db := util.Db.Table("product")
	if keyword != "" {
		db = db.Where("name LIKE ? OR subtitle LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if category != nil && *category > 0 {
		db = db.Where("category_id = ?", *category)
	}
	if status != nil {
		db = db.Where("product.status = ?", *status)
	}
	if err := db.Count(&total).Error; err != nil {
		logger.Log.Error("管理员统计商品列表失败", zap.Error(err))
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
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
	).Joins("LEFT JOIN product_sku ON product_sku.product_id = product.id").
		Group("product.id").
		Order("product.id desc").
		Offset(offset).
		Limit(pageSize).
		Scan(&products).Error
	if err != nil {
		logger.Log.Error("管理员查询商品列表失败", zap.Error(err))
		return nil, 0, err
	}
	return products, total, nil
}

func AdminUpdateProduct(productId uint, updates map[string]interface{}) error {
	result := util.Db.Table("product").Where("id = ?", productId).Updates(updates)
	if result.Error != nil {
		logger.Log.Error("管理员修改商品失败",
			zap.Error(result.Error),
			zap.Uint("product_id", productId),
			zap.Int("field_count", len(updates)),
		)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("商品不存在或没有任何修改")
	}
	return nil
}

func AdminDeleteProduct(productId uint) error {
	return util.Db.Transaction(func(tx *gorm.DB) error {
		result := tx.Table("product").Where("id = ?", productId).Updates(map[string]interface{}{
			"status":     0,
			"updated_at": gorm.Expr("NOW()"),
		})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("商品不存在")
		}
		if err := tx.Table("product_sku").Where("product_id = ?", productId).Updates(map[string]interface{}{
			"status":     0,
			"updated_at": gorm.Expr("NOW()"),
		}).Error; err != nil {
			return err
		}
		return nil
	})
}

func AdminCreateProductSku(productSku *models.ProductSku) error {
	err := util.Db.Table("product_sku").Create(productSku).Error
	if err != nil {
		logger.Log.Error("新增商品SKU失败",
			zap.Error(err),
			zap.Uint("product_id", productSku.ProductId),
			zap.String("sku_code", productSku.SkuCode),
		)
		return err
	}
	return nil
}

func AdminCountSkuCode(productId uint, skuCode string) (int64, error) {
	var count int64
	err := util.Db.Table("product_sku").Where("product_id = ? AND sku_code = ?", productId, skuCode).Count(&count).Error
	if err != nil {
		logger.Log.Error("管理员统计SKU编码失败", zap.Error(err), zap.Uint("product_id", productId), zap.String("sku_code", skuCode))
		return 0, err
	}
	return count, nil
}

func AdminCountSkuCodeExcludeId(productId uint, skuCode string, excludeSkuId uint) (int64, error) {
	var count int64
	err := util.Db.Table("product_sku").Where("product_id = ? AND sku_code = ? AND id <> ?", productId, skuCode, excludeSkuId).Count(&count).Error
	if err != nil {
		logger.Log.Error("管理员统计SKU编码失败", zap.Error(err), zap.Uint("product_id", productId), zap.Uint("sku_id", excludeSkuId), zap.String("sku_code", skuCode))
		return 0, err
	}
	return count, nil
}

func AdminGetSkuById(skuId uint) (*models.ProductSku, error) {
	var productSku models.ProductSku
	err := util.Db.Table("product_sku").Select(
		"id",
		"product_id",
		"sku_code",
		"sku_name",
		"price",
		"stock",
		"image",
		"status",
		"created_at",
		"updated_at",
	).Where("id = ?", skuId).Limit(1).Scan(&productSku).Error
	if err != nil {
		logger.Log.Error("管理员查询SKU失败", zap.Error(err), zap.Uint("sku_id", skuId))
		return nil, err
	}
	if productSku.Id == 0 {
		return nil, nil
	}
	return &productSku, nil
}

func AdminGetSkuListByProductId(productId uint) ([]models.ProductSku, error) {
	var list []models.ProductSku
	err := util.Db.Table("product_sku").Select(
		"id",
		"product_id",
		"sku_code",
		"sku_name",
		"price",
		"stock",
		"image",
		"status",
		"created_at",
		"updated_at",
	).Where("product_id = ?", productId).Order("id ASC").Scan(&list).Error
	if err != nil {
		logger.Log.Error("管理员查询SKU列表失败", zap.Error(err), zap.Uint("product_id", productId))
		return nil, err
	}
	return list, nil
}

func AdminUpdateSku(skuId uint, updates map[string]interface{}) error {
	result := util.Db.Table("product_sku").Where("id = ?", skuId).Updates(updates)
	if result.Error != nil {
		logger.Log.Error("管理员修改SKU失败",
			zap.Error(result.Error),
			zap.Uint("sku_id", skuId),
			zap.Int("field_count", len(updates)),
		)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("SKU不存在或没有任何修改")
	}
	return nil
}

func AdminDeleteSku(skuId uint) error {
	result := util.Db.Table("product_sku").Where("id = ?", skuId).Updates(map[string]interface{}{
		"status":     0,
		"updated_at": gorm.Expr("NOW()"),
	})
	if result.Error != nil {
		logger.Log.Error("管理员删除SKU失败", zap.Error(result.Error), zap.Uint("sku_id", skuId))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("SKU不存在")
	}
	return nil
}
