package main

import (
	"MYshop/controller"
	_ "MYshop/util"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/send_code", controller.SendRegister)
		userGroup.POST("/register", controller.Register)
		userGroup.POST("/login", controller.Login)
	}
	r.Run(":8080")
}
