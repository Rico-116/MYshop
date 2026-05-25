package Service

import (
	"MYshop/dao"
	"MYshop/models"
	"errors"
	"strings"
)

func SearchProducts(req models.SearchProductRequest) (*models.SearchProductResult, error) {
	req.Keyword = strings.TrimSpace(req.Keyword)
	req.Sort = strings.TrimSpace(req.Sort)

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 50 {
		req.PageSize = 50
	}
	if req.MinPrice != nil && *req.MinPrice < 0 {
		return nil, errors.New("最低价格不能小于0")
	}
	if req.MaxPrice != nil && *req.MaxPrice < 0 {
		return nil, errors.New("最高价格不能小于0")
	}
	if req.MinPrice != nil && req.MaxPrice != nil && *req.MinPrice > *req.MaxPrice {
		return nil, errors.New("最低价格不能大于最高价格")
	}

	switch req.Sort {
	case "", "price_asc", "price_desc", "hot", "newest":
	default:
		return nil, errors.New("排序参数错误")
	}

	list, total, err := dao.SearchProductList(req)
	if err != nil {
		return nil, err
	}

	return &models.SearchProductResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	}, nil
}
