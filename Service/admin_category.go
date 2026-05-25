package Service

import (
	"MYshop/dao"
	"MYshop/models"
	"errors"
	"gorm.io/gorm"
	"strings"
	"time"
)

func GetAdminCategoryList(status *uint, level int) ([]models.Category, error) {
	if status != nil && *status != 0 && *status != 1 {
		return nil, errors.New("分类状态只能是0或1")
	}
	if level != 0 && level != 1 && level != 2 {
		return nil, errors.New("分类层级参数错误")
	}
	return dao.AdminGetCategoryList(status, level)
}

func CreateAdminCategory(req models.AdminCreateCategoryRequest) (*models.Category, error) {
	name := strings.TrimSpace(req.Name)
	icon := strings.TrimSpace(req.Icon)
	if name == "" {
		return nil, errors.New("分类名称不能为空")
	}
	if req.ParentId == 0 {
		if icon == "" {
			return nil, errors.New("一级分类必须上传图标")
		}
	} else {
		parent, err := dao.AdminGetCategoryById(req.ParentId)
		if err != nil {
			return nil, err
		}
		if parent == nil || parent.Status != 1 {
			return nil, errors.New("父分类不存在或已禁用")
		}
		if parent.ParentId != 0 {
			return nil, errors.New("只能创建二级分类，不能创建三级分类")
		}
		icon = ""
	}
	status := uint(1)
	if req.Status != nil {
		status = *req.Status
	}
	if status != 0 && status != 1 {
		return nil, errors.New("分类状态只能是0或1")
	}
	category := &models.Category{
		Name:      name,
		ParentId:  req.ParentId,
		Sort:      req.Sort,
		Icon:      icon,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := dao.AdminCreateCategory(category); err != nil {
		return nil, err
	}
	return category, nil
}

func UpdateAdminCategory(categoryId uint, req models.AdminUpdateCategoryRequest) (*models.Category, error) {
	if categoryId == 0 {
		return nil, errors.New("分类id不能为空")
	}
	category, err := dao.AdminGetCategoryById(categoryId)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("分类不存在")
	}

	parentId := category.ParentId
	icon := strings.TrimSpace(category.Icon)

	updates := make(map[string]interface{})
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.New("分类名称不能为空")
		}
		updates["name"] = name
	}
	if req.ParentId != nil {
		if *req.ParentId == categoryId {
			return nil, errors.New("父分类不能选择自己")
		}
		if *req.ParentId > 0 {
			parent, err := dao.AdminGetCategoryById(*req.ParentId)
			if err != nil {
				return nil, err
			}
			if parent == nil || parent.Status != 1 {
				return nil, errors.New("父分类不存在或已禁用")
			}
			if parent.ParentId != 0 {
				return nil, errors.New("只能选择一级分类作为父分类")
			}
		}
		parentId = *req.ParentId
		updates["parent_id"] = *req.ParentId
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.Icon != nil {
		icon = strings.TrimSpace(*req.Icon)
		updates["icon"] = icon
	}
	if req.Status != nil {
		if *req.Status != 0 && *req.Status != 1 {
			return nil, errors.New("分类状态只能是0或1")
		}
		updates["status"] = *req.Status
	}
	if len(updates) == 0 {
		return nil, errors.New("没有需要修改的分类信息")
	}
	if parentId == 0 {
		if icon == "" {
			return nil, errors.New("一级分类必须上传图标")
		}
	} else {
		updates["icon"] = ""
	}
	updates["updated_at"] = gorm.Expr("NOW()")
	if err := dao.AdminUpdateCategory(categoryId, updates); err != nil {
		return nil, err
	}
	return dao.AdminGetCategoryById(categoryId)
}

func DeleteAdminCategory(categoryId uint) error {
	if categoryId == 0 {
		return errors.New("分类id不能为空")
	}
	category, err := dao.AdminGetCategoryById(categoryId)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("分类不存在")
	}
	childCount, err := dao.CountChildCategory(int(categoryId))
	if err != nil {
		return err
	}
	if childCount > 0 {
		return errors.New("该分类下还有子分类，不能删除")
	}
	productCount, err := dao.AdminCountProductByCategoryId(categoryId)
	if err != nil {
		return err
	}
	if productCount > 0 {
		return errors.New("该分类下还有上架商品，不能删除")
	}
	return dao.AdminDeleteCategory(categoryId)
}
