package controller

import (
	"MYshop/Service"
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateProductComment(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}

	var req models.CreateProductCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("创建商品评论参数错误",
			zap.Error(err),
			zap.Uint("user_id", userId),
		)
		util.Fail(c, 400, "参数错误")
		return
	}

	comment, err := Service.CreateProductComment(userId, req)
	if err != nil {
		logger.Log.Warn("创建商品评论失败",
			zap.Error(err),
			zap.Uint("user_id", userId),
			zap.Uint("product_id", req.ProductId),
		)
		util.Fail(c, 400, err.Error())
		return
	}

	util.Success(c, "评论成功", comment)
}

func GetProductCommentTree(c *gin.Context) {
	productId64, err := strconv.ParseUint(c.Query("product_id"), 10, 64)
	if err != nil || productId64 == 0 {
		util.Fail(c, 400, "商品参数错误")
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		util.Fail(c, 400, "分页参数错误")
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		util.Fail(c, 400, "分页大小参数错误")
		return
	}

	result, err := Service.GetProductCommentTree(uint(productId64), page, pageSize)
	if err != nil {
		logger.Log.Warn("获取商品评论失败",
			zap.Error(err),
			zap.Uint64("product_id", productId64),
		)
		util.Fail(c, 400, err.Error())
		return
	}

	util.Success(c, "获取商品评论成功", result)
}
