package controller

import (
	"MYshop/Service"
	"MYshop/models"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AddAddress(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}
	var req models.AddAddressRequest
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}
	if err := Service.AddAddress(userId, req); err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "新增地址成功", nil)
}
func GetAddressList(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}
	list, err := Service.GetAddressList(userId)
	if err != nil {
		util.Fail(c, 500, "获取地址列表失败")
		return
	}
	util.Success(c, "获取地址列表成功", gin.H{"list": list})
}
func SetDefaultAddress(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.Fail(c, 400, "地址ID错误")
		return
	}
	if err := Service.SetDefaultAddress(userId, uint(id)); err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "设置默认地址成功", nil)
}
