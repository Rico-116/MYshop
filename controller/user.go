package controller

import (
	"MYshop/Service"
	"MYshop/dao"
	"MYshop/models"
	"MYshop/util"
	"github.com/gin-gonic/gin"
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
		util.Fail(c, 400, "参数错误")
		return
	}
	user, err := dao.GetByUsername(req.Username)
	if err != nil {
		util.Fail(c, 400, "查询用户失败")
		return
	}
	if user == nil {
		util.Fail(c, 400, "用户不存在")
		return
	}
	if !util.CheckPasswordHash(user.Password, req.Password) {
		util.Fail(c, 400, "密码错误")
		return
	}
	token, err := util.GenerateToken(user.UserId, user.Username)
	if err != nil {
		util.Fail(c, 400, "生成token失败")
		return
	}
	util.Success(c, "登录成功", gin.H{
		"token":    token,
		"user_id":  user.UserId,
		"username": user.Username,
	})
}
func SendLoginCode(c *gin.Context) {
	var req models.SendLoginCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}
	err := Service.SendLoginCode(req)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "验证码发送成功", gin.H{
		"expire_seconds": 300,
	})
}
func EmailLogin(c *gin.Context) {
	var req models.EmailLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}
	token, user, err := Service.EmailLogin(req)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "登录成功", gin.H{
		"token":    token,
		"user_id":  user.UserId,
		"username": user.Username,
		"email":    user.Email,
	})
}
func DeleteUser(c *gin.Context) {

}
