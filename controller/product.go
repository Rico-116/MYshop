package controller

import (
	"MYshop/dao"
	"MYshop/package/logger"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetProductList(c *gin.Context) {
	list, err := dao.GetProductList()
	if err != nil {
		logger.Log.Error("获取商品列表失败", zap.Error(err))
		util.Fail(c, 500, "获取商品列表失败")
		return
	}
	logger.Log.Info("获取商品列表成功", zap.Any("list", list))
	util.Success(c, "获取商品列表成功", gin.H{"list": list})

}
