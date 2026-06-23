package controller

// @Summary 发送注册验证码
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body models.SendRegisterCodeRequest true "邮箱信息"
// @Success 200 {object} util.Response
// @Router /user/send_code [post]
func swaggerSendRegister() {}

// @Summary 用户注册
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body models.RegisterRequest true "注册信息"
// @Success 200 {object} util.Response
// @Router /user/register [post]
func swaggerRegister() {}

// @Summary 用户账号密码登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body models.LoginRequest true "登录信息"
// @Success 200 {object} util.Response
// @Router /user/login [post]
func swaggerLogin() {}

// @Summary 发送邮箱登录验证码
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body models.SendLoginCodeRequest true "邮箱信息"
// @Success 200 {object} util.Response
// @Router /user/send_login_code [post]
func swaggerSendLoginCode() {}

// @Summary 邮箱验证码登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body models.EmailLoginRequest true "邮箱登录信息"
// @Success 200 {object} util.Response
// @Router /user/email_login [post]
func swaggerEmailLogin() {}

// @Summary 发送重置密码验证码
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body models.SendResetPasswordRequest true "邮箱信息"
// @Success 200 {object} util.Response
// @Router /user/send_reset_password_code [post]
func swaggerSendResetPasswordCode() {}

// @Summary 重置密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body models.ResetPasswordRequest true "重置密码信息"
// @Success 200 {object} util.Response
// @Router /user/reset_password [post]
func swaggerResetPassword() {}

// @Summary 获取首页轮播图
// @Tags 首页
// @Produce json
// @Success 200 {object} util.Response
// @Router /index/banners [get]
func swaggerGetBannerList() {}

// @Summary 获取分类树
// @Tags 首页
// @Produce json
// @Success 200 {object} util.Response
// @Router /index/categories [get]
func swaggerGetCategoryTree() {}

// @Summary 获取商品列表
// @Tags 首页
// @Produce json
// @Success 200 {object} util.Response
// @Router /index/products [get]
func swaggerGetProductList() {}

// @Summary 获取商品详情
// @Tags 首页
// @Produce json
// @Param id query int true "商品ID"
// @Success 200 {object} util.Response
// @Router /index/product/detail [get]
func swaggerGetProductDetail() {}

// @Summary 按分类获取商品列表
// @Tags 首页
// @Produce json
// @Param category_id query int true "分类ID"
// @Success 200 {object} util.Response
// @Router /index/products/category [get]
func swaggerGetProductListByCategory() {}

// @Summary 获取分类展示
// @Tags 首页
// @Produce json
// @Param category_id query int true "分类ID"
// @Success 200 {object} util.Response
// @Router /index/category/display [get]
func swaggerGetCategoryDisplay() {}

// @Summary 获取热门商品
// @Tags 首页
// @Produce json
// @Success 200 {object} util.Response
// @Router /index/products/hot [get]
func swaggerGetHotProductList() {}

// @Summary 获取商品评论树
// @Tags 首页
// @Produce json
// @Param product_id query int true "商品ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} util.Response
// @Router /index/product/comments [get]
func swaggerGetProductCommentTree() {}

// @Summary 搜索商品
// @Tags 首页
// @Produce json
// @Param keyword query string false "关键词"
// @Param category_id query int false "分类ID"
// @Param min_price query number false "最低价格"
// @Param max_price query number false "最高价格"
// @Param sort query string false "排序方式"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} util.Response
// @Router /index/search [get]
func swaggerSearchProducts() {}

// @Summary 获取用户端秒杀活动列表
// @Tags 秒杀
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} util.Response
// @Router /index/seckill/list [get]
func swaggerGetUserSeckillList() {}

// @Summary 加入购物车
// @Tags 购物车
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.AddCartRequest true "购物车信息"
// @Success 200 {object} util.Response
// @Router /auth/cart/add [post]
func swaggerAddCart() {}

// @Summary 获取购物车列表
// @Tags 购物车
// @Security BearerAuth
// @Produce json
// @Success 200 {object} util.Response
// @Router /auth/cart/list [get]
func swaggerGetCartList() {}

// @Summary 修改购物车数量
// @Tags 购物车
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.UpdateCartQuantityRequest true "数量信息"
// @Success 200 {object} util.Response
// @Router /auth/cart/quantity [put]
func swaggerUpdateCartQuantity() {}

// @Summary 修改购物车选中状态
// @Tags 购物车
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.UpdateCartCheckRequest true "选中状态"
// @Success 200 {object} util.Response
// @Router /auth/cart/check [put]
func swaggerUpdateCartChecked() {}

// @Summary 删除购物车商品
// @Tags 购物车
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.DeleteCartRequest true "购物车项信息"
// @Success 200 {object} util.Response
// @Router /auth/cart/delete [delete]
func swaggerDeleteCart() {}

// @Summary 获取用户资料
// @Tags 用户资料
// @Security BearerAuth
// @Produce json
// @Success 200 {object} util.Response
// @Router /auth/user/display/profile [get]
func swaggerGetUserProfile() {}

// @Summary 修改用户资料
// @Tags 用户资料
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.UpdateUserProfileRequest true "用户资料"
// @Success 200 {object} util.Response
// @Router /auth/user/profile [put]
func swaggerUpdateUserProfile() {}

// @Summary 新增收货地址
// @Tags 地址
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.AddAddressRequest true "地址信息"
// @Success 200 {object} util.Response
// @Router /auth/address/add [post]
func swaggerAddAddress() {}

// @Summary 获取收货地址列表
// @Tags 地址
// @Security BearerAuth
// @Produce json
// @Success 200 {object} util.Response
// @Router /auth/address/list [get]
func swaggerGetAddressList() {}

// @Summary 设置默认地址
// @Tags 地址
// @Security BearerAuth
// @Produce json
// @Param id query int true "地址ID"
// @Success 200 {object} util.Response
// @Router /auth/address/default [put]
func swaggerSetDefaultAddress() {}

// @Summary 修改收货地址
// @Tags 地址
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "地址ID"
// @Param data body models.UpdateAddressRequest true "地址信息"
// @Success 200 {object} util.Response
// @Router /auth/address/update/{id} [put]
func swaggerUpdateAddress() {}

// @Summary 删除收货地址
// @Tags 地址
// @Security BearerAuth
// @Produce json
// @Param id path int true "地址ID"
// @Success 200 {object} util.Response
// @Router /auth/address/delete/{id} [delete]
func swaggerDeleteAddress() {}

// @Summary 预览订单
// @Tags 订单
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.OrderPreviewRequest true "订单预览信息"
// @Success 200 {object} util.Response
// @Router /auth/order/preview [post]
func swaggerPreviewOrder() {}

// @Summary 创建订单
// @Tags 订单
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.CreateOrderRequest true "创建订单信息"
// @Success 200 {object} util.Response
// @Router /auth/order/create [post]
func swaggerCreateOrder() {}

// @Summary 获取订单列表
// @Tags 订单
// @Security BearerAuth
// @Produce json
// @Param status query int false "订单状态，-1表示全部" default(-1)
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} util.Response
// @Router /auth/order/list [get]
func swaggerGetOrderList() {}

// @Summary 获取订单详情
// @Tags 订单
// @Security BearerAuth
// @Produce json
// @Param order_no query string true "订单号"
// @Success 200 {object} util.Response
// @Router /auth/order/detail [get]
func swaggerGetOrderDetail() {}

// @Summary 支付订单
// @Tags 订单
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.PayOrderRequest true "支付信息"
// @Success 200 {object} util.Response
// @Router /auth/order/pay [post]
func swaggerPayOrder() {}

// @Summary 获取支付页
// @Tags 订单
// @Security BearerAuth
// @Produce json
// @Param order_no query string true "订单号"
// @Success 200 {object} util.Response
// @Router /auth/order/pay/page [get]
func swaggerGetPayPage() {}

// @Summary 取消订单
// @Tags 订单
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.CancelOrderRequest true "取消订单信息"
// @Success 200 {object} util.Response
// @Router /auth/order/cancel [post]
func swaggerCancelUserOrder() {}

// @Summary 确认收货
// @Tags 订单
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.ConfirmReceiveRequest true "确认收货信息"
// @Success 200 {object} util.Response
// @Router /auth/order/confirm [post]
func swaggerConfirmReceive() {}

// @Summary 删除订单
// @Tags 订单
// @Security BearerAuth
// @Produce json
// @Param order_no query string true "订单号"
// @Success 200 {object} util.Response
// @Router /auth/order/delete [delete]
func swaggerDeleteUserOrder() {}

// @Summary 创建商品评论
// @Tags 评论
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.CreateProductCommentRequest true "评论信息"
// @Success 200 {object} util.Response
// @Router /auth/product/comment [post]
func swaggerCreateProductComment() {}

// @Summary 提交秒杀
// @Tags 秒杀
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.SeckillSubmitRequest true "秒杀信息"
// @Success 200 {object} util.Response
// @Router /auth/seckill/submit [post]
func swaggerSubmitSeckill() {}

// @Summary 查询秒杀结果
// @Tags 秒杀
// @Security BearerAuth
// @Produce json
// @Param activity_id query int true "秒杀活动ID"
// @Success 200 {object} util.Response
// @Router /auth/seckill/result [get]
func swaggerGetSeckillResult() {}

// @Summary 管理员登录
// @Tags 管理员
// @Accept json
// @Produce json
// @Param data body AdminLoginRequest true "管理员登录信息"
// @Success 200 {object} util.Response
// @Router /admin/login [post]
func swaggerAdminLogin() {}

// @Summary 管理端商品列表
// @Tags 管理员-商品
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "关键词"
// @Param category_id query int false "分类ID"
// @Param status query int false "状态"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} util.Response
// @Router /admin/product/list [get]
func swaggerAdminGetProductList() {}

// @Summary 管理端新增商品
// @Tags 管理员-商品
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.AdminCreateProductRequest true "商品信息"
// @Success 200 {object} util.Response
// @Router /admin/product/add [post]
func swaggerAdminCreateProduct() {}

// @Summary 管理端修改商品
// @Tags 管理员-商品
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "商品ID"
// @Param data body models.AdminUpdateProductRequest true "商品信息"
// @Success 200 {object} util.Response
// @Router /admin/product/update/{id} [put]
func swaggerAdminUpdateProduct() {}

// @Summary 管理端删除商品
// @Tags 管理员-商品
// @Security BearerAuth
// @Produce json
// @Param id path int true "商品ID"
// @Success 200 {object} util.Response
// @Router /admin/product/delete/{id} [delete]
func swaggerAdminDeleteProduct() {}

// @Summary 管理端新增商品SKU
// @Tags 管理员-商品
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.AdminCreateSkuRequest true "SKU信息"
// @Success 200 {object} util.Response
// @Router /admin/product/sku/add [post]
func swaggerAdminCreateProductSku() {}

// @Summary 管理端修改商品SKU
// @Tags 管理员-商品
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "SKU ID"
// @Param data body models.AdminUpdateSkuRequest true "SKU信息"
// @Success 200 {object} util.Response
// @Router /admin/product/sku/update/{id} [put]
func swaggerAdminUpdateProductSku() {}

// @Summary 管理端删除商品SKU
// @Tags 管理员-商品
// @Security BearerAuth
// @Produce json
// @Param id path int true "SKU ID"
// @Success 200 {object} util.Response
// @Router /admin/product/sku/delete/{id} [delete]
func swaggerAdminDeleteProductSku() {}

// @Summary 管理端订单列表
// @Tags 管理员-订单
// @Security BearerAuth
// @Produce json
// @Param status query int false "订单状态，-1表示全部" default(-1)
// @Param order_no query string false "订单号"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} util.Response
// @Router /admin/order/list [get]
func swaggerAdminGetOrderList() {}

// @Summary 管理端订单详情
// @Tags 管理员-订单
// @Security BearerAuth
// @Produce json
// @Param order_no query string true "订单号"
// @Success 200 {object} util.Response
// @Router /admin/order/detail [get]
func swaggerAdminGetOrderDetail() {}

// @Summary 管理端订单发货
// @Tags 管理员-订单
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.AdminShipOrderRequest true "发货信息"
// @Success 200 {object} util.Response
// @Router /admin/order/ship [post]
func swaggerAdminShipOrder() {}

// @Summary 管理端取消订单
// @Tags 管理员-订单
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.AdminCancelOrderRequest true "取消订单信息"
// @Success 200 {object} util.Response
// @Router /admin/order/cancel [post]
func swaggerAdminCancelOrder() {}

// @Summary 管理端用户列表
// @Tags 管理员-用户
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "关键词"
// @Param status query int false "用户状态"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} util.Response
// @Router /admin/user/list [get]
func swaggerAdminGetUserList() {}

// @Summary 管理端修改用户状态
// @Tags 管理员-用户
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.AdminUpdateUserStatusRequest true "用户状态信息"
// @Success 200 {object} util.Response
// @Router /admin/user/status [put]
func swaggerAdminUpdateUserStatus() {}

// @Summary 管理端分类列表
// @Tags 管理员-分类
// @Security BearerAuth
// @Produce json
// @Param status query int false "分类状态"
// @Param level query int false "分类层级"
// @Success 200 {object} util.Response
// @Router /admin/category/list [get]
func swaggerAdminGetCategoryList() {}

// @Summary 管理端新增分类
// @Tags 管理员-分类
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.AdminCreateCategoryRequest true "分类信息"
// @Success 200 {object} util.Response
// @Router /admin/category/add [post]
func swaggerAdminCreateCategory() {}

// @Summary 管理端修改分类
// @Tags 管理员-分类
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Param data body models.AdminUpdateCategoryRequest true "分类信息"
// @Success 200 {object} util.Response
// @Router /admin/category/update/{id} [put]
func swaggerAdminUpdateCategory() {}

// @Summary 管理端删除分类
// @Tags 管理员-分类
// @Security BearerAuth
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} util.Response
// @Router /admin/category/delete/{id} [delete]
func swaggerAdminDeleteCategory() {}

// @Summary 管理端创建秒杀活动
// @Tags 管理员-秒杀
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.AdminCreateSeckillActivityRequest true "秒杀活动信息"
// @Success 200 {object} util.Response
// @Router /admin/seckill/create [post]
func swaggerAdminCreateSeckillActivity() {}

// @Summary 管理端秒杀活动列表
// @Tags 管理员-秒杀
// @Security BearerAuth
// @Produce json
// @Param status query int false "活动状态，-1表示全部" default(-1)
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} util.Response
// @Router /admin/seckill/list [get]
func swaggerAdminGetSeckillActivityList() {}

// @Summary 管理端修改秒杀活动
// @Tags 管理员-秒杀
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.AdminUpdateSeckillActivityRequest true "秒杀活动信息"
// @Success 200 {object} util.Response
// @Router /admin/seckill/update [put]
func swaggerAdminUpdateSeckillActivity() {}

// @Summary 管理端关闭秒杀活动
// @Tags 管理员-秒杀
// @Security BearerAuth
// @Produce json
// @Param activity_id query int true "秒杀活动ID"
// @Success 200 {object} util.Response
// @Router /admin/seckill/close [post]
func swaggerAdminCloseSeckillActivity() {}

// @Summary 管理端初始化秒杀库存
// @Tags 管理员-秒杀
// @Security BearerAuth
// @Produce json
// @Param activity_id query int true "秒杀活动ID"
// @Success 200 {object} util.Response
// @Router /admin/seckill/init [post]
func swaggerAdminInitSeckillStock() {}
