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
func GetCategoryById(id int) (*models.Category, error) {
	var category models.Category
	sql := "select id,name,parent_id,sort,icon,status,created_at FROM category WHERE id=? AND status=1 LIMIT 1"
	err := util.Db.Raw(sql, id).Scan(&category).Error
	if err != nil {
		logger.Log.Error("查询分类失败", zap.Error(err))
		return nil, err
	}
	return &category, nil
}
func GetChildCategoryList(parentId int) ([]models.Category, error) {
	var list []models.Category
	sql := `
		SELECT
			id,
			name,
			parent_id,
			sort,
			status,
			icon,
			created_at,
			updated_at
		FROM category
		WHERE parent_id = ? AND status = 1
		ORDER BY sort ASC, id ASC
	`
	err := util.Db.Raw(sql, parentId).Scan(&list).Error
	if err != nil {
		logger.Log.Error("查询子分类失败", zap.Error(err))
		return nil, err
	}
	return list, nil
}
func CountChildCategory(parentId int) (int, error) {
	var count int
	sql := "SELECT COUNT(1) FROM category WHERE parent_id = ? AND status = 1"
	err := util.Db.Raw(sql, parentId).Scan(&count).Error
	if err != nil {
		logger.Log.Error("统计子分类的数量失败", zap.Error(err))
		return 0, err
	}
	return count, nil
}
func IsLeafCategory(categoryId int) (bool, error) {
	count, err := CountChildCategory(categoryId)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
func GetSkuCategoryById(id uint) (*models.Category, error) {
	sql := `
		SELECT
			id,
			name,
			parent_id,
			sort,
			status,
			icon,
			created_at,
			updated_at
		FROM category
		where id=? AND status=1 LIMIT 1
	`
	var category models.Category
	err := util.Db.Raw(sql, id).Scan(&category).Error
	if err != nil {
		logger.Log.Error("获取商品类别", zap.Error(err))
		return nil, err
	}
	return &category, nil
}
