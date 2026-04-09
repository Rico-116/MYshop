package controller

import (
	"MYshop/Service"
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddCart(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}
	var req models.AddCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("加入购物车参数错误", zap.Error(err))
		util.Fail(c, 400, "参数错误")
		return
	}
	err := Service.AddCart(userId, req)
	if err != nil {
		logger.Log.Error("加入购物车失败", zap.Error(err), zap.Error(err), zap.Uint("user_id", userId), zap.Uint("sku_id", req.SkuId))
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "加入购物车成功", nil)
}
