package controller

import (
	"MYshop/dao"
	"MYshop/package/logger"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetBannerList(c *gin.Context) {
	list, err := dao.GetBannerList()
	if err != nil {
		logger.Log.Warn("获取轮播图失败", zap.Error(err))
		util.Fail(c, 500, "轮播图获取失败")
		return
	}
	logger.Log.Info("轮播图获取成功", zap.Any("list", list))
	util.Success(c, "获取轮播图成功", gin.H{"list": list})
}
