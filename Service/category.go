package Service

import (
	"MYshop/dao"
	"MYshop/models"
)

func GetCategoryTree() ([]models.CategoryTree, error) {
	list, err := dao.GetCategoryList()
	if err != nil {
		return nil, err
	}
	var result []models.CategoryTree
	childrenMap := make(map[uint][]models.CategoryTree)
	for _, item := range list {
		if item.ParentId != 0 {
			childrenMap[item.ParentId] = append(childrenMap[item.ParentId], models.CategoryTree{
				Id:   item.Id,
				Name: item.Name,
			})
		}
	}
	for _, item := range list {
		if item.ParentId == 0 {
			result = append(result, models.CategoryTree{
				Id:       item.Id,
				Icon:     item.Icon,
				Name:     item.Name,
				Children: childrenMap[item.Id],
			})
		}
	}
	return result, nil
}
