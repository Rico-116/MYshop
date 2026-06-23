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
			p.subtitle AS product_name,
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
			p.subtitle AS product_name,
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
		 receiver_detail_address, remark, user_deleted, created_at, updated_at)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, NOW(), NOW())
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
		logger.Log.Error("创建订单主表失败",
			zap.Error(result.Error),
			zap.String("order_no", order.OrderNo),
			zap.Uint("user_id", order.UserId),
			zap.Float64("pay_amount", order.PayAmount),
		)
		return 0, result.Error
	}
	var orderId uint
	idSql := `SELECT id FROM orders WHERE order_no = ?`
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
		logger.Log.Error("创建订单明细失败",
			zap.Error(err),
			zap.Uint("order_id", item.OrderId),
			zap.String("order_no", item.OrderNo),
			zap.Uint("product_id", item.ProductId),
			zap.Uint("sku_id", item.SkuId),
		)
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
	sql := `DELETE FROM cart WHERE id IN ? AND user_id = ?`
	err := tx.Exec(sql, cartsIds, userId).Error
	if err != nil {
		logger.Log.Error("删除已下单购物车项失败", zap.Error(err), zap.Uint("user_id", userId))
		return err
	}
	return nil
}
func UpdateOrderToCanceledIfUnpaidTx(tx *gorm.DB, orderNo string) (int64, error) {
	sql := `UPDATE orders SET status=?,close_time=NOW(),updated_at=NOW() WHERE order_no=? AND status=?`
	res := tx.Exec(sql, 0, orderNo, 1)
	return res.RowsAffected, res.Error
}
func GetOrderItemsByOrderNo(orderNo string) ([]models.OrderItemVO, error) {
	var items []models.OrderItemVO
	sql := `
		SELECT
			order_no,
			product_id,
			sku_id,
			product_name,
			product_image,
			sku_name,
			price,
			quantity,
			total_amount
		FROM order_item
		WHERE order_no = ?
		ORDER BY id ASC
	`
	err := util.Db.Raw(sql, orderNo).Scan(&items).Error
	return items, err
}
func RestoreSkuStockTx(tx *gorm.DB, skuId uint, quantity int) error {
	sql := `UPDATE product_sku
		SET stock = stock + ?, updated_at = NOW()
		WHERE id = ?`
	return tx.Exec(sql, quantity, skuId).Error

}
func GetOrderListByUserId(userId uint, status int, page int, pageSize int) ([]models.OrderListVO, error) {
	var list []models.OrderListVO
	offset := (page - 1) * pageSize
	db := util.Db.Table("orders AS o").Select(
		"o.id",
		"o.order_no",
		"o.user_id",
		"o.status",
		"o.total_amount",
		"o.pay_amount",
		"o.freight_amount",
		"o.receiver_name",
		"o.receiver_phone",
		"o.receiver_province",
		"o.receiver_city",
		"o.receiver_district",
		"o.receiver_detail_address",
		"o.remark",
		"o.pay_time",
		"o.delivery_time",
		"o.finish_time",
		"o.close_time",
		"o.created_at",
		"o.updated_at",
		"CASE WHEN so.id IS NULL THEN false ELSE true END AS is_seckill",
		"COALESCE(so.activity_id, 0) AS seckill_activity_id",
	).Joins("LEFT JOIN seckill_order AS so ON so.order_no = o.order_no").Where("o.user_id = ? AND COALESCE(o.user_deleted, 0) = 0", userId)
	if status != models.OrderStatusAll {
		db = db.Where("o.status = ?", status)
	}
	err := db.Order("o.created_at desc").Offset(offset).Limit(pageSize).Scan(&list).Error
	if err != nil {
		logger.Log.Error("查询订单列表失败",
			zap.Error(err),
			zap.Uint("user_id", userId),
			zap.Int("status", status),
		)
		return nil, err
	}
	return list, nil
}

func CountOrderListByUserId(userId uint, status int) (int, error) {
	var total int64
	db := util.Db.Table("orders").Where("user_id = ? AND COALESCE(user_deleted, 0) = 0", userId)
	if status != models.OrderStatusAll {
		db = db.Where("status = ?", status)
	}
	err := db.Count(&total).Error
	if err != nil {
		logger.Log.Error("统计订单数量失败",
			zap.Error(err),
			zap.Uint("user_id", userId),
			zap.Int("status", status),
		)
		return 0, err
	}
	return int(total), nil
}

func GetMergedOrderListByUserId(userId uint, status int, page int, pageSize int) ([]models.OrderListVO, error) {
	var list []models.OrderListVO
	offset := (page - 1) * pageSize

	normalStatusWhere := ""
	seckillStatusWhere := ""
	args := []interface{}{userId}
	if status != models.OrderStatusAll {
		normalStatusWhere = " AND o.status = ?"
		seckillStatusWhere = " AND COALESCE(o.status, CASE WHEN so.status = 0 THEN 0 ELSE 1 END) = ?"
		args = append(args, status)
	}

	sql := `
		SELECT *
		FROM (
			SELECT
				o.id,
				o.order_no,
				o.user_id,
				o.status,
				o.total_amount,
				o.pay_amount,
				o.freight_amount,
				o.coupon_amount,
				o.receiver_name,
				o.receiver_phone,
				o.receiver_province,
				o.receiver_city,
				o.receiver_district,
				o.receiver_detail_address,
				o.remark,
				o.pay_time,
				o.delivery_time,
				o.finish_time,
				o.close_time,
				o.created_at,
				o.updated_at,
				0 AS is_seckill,
				0 AS seckill_activity_id
			FROM orders AS o
			WHERE o.user_id = ?
				AND COALESCE(o.user_deleted, 0) = 0
				AND NOT EXISTS (
					SELECT 1 FROM seckill_order AS sox WHERE sox.order_no = o.order_no
				)` + normalStatusWhere + `

			UNION ALL

			SELECT
				COALESCE(o.id, so.id) AS id,
				so.order_no,
				so.user_id,
				COALESCE(o.status, CASE WHEN so.status = 0 THEN 0 ELSE 1 END) AS status,
				COALESCE(o.total_amount, a.seckill_price, 0) AS total_amount,
				COALESCE(o.pay_amount, a.seckill_price, 0) AS pay_amount,
				COALESCE(o.freight_amount, 0) AS freight_amount,
				COALESCE(o.coupon_amount, 0) AS coupon_amount,
				COALESCE(o.receiver_name, '') AS receiver_name,
				COALESCE(o.receiver_phone, '') AS receiver_phone,
				COALESCE(o.receiver_province, '') AS receiver_province,
				COALESCE(o.receiver_city, '') AS receiver_city,
				COALESCE(o.receiver_district, '') AS receiver_district,
				COALESCE(o.receiver_detail_address, '') AS receiver_detail_address,
				COALESCE(o.remark, '秒杀订单') AS remark,
				o.pay_time,
				o.delivery_time,
				o.finish_time,
				o.close_time,
				COALESCE(o.created_at, so.created_at) AS created_at,
				COALESCE(o.updated_at, so.updated_at) AS updated_at,
				1 AS is_seckill,
				so.activity_id AS seckill_activity_id
			FROM seckill_order AS so
			LEFT JOIN orders AS o ON o.order_no = so.order_no
			LEFT JOIN seckill_activity AS a ON a.id = so.activity_id
			WHERE so.user_id = ?
				AND (o.id IS NULL OR COALESCE(o.user_deleted, 0) = 0)` + seckillStatusWhere + `
		) AS merged_orders
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	args = append(args, userId)
	if status != models.OrderStatusAll {
		args = append(args, status)
	}
	args = append(args, pageSize, offset)

	if err := util.Db.Raw(sql, args...).Scan(&list).Error; err != nil {
		logger.Log.Error("查询合并订单列表失败",
			zap.Error(err),
			zap.Uint("user_id", userId),
			zap.Int("status", status),
		)
		return nil, err
	}
	return list, nil
}

func CountMergedOrderListByUserId(userId uint, status int) (int, error) {
	var total int64

	normalStatusWhere := ""
	seckillStatusWhere := ""
	args := []interface{}{userId}
	if status != models.OrderStatusAll {
		normalStatusWhere = " AND o.status = ?"
		seckillStatusWhere = " AND COALESCE(o.status, CASE WHEN so.status = 0 THEN 0 ELSE 1 END) = ?"
		args = append(args, status)
	}

	sql := `
		SELECT COUNT(*) AS total
		FROM (
			SELECT o.order_no
			FROM orders AS o
			WHERE o.user_id = ?
				AND COALESCE(o.user_deleted, 0) = 0
				AND NOT EXISTS (
					SELECT 1 FROM seckill_order AS sox WHERE sox.order_no = o.order_no
				)` + normalStatusWhere + `

			UNION ALL

			SELECT so.order_no
			FROM seckill_order AS so
			LEFT JOIN orders AS o ON o.order_no = so.order_no
			WHERE so.user_id = ?
				AND (o.id IS NULL OR COALESCE(o.user_deleted, 0) = 0)` + seckillStatusWhere + `
		) AS merged_orders
	`

	args = append(args, userId)
	if status != models.OrderStatusAll {
		args = append(args, status)
	}

	if err := util.Db.Raw(sql, args...).Scan(&total).Error; err != nil {
		logger.Log.Error("统计合并订单数量失败",
			zap.Error(err),
			zap.Uint("user_id", userId),
			zap.Int("status", status),
		)
		return 0, err
	}
	return int(total), nil
}
func GetOrderItemsByOrderNos(orderNos []string) ([]models.OrderItemVO, error) {
	var list []models.OrderItemVO
	if len(orderNos) == 0 {
		return list, nil
	}

	err := util.Db.Table("order_item").
		Select(
			"order_no",
			"product_id",
			"sku_id",
			"product_name",
			"product_image",
			"sku_name",
			"price",
			"quantity",
			"total_amount",
		).
		Where("order_no IN ?", orderNos).
		Order("created_at desc").
		Find(&list).Error

	if err != nil {
		logger.Log.Error("批量查询订单商品失败",
			zap.Error(err),
			zap.Int("order_no_count", len(orderNos)))
		return nil, err
	}

	return list, nil
}

func GetOrderByOrderNoAndUserId(userId uint, orderNo string) (*models.Order, error) {
	var order models.Order
	db := util.Db.Table("orders")
	db = db.Select("id", "order_no", "user_id", "status", "total_amount", "pay_amount", "coupon_amount", "freight_amount", "receiver_name", "receiver_phone", "receiver_province", "receiver_city", "receiver_district", "receiver_detail_address", "remark", "pay_time", "delivery_time", "finish_time", "close_time", "user_deleted", "created_at", "updated_at").Where("user_id = ? AND order_no = ? AND COALESCE(user_deleted, 0) = 0", userId, orderNo)
	err := db.First(&order).Error
	if err != nil {
		logger.Log.Error("订单查询失败", zap.Error(err), zap.Uint("user_id", userId), zap.String("order_no", orderNo))
		return nil, err
	}
	if order.Id == 0 {
		return nil, nil
	}
	return &order, nil
}
func UpdateOrderToPaidTx(tx *gorm.DB, userId uint, orderNo string) (int64, error) {
	result := tx.Table("orders").Where("user_id = ? AND order_no = ? AND status=?", userId, orderNo, models.OrderStatusUnpaid).Updates(map[string]interface{}{"status": models.OrderStatusPaid, "pay_time": gorm.Expr("NOW()"), "updated_at": gorm.Expr("NOW()")})
	if result.Error != nil {
		logger.Log.Error("更新订单为已支付失败",
			zap.Error(result.Error),
			zap.Uint("user_id", userId),
			zap.String("order_no", orderNo),
		)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func AdminGetOrderList(status int, orderNo string, page, pageSize int) ([]models.OrderListVO, int64, error) {
	var list []models.OrderListVO
	var total int64

	db := util.Db.Table("orders AS o").Select(
		"o.id",
		"o.order_no",
		"o.user_id",
		"o.status",
		"o.total_amount",
		"o.pay_amount",
		"o.freight_amount",
		"o.coupon_amount",
		"o.receiver_name",
		"o.receiver_phone",
		"o.receiver_province",
		"o.receiver_city",
		"o.receiver_district",
		"o.receiver_detail_address",
		"o.remark",
		"o.pay_time",
		"o.delivery_time",
		"o.finish_time",
		"o.close_time",
		"o.created_at",
		"o.updated_at",
		"CASE WHEN so.id IS NULL THEN false ELSE true END AS is_seckill",
		"COALESCE(so.activity_id, 0) AS seckill_activity_id",
	).Joins("LEFT JOIN seckill_order AS so ON so.order_no = o.order_no")
	if status != models.OrderStatusAll {
		db = db.Where("o.status = ?", status)
	}
	if orderNo != "" {
		db = db.Where("o.order_no LIKE ?", "%"+orderNo+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := db.Order("o.created_at desc").Offset(offset).Limit(pageSize).Scan(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func AdminGetMergedOrderList(status int, orderNo string, page, pageSize int) ([]models.OrderListVO, int64, error) {
	var list []models.OrderListVO
	var total int64
	offset := (page - 1) * pageSize

	normalStatusWhere := ""
	seckillStatusWhere := ""
	normalOrderNoWhere := ""
	seckillOrderNoWhere := ""
	listArgs := make([]interface{}, 0)
	countArgs := make([]interface{}, 0)

	if status != models.OrderStatusAll {
		normalStatusWhere = " AND o.status = ?"
		seckillStatusWhere = " AND COALESCE(o.status, CASE WHEN so.status = 0 THEN 0 ELSE 1 END) = ?"
		listArgs = append(listArgs, status)
		countArgs = append(countArgs, status)
	}
	if orderNo != "" {
		normalOrderNoWhere = " AND o.order_no LIKE ?"
		seckillOrderNoWhere = " AND so.order_no LIKE ?"
		likeOrderNo := "%" + orderNo + "%"
		listArgs = append(listArgs, likeOrderNo)
		countArgs = append(countArgs, likeOrderNo)
	}

	listSQL := `
		SELECT *
		FROM (
			SELECT
				o.id,
				o.order_no,
				o.user_id,
				o.status,
				o.total_amount,
				o.pay_amount,
				o.freight_amount,
				o.coupon_amount,
				o.receiver_name,
				o.receiver_phone,
				o.receiver_province,
				o.receiver_city,
				o.receiver_district,
				o.receiver_detail_address,
				o.remark,
				o.pay_time,
				o.delivery_time,
				o.finish_time,
				o.close_time,
				o.created_at,
				o.updated_at,
				0 AS is_seckill,
				0 AS seckill_activity_id
			FROM orders AS o
			WHERE NOT EXISTS (
				SELECT 1 FROM seckill_order AS sox WHERE sox.order_no = o.order_no
			)` + normalStatusWhere + normalOrderNoWhere + `

			UNION ALL

			SELECT
				COALESCE(o.id, so.id) AS id,
				so.order_no,
				so.user_id,
				COALESCE(o.status, CASE WHEN so.status = 0 THEN 0 ELSE 1 END) AS status,
				COALESCE(o.total_amount, a.seckill_price, 0) AS total_amount,
				COALESCE(o.pay_amount, a.seckill_price, 0) AS pay_amount,
				COALESCE(o.freight_amount, 0) AS freight_amount,
				COALESCE(o.coupon_amount, 0) AS coupon_amount,
				COALESCE(o.receiver_name, '') AS receiver_name,
				COALESCE(o.receiver_phone, '') AS receiver_phone,
				COALESCE(o.receiver_province, '') AS receiver_province,
				COALESCE(o.receiver_city, '') AS receiver_city,
				COALESCE(o.receiver_district, '') AS receiver_district,
				COALESCE(o.receiver_detail_address, '') AS receiver_detail_address,
				COALESCE(o.remark, '秒杀订单') AS remark,
				o.pay_time,
				o.delivery_time,
				o.finish_time,
				o.close_time,
				COALESCE(o.created_at, so.created_at) AS created_at,
				COALESCE(o.updated_at, so.updated_at) AS updated_at,
				1 AS is_seckill,
				so.activity_id AS seckill_activity_id
			FROM seckill_order AS so
			LEFT JOIN orders AS o ON o.order_no = so.order_no
			LEFT JOIN seckill_activity AS a ON a.id = so.activity_id
			WHERE 1 = 1` + seckillStatusWhere + seckillOrderNoWhere + `
		) AS merged_orders
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	if status != models.OrderStatusAll {
		listArgs = append(listArgs, status)
	}
	if orderNo != "" {
		listArgs = append(listArgs, "%"+orderNo+"%")
	}
	listArgs = append(listArgs, pageSize, offset)

	if err := util.Db.Raw(listSQL, listArgs...).Scan(&list).Error; err != nil {
		return nil, 0, err
	}

	countSQL := `
		SELECT COUNT(*) AS total
		FROM (
			SELECT o.order_no
			FROM orders AS o
			WHERE NOT EXISTS (
				SELECT 1 FROM seckill_order AS sox WHERE sox.order_no = o.order_no
			)` + normalStatusWhere + normalOrderNoWhere + `

			UNION ALL

			SELECT so.order_no
			FROM seckill_order AS so
			LEFT JOIN orders AS o ON o.order_no = so.order_no
			WHERE 1 = 1` + seckillStatusWhere + seckillOrderNoWhere + `
		) AS merged_orders
	`

	if status != models.OrderStatusAll {
		countArgs = append(countArgs, status)
	}
	if orderNo != "" {
		countArgs = append(countArgs, "%"+orderNo+"%")
	}

	if err := util.Db.Raw(countSQL, countArgs...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func AdminGetOrderByOrderNo(orderNo string) (*models.Order, error) {
	var order models.Order
	err := util.Db.Table("orders").Where("order_no = ?", orderNo).Limit(1).Scan(&order).Error
	if err != nil {
		return nil, err
	}
	if order.Id == 0 {
		return nil, nil
	}
	return &order, nil
}

func AdminUpdateOrderToShipped(orderNo string) (int64, error) {
	result := util.Db.Table("orders").Where("order_no = ? AND status = ?", orderNo, models.OrderStatusPaid).Updates(map[string]interface{}{
		"status":        models.OrderStatusShipped,
		"delivery_time": gorm.Expr("NOW()"),
		"updated_at":    gorm.Expr("NOW()"),
	})
	return result.RowsAffected, result.Error
}

func AdminUpdateOrderToCanceledTx(tx *gorm.DB, orderNo string) (int64, error) {
	result := tx.Table("orders").
		Where("order_no = ? AND status IN ?", orderNo, []int{models.OrderStatusUnpaid, models.OrderStatusPaid}).
		Updates(map[string]interface{}{
			"status":     models.OrderStatusCanceled,
			"close_time": gorm.Expr("NOW()"),
			"updated_at": gorm.Expr("NOW()"),
		})
	return result.RowsAffected, result.Error
}
func UpdateOrderUserDeleted(userId uint, orderNo string) (int64, error) {
	result := util.Db.Table("orders").
		Where("user_id = ? AND order_no = ? AND COALESCE(user_deleted, 0) = 0 AND status IN ?", userId, orderNo, []int{models.OrderStatusCanceled, models.OrderStatusFinished}).
		Updates(map[string]interface{}{
			"user_deleted": 1,
			"updated_at":   gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		logger.Log.Error("用户删除订单失败",
			zap.Error(result.Error),
			zap.Uint("user_id", userId),
			zap.String("order_no", orderNo),
		)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func UpdateUserOrderToCanceledTx(tx *gorm.DB, userId uint, orderNo string) (int64, error) {
	result := tx.Table("orders").
		Where("user_id = ? AND order_no = ? AND COALESCE(user_deleted, 0) = 0 AND status IN ?", userId, orderNo, []int{models.OrderStatusUnpaid, models.OrderStatusPaid}).
		Updates(map[string]interface{}{
			"status":     models.OrderStatusCanceled,
			"close_time": gorm.Expr("NOW()"),
			"updated_at": gorm.Expr("NOW()"),
		})
	return result.RowsAffected, result.Error
}

func UpdateOrderToFinished(userId uint, orderNo string) (int64, error) {
	result := util.Db.Table("orders").
		Where("user_id = ? AND order_no = ? AND status = ? AND COALESCE(user_deleted, 0) = 0", userId, orderNo, models.OrderStatusShipped).
		Updates(map[string]interface{}{
			"status":      models.OrderStatusFinished,
			"finish_time": gorm.Expr("NOW()"),
			"updated_at":  gorm.Expr("NOW()"),
		})
	return result.RowsAffected, result.Error
}

func GetOrderItemByOrderProduct(userId uint, orderNo string, productId uint, skuId uint) (*models.OrderItem, error) {
	var item models.OrderItem
	db := util.Db.Table("order_item").
		Where("user_id = ? AND order_no = ? AND product_id = ?", userId, orderNo, productId)
	if skuId > 0 {
		db = db.Where("sku_id = ?", skuId)
	}
	err := db.Order("id asc").Limit(1).Scan(&item).Error
	if err != nil {
		return nil, err
	}
	if item.Id == 0 {
		return nil, nil
	}
	return &item, nil
}
