package dao

import (
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CheckProductSkuExist 校验商品SKU是否存在
func CheckProductSkuExist(productId uint, skuId uint) error {
	var count int64
	err := util.Db.Table("product_sku").
		Where("id = ? AND product_id = ?", skuId, productId).
		Count(&count).Error

	if err != nil {
		logger.Log.Error("校验商品SKU失败",
			zap.Error(err),
			zap.Uint("product_id", productId),
			zap.Uint("sku_id", skuId),
		)
		return err
	}

	if count == 0 {
		return errors.New("商品SKU不存在或不属于该商品")
	}

	return nil
}

func CreateSeckillActivity(activity *models.SeckillActivity) error {
	activity.SoldCount = 0
	if err := util.Db.Table("seckill_activity").Create(activity).Error; err != nil {
		logger.Log.Error("创建秒杀活动失败",
			zap.Error(err),
			zap.Uint("product_id", activity.ProductId),
			zap.Uint("sku_id", activity.SkuId),
			zap.Int("stock", activity.Stock),
		)
		return err
	}
	return nil
}

func GetSeckillActivityById(activityId uint) (*models.SeckillActivity, error) {
	var activity models.SeckillActivity
	err := util.Db.Table("seckill_activity").Where("id = ? ", activityId).First(&activity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Log.Error("查询秒杀活动失败", zap.Error(err), zap.Uint("activity_id", activityId))
		return nil, err
	}
	return &activity, nil
}

func GetUnfinishedSeckillActivities() ([]models.SeckillActivity, error) {
	var activities []models.SeckillActivity
	err := util.Db.Table("seckill_activity").
		Where("status = ? AND end_time > NOW()", models.SeckillActivityStatusOn).
		Find(&activities).Error
	if err != nil {
		logger.Log.Error("查询未结束秒杀活动失败", zap.Error(err))
		return nil, err
	}
	return activities, nil
}

func GetAdminSeckillActivityList(page int, pageSize int, status int) ([]models.AdminSeckillActivityVO, int64, error) {
	var list []models.AdminSeckillActivityVO
	var total int64
	query := util.Db.Table("seckill_activity AS a").Joins("LEFT JOIN product As p ON p.id = a.product_id").
		Joins("LEFT JOIN product_sku AS s ON s.id = a.sku_id")
	switch status {
	case models.AdminSeckillListStatusClosed:
		query = query.Where("a.status = ?", models.SeckillActivityStatusOff)
	case models.AdminSeckillListStatusNotStarted:
		query = query.Where("a.status = ? AND a.start_time > NOW()", models.SeckillActivityStatusOn)
	case models.AdminSeckillListStatusRunning:
		query = query.Where("a.status = ? AND a.start_time <= NOW() AND a.end_time > NOW()", models.SeckillActivityStatusOn)
	case models.AdminSeckillListStatusEnded:
		query = query.Where("a.status = ? AND a.end_time <= NOW()", models.SeckillActivityStatusOn)
	}
	if err := query.Count(&total).Error; err != nil {
		logger.Log.Error("统计秒杀活动数量失败", zap.Error(err))
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Select(`a.id,
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
			DATE_FORMAT(a.end_time, '%Y-%m-%d %H:%i:%s') AS end_time,
			DATE_FORMAT(a.created_at, '%Y-%m-%d %H:%i:%s') AS created_at`).Order("a.id desc").Limit(pageSize).Offset(offset).Scan(&list).Error
	if err != nil {
		logger.Log.Error("查询秒杀活动列表失败", zap.Error(err))
		return nil, 0, err
	}
	return list, total, nil
}

func CloseSeckillActivity(activityId uint) error {
	err := util.Db.Table("seckill_activity").
		Where("id = ?", activityId).
		Updates(map[string]interface{}{
			"status":     models.SeckillActivityStatusOff,
			"updated_at": gorm.Expr("NOW()"),
		}).Error

	if err != nil {
		logger.Log.Error("关闭秒杀活动失败",
			zap.Error(err),
			zap.Uint("activity_id", activityId),
		)
		return err
	}

	return nil
}

func UpdateSeckillActivity(activity *models.SeckillActivity) error {
	sql := `
		UPDATE seckill_activity
		SET
			seckill_price = ?,
			stock = ?,
			start_time = ?,
			end_time = ?,
			status = ?,
			updated_at = NOW()
		WHERE id = ?
	`

	err := util.Db.Exec(
		sql,
		activity.SeckillPrice,
		activity.Stock,
		activity.StartTime,
		activity.EndTime,
		activity.Status,
		activity.Id,
	).Error

	if err != nil {
		logger.Log.Error("修改秒杀活动失败",
			zap.Error(err),
			zap.Uint("activity_id", activity.Id),
		)
		return err
	}

	return nil
}
