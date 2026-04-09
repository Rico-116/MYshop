package dao

import (
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"go.uber.org/zap"
)

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
	err := util.Db.Exec(sql,
		cart.UserId,
		cart.ProductId,
		cart.SkuId,
		cart.Quantity,
		cart.Checked).Error
	if err != nil {
		logger.Log.Error("新增购物车失败", zap.Error(err), zap.Any("cart", cart))
		return err
	}
	return nil
}

func UpdateCartQuantity(cartId uint, quantity int) error {
	sql := `UPDATE cart
		SET quantity = ?, updated_at = NOW()
		WHERE id = ?`
	err := util.Db.Exec(sql, quantity, cartId).Error
	if err != nil {
		logger.Log.Error("更新购物车数量失败", zap.Error(err), zap.Uint("cart_id", cartId), zap.Int("quantity", quantity))
		return err
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
		WHERE c.user_id = ?
		ORDER BY c.id DESC`
	err := util.Db.Raw(sql, userId).Scan(&list).Error
	if err != nil {
		logger.Log.Error("查询购物车列表失败", zap.Error(err), zap.Uint("user_id", userId))
		return nil, err
	}
	return list, nil
}
func GetCartById(cartId, userId uint) (*models.Cart, error) {
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
	return &cart, nil
}
func UpdateCartChecked(cartId uint, checked int) error {
	sql := `UPDATE cart
		SET checked = ?, updated_at = NOW()
		WHERE id = ?`
	err := util.Db.Exec(sql, checked, cartId).Error
	if err != nil {
		logger.Log.Error("更新购物车勾选状态失败", zap.Error(err), zap.Uint("cart_id", cartId), zap.Int("checked", checked))
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
