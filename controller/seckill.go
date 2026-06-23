package controller

import (
	"MYshop/Service"
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 用户端秒杀活动列表
func GetUserSeckillList(c *gin.Context) {
	page := Service.ParseIntQuery(c.Query("page"), 1)
	pageSize := Service.ParseIntQuery(c.Query("pageSize"), 10)

	result, err := Service.GetUserSeckillActivityList(page, pageSize)
	if err != nil {
		logger.Log.Warn("查询用户端秒杀列表失败", zap.Error(err))
		util.Fail(c, 400, err.Error())
		return
	}

	util.Success(c, "查询秒杀活动成功", result)
}

// 用户点击立即秒杀
func SubmitSeckill(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}

	var req models.SeckillSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}

	result, err := Service.SubmitSeckill(userId, req)
	if err != nil {
		logger.Log.Warn("提交秒杀失败",
			zap.Error(err),
			zap.Uint("user_id", userId),
			zap.String("activity_id", req.ActivityId),
		)
		util.Fail(c, 400, err.Error())
		return
	}

	util.Success(c, "秒杀请求已提交，排队中", result)
}

// 查询秒杀结果
func GetSeckillResult(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		util.Fail(c, 401, "请先登录")
		return
	}

	activityId := Service.ParseUintQuery(c.Query("activity_id"))
	if activityId == 0 {
		util.Fail(c, 400, "秒杀活动ID不能为空")
		return
	}

	result, err := Service.GetSeckillResult(userId, activityId)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}

	util.Success(c, "查询秒杀结果成功", result)
}
