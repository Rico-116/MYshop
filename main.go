package main

// @title           My API
// @version         1.0
// @description     这是我的 Go 服务接口文档
// @host            localhost:8080
// @BasePath        /api/v1
import (
	"MYshop/controller"
	"MYshop/package/logger"
	_ "MYshop/util"
	"github.com/gin-gonic/gin"
)

func main() {
	//hash, _ := util.HashPassword("123456")
	//fmt.Println("hash =", hash)
	//fmt.Println(util.CheckPasswordHash("123456", hash))  // 应该是 true
	//fmt.Println(util.CheckPasswordHash("1234567", hash)) // 应该是 false
	if err := logger.Init("dev"); err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Log.Info("服务启动成功")
	logger.Sugar.Infof("服务已启动, port=%d", 8080)
	r := gin.Default()
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/send_code", controller.SendRegister)
		userGroup.POST("/register", controller.Register)
		userGroup.POST("/login", controller.Login)
		userGroup.POST("/send_login_code", controller.SendLoginCode)
		userGroup.POST("/email_login", controller.EmailLogin)
		//userGroup.POST("/del",controller.)
		//userGroup.POST("/index",controller.)
	}
	//authApi:= r.Group("/api/auth")
	//authApi.Use(middleware.JWTAuthMiddleware()){}
	r.Run(":8080")
}
