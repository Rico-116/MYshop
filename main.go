package main

// @title           My API
// @version         1.0
// @description     这是我的 Go 服务接口文档
// @host            localhost:8080
// @BasePath        /api/v1

import (
	"MYshop/Service"
	"MYshop/controller"
	"MYshop/middleware"
	"MYshop/package/logger"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	// 初始化日志
	if err := logger.Init("dev"); err != nil {
		panic(err)
	}
	defer logger.Sync()

	// 初始化 RabbitMQ
	if err := util.InitRabbitMQ(); err != nil {
		log.Fatalf("init rabbitmq error: %v", err)
	}
	defer util.CloseRabbitMQ()

	// 启动订单超时关单消费者
	if err := Service.StartOrderCloseConsumer(); err != nil {
		log.Fatalf("启动消费者失败: %v", err)
	}

	logger.Log.Info("服务启动成功")
	logger.Sugar.Infof("服务已启动, port=%d", 8080)

	// 创建 Gin 实例
	r := gin.Default()

	// 跨域配置：一定要放在所有路由注册之前
	r.Use(corsMiddleware())

	// 启动商品热度定时写回任务
	Service.StartProductHotWriteBackWorker(5 * time.Minute)

	// 静态资源
	r.Static("/static", "./static")

	// 用户相关接口
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/send_code", controller.SendRegister)
		userGroup.POST("/register", controller.Register)
		userGroup.POST("/login", controller.Login)
		userGroup.POST("/send_login_code", controller.SendLoginCode)
		userGroup.POST("/email_login", controller.EmailLogin)
		userGroup.POST("/send_reset_password_code", controller.SendResetPasswordCode)
		userGroup.POST("/reset_password", controller.ResetPassword)
	}

	// 首页、商品、分类相关接口
	indexGroup := r.Group("/api/index")
	{
		indexGroup.GET("/banners", controller.GetBannerList)
		indexGroup.GET("/categories", controller.GetCategoryTree)
		indexGroup.GET("/products", controller.GetProductList)
		indexGroup.GET("/product/detail", controller.GetProductDetail)
		indexGroup.GET("/products/category", controller.GetProductListByCategory)
		indexGroup.GET("/category/display", controller.GetCategoryDisplay)
		indexGroup.GET("/products/hot", controller.GetHotProductList)
		indexGroup.GET("/search", controller.SearchProducts)
	}

	// 需要登录的接口
	authGroup := r.Group("/api/auth")
	authGroup.Use(middleware.JWTAuthMiddleware())
	{
		// 购物车
		authGroup.POST("/cart/add", controller.AddCart)
		authGroup.GET("/cart/list", controller.GetCartList)
		authGroup.PUT("/cart/quantity", controller.UpdateCartQuantity)
		authGroup.PUT("/cart/check", controller.UpdateCartChecked)
		authGroup.DELETE("/cart/delete", controller.DeleteCart)

		// 地址
		authGroup.POST("/address/add", controller.AddAddress)
		authGroup.GET("/address/list", controller.GetAddressList)
		authGroup.PUT("/address/default", controller.SetDefaultAddress)
		authGroup.PUT("/address/update/:id", controller.UpdateAddress)
		authGroup.DELETE("/address/delete/:id", controller.DeleteAddress)

		// 订单
		authGroup.POST("/order/preview", controller.PreviewOrder)
		authGroup.POST("/order/create", controller.CreateOrder)
		authGroup.GET("/order/list", controller.GetOrderList)
		authGroup.GET("/order/detail", controller.GetOrderDetail)
		authGroup.POST("/order/pay", controller.PayOrder)
		authGroup.GET("/order/pay/page", controller.GetPayPage)
		authGroup.DELETE("/order/delete", controller.DeleteUserOrder)
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
		}
	}

	// 启动服务
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务启动失败: %v", err)
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
