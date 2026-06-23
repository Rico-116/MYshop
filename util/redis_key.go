package util

import "fmt"

const (
	HotProductZetKey             = "shop:product:hot:zet"         //热门商品排行榜
	ProductClickWriteBackHK      = "shop:product:click:writeback" //待回写到MySQL的点击量：Hash
	ProductAntiBrushPrefix       = "shop:product:view:dedup"      //防刷前缀：String
	ProductListKeyPrefix         = "shop:product:list:"
	ProductDetailKeyPrefix       = "shop:product:detail"
	ProductSkuListKeyPrefix      = "shop:product:sku:list:"
	ProductCategoryListKeyPrefix = "shop:product:category:list:"
	CategoryDetailKeyPrefix      = "shop:category:detail:"
	ProductDetailLockKeyPrefix   = "shop:product:detail:lock:"
	CacheNullValue               = "null"
	CartListKeyPrefix            = "shop:cart:list:user"
	SeckillStockKeyPrefix        = "shop:seckill:stock"
	SeckillUserKeyPrefix         = "shop:seckill:user"
	SeckillResultKeyPrefix       = "shop:seckill:result"
	SeckillListKeyPrefix         = "shop:seckill:list"
	SeckillListLockKeyPrefix     = "shop:seckill:list:lock"
)

func SeckillUserKey(activityId uint, userId uint) string {
	return fmt.Sprintf("%s:%d:%d", SeckillUserKeyPrefix, activityId, userId)
}

func SeckillResultKey(activityId uint, userId uint) string {
	return fmt.Sprintf("%s:%d:%d", SeckillResultKeyPrefix, activityId, userId)
}

func SeckillStockKey(activityId uint) string {
	return fmt.Sprintf("%s:%d", SeckillStockKeyPrefix, activityId)
}

func SeckillListKey(page int, pageSize int) string {
	return fmt.Sprintf("%s:%d:%d", SeckillListKeyPrefix, page, pageSize)
}

func SeckillListLockKey(page int, pageSize int) string {
	return fmt.Sprintf("%s:%d:%d", SeckillListLockKeyPrefix, page, pageSize)
}

func ProductAntiBrushKey(ProductID int, identity string) string {
	return fmt.Sprintf("%s:%d:%s", ProductAntiBrushPrefix, ProductID, identity)
}
func ProductListKey() string {
	return ProductListKeyPrefix + ":all"
}
func ProductDetailKey(productID int) string {
	return fmt.Sprintf("%s:%d", ProductDetailKeyPrefix, productID)
}
func ProductSkuListKey(productID int) string {
	return fmt.Sprintf("%s:%d", ProductSkuListKeyPrefix, productID)
}
func ProductCategoryListKey(productID int) string {
	return fmt.Sprintf("%s:%d", ProductCategoryListKeyPrefix, productID)
}
func CategoryDetailKey(categoryID uint) string {
	return fmt.Sprintf("%s%d", CategoryDetailKeyPrefix, categoryID)
}
func ProductDetailLockKey(productID int) string {
	return fmt.Sprintf("%s:%d", ProductDetailLockKeyPrefix, productID)
}
func CartListKey(UserID string) string {
	return fmt.Sprintf("%s:%s", CartListKeyPrefix, UserID)
}
