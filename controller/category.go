package controller

import (
	"MYshop/Service"
	"MYshop/package/logger"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func GetCategoryTree(c *gin.Context) {
	list, err := Service.GetCategoryTree()
	if err != nil {
		logger.Log.Warn("获取分类树失败")
		util.Fail(c, 500, "获取分类树失败")
		return
	}
	logger.Log.Info("获取分类成功", zap.Any("list", list))
	util.Success(c, "获取分类成功", gin.H{"list": list})
}
func GetCategoryDisplay(c *gin.Context) {
	categoryStr := c.Query("category_id")
	categoryId, err := strconv.Atoi(categoryStr)
	if err != nil || categoryId <= 0 {
		logger.Log.Warn("分类id参数错误", zap.Error(err))
		util.Fail(c, 400, "分类id参数错误")
		return
	}
	data, err := Service.GetCategoryDisplay(categoryId)
	if err != nil {
		logger.Log.Warn("获取分类展示失败", zap.Error(err))
		util.Fail(c, 500, "获取分类展示失败")
		return
	}
	util.Success(c, "获取分类展示成功", gin.H{
		"current_category": data.CurrentCategory,
		"sub_categories":   data.SubCategories,
		"product_list":     data.ProductList,
		"is_leaf":          data.IsLeaf,
	})
}
