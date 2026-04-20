package controller

import (
	"MYshop/Service"
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateOrder(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}

	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("创建订单参数错误", zap.Error(err))
		util.Fail(c, 400, "参数错误")
		return
	}

	result, err := Service.CreateOrder(userId, req)
	if err != nil {
		logger.Log.Warn("创建订单失败", zap.Error(err), zap.Uint("user_id", userId))
		util.Fail(c, 400, err.Error())
		return
	}

	util.Success(c, "创建订单成功", gin.H{
		"order": result,
	})
}
