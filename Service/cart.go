package Service

import (
	"MYshop/dao"
	"MYshop/models"
	"MYshop/package/logger"
	"errors"
	"go.uber.org/zap"
)

func AddCart(userId uint, req models.AddCartRequest) error {
	if userId == 0 {
		return errors.New("请先登录")
	}
	if req.SkuId == 0 {
		return errors.New("sku_id不能为空")
	}
	if req.Quantity <= 0 {
		return errors.New("加入购物车数量必须大于0")
	}
	//查sku
	sku, err := dao.GetSkuById(req.SkuId)
	if err != nil {
		logger.Log.Error("加入购物车失败：查询sku失败", zap.Error(err), zap.Uint("sku_id", req.SkuId))
		return errors.New("查询商品规格失败")
	}
	if sku == nil || sku.Id == 0 {
		return errors.New("商品规格不存在或者已下架")
	}
	product, err := dao.GetProductById(int(sku.ProductId))
	if err != nil {
		logger.Log.Error("加入购物车失败：查询商品失败", zap.Error(err), zap.Uint("product_id", sku.ProductId))
		return errors.New("查询商品失败")
	}
	if product == nil || product.Id == 0 {
		return errors.New("商品不存在或者已下架")
	}
	cart, err := dao.GetCartByUserAndSku(userId, req.SkuId)
	if err != nil {
		logger.Log.Error("加入购物车失败：查询购物车失败", zap.Error(err), zap.Uint("user_id", userId), zap.Uint("sku_id", req.SkuId))
		return errors.New("查询购物车失败")
	}

	finalQty := req.Quantity
	if cart != nil && cart.Id > 0 {
		finalQty = cart.Quantity + req.Quantity
	}
	if int(finalQty) > sku.Stock {
		return errors.New("加入购物车失败，库存不足")
	}
	if cart != nil && cart.Id > 0 {
		err = dao.UpdateCartQuantity(cart.Id, int(finalQty))
		if err != nil {
			logger.Log.Error("加入购物车失败：更新购物车失败", zap.Error(err), zap.Uint("cart_id", cart.Id))
			return errors.New("加入购物车失败")
		}
	} else {
		newCart := &models.Cart{
			UserId:    userId,
			ProductId: sku.ProductId,
			SkuId:     sku.Id,
			Quantity:  req.Quantity,
			Checked:   1,
		}
		err = dao.CreateCart(newCart)
		if err != nil {
			logger.Log.Error("加入购物车失败：新增购物车失败", zap.Error(err), zap.Any("cart", newCart))
			return errors.New("加入购物车失败")
		}
	}
	return nil
}
func GetCartList(userId uint) (map[string]interface{}, error) {
	list, err := dao.GetCartListByUserId(userId)
	if err != nil {
		return nil, errors.New("获取购物车列表失败")
	}
	checkedCount := 0
	checkedAmount := 0.0
	for i := range list {
		list[i].Subtotal = list[i].Price * float64(list[i].Quantity)
		if list[i].ProductId == 0 || list[i].SkuId == 0 || list[i].ProductStatus != 1 || list[i].SkuStatus != 1 || list[i].Stock <= 0 {
			list[i].Invalid = 1
		}
		if list[i].Checked == 1 && list[i].Invalid == 0 {
			checkedCount++
			checkedAmount += list[i].Subtotal
		}
	}
	return map[string]interface{}{
		"list":           list,
		"total_count":    len(list),
		"checked_count":  checkedCount,
		"checked_amount": checkedAmount,
	}, nil
}
func ChangeCartQuantity(userId uint, req models.UpdateCartQuantityRequest) error {
	if userId == 0 {
		return errors.New("请先登录")
	}
	if req.CartId == 0 {
		return errors.New("cart_id不能为空")
	}
	if req.Quantity == 0 {
		return errors.New("数量必须大于0")
	}
	cart, err := dao.GetCartById(userId, req.CartId)
	if err != nil {
		return errors.New("查询本商品信息失败")
	}
	if cart == nil || cart.Id == 0 {
		return errors.New("该条记录不存在")
	}
	sku, err := dao.GetSkuById(cart.SkuId)
	if err != nil {
		return errors.New("查询商品规格失败")
	}
	if sku == nil || sku.Id == 0 {
		return errors.New("商品规格不存在或者已下架")
	}
	if int(req.Quantity) > sku.Stock {
		return errors.New("库存不足")
	}
	err = dao.UpdateCartQuantity(req.CartId, int(req.Quantity))
	if err != nil {
		return errors.New("更新购物车数量失败")
	}
	return nil
}
func ChangeCartChecked(userId uint, req models.UpdateCartCheckRequest) error {
	if userId == 0 {
		return errors.New("请先登录")
	}
	if req.CartId == 0 {
		return errors.New("cart_id不能为空")
	}
	if req.Checked != 0 && req.Checked != 1 {
		return errors.New("商品状态参数错误")
	}
	cart, err := dao.GetCartById(userId, req.CartId)
	if err != nil {
		return errors.New("查询购物车失败")
	}
	if cart == nil || cart.Id == 0 {
		return errors.New("购物车记录不存在")
	}
	err = dao.UpdateCartChecked(req.CartId, int(req.Checked))
	if err != nil {
		return errors.New("更新勾选状态失败")
	}
	return nil
}
func DeleteCart(userId uint, cartId uint) error {
	if userId == 0 {
		return errors.New("请先登录")
	}
	if cartId == 0 {
		return errors.New("cart_id不能为空")
	}
	cart, err := dao.GetCartById(userId, cartId)
	if err != nil {
		return errors.New("查询购物车失败")
	}
	if cart == nil || cart.Id == 0 {
		return errors.New("购物车记录不存在")
	}
	err = dao.DeleteCartById(cartId)
	if err != nil {
		return errors.New("删除购物车失败")
	}
	return nil
}
