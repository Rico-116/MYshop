package Service

import (
	"MYshop/dao"
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"math/rand"
	//"strconv"
	"time"
)

const (
	ProductListTTL     = 10 * time.Minute
	ProductDetailTTL   = 30 * time.Minute
	ProductCategoryTTL = 10 * time.Minute
	ProductSkuTTL      = 30 * time.Minute
	NullCacheTTL       = 2 * time.Minute
	LockTTL            = 10 * time.Minute
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func withJitter(base time.Duration, maxJitterSeconds int) time.Duration {
	if maxJitterSeconds <= 0 {
		return base
	}
	return base + time.Duration(rand.Intn(maxJitterSeconds))*time.Second
}

// GetProductListWithCache 商品列表：先查redis，没找到的话去查mysql
func GetProductListWithCache() ([]models.Product, error) {
	key := util.ProductListKey()
	val, err := util.RDB.Get(util.Ctx, key).Result()
	if err == nil {
		var list []models.Product
		if e := json.Unmarshal([]byte(val), &list); e == nil {
			//反序列化，将json模式变回来，弄到list里面
			return list, nil
		} else {
			logger.Log.Warn("商品列表缓存反序列化失败", zap.Error(err))
		}

	}
	if err != nil && errors.Is(err, redis.Nil) {
		logger.Log.Warn("读取商品列表缓存失败", zap.Error(err))
	}
	list, err := dao.GetProductList()
	if err != nil {
		logger.Log.Warn("在mysql中读取失败", zap.Error(err))
		return nil, err
	}
	data, _ := json.Marshal(list) //将数据序列化
	if e := util.RDB.Set(util.Ctx, key, data, ProductListTTL).Err(); e != nil {
		logger.Log.Warn("写入商品缓存列表失败", zap.Error(e))
	}
	return list, nil

}
func getProductDetailWithLock(productID int) (*models.Product, error) {
	key := util.ProductDetailKey(productID)
	lockKey := util.ProductDetailLockKey(productID)
	locked, err := util.RDB.SetNX(util.Ctx, lockKey, "1", LockTTL).Result()
	if err != nil {
		logger.Log.Warn("商品详情加锁失败", zap.Error(err), zap.Int("product_id", productID))

	}
	if locked {
		defer func() {
			if err := util.RDB.Del(util.Ctx, lockKey).Err(); err != nil {
				logger.Log.Warn("商品解锁失败", zap.Error(err))
			}
		}()
		// 双检：防止拿到锁前别的请求已经重建好缓存
		val, err := util.RDB.Get(util.Ctx, key).Result()
		if err == nil {
			if val == util.CacheNullValue {
				return nil, nil
			}
			var product models.Product
			if e := json.Unmarshal([]byte(val), &product); e == nil {
				return &product, nil
			}
		}
		product, err := dao.GetProductById(productID)
		if err != nil {
			return nil, err
		}
		//商品不存在写控制防止击穿
		if product == nil || product.Id == 0 {
			if e := util.RDB.Set(util.Ctx, key, util.CacheNullValue, withJitter(NullCacheTTL, 60)).Err(); e != nil {
				logger.Log.Warn("写入商品空值缓存失败", zap.Error(e))
			}
			return nil, nil
		}
		data, _ := json.Marshal(product)
		if e := util.RDB.Set(util.Ctx, key, data, withJitter(ProductDetailTTL, 300)).Err(); e != nil {
			logger.Log.Warn("写入商品详情缓存失败", zap.Error(e))
		}
		return product, nil
	}
	time.Sleep(80 * time.Millisecond)
	val, err := util.RDB.Get(util.Ctx, key).Result()
	if err == nil {
		if val == util.CacheNullValue {
			return nil, nil
		}
		var product models.Product
		if e := json.Unmarshal([]byte(val), &product); e == nil {
			return &product, nil
		}
	}
	product, err := dao.GetProductById(productID)
	if err != nil {
		return nil, err
	}
	if product == nil || product.Id == 0 {
		return nil, nil
	}
	return product, nil
}

// GetProductDetailWithCache 商品详情：先查Redis，没命中再查MySQL
func GetProductDetailWithCache(ProductID int) (*models.Product, error) {
	key := util.ProductDetailKey(ProductID)
	val, err := util.RDB.Get(util.Ctx, key).Result()
	if err == nil {
		if val == util.CacheNullValue {
			return nil, nil
		}
		var product models.Product
		if e := json.Unmarshal([]byte(val), &product); e == nil {
			return &product, nil
		} else {
			logger.Log.Warn("商品详情缓存反序列化失败", zap.Error(e))
		}
	}
	if err != nil && errors.Is(err, redis.Nil) {
		logger.Log.Warn("读取商品详情缓存失败", zap.Error(err))
	}
	return getProductDetailWithLock(ProductID)
}

// GetProductSkuListWithCache 商品SKU列表缓存
func GetProductSkuListWithCache(productID int) ([]models.ProductSku, error) {
	key := util.ProductCategoryListKey(productID)
	val, err := util.RDB.Get(util.Ctx, key).Result()
	if err == nil {
		var list []models.ProductSku
		if e := json.Unmarshal([]byte(val), &list); e == nil {
			return list, nil
		} else {
			logger.Log.Warn("商品SKU缓存反序列化失败", zap.Error(e))
		}
	}
	list, err := dao.GetProductSkuListByProductId(productID)
	if err != nil {
		return list, err
	}
	data, _ := json.Marshal(list)
	if e := util.RDB.Set(util.Ctx, key, data, ProductSkuTTL).Err(); e != nil {
		logger.Log.Warn("写入商品SKU缓存失败", zap.Error(e))
	}
	return list, nil
}
func GetProductListByCategoryWithCache(categoryID int) ([]models.Product, error) {
	key := util.ProductCategoryListKey(categoryID)
	val, err := util.RDB.Get(util.Ctx, key).Result()
	if err == nil {
		var list []models.Product
		if e := json.Unmarshal([]byte(val), &list); e == nil {
			return list, nil
		} else {
			logger.Log.Warn("分类商品反序列化失败", zap.Error(e))
		}
	}
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.Log.Warn("读取分类商品缓存失败", zap.Error(err))
	}
	list, err := dao.GetProductByCategoryId(categoryID)
	if err != nil {
		return list, err
	}
	data, _ := json.Marshal(list)
	if e := util.RDB.Set(util.Ctx, key, data, ProductListTTL).Err(); e != nil {
		logger.Log.Warn("写入分类商品缓存失败", zap.Error(e))
	}
	return list, nil
}
func DeleteProductCache(productID int, categoryID int) {
	_ = util.RDB.Del(util.Ctx,
		util.ProductDetailKey(productID),
		util.ProductSkuListKey(productID),
		util.ProductListKey(),
		util.ProductCategoryListKey(categoryID),
	).Err()
}
