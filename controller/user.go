package controller

import (
	"MYshop/Service"
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

}
