package dao

import (
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

func AdminGetCategoryById(id uint) (*models.Category, error) {
	var category models.Category
	err := util.Db.Table("category").
		Select("id", "name", "parent_id", "sort", "icon", "status", "created_at", "updated_at").
		Where("id = ?", id).
		Limit(1).
		Scan(&category).Error
	if err != nil {
		return nil, err
	}
	if category.Id == 0 {
		return nil, nil
	}
	return &category, nil
}

func AdminGetCategoryList(status *uint, level int) ([]models.Category, error) {
	var list []models.Category
	db := util.Db.Table("category").Select("id", "name", "parent_id", "sort", "icon", "status", "created_at", "updated_at")
	if status != nil {
		db = db.Where("status = ?", *status)
	}
	if level == 1 {
		db = db.Where("parent_id = 0")
	}
	if level == 2 {
		db = db.Where("parent_id <> 0")
	}
	err := db.Order("sort ASC, id ASC").Scan(&list).Error
	return list, err
}

func AdminCreateCategory(category *models.Category) error {
	return util.Db.Table("category").Create(category).Error
}

func AdminUpdateCategory(categoryId uint, updates map[string]interface{}) error {
	result := util.Db.Table("category").Where("id = ?", categoryId).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("分类不存在或没有任何修改")
	}
	return nil
}

func AdminDeleteCategory(categoryId uint) error {
	result := util.Db.Table("category").Where("id = ?", categoryId).Updates(map[string]interface{}{
		"status":     0,
		"updated_at": gorm.Expr("NOW()"),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("分类不存在")
	}
	return nil
}

func AdminCountProductByCategoryId(categoryId uint) (int64, error) {
	var count int64
	err := util.Db.Table("product").Where("category_id = ? AND status = 1", categoryId).Count(&count).Error
	return count, err
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
