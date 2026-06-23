package main

import (
	"MYshop/Service"
	"MYshop/controller"
	_ "MYshop/docs"
	"MYshop/middleware"
	"MYshop/package/logger"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"os"
	"time"
)

// @title MYshop 商城系统接口文档
// @version 1.0
// @description 基于 Go + Gin + MySQL + Redis + RabbitMQ 的商城项目接口文档
// @host localhost:8080
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 请输入 Bearer token，例如：Bearer eyJhbGciOiJIUzI1NiIs...
func main() {
	// 初始化日志
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	if err := logger.Init(env); err != nil {
		panic(err)
	}
	defer logger.Sync()

	if env == "prod" || env == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	if err := util.InitMySQL(); err != nil {
		logger.Log.Fatal("mysql init failed", zap.Error(err))
	}
	defer util.CloseMySQL()

	if err := util.InitRedis(); err != nil {
		logger.Log.Fatal("redis init failed", zap.Error(err))
	}
	defer util.CloseRedis()

	// 初始化 RabbitMQ
	if err := util.InitRabbitMQ(); err != nil {
		logger.Log.Fatal("rabbitmq init failed", zap.Error(err))
	}
	defer util.CloseRabbitMQ()

	// 启动订单超时关单消费者
	if err := Service.StartOrderCloseConsumer(); err != nil {
		logger.Log.Fatal("order close consumer start failed", zap.Error(err))
	}
	if err := Service.StartSeckillOrderConsumer(); err != nil {
		logger.Log.Fatal("seckill order consumer start failed", zap.Error(err))
	}

	// 创建 Gin 实例
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 跨域配置：一定要放在所有路由注册之前
	r.Use(corsMiddleware())

	// 启动商品热度定时写回任务
	Service.StartProductHotWriteBackWorker(5 * time.Minute)

	if err := Service.LoadUnfinishedSeckillStockToRedis(); err != nil {
		logger.Log.Fatal("seckill stock preload failed", zap.Error(err))
	}

	// 静态资源
	r.Static("/static", "./static")

	// 用户相关接口
	userGroup := r.Group("/api/user")
	{
		// @Router /user/send_code [post]
		userGroup.POST("/send_code", controller.SendRegister)
		// @Router /user/register [post]
		userGroup.POST("/register", controller.Register)
		// @Router /user/login [post]
		userGroup.POST("/login", controller.Login)
		// @Router /user/send_login_code [post]
		userGroup.POST("/send_login_code", controller.SendLoginCode)
		// @Router /user/email/login [post]
		userGroup.POST("/email_login", controller.EmailLogin)
		// @Router /user/send_reset_password_code [post]
		userGroup.POST("/send_reset_password_code", controller.SendResetPasswordCode)
		// @Router /user/reset_password [post]
		userGroup.POST("/reset_password", controller.ResetPassword)

	}

	// 首页、商品、分类相关接口
	indexGroup := r.Group("/api/index")
	{
		// @Router /index/banners [get]
		indexGroup.GET("/banners", controller.GetBannerList)
		// @Router /index/categories [get]
		indexGroup.GET("/categories", controller.GetCategoryTree)
		// @Router /index/products [get]
		indexGroup.GET("/products", controller.GetProductList)
		// @Router /index/product/detail [get]
		indexGroup.GET("/product/detail", controller.GetProductDetail)
		// @Router /index/products/category [get]
		indexGroup.GET("/products/category", controller.GetProductListByCategory)
		// @Router /index/category/display [get]
		indexGroup.GET("/category/display", controller.GetCategoryDisplay)
		// @Router /index/products/hot [get]
		indexGroup.GET("/products/hot", controller.GetHotProductList)
		// @Router /index/product/comments [get]
		indexGroup.GET("/product/comments", controller.GetProductCommentTree)
		// @Router /index/search [get]
		indexGroup.GET("/search", controller.SearchProducts)
		// @Router /index/seckill/list [get]
		indexGroup.GET("/seckill/list", controller.GetUserSeckillList)
	}

	// 需要登录的接口
	authGroup := r.Group("/api/auth")
	authGroup.Use(middleware.JWTAuthMiddleware())
	{
		// 购物车
		// @Router /auth/cart/add [post]
		authGroup.POST("/cart/add", controller.AddCart)
		// @Router /auth/cart/list [get]
		authGroup.GET("/cart/list", controller.GetCartList)
		// @Router /auth/cart/quantity [put]
		authGroup.PUT("/cart/quantity", controller.UpdateCartQuantity)
		// @Router /auth/cart/check [put]
		authGroup.PUT("/cart/check", controller.UpdateCartChecked)
		// @Router /auth/cart/delete [delete]
		authGroup.DELETE("/cart/delete", controller.DeleteCart)
		// @Router /auth/user/display/profile [get]
		authGroup.GET("/user/display/profile", controller.GetUserProfile)
		// @Router /auth/user/profile [put]
		authGroup.PUT("/user/profile", controller.UpdateUserProfile)

		// 地址
		// @Router /auth/address/add [post]
		authGroup.POST("/address/add", controller.AddAddress)
		// @Router /auth/address/list [get]
		authGroup.GET("/address/list", controller.GetAddressList)
		// @Router /auth/address/default [put]
		authGroup.PUT("/address/default", controller.SetDefaultAddress)
		// @Router /auth/address/update/{id} [put]
		authGroup.PUT("/address/update/:id", controller.UpdateAddress)
		authGroup.DELETE("/address/delete/:id", controller.DeleteAddress)

		// 订单
		authGroup.POST("/order/preview", controller.PreviewOrder)
		authGroup.POST("/order/create", controller.CreateOrder)
		authGroup.GET("/order/list", controller.GetOrderList)
		authGroup.GET("/order/detail", controller.GetOrderDetail)
		authGroup.POST("/order/pay", controller.PayOrder)
		authGroup.GET("/order/pay/page", controller.GetPayPage)
		authGroup.POST("/order/cancel", controller.CancelUserOrder)
		authGroup.POST("/order/confirm", controller.ConfirmReceive)
		authGroup.DELETE("/order/delete", controller.DeleteUserOrder)
		authGroup.POST("/product/comment", controller.CreateProductComment)

		authGroup.POST("/seckill/submit", controller.SubmitSeckill)

		authGroup.GET("/seckill/result", controller.GetSeckillResult)
	}

	// 管理员接口
	adminGroup := r.Group("/api/admin")
	{
		adminGroup.POST("/login", controller.AdminLogin)
		adminAuthGroup := adminGroup.Group("")
		adminAuthGroup.Use(middleware.AdminJWTAuthMiddleware())
		{
			adminAuthGroup.GET("/product/list", controller.AdminGetProductList)
			adminAuthGroup.POST("/product/add", controller.AdminCreateProduct)
			adminAuthGroup.PUT("/product/update/:id", controller.AdminUpdateProduct)
			adminAuthGroup.DELETE("/product/delete/:id", controller.AdminDeleteProduct)
			adminAuthGroup.POST("/product/sku/add", controller.AdminCreateProductSku)
			adminAuthGroup.PUT("/product/sku/update/:id", controller.AdminUpdateProductSku)
			adminAuthGroup.DELETE("/product/sku/delete/:id", controller.AdminDeleteProductSku)

			adminAuthGroup.GET("/order/list", controller.AdminGetOrderList)
			adminAuthGroup.GET("/order/detail", controller.AdminGetOrderDetail)
			adminAuthGroup.POST("/order/ship", controller.AdminShipOrder)
			adminAuthGroup.POST("/order/cancel", controller.AdminCancelOrder)

			adminAuthGroup.GET("/user/list", controller.AdminGetUserList)
			adminAuthGroup.PUT("/user/status", controller.AdminUpdateUserStatus)

			adminAuthGroup.GET("/category/list", controller.AdminGetCategoryList)
			adminAuthGroup.POST("/category/add", controller.AdminCreateCategory)
			adminAuthGroup.PUT("/category/update/:id", controller.AdminUpdateCategory)
			adminAuthGroup.DELETE("/category/delete/:id", controller.AdminDeleteCategory)
			adminAuthGroup.POST("/seckill/create", controller.AdminCreateSeckillActivity)
			adminAuthGroup.GET("/seckill/list", controller.AdminGetSeckillActivityList)
			adminAuthGroup.PUT("/seckill/update", controller.AdminUpdateSeckillActivity)
			adminAuthGroup.POST("/seckill/close", controller.AdminCloseSeckillActivity)
			adminAuthGroup.POST("/seckill/init", controller.AdminInitSeckillStock)
		}
	}

	// 启动服务
	logger.Log.Info("server starting", zap.Int("port", 8080))
	if err := r.Run(":8080"); err != nil {
		logger.Log.Fatal("server start failed", zap.Error(err))
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Content-Length,Accept,Authorization,token,Token,X-Requested-With,X-CSRF-Token,Access-Control-Request-Private-Network")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Authorization,token")
		c.Header("Access-Control-Allow-Private-Network", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
