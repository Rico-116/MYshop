package dao

import (
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"go.uber.org/zap"
)

func GetCategoryList() ([]models.Category, error) {
	var category []models.Category
	sql := "SELECT id,name,parent_id,sort,icon,status,created_at,updated_at FROM category WHERE status=1 ORDER BY sort ASC,id ASC"
	err := util.Db.Raw(sql).Scan(&category).Error
	if err != nil {
		logger.Log.Error("查询数据库失败", zap.Error(err))
		return category, err
	}
	return category, nil
}
