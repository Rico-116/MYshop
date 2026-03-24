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

func SendRegister(c *gin.Context) {
	//var req models.SendRegisterCodeRequest
	//if err := c.ShouldBindJSON(&req); err != nil {
	//	util.Fail(c, 400, err.Error())
	//	return
	//}
	err := Service.SendEmailCode(c)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "验证码发送成功", gin.H{
		"expire_seconds": 300,
	})
}
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("注册参数错误",
			zap.Error(err))
		util.Fail(c, 400, "参数错误")
		return
	}
	err := Service.Register(req)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "注册成功", nil)
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("登录失败：参数错误",
			zap.Error(err),
			zap.String("ip", c.ClientIP()),
		)
		util.Fail(c, 400, "参数错误")
		return
	}

	if req.Username == "" {
		logger.Log.Warn("登录失败：用户名不能为空",
			zap.String("ip", c.ClientIP()),
		)
		util.Fail(c, 400, "用户名不能为空")
		return
	}

	if req.Password == "" {
		logger.Log.Warn("登录失败：密码不能为空",
			zap.String("username", req.Username),
			zap.String("ip", c.ClientIP()),
		)
		util.Fail(c, 400, "密码不能为空")
		return
	}

	user, err := dao.GetByUsername(req.Username)
	if err != nil {
		logger.Log.Error("登录失败：查询用户失败",
			zap.String("username", req.Username),
			zap.String("ip", c.ClientIP()),
			zap.Error(err),
		)
		util.Fail(c, 500, "查询用户失败")
		return
	}

	if user == nil {
		logger.Log.Warn("登录失败：用户不存在",
			zap.String("username", req.Username),
			zap.String("ip", c.ClientIP()),
		)
		util.Fail(c, 400, "用户不存在")
		return
	}

	if !util.CheckPasswordHash(user.Password, req.Password) {
		logger.Log.Warn("登录失败：密码错误",
			zap.String("username", req.Username),
			zap.Uint("user_id", user.UserId),
			zap.String("ip", c.ClientIP()),
		)
		util.Fail(c, 400, "用户名或密码错误")
		return
	}

	token, err := util.GenerateToken(user.UserId, user.Username)
	if err != nil {
		logger.Log.Error("登录失败：生成 token 失败",
			zap.Uint("user_id", user.UserId),
			zap.String("username", user.Username),
			zap.String("ip", c.ClientIP()),
			zap.Error(err),
		)
		util.Fail(c, 500, "登录失败")
		return
	}

	logger.Log.Info("用户登录成功",
		zap.Uint("user_id", user.UserId),
		zap.String("username", user.Username),
		zap.String("ip", c.ClientIP()),
	)

	util.Success(c, "登录成功", gin.H{
		"token":    token,
		"user_id":  user.UserId,
		"username": user.Username,
	})
}
func SendLoginCode(c *gin.Context) {
	var req models.SendLoginCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("参数错误",
			zap.Error(err))
		util.Fail(c, 400, "参数错误")
		return
	}
	err := Service.SendLoginCode(req)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	logger.Log.Info("验证码发送成功",
		zap.String("ip", c.ClientIP()),
		zap.String("Email", req.Email))
	util.Success(c, "验证码发送成功", gin.H{
		"expire_seconds": 300,
	})
}
func EmailLogin(c *gin.Context) {
	var req models.EmailLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("参数错误",
			zap.Error(err))
		util.Fail(c, 400, "参数错误")
		return
	}
	token, user, err := Service.EmailLogin(req)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	logger.Log.Info("登录成功",
		zap.Uint("user_id", user.UserId),
		zap.String("email", req.Email),
		zap.String("token", token),
		zap.String("username", user.Username))
	util.Success(c, "登录成功", gin.H{
		"token":    token,
		"user_id":  user.UserId,
		"username": user.Username,
		"email":    user.Email,
	})
}
func DeleteUser(c *gin.Context) {

}
