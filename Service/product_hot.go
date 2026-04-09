package Service

import (
	"MYshop/dao"
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

var hotCtx = context.Background()

func BuildUserIdentity(userID interface{}, clientIP, userAgent string) string {
	if userID != nil {
		uid := strings.TrimSpace(fmt.Sprintf("%v", userID))
		if uid != "" {
			return "user:" + uid
		}
	}
	return fmt.Sprintf("ip:%s:us:%s", clientIP, userAgent)
}
func RecordProductView(ProductID int, identity string) error {
	if ProductID <= 0 || identity == "" {
		return errors.New("invalid product view params")
	}
	//防刷机制：10min内同一用户点击同一商品，无法上涨浏览量
	antiBrushKey := util.ProductAntiBrushKey(ProductID, identity)
	ok, err := util.RDB.SetNX(hotCtx, antiBrushKey, 1, 10*time.Minute).Result()
	//十分钟内同一个人再点这个商品，SetNX会失败，失败就return nil
	if err != nil {
		return err
	}
	if !ok {
		return nil //命中防刷，不重复计数，这里十分钟内点了同一个
	}
	//ZSet：实时排行榜
	//Hash：待回写MySQL的增量
	//先防刷才能进行这两步
	pipe := util.RDB.TxPipeline()
	pipe.ZIncrBy(hotCtx, util.HotProductZetKey, 1, strconv.Itoa(ProductID))
	pipe.HIncrBy(hotCtx, util.ProductClickWriteBackHK, strconv.Itoa(ProductID), 1)
	_, err = pipe.Exec(hotCtx)
	if err != nil {
		return err
	}
	return nil
}
func GetHotProducts(limit int) ([]models.Product, error) {
	if limit <= 0 {
		limit = 10
	}
	//logger.Log.Info("<UNK>", zap.Int("limit", limit))
	//fmt.Println(limit)
	// 1. 先从 Redis 热榜取
	idStrs, err := util.RDB.ZRevRange(hotCtx, util.HotProductZetKey, 0, int64(limit-1)).Result()
	if err == nil && len(idStrs) > 0 {
		productIDs := make([]int, 0, len(idStrs))
		for _, idStr := range idStrs {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				continue
			}
			productIDs = append(productIDs, id)
		}

		if len(productIDs) > 0 {
			products, err := dao.GetProductByIDs(productIDs)
			if err != nil {
				return nil, err
			}

			productMap := make(map[int]models.Product, len(products))
			for _, product := range products {
				productMap[int(product.Id)] = product
			}

			ordered := make([]models.Product, 0, len(productIDs))
			for _, id := range productIDs {
				if product, ok := productMap[id]; ok {
					ordered = append(ordered, product)
				}
			}
			if len(ordered) > 0 {
				return ordered, nil
			}
			//return products, nil
		}

	}

	// 2. Redis 为空时，走数据库兜底
	return dao.GetDefaultHotProducts(limit)
}
func FlushProductHotClickToDB() error { //将待回写的点击量刷入MySQL
	clickMapStr, err := util.RDB.HGetAll(hotCtx, util.ProductClickWriteBackHK).Result()
	if err != nil {
		logger.Log.Error("读取Redis待回写热度失败", zap.Error(err))
	}
	if len(clickMapStr) == 0 {
		return nil
	}
	clickMap := make(map[int]int64, len(clickMapStr))
	fields := make([]string, 0, len(clickMapStr))
	for productIDstr, deltaStr := range clickMapStr {
		productID, err1 := strconv.Atoi(productIDstr)
		delta, err2 := strconv.ParseInt(deltaStr, 10, 64)
		if err1 != nil || err2 != nil {
			logger.Log.Warn("待回写热度数据格式异常", zap.String("product_id", productIDstr), zap.String("delta", deltaStr))
			continue
		}
		if productID <= 0 || delta <= 0 {
			continue
		}
		clickMap[productID] = delta
		fields = append(fields, productIDstr)
	}
	if len(clickMap) == 0 {
		return nil
	}
	err = dao.BatchAddProductClickCount(clickMap)
	if err != nil {
		logger.Log.Error("回写商品到MYSQL失败", zap.Error(err))
		return err
	}
	err = util.RDB.HDel(hotCtx, util.ProductClickWriteBackHK, fields...).Err()
	if err != nil {
		logger.Log.Error("删除Redis已回写热度字段失败", zap.Error(err))
		return err
	}
	logger.Log.Info("商品热度回写MySQl成功", zap.Int("count", len(clickMap)))
	return nil
}

// 定时回写入口
func StartProductHotWriteBackWorker(interval time.Duration) {
	if interval <= 0 {
		interval = time.Minute * 5
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			if err := FlushProductHotClickToDB(); err != nil {
				logger.Log.Error("定时回写商品热度失败", zap.Error(err))
			}
		}
	}()
}
