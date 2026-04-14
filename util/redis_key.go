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
	ProductDetailLockKeyPrefix   = "shop:product:detail:lock:"
	CacheNullValue               = "null"
	CartListKeyPrefix            = "shop:cart:list:user"
)

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
func ProductDetailLockKey(productID int) string {
	return fmt.Sprintf("%s:%d", ProductDetailLockKeyPrefix, productID)
}
func CartListKey(UserID string) string {
	return fmt.Sprintf("%s:%s", CartListKeyPrefix, UserID)
}
