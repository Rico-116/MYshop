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
