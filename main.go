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
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
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
	//r.GET("/ping", func(c *gin.Context) {
	//	c.JSON(200, gin.H{"msg": "pong"})
	//})
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	//r.Use(func(c *gin.Context) {
	//	println("method =", c.Request.Method, "path =", c.Request.URL.Path)
	//	c.Next()
	//})
	//
	//r.OPTIONS("/*path", func(c *gin.Context) {
	//	c.Status(204)
	//})
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/send_code", controller.SendRegister)
		userGroup.POST("/register", controller.Register)
		userGroup.POST("/login", controller.Login)
		userGroup.POST("/send_login_code", controller.SendLoginCode)
		userGroup.POST("/email_login", controller.EmailLogin)
		userGroup.POST("/send_reset_password_code", controller.SendResetPasswordCode)
		userGroup.POST("/reset_password", controller.ResetPassword)
		//userGroup.POST("/del",controller.)
		//userGroup.POST("/index", controller.Index)
	}
	indexGroup := r.Group("/api/index")
	{
		indexGroup.GET("/banners", controller.GetBannerList)
		indexGroup.GET("/categories", controller.GetCategoryTree)
	}
	//authApi:= r.Group("/api/auth")
	//authApi.Use(middleware.JWTAuthMiddleware()){}
	r.Run(":8080")
}
