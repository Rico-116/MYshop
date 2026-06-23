package controller

import (
	"MYshop/Service"
	"MYshop/dao"
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
	if req.SkuId == 0 || req.Quantity <= 0 {
		util.Fail(c, 400, "参数错误")
		return
	}
	productSku, err := dao.GetSkuById(req.SkuId)
	if err != nil {
		util.Fail(c, 400, "获取商品项失败")
		return
	}
	err = Service.AddCart(productSku.ProductId, userId, req.SkuId, int(req.Quantity))
	if err != nil {
		logger.Log.Error("加入购物车失败", zap.Error(err), zap.Uint("user_id", userId), zap.Uint("sku_id", req.SkuId))
		util.Fail(c, 500, err.Error())
		return
	}
	util.Success(c, "加入购物车成功", nil)
}
func GetCartList(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "用户未登录，请先登录")
		return
	}

	list, err := Service.GetCartDisplayList(uint(userId))
	if err != nil {
		logger.Log.Warn("获取购物车列表失败", zap.Error(err), zap.Uint("user_id", userId))
		util.Fail(c, 500, "获取购物车列表失败")
		return
	}
	selectedCount := 0
	invalidCount := 0
	selectedAmount := 0.0
	for _, item := range list {
		if item.Invalid == 1 {
			invalidCount++
		}
		if item.CanCheckout == 1 && item.Checked == 1 {
			selectedCount++
			selectedAmount += item.TotalAmount
		}
	}
	logger.Log.Info("获取购物车列表成功",
		zap.Uint("user_id", userId),
		zap.Int("count", len(list)),
		zap.Int("selected_count", selectedCount),
		zap.Int("invalid_count", invalidCount),
	)
	util.Success(c, "获取购物车列表成功", gin.H{
		"list": list,
		"summary": gin.H{
			"selected_count":  selectedCount,
			"invalid_count":   invalidCount,
			"selected_amount": selectedAmount,
		},
	})
}

func UpdateCartQuantity(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "用户未登录")
		return
	}
	var req models.UpdateCartQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}
	if req.CartId == 0 || req.Quantity <= 0 {
		util.Fail(c, 400, "参数错误")
		return
	}
	if err := Service.UpdateCartQuantity(userId, req.CartId, int(req.Quantity)); err != nil {
		logger.Log.Warn("修改购物车数量失败", zap.Error(err))
		util.Fail(c, 500, err.Error())
		return
	}
	util.Success(c, "修改购物车数量成功", nil)
}
func UpdateCartChecked(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "用户未登录")
		return
	}
	var req models.UpdateCartCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("修改购物车勾选状态参数错误", zap.Error(err), zap.Uint("user_id", userId))
		util.Fail(c, 400, "参数错误")
		return
	}
	if req.CartId == 0 && len(req.CartIds) == 0 {
		util.Fail(c, 400, "参数错误")
		return
	}
	if err := Service.UpdateCartChecked(userId, req.CartId, req.CartIds, *req.Checked == 1); err != nil {
		logger.Log.Warn("修改购物车勾选状态失败", zap.Error(err))
		util.Fail(c, 500, err.Error())
		return
	}
	util.Success(c, "修改购物车勾选状态成功", nil)
}

func DeleteCart(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "用户未登录")
		return
	}
	var req models.DeleteCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}
	if req.CartId == 0 {
		util.Fail(c, 400, "参数错误")
		return
	}
	if err := Service.DeleteCartItem(userId, req.CartId); err != nil {
		logger.Log.Warn("删除购物车失败", zap.Error(err))
		util.Fail(c, 500, err.Error())
		return
	}
	util.Success(c, "删除购物车成功", nil)
}
