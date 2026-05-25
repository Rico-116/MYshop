package controller

import (
	"MYshop/Service"
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
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
	logger.Log.Info("创建订单成功",
		zap.Uint("user_id", userId),
		zap.String("order_no", result.OrderNo),
		zap.Uint("order_id", result.OrderId),
	)
	util.Success(c, "创建订单成功", gin.H{
		"order": result,
	})
}
func GetOrderList(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}
	statusStr := c.DefaultQuery("status", "-1")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		util.Fail(c, 400, "订单状态参数错误")
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		util.Fail(c, 400, "分页参数错误")
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		util.Fail(c, 400, "分页大小参数错误")
		return
	}
	req := models.OrderListRequest{
		Status:   status,
		Page:     page,
		PageSize: pageSize,
	}
	result, err := Service.GetOrderList(userId, req)
	if err != nil {
		logger.Log.Warn("查询订单列表失败",
			zap.Error(err),
			zap.Uint("user_id", userId),
			zap.Int("status", status),
		)
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "查询订单列表成功", result)
}
func PreviewOrder(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}
	var req models.OrderPreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("订单确认页参数错误", zap.Error(err))
		util.Fail(c, 400, "参数错误")
		return
	}
	result, err := Service.PreviewOrder(userId, req)
	if err != nil {
		logger.Log.Warn("生成订单确认页面失败", zap.Error(err), zap.Uint("user_id", userId))
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "获取订单确认页成功", result)
}
func GetPayPage(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}
	OrderNo := c.Query("order_no")
	if OrderNo == "" {
		util.Fail(c, 400, "订单号不能为空")
		return
	}
	result, err := Service.GetPayPage(userId, OrderNo)
	if err != nil {
		logger.Log.Warn("获取支付页失败",
			zap.Error(err),
			zap.Uint("user_id", userId),
			zap.String("order_no", OrderNo),
		)
		util.Fail(c, 400, err.Error())
		return
	}

	util.Success(c, "获取支付页成功", result)
}
func PayOrder(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}
	var req models.PayOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("支付订单参数错误", zap.Error(err))
		util.Fail(c, 400, "参数错误")
		return
	}
	result, err := Service.PayOrder(userId, req)
	if err != nil {
		logger.Log.Warn("支付订单失败",
			zap.Error(err),
			zap.Uint("user_id", userId),
			zap.String("order_no", req.OrderNO),
		)
		util.Fail(c, 400, err.Error())
		return
	}

	util.Success(c, "支付成功", result)
}

func GetOrderDetail(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}
	orderNo := c.Query("order_no")
	if orderNo == "" {
		util.Fail(c, 400, "订单号不能为空")
		return
	}
	result, err := Service.GetOrderDetail(userId, orderNo)
	if err != nil {
		logger.Log.Warn("获取订单详情失败",
			zap.Error(err),
			zap.Uint("user_id", userId),
			zap.String("order_no", orderNo),
		)
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "获取订单详情成功", result)
}

func DeleteUserOrder(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}
	orderNo := c.Query("order_no")
	if orderNo == "" {
		util.Fail(c, 400, "订单号不能为空")
		return
	}
	if err := Service.DeleteUserOrder(userId, orderNo); err != nil {
		logger.Log.Warn("删除订单失败",
			zap.Error(err),
			zap.Uint("user_id", userId),
			zap.String("order_no", orderNo),
		)
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "删除订单成功", nil)
}
