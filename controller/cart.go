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
func GetCartList(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}
	data, err := Service.GetCartList(userId)
	if err != nil {
		util.Fail(c, 500, err.Error())
		return
	}
	util.Success(c, "获取购物车列表成功", data)
}
func ChangeCartQuantity(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}
	var req models.UpdateCartQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}
	err := Service.ChangeCartQuantity(userId, req)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "修改数量成功", nil)
}
func ChangeCartChecked(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}
	var req models.UpdateCartCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}
	err := Service.ChangeCartChecked(userId, req)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "修改购物车勾选状态成功", nil)
}
func DeleteCart(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}

	var req models.DeleteCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}

	err := Service.DeleteCart(userId, req.CartId)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "删除购物车成功", nil)
}
