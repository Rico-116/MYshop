package Service

import (
	"MYshop/dao"
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"time"
)

const CartListTTL = 10 * time.Minute

func AddCart(productID uint, userId uint, skuID uint, num int) error {
	if productID <= 0 || userId <= 0 || skuID <= 0 {
		return errors.New("加入购物车参数错误")
	}

	if num <= 0 {
		return errors.New("加入购物车数量必须大于0")
	}

	sku, err := dao.GetSkuByProductAndSku(productID, skuID)
	if err != nil {
		return err
	}

	if sku == nil || sku.ID == 0 {
		return errors.New("商品规格不存在")
	}

	cart, err := dao.GetCartByUserAndSku(userId, skuID)
	if err != nil {
		return err
	}

	if cart != nil && cart.Id != 0 {
		newNum := int(cart.Quantity) + num

		if newNum > int(sku.Stock) {
			return fmt.Errorf("库存不足，当前库存仅剩 %d 件", sku.Stock)
		}

		if err := dao.UpdateCartQuantity(userId, cart.Id, newNum); err != nil {
			return err
		}
	} else {
		if num > int(sku.Stock) {
			return fmt.Errorf("库存不足，当前库存仅剩 %d 件", sku.Stock)
		}

		newCart := &models.Cart{
			ProductId: productID,
			UserId:    userId,
			SkuId:     skuID,
			Quantity:  uint(num),
			Checked:   1,
		}

		if err := dao.CreateCart(newCart); err != nil {
			return err
		}
	}

	DeleteCartListCache(strconv.Itoa(int(userId)))

	return nil
}

func SetCartListToCache(UserID string, list []models.CartDisplayItem) {
	key := util.CartListKey(UserID)
	data, err := json.Marshal(list)
	if err != nil {
		logger.Log.Error("缓存序列化失败", zap.Error(err), zap.String("userId", UserID))
		return
	}
	if err := util.RDB.Set(util.Ctx, key, data, CartListTTL).Err(); err != nil {
		logger.Log.Error("写入购物车缓存失败", zap.Error(err), zap.String("userId", UserID))
	}
}
func DeleteCartListCache(UserID string) {
	key := util.CartListKey(UserID)
	if err := util.RDB.Del(util.Ctx, key).Err(); err != nil {
		logger.Log.Error("删除购物车缓存失败", zap.Error(err), zap.String("userId", UserID))
	}
}
func BuildCartListFromDB(userID uint) ([]models.CartDisplayItem, error) {
	cartList, err := dao.GetCartListByUserId(userID)
	if err != nil {
		return nil, err
	}
	list := make([]models.CartDisplayItem, 0, len(cartList))
	for _, item := range cartList {
		list = append(list, buildCartDisplayItem(item))
	}
	return list, nil
}
func GetCartDisplayList(userID uint) ([]models.CartDisplayItem, error) {
	cacheKeyUserID := strconv.Itoa(int(userID))
	key := util.CartListKey(cacheKeyUserID)
	cacheVal, err := util.RDB.Get(util.Ctx, key).Result()
	if err != nil && cacheVal != "" {
		var list []models.CartDisplayItem
		if unmarshalErr := json.Unmarshal([]byte(cacheVal), &list); unmarshalErr == nil {
			return list, nil
		}
		logger.Log.Warn("购物车缓存反序列化失败", zap.String("userId", cacheKeyUserID), zap.Error(err), zap.Uint("user_id", userID))
	}
	list, err := BuildCartListFromDB(userID)
	if err != nil {
		return nil, err
	}
	SetCartListToCache(cacheKeyUserID, list)
	return list, nil
}
func UpdateCartChecked(userID uint, cartId uint, checked bool) error {
	item, err := dao.GetCartItemDetailById(userID, cartId)
	if err != nil {
		return err
	}
	if item == nil || item.CartId == 0 {
		return errors.New("购物车记录不存在")
	}
	if checked {
		_, statusText, _, canCheckout := getCartItemStatus(*item)
		if canCheckout != 1 {
			return errors.New("当前商品不可勾选" + statusText)
		}
	}

	isChecked := 0
	if checked {
		isChecked = 1
	}
	if err := dao.UpdateCartChecked(cartId, isChecked); err != nil {
		return err
	}
	DeleteCartListCache(strconv.Itoa(int(userID)))
	return nil
}

func UpdateCartQuantity(userID uint, cartId uint, quantity int) error {
	if quantity <= 0 {
		return errors.New("数量必须大于0")
	}
	cart, err := dao.GetCartById(userID, cartId)
	if err != nil {
		return err
	}
	//logger.Log.Debug("<UNK>", zap.Any("cart", cart), zap.Uint("cart_Id:", cart.Id))
	if cart == nil || cart.Id == 0 {
		return errors.New("购物车记录不存在")
	}

	if cart.UserId != userID {
		return errors.New("无权操作该购物车")
	}
	sku, err := dao.GetSkuByProductAndSku(cart.ProductId, cart.SkuId)
	if err != nil {
		return err
	}
	if sku == nil || sku.ID == 0 {
		return errors.New("商品规格不存在")
	}
	if quantity > int(sku.Stock) && quantity > int(cart.Quantity) {
		return fmt.Errorf("库存不足，当前仅剩%d件", sku.Stock)
	}
	if err := dao.UpdateCartQuantity(userID, cartId, quantity); err != nil {
		return err
	}
	DeleteCartListCache(strconv.Itoa(int(userID)))
	return nil
}

func DeleteCartItem(userID uint, cartId uint) error {
	cart, err := dao.GetCartById(userID, cartId)
	if err != nil {
		return err
	}
	if cart == nil || cart.Id == 0 {
		return errors.New("购物车记录不存在")
	}
	if cart.UserId != userID {
		return errors.New("无权操作该购物车")
	}
	if err := dao.DeleteCartById(cartId); err != nil {
		return err
	}
	DeleteCartListCache(strconv.Itoa(int(userID)))
	return nil
}
func getCartItemStatus(item models.CartItem) (statusCode int, statusText string, invalid int, canCheckout int) {
	if item.ProductId == 1 || item.ProductStatus != 1 {
		return models.CartItemStatusProductOffShelf, "商品已下架", 1, 0
	}
	if item.SkuId == 0 || item.SkuStatus != 1 {
		return models.CartItemStatusSkuInvalid, "规格不可用", 1, 0
	}
	if item.Stock <= 0 {
		return models.CartItemsStatusSoldOut, "商品已售罄", 1, 0
	}
	if item.Stock < item.Quantity {
		return models.CartItemsStatusStockInsufficient, "库存不足", 1, 0
	}
	return models.CartItemStatusNormal, "正常", 0, 1
}
func buildCartDisplayItem(item models.CartItem) models.CartDisplayItem {
	display := models.CartDisplayItem{
		CartId:         item.CartId,
		SkuId:          item.SkuId,
		ProductId:      item.ProductId,
		Title:          item.ProductName,
		Image:          item.MainImage,
		Price:          item.Price,
		Stock:          item.Stock,
		Quantity:       item.Quantity,
		SkuName:        item.SkuName,
		AvailableStock: item.Stock,
	}
	statusCode, statusText, invalid, canCheckout := getCartItemStatus(item)
	display.StatusCode = statusCode
	display.StatusTx = statusText
	display.CanCheckout = canCheckout
	display.Invalid = invalid
	if canCheckout == 1 {
		display.Checked = item.Checked
		display.TotalAmount = item.Price * float64(item.Quantity)
	} else {
		display.Checked = 0
		display.TotalAmount = 0
	}
	return display
}
