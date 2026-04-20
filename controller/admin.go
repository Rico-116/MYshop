package controller

import (
	"MYshop/dao"
	"MYshop/package/logger"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AdminLogin(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}

	if req.Username == "" {
		util.Fail(c, 400, "用户名不能为空")
		return
	}
	if req.Password == "" {
		util.Fail(c, 400, "密码不能为空")
		return
	}

	admin, err := dao.GetAdminByUsername(req.Username)
	if err != nil {
		logger.Log.Error("管理员登录失败：查询管理员失败",
			zap.String("username", req.Username),
			zap.String("ip", c.ClientIP()),
			zap.Error(err),
		)
		util.Fail(c, 500, "查询管理员失败")
		return
	}

	if admin == nil {
		util.Fail(c, 400, "管理员不存在")
		return
	}

	if admin.Status != 1 {
		util.Fail(c, 400, "管理员账号已被禁用")
		return
	}

	if !util.CheckPasswordHash(admin.Password, req.Password) {
		util.Fail(c, 400, "用户名或密码错误")
		return
	}

	token, err := util.GenerateToken(admin.Id, admin.Username, "admin")
	if err != nil {
		logger.Log.Error("管理员登录失败：生成 token 失败",
			zap.Uint("admin_id", admin.Id),
			zap.String("username", admin.Username),
			zap.String("ip", c.ClientIP()),
			zap.Error(err),
		)
		util.Fail(c, 500, "登录失败")
		return
	}

	logger.Log.Info("管理员登录成功",
		zap.Uint("admin_id", admin.Id),
		zap.String("username", admin.Username),
		zap.String("ip", c.ClientIP()),
	)

	util.Success(c, "管理员登录成功", gin.H{
		"token":    token,
		"admin_id": admin.Id,
		"username": admin.Username,
		"nickname": admin.Nickname,
		"role":     "admin",
	})
}
