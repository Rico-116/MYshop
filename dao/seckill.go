package dao

import (
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"errors"

	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ErrSeckillOrderDuplicated = errors.New("seckill order already exists")

func GetUserSeckillActivityList(page int, pageSize int) ([]models.UserSeckillActivityVO, int64, error) {
	var list []models.UserSeckillActivityVO
	var total int64

	db := util.Db.Table("seckill_activity AS a").
		Joins("LEFT JOIN product AS p ON p.id = a.product_id").
		Joins("LEFT JOIN product_sku AS s ON s.id = a.sku_id").
		Where("a.status = ? AND a.end_time > NOW()", models.SeckillActivityStatusOn)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	err := db.Select(`
		a.id,
		a.product_id,
		a.sku_id,
		p.subtitle AS product_name,
		p.main_image AS product_image,
		s.sku_name,
		s.price AS original_price,
		a.seckill_price,
		a.stock,
		a.sold_count,
		(a.stock - a.sold_count) AS remain_stock,
		a.status,
		DATE_FORMAT(a.start_time, '%Y-%m-%d %H:%i:%s') AS start_time,
		DATE_FORMAT(a.end_time, '%Y-%m-%d %H:%i:%s') AS end_time
	`).
		Order("a.start_time ASC, a.id DESC").
		Limit(pageSize).
		Offset(offset).
		Scan(&list).Error

	return list, total, err
}

func GetSeckillOrderBuyItem(activityId uint) (*models.SeckillOrderBuyItem, error) {
	var item models.SeckillOrderBuyItem

	err := util.Db.Table("seckill_activity AS a").
		Select(`
			a.id AS activity_id,
			a.product_id,
			a.sku_id,
			p.name AS product_name,
			p.main_image AS product_image,
			p.status AS product_status,
			s.sku_name,
			s.status AS sku_status,
			a.seckill_price AS price,
			s.stock
		`).
		Joins("LEFT JOIN product AS p ON p.id = a.product_id").
		Joins("LEFT JOIN product_sku AS s ON s.id = a.sku_id").
		Where("a.id = ?", activityId).
		Limit(1).
		Scan(&item).Error

	if err != nil {
		logger.Log.Error("查询秒杀下单商品失败",
			zap.Error(err),
			zap.Uint("activity_id", activityId),
		)
		return nil, err
	}

	if item.ActivityId == 0 {
		return nil, nil
	}

	return &item, nil
}
func CheckSeckillOrderExists(activityId uint, userId uint) (bool, string, error) {
	var order models.SeckillOrder

	err := util.Db.Table("seckill_order AS so").
		Select("so.id, so.order_no").
		Joins("LEFT JOIN orders AS o ON o.order_no = so.order_no").
		Where("so.activity_id = ? AND so.user_id = ? AND so.status = 1", activityId, userId).
		Where("o.id IS NULL OR o.status <> ?", models.OrderStatusCanceled).
		Limit(1).
		Scan(&order).Error

	if err != nil {
		return false, "", err
	}

	return order.Id != 0, order.OrderNo, nil
}

func GetSeckillOrderByOrderNo(orderNo string) (*models.SeckillOrder, error) {
	var order models.SeckillOrder

	err := util.Db.Table("seckill_order").
		Select("id, activity_id, user_id, order_no, status, created_at, updated_at").
		Where("order_no = ?", orderNo).
		Limit(1).
		Scan(&order).Error
	if err != nil {
		return nil, err
	}
	if order.Id == 0 {
		return nil, nil
	}
	return &order, nil
}

func CancelSeckillOrderAndRestoreActivityStockTx(tx *gorm.DB, orderNo string) (*models.SeckillOrder, error) {
	var order models.SeckillOrder
	err := tx.Table("seckill_order").
		Select("id, activity_id, user_id, order_no, status").
		Where("order_no = ? AND status = 1", orderNo).
		Limit(1).
		Scan(&order).Error
	if err != nil {
		return nil, err
	}
	if order.Id == 0 {
		return nil, nil
	}

	res := tx.Table("seckill_order").
		Where("id = ? AND status = 1", order.Id).
		Updates(map[string]interface{}{
			"status":     0,
			"updated_at": gorm.Expr("NOW()"),
		})
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}

	res = tx.Exec(`
		UPDATE seckill_activity
		SET sold_count = sold_count - 1, updated_at = NOW()
		WHERE id = ? AND sold_count > 0
	`, order.ActivityId)
	if res.Error != nil {
		return nil, res.Error
	}

	return &order, nil
}

func CreateSeckillOrderTx(tx *gorm.DB, activityId uint, userId uint, orderNo string) error {
	sql := `
		INSERT INTO seckill_order(activity_id, user_id, order_no, status, created_at, updated_at)
		VALUES (?, ?, ?, 1, NOW(), NOW())
	`
	err := tx.Exec(sql, activityId, userId, orderNo).Error
	if err == nil {
		return nil
	}

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return ErrSeckillOrderDuplicated
	}

	return err
}
func DeductSeckillActivityStockTx(tx *gorm.DB, activityId uint) error {
	sql := `
		UPDATE seckill_activity
		SET sold_count = sold_count + 1, updated_at = NOW()
		WHERE id = ? AND status = ? AND sold_count < stock
	`

	res := tx.Exec(sql, activityId, models.SeckillActivityStatusOn)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("秒杀活动库存不足")
	}

	return nil
}
