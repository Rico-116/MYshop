package controller

import (
	"MYshop/Service"
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func ParseUintQuery(val string) uint {
	val = strings.TrimSpace(val)

	if val == "" {
		return 0
	}

	n, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}

	return uint(n)
}
func AdminCreateSeckillActivity(c *gin.Context) {
	var req models.AdminCreateSeckillActivityRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("创建秒杀活动参数错误", zap.Error(err))
		util.Fail(c, 400, "参数错误")
		return
	}

	if err := Service.AdminCreateSeckillActivity(req); err != nil {
		logger.Log.Warn("创建秒杀活动失败", zap.Error(err))
		util.Fail(c, 400, err.Error())
		return
	}

	util.Success(c, "创建秒杀活动成功", nil)
}
func AdminGetSeckillActivityList(c *gin.Context) {
	page := Service.ParseIntQuery(c.Query("page"), 1)
	pageSize := Service.ParseIntQuery(c.Query("pageSize"), 10)

	// status = -1 表示查询全部
	status := Service.ParseIntQuery(c.Query("status"), -1)

	result, err := Service.AdminGetSeckillActivityList(page, pageSize, status)
	if err != nil {
		logger.Log.Warn("查询秒杀活动列表失败", zap.Error(err))
		util.Fail(c, 400, err.Error())
		return
	}

	util.Success(c, "查询秒杀活动列表成功", result)
}
func AdminUpdateSeckillActivity(c *gin.Context) {
	var req models.AdminUpdateSeckillActivityRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("修改秒杀活动参数错误", zap.Error(err))
		util.Fail(c, 400, "参数错误")
		return
	}

	if err := Service.AdminUpdateSeckillActivity(req); err != nil {
		logger.Log.Warn("修改秒杀活动失败", zap.Error(err))
		util.Fail(c, 400, err.Error())
		return
	}

	util.Success(c, "修改秒杀活动成功", nil)
}
func AdminCloseSeckillActivity(c *gin.Context) {
	activityId := Service.ParseUintQuery(c.Query("activity_id"))

	if activityId == 0 {
		util.Fail(c, 400, "秒杀活动ID不能为空")
		return
	}

	if err := Service.AdminCloseSeckillActivity(activityId); err != nil {
		logger.Log.Warn("关闭秒杀活动失败",
			zap.Error(err),
			zap.Uint("activity_id", activityId),
		)
		util.Fail(c, 400, err.Error())
		return
	}

	util.Success(c, "关闭秒杀活动成功", nil)
}
func AdminInitSeckillStock(c *gin.Context) {
	activityId := Service.ParseUintQuery(c.Query("activity_id"))

	if activityId == 0 {
		util.Fail(c, 400, "秒杀活动ID不能为空")
		return
	}

	if err := Service.AdminInitSeckillStock(activityId); err != nil {
		logger.Log.Warn("初始化秒杀库存失败",
			zap.Error(err),
			zap.Uint("activity_id", activityId),
		)
		util.Fail(c, 400, err.Error())
		return
	}

	util.Success(c, "初始化秒杀库存成功", gin.H{
		"activity_id": activityId,
	})
}
