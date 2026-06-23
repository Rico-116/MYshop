package controller

import (
	"MYshop/Service"
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ConfirmReceive(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}

	var req models.ConfirmReceiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("确认收货参数错误",
			zap.Error(err),
			zap.Uint("user_id", userId),
		)
		util.Fail(c, 400, "参数错误")
		return
	}

	result, err := Service.ConfirmReceive(userId, req)
	if err != nil {
		logger.Log.Warn("确认收货失败",
			zap.Error(err),
			zap.Uint("user_id", userId),
			zap.String("order_no", req.OrderNo),
		)
		util.Fail(c, 400, err.Error())
		return
	}

	util.Success(c, "确认收货成功", result)
}
