package dao

import (
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func GetCartOrderItemsByIds(userId uint, cartIds []uint) ([]models.OrderBuyItem, error) {
	var list []models.OrderBuyItem
	sql := `
		SELECT
			c.id AS cart_id,
			c.product_id,
			c.sku_id,
			c.quantity,
			p.name AS product_name,
			p.main_image AS product_image,
			p.status AS product_status,
			s.sku_name,
			s.price,
			s.stock,
			s.status AS sku_status
		FROM cart c
		LEFT JOIN product p ON c.product_id = p.id
		LEFT JOIN product_sku s ON c.sku_id = s.id
		WHERE c.user_id = ? AND c.id IN ?
		ORDER BY c.id DESC
	`
	err := util.Db.Raw(sql, userId, cartIds).Scan(&list).Error
	if err != nil {
		logger.Log.Error("查询购物车下单商品失败", zap.Error(err), zap.Uint("user_id", userId))
		return nil, err
	}
	return list, nil
}
func GetDirectOrderItemBySkuId(skuId uint) (*models.OrderBuyItem, error) {
	var item models.OrderBuyItem
	sql := `
		SELECT
			p.id AS product_id,
			s.id AS sku_id,
			p.name AS product_name,
			p.main_image AS product_image,
			p.status AS product_status,
			s.sku_name,
			s.price,
			s.stock,
			s.status AS sku_status
		FROM product_sku s
		LEFT JOIN product p ON s.product_id = p.id
		WHERE s.id = ?
		LIMIT 1
	`
	err := util.Db.Raw(sql, skuId).Scan(&item).Error
	if err != nil {
		logger.Log.Error("查询直接购买商品失败", zap.Error(err), zap.Uint("sku_id", skuId))
		return nil, err
	}
	if item.SkuId == 0 {
		return nil, nil
	}
	return &item, nil
}
func CreateOrderTx(tx *gorm.DB, order *models.Order) (uint, error) {
	sql := `
		INSERT INTO orders
		(order_no, user_id, status, total_amount, pay_amount, coupon_amount, freight_amount,
		 receiver_name, receiver_phone, receiver_province, receiver_city, receiver_district,
		 receiver_detail_address, remark, created_at, updated_at)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	result := tx.Exec(sql,
		order.OrderNo,
		order.UserId,
		order.Status,
		order.TotalAmount,
		order.PayAmount,
		order.CouponAmount,
		order.FreightAmount,
		order.ReceiverName,
		order.ReceiverPhone,
		order.ReceiverProvince,
		order.ReceiverCity,
		order.ReceiverDistrict,
		order.ReceiverDetailAddress,
		order.Remark,
	)
	if result.Error != nil {
		logger.Log.Error("创建订单主表失败", zap.Error(result.Error), zap.Any("order", order))
		return 0, result.Error
	}
	var orderId uint
	idSql := `SELECT order_id FROM orders WHERE order_no = ?`
	if err := tx.Raw(idSql, order.OrderNo).Scan(&orderId).Error; err != nil {
		logger.Log.Error("查询订单ID失败", zap.Error(err), zap.String("order_no", order.OrderNo))
		return 0, err
	}
	if orderId == 0 {
		return 0, errors.New("创建订单失败")
	}
	return orderId, nil
}
func CreateOrderItemTx(tx *gorm.DB, item *models.OrderItem) error {
	sql := `
		INSERT INTO order_item
		(order_id, order_no, user_id, product_id, sku_id, product_name, product_image, sku_name,
		 price, quantity, total_amount, created_at, updated_at)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	err := tx.Exec(sql,
		item.OrderId,
		item.OrderNo,
		item.UserId,
		item.ProductId,
		item.SkuId,
		item.ProductName,
		item.ProductImage,
		item.SkuName,
		item.Price,
		item.Quantity,
		item.TotalAmount,
	).Error
	if err != nil {
		logger.Log.Error("创建订单明细失败", zap.Error(err), zap.Any("item", item))
		return err
	}
	return nil
}
func DeductSkuStockTx(tx *gorm.DB, skuId uint, quantity int) error {
	sql := `UPDATE product_sku SET stock=stock-?,updated_at=NOW() WHERE id = ? AND status=1 AND stock>=?`
	result := tx.Exec(sql, quantity, skuId, quantity)
	if result.Error != nil {
		logger.Log.Error("扣减库存失败", zap.Uint("sku_id", skuId), zap.Int("quantity", quantity))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("商品库存不足")
	}
	return nil
}
func DeleteCartItemByIdsTx(tx *gorm.DB, userId uint, cartsIds []uint) error {
	sql := `DELETE FROM cart WHERE id = ? AND user_id = ?`
	err := tx.Exec(sql, cartsIds, userId).Error
	if err != nil {
		logger.Log.Error("删除已下单购物车项失败", zap.Error(err), zap.Uint("user_id", userId))
		return err
	}
	return nil
}
