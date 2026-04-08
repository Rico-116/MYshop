package util

import "fmt"

const (
	HotProductZetKey        = "shop:product:hot:zet"         //热门商品排行榜
	ProductClickWriteBackHK = "shop:product:click:writeback" //待回写到MySQL的点击量：Hash
	ProductAntiBrushPrefix  = "shop:product:view:dedup"      //防刷前缀：String
)

func ProductAntiBrushKey(ProductID int, identity string) string {
	return fmt.Sprintf("%s:%d:%s", ProductAntiBrushPrefix, ProductID, identity)
}
