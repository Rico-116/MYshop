package Service

import (
	"MYshop/dao"
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"strconv"
	"time"
)

const CartListTTL = 10 * time.Minute

func AddCart(productID uint, userId uint, skuID uint, num int) error {
	cart, err := dao.GetCartByUserAndSku(userId, skuID)
	if err != nil {
		return err
	}
	if cart != nil && cart.Id != 0 {
		newNum := cart.Quantity + uint(num)
		if err := dao.UpdateCartQuantity(userId, cart.Id, int(newNum)); err != nil {
			return err
		}
	} else {
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
		sku, err := dao.GetSkuById(item.SkuId)
		if err != nil {
			return nil, err
		}
		if sku == nil || sku.Id == 0 {
			continue
		}
		product, err := dao.GetProductById(int(item.ProductId))
		if err != nil {
			return nil, err
		}
		if product == nil || product.Id == 0 {
			continue
		}
		price := sku.Price
		quantity := item.Quantity
		totalAmount := price * float64(quantity)
		list = append(list, models.CartDisplayItem{
			CartId:      item.CartId,
			SkuId:       sku.Id,
			SkuName:     sku.SkuName,
			ProductId:   product.Id,
			Title:       product.Subtitle,
			Image:       product.MainImage,
			Price:       price,
			Stock:       sku.Stock,
			Quantity:    quantity,
			Checked:     item.Checked,
			TotalAmount: totalAmount,
		})

	}
	logger.Log.Debug("<UNK>", zap.Any("list", list))
	return list, nil
}
func GetCartDisplayList(userID uint) ([]models.CartDisplayItem, error) {
	cacheKeyUserID := strconv.Itoa(int(userID))
	list, err := BuildCartListFromDB(userID)
	if err != nil {
		return nil, err
	}
	if list != nil {
		return list, nil
	}
	list, err = BuildCartListFromDB(userID)
	if err != nil {
		return nil, err
	}
	SetCartListToCache(cacheKeyUserID, list)
	return list, nil
}
func UpdateCartChecked(userID uint, cartId uint, checked bool) error {
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
