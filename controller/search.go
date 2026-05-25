package controller

import (
	"MYshop/Service"
	"MYshop/models"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

func SearchProducts(c *gin.Context) {
	req := models.SearchProductRequest{
		Keyword: c.Query("keyword"),
		Sort:    c.Query("sort"),
	}

	if categoryStr := c.Query("category_id"); categoryStr != "" {
		categoryId, err := strconv.Atoi(categoryStr)
		if err != nil || categoryId < 0 {
			util.Fail(c, 400, "分类参数错误")
			return
		}
		value := uint(categoryId)
		req.CategoryId = &value
	}
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			util.Fail(c, 400, "最低价格参数错误")
			return
		}
		req.MinPrice = &minPrice
	}
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			util.Fail(c, 400, "最高价格参数错误")
			return
		}
		req.MaxPrice = &maxPrice
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	req.Page = page
	req.PageSize = pageSize

	result, err := Service.SearchProducts(req)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "搜索商品成功", result)
}
