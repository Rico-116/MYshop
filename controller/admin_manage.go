package controller

import (
	"MYshop/Service"
	"MYshop/models"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AdminGetOrderList(c *gin.Context) {
	status, err := strconv.Atoi(c.DefaultQuery("status", "-1"))
	if err != nil {
		util.Fail(c, 400, "订单状态参数错误")
		return
	}
	orderNo := c.Query("order_no")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := Service.GetAdminOrderList(status, orderNo, page, pageSize)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "查询订单列表成功", result)
}

func AdminGetOrderDetail(c *gin.Context) {
	orderNo := c.Query("order_no")
	result, err := Service.GetAdminOrderDetail(orderNo)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "查询订单详情成功", result)
}

func AdminShipOrder(c *gin.Context) {
	var req models.AdminShipOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}
	result, err := Service.ShipAdminOrder(req)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "发货成功", result)
}

func AdminCancelOrder(c *gin.Context) {
	var req models.AdminCancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}
	result, err := Service.CancelAdminOrder(req)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "取消订单成功", result)
}

func AdminGetUserList(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		value, err := strconv.Atoi(statusStr)
		if err != nil {
			util.Fail(c, 400, "用户状态参数错误")
			return
		}
		status = &value
	}

	list, total, err := Service.GetAdminUserList(keyword, status, page, pageSize)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "查询用户列表成功", gin.H{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func AdminUpdateUserStatus(c *gin.Context) {
	var req models.AdminUpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}
	user, err := Service.UpdateAdminUserStatus(req)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	msg := "恢复用户成功"
	if req.Status == 0 {
		msg = "拉黑用户成功"
	}
	util.Success(c, msg, gin.H{"user": user})
}
