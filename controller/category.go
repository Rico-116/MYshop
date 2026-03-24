package controller

import (
	"MYshop/Service"
	"MYshop/package/logger"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
