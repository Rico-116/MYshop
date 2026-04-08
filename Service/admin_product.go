package Service

import (
	"MYshop/dao"
	"errors"
)

func ValidateProductCategory(categoryId int) error { //目前是一个半成品
	category, err := dao.GetCategoryById(categoryId)
	if err != nil {
		return err
	}
	if (category.Id == 0) || (category == nil) {
		return errors.New("分类不存在")
	}
	isLeaf, err := dao.IsLeafCategory(categoryId)
	if err != nil {
		return err
	}
	if !isLeaf {
		return errors.New("只能选择最小分类添加")
	}
	return nil
}
