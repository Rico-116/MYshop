package controller

import (
	"MYshop/Service"
	"MYshop/package/logger"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func GetProductList(c *gin.Context) {
	list, err := Service.GetProductListWithCache()
	if err != nil {
		logger.Log.Error("获取商品列表失败", zap.Error(err))
		util.Fail(c, 500, "获取商品列表失败")
		return
	}
	logger.Log.Info("获取商品列表成功", zap.Int("count", len(list)))
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
	product, err := Service.GetProductDetailWithCache(int(id64))
	if err != nil {
		logger.Log.Warn("获取商品详情失败", zap.Error(err))
		util.Fail(c, 500, "获取商品详情失败")
		return
	}
	if product == nil || product.Id == 0 {
		logger.Log.Warn("商品不存在", zap.Uint64("product_id", id64))
		util.Fail(c, 404, "商品不存在")
		return
	}
	skuList, err := Service.GetProductSkuListWithCache(int(id64))
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
	Category, err := Service.GetCategoryDetailWithCache(product.CategoryId)
	if err != nil {
		logger.Log.Warn("获取商品类别失败", zap.Error(err), zap.Uint("category_id", product.CategoryId))
		util.Fail(c, 500, "获取商品类别失败，请稍后再试")
		return
	}
	if Category == nil || Category.Id == 0 {
		logger.Log.Warn("商品类别不存在", zap.Uint("category_id", product.CategoryId))
		util.Fail(c, 404, "商品类别不存在")
		return
	}
	util.Success(c, "获取商品详情成功", gin.H{
		"detail":         product,
		"sku_list":       skuList,
		"sku_category":   Category.Name,
		"category_id":    Category.Id,
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
	list, err := Service.GetProductListByCategoryWithCache(int(categoryId))
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
	limitStr := "10"
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {

		limit = 10
	}
	list, err := Service.GetHotProducts(limit)
	if err != nil {
		logger.Log.Warn("获取热门推荐失败", zap.Error(err))
		util.Fail(c, 500, "获取热门推荐失败")
		return
	}
	logger.Log.Info("获取热门推荐成功", zap.Int("count", len(list)), zap.Int("limit", limit))
	util.Success(c, "获取热门推荐成功", gin.H{
		"list": list,
	})

}
