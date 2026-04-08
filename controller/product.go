package controller

import (
	"MYshop/Service"
	"MYshop/dao"
	"MYshop/package/logger"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func GetProductList(c *gin.Context) {
	list, err := dao.GetProductList()
	if err != nil {
		logger.Log.Error("获取商品列表失败", zap.Error(err))
		util.Fail(c, 500, "获取商品列表失败")
		return
	}
	logger.Log.Info("获取商品列表成功", zap.Any("list", list))
	util.Success(c, "获取商品列表成功", gin.H{"list": list})

}
func GetProductDetail(c *gin.Context) {
	idStr := c.Query("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		logger.Log.Warn("商品id参数错误", zap.Error(err))
		util.Fail(c, 400, "商品id参数错误")
		return
	}
	product, err := dao.GetProductById(int(id64))
	if err != nil {
		logger.Log.Warn("获取商品详情失败", zap.Error(err))
		util.Fail(c, 500, "获取商品详情失败")
		return
	}
	if product == nil || product.Id == 0 {
		logger.Log.Warn("商品不存在", zap.Any("product", product))
		util.Fail(c, 404, "商品不存在")
		return
	}
	skuList, err := dao.GetProductSkuListByProductId(int(id64))
	if err != nil {
		logger.Log.Warn("获取商品sku失败", zap.Error(err))
		util.Fail(c, 500, "获取商品sku失败")
		return
	}
	var defaultSkuId uint = 0
	for _, sku := range skuList {
		if sku.Stock > 0 {
			defaultSkuId = sku.Id
			break
		}
	}
	if defaultSkuId == 0 && len(skuList) > 0 {
		defaultSkuId = skuList[0].Id
	}
	identity := Service.BuildUserIdentity(
		c.GetString("user_id"),
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)
	if err = Service.RecordProductView(int(id64), identity); err != nil {
		logger.Log.Warn("记录商品热度失败", zap.Error(err), zap.Uint64("product_id", id64))
	}
	util.Success(c, "获取商品详情成功", gin.H{
		"detail":         product,
		"sku_list":       skuList,
		"default_sku_id": defaultSkuId,
	})
}
func GetProductListByCategory(c *gin.Context) {
	categoryStr := c.Query("category_id")
	categoryId, err := strconv.ParseUint(categoryStr, 10, 64)
	if err != nil || categoryId == 0 {
		logger.Log.Warn("分类id参数错误", zap.Error(err))
		util.Fail(c, 400, "分类id参数错误")
		return
	}
	list, err := dao.GetProductByCategoryId(int(categoryId))
	if err != nil {
		logger.Log.Warn("按分类获取商品失败", zap.Error(err))
		util.Fail(c, 500, "按分类获取商品失败")
		return
	}
	util.Success(c, "按分类获取商品成功", gin.H{
		"list": list,
	})
}
func GetHotProductList(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 8
	}
	list, err := Service.GetHotProducts(limit)
	if err != nil {
		logger.Log.Warn("获取热门推荐失败", zap.Error(err))
		util.Fail(c, 500, "获取热门推荐失败")
		return
	}
	util.Success(c, "获取热门推荐成功", gin.H{
		"list": list,
	})

}
