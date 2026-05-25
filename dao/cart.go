package dao

import (
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"errors"
	"go.uber.org/zap"
	"time"
)

type CartSkuInfo struct {
	ID        uint
	ProductID uint
	Stock     uint
}

func GetSkuByProductAndSku(productID uint, skuID uint) (*CartSkuInfo, error) {
	var sku CartSkuInfo

	sql := `
		SELECT 
			id,
			product_id,
			stock
		FROM product_sku
		WHERE id = ?
		  AND product_id = ?
		LIMIT 1
	`

	err := util.Db.Raw(sql, skuID, productID).Scan(&sku).Error
	if err != nil {
		logger.Log.Error("查询SKU库存失败",
			zap.Error(err),
			zap.Uint("product_id", productID),
			zap.Uint("sku_id", skuID),
		)
		return nil, err
	}

	if sku.ID == 0 {
		return nil, nil
	}

	return &sku, nil
}
func GetSkuById(skuId uint) (*models.ProductSku, error) {
	var sku models.ProductSku
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
		WHERE id = ? AND status = 1
		LIMIT 1
	`
	err := util.Db.Raw(sql, skuId).Scan(&sku).Error
	if err != nil {
		logger.Log.Error("查询sku失败", zap.Error(err), zap.Uint("sku_id", skuId))
		return nil, err
	}
	return &sku, nil
}
func GetCartByUserAndSku(userId uint, skuId uint) (*models.Cart, error) {
	var cart models.Cart
	sql := `
		SELECT
			id,
			user_id,
			product_id,
			sku_id,
			quantity,
			checked,
			created_at,
			updated_at
		FROM cart
		WHERE user_id = ? AND sku_id = ?
		LIMIT 1
	`
	err := util.Db.Raw(sql, userId, skuId).Scan(&cart).Error
	if err != nil {
		logger.Log.Error("查询购物车失败", zap.Error(err), zap.Uint("user_id", userId), zap.Uint("sku_id", skuId))
		return nil, err
	}
	return &cart, nil
}
func CreateCart(cart *models.Cart) error {
	sql := `
		INSERT INTO cart
		(user_id, product_id, sku_id, quantity, checked, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`
	result := util.Db.Exec(sql,
		cart.UserId,
		cart.ProductId,
		cart.SkuId,
		cart.Quantity,
		cart.Checked)
	if result.Error != nil {
		logger.Log.Error("新增购物车失败", zap.Error(result.Error), zap.Any("cart", cart))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("新增购物车失败")
	}
	return nil
}

func UpdateCartQuantity(userId uint, cartId uint, quantity int) error {
	sql := `UPDATE cart SET quantity = ?,updated_at=? WHERE id = ? AND user_id = ?`
	result := util.Db.Exec(sql, quantity, time.Now(), cartId, userId)
	if result.Error != nil {
		logger.Log.Error("更新购物车数量失败",
			zap.Error(result.Error),
			zap.Uint("user_id", userId),
			zap.Uint("cart_id", cartId),
			zap.Int("quantity", quantity),
		)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("购物车项不存在或无权限修改")
	}
	return nil
}

func GetCartListByUserId(userId uint) ([]models.CartItem, error) {
	var list []models.CartItem
	sql := `SELECT
			c.id AS cart_id,
			c.product_id,
			c.sku_id,
			c.quantity,
			c.checked,
			p.subtitle AS product_name,
			p.main_image,
			p.status AS product_status,
			s.sku_name,
			s.price,
			s.stock,
			s.status AS sku_status
		FROM cart c
		LEFT JOIN product p ON c.product_id = p.id
		LEFT JOIN product_sku s ON c.sku_id = s.id
		WHERE c.user_id = ?
		ORDER BY c.id DESC`
	err := util.Db.Raw(sql, userId).Scan(&list).Error
	if err != nil {
		logger.Log.Error("查询购物车列表失败", zap.Error(err), zap.Uint("user_id", userId))
		return nil, err
	}
	//logger.Log.Debug("<UNK>", zap.Any("list", list))
	return list, nil
}
func GetCartById(userId, cartId uint) (*models.Cart, error) {
	var cart models.Cart
	sql := `
		SELECT
			id,
			user_id,
			product_id,
			sku_id,
			quantity,
			checked,
			created_at,
			updated_at
		FROM cart
		WHERE id = ? AND user_id = ?
		LIMIT 1
	`
	err := util.Db.Raw(sql, cartId, userId).Scan(&cart).Error
	if err != nil {
		logger.Log.Error("查询购物车详情失败", zap.Error(err), zap.Uint("user_id", userId), zap.Uint("cart_id", cartId))
		return nil, err
	}
	if cart.Id == 0 {
		return nil, nil
	}
	//logger.Log.Debug("", zap.Any("cart", cart))
	return &cart, nil
}
func UpdateCartChecked(userId uint, cartId uint, checked int) error {
	sql := `UPDATE cart
		SET checked = ?, updated_at = NOW()
		WHERE id = ? AND user_id = ?`
	result := util.Db.Exec(sql, checked, cartId, userId)
	if result.Error != nil {
		logger.Log.Error("更新购物车勾选状态失败", zap.Error(result.Error), zap.Uint("cart_id", cartId), zap.Int("checked", checked))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("购物车项不存在或无权限修改")
	}
	return nil
}

func UpdateCartCheckedBatch(userId uint, cartIds []uint, checked int) error {
	sql := `UPDATE cart
		SET checked = ?, updated_at = NOW()
		WHERE user_id = ? AND id IN ?`
	err := util.Db.Exec(sql, checked, userId, cartIds).Error
	if err != nil {
		logger.Log.Error("批量更新购物车勾选状态失败", zap.Error(err), zap.Uint("user_id", userId), zap.Any("cart_ids", cartIds), zap.Int("checked", checked))
		return err
	}
	return nil
}
func DeleteCartById(cartId uint) error {
	sql := `DELETE FROM cart WHERE id = ?`
	err := util.Db.Exec(sql, cartId).Error
	if err != nil {
		logger.Log.Error("删除购物车失败", zap.Error(err), zap.Uint("cart_id", cartId))
		return err
	}
	return nil
}
func GetCartItemDetailById(userId uint, cartId uint) (*models.CartItem, error) {
	var item models.CartItem
	sql := `
		SELECT
			c.id AS cart_id,
			c.product_id,
			c.sku_id,
			c.quantity,
			c.checked,
			p.name AS product_name,
			p.main_image,
			p.status AS product_status,
			s.sku_name,
			s.price,
			s.stock,
			s.status AS sku_status
		FROM cart c
		LEFT JOIN product p ON c.product_id = p.id
		LEFT JOIN product_sku s ON c.sku_id = s.id
		WHERE c.user_id = ? AND c.id = ?
		LIMIT 1
	`
	err := util.Db.Raw(sql, userId, cartId).Scan(&item).Error
	if err != nil {
		logger.Log.Error("查询购物车商品详情失败", zap.Error(err), zap.Uint("user_id", userId), zap.Uint("cart_id", cartId))
		return nil, err
	}
	if item.CartId == 0 {
		return nil, nil
	}
	return &item, nil
}
