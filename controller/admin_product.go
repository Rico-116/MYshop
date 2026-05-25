package controller

import (
	"MYshop/Service"
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func AdminCreateProduct(c *gin.Context) {
	var req models.AdminCreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}
	product, err := Service.CreateAdminProduct(req)
	if err != nil {
		logger.Log.Warn("管理员新增商品失败", zap.Error(err), zap.Any("req", req))
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "新增商品成功", gin.H{
		"product_id": product.Id,
		"product":    product,
	})
}

func AdminGetProductList(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var categoryId *uint
	if categoryStr := c.Query("category_id"); categoryStr != "" {
		categoryInt, err := strconv.Atoi(categoryStr)
		if err != nil || categoryInt < 0 {
			util.Fail(c, 400, "分类参数错误")
			return
		}
		value := uint(categoryInt)
		categoryId = &value
	}

	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		value, err := strconv.Atoi(statusStr)
		if err != nil {
			util.Fail(c, 400, "状态参数错误")
			return
		}
		status = &value
	}

	list, total, err := Service.GetAdminProductList(keyword, categoryId, status, page, pageSize)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "查询商品列表成功", gin.H{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func AdminGetProductDetail(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil || productId <= 0 {
		util.Fail(c, 400, "商品id错误")
		return
	}
	product, skuList, err := Service.GetAdminProductDetail(uint(productId))
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "查询商品详情成功", gin.H{
		"product": product,
		"skus":    skuList,
	})
}

func AdminUpdateProduct(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil || productId <= 0 {
		util.Fail(c, 400, "商品id错误")
		return
	}
	var req models.AdminUpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}
	product, err := Service.UpdateAdminProduct(uint(productId), req)
	if err != nil {
		logger.Log.Warn("管理员修改商品失败", zap.Error(err), zap.Any("req", req))
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "修改商品成功", gin.H{"product": product})
}

func AdminDeleteProduct(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil || productId <= 0 {
		util.Fail(c, 400, "商品id错误")
		return
	}
	if err := Service.DeleteAdminProduct(uint(productId)); err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "删除商品成功", nil)
}

func AdminCreateProductSku(c *gin.Context) {
	var req models.AdminCreateSkuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}
	sku, err := Service.CreateAdminProductSku(req)
	if err != nil {
		logger.Log.Warn("管理员新增商品SKU失败", zap.Error(err), zap.Any("req", req))
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "新增商品SKU成功", gin.H{
		"sku_id": sku.Id,
		"sku":    sku,
	})
}

func AdminUpdateProductSku(c *gin.Context) {
	skuId, err := strconv.Atoi(c.Param("id"))
	if err != nil || skuId <= 0 {
		util.Fail(c, 400, "SKU id错误")
		return
	}
	var req models.AdminUpdateSkuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}
	sku, err := Service.UpdateAdminProductSku(uint(skuId), req)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "修改商品SKU成功", gin.H{"sku": sku})
}

func AdminDeleteProductSku(c *gin.Context) {
	skuId, err := strconv.Atoi(c.Param("id"))
	if err != nil || skuId <= 0 {
		util.Fail(c, 400, "SKU id错误")
		return
	}
	if err := Service.DeleteAdminProductSku(uint(skuId)); err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "删除商品SKU成功", nil)
}
