package Service

import (
	"MYshop/dao"
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func CreateOrder(userId uint, req models.CreateOrderRequest) (*models.CreateOrderResult, error) {
	if userId == 0 {
		return nil, errors.New("用户未登录")
	}
	if req.AddressId == 0 {
		return nil, errors.New("请选择收货地址")
	}
	if req.SourceType != "cart" && req.SourceType != "direct" {
		return nil, errors.New("下单来源错误")
	}
	addr, err := dao.GetAddressByIdAndUserId(req.AddressId, userId)
	if err != nil {
		return nil, err
	}
	if addr == nil || addr.Id == 0 {
		return nil, errors.New("收货地址不存在")
	}
	items, totalAmount, err := buildOrderItemsBySource(userId, req)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New("没有可下单的商品")
	}
	orderNo := generateOrderNo(userId)
	freightAmount := 0.0
	couponAmount := 0.0
	payAmount := totalAmount + freightAmount - couponAmount
	tx := util.Db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	for _, item := range items {
		if err := dao.DeductSkuStockTx(tx, item.SkuId, item.Quantity); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	order := &models.Order{
		OrderNo:               orderNo,
		UserId:                userId,
		Status:                models.OrderStatusUnpaid,
		TotalAmount:           totalAmount,
		PayAmount:             payAmount,
		CouponAmount:          couponAmount,
		FreightAmount:         freightAmount,
		ReceiverName:          addr.ReceiverName,
		ReceiverPhone:         addr.ReceiverPhone,
		ReceiverProvince:      addr.Province,
		ReceiverCity:          addr.City,
		ReceiverDistrict:      addr.District,
		ReceiverDetailAddress: addr.DetailAddress,
		Remark:                req.Remark,
	}
	orderId, err := dao.CreateOrderTx(tx, order)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, item := range items {
		orderItem := &models.OrderItem{
			OrderId:      orderId,
			OrderNo:      orderNo,
			UserId:       userId,
			ProductId:    item.ProductId,
			SkuId:        item.SkuId,
			ProductName:  item.ProductName,
			ProductImage: item.ProductImage,
			SkuName:      item.SkuName,
			Price:        item.Price,
			Quantity:     item.Quantity,
			TotalAmount:  item.TotalAmount,
		}
		if err := dao.CreateOrderItemTx(tx, orderItem); err != nil {
			tx.Rollback()
			return nil, err
		}
		if req.SourceType == "cart" && len(req.CartIds) > 0 {
			if err := dao.DeleteCartItemByIdsTx(tx, userId, req.CartIds); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		if err := tx.Commit().Error; err != nil {
			return nil, err
		}
		if req.SourceType == "cart" {
			DeleteCartListCache(strconv.Itoa(int(userId)))
		}
		logger.Log.Info("创建订单成功",
			zap.Uint("user_id", userId),
			zap.String("order_no", orderNo),
			zap.Uint("order_id", orderId),
		)
	}
	return &models.CreateOrderResult{
		OrderId:     orderId,
		OrderNo:     orderNo,
		TotalAmount: totalAmount,
		ItemCount:   len(items),
		Status:      models.OrderStatusUnpaid,
	}, nil
}
func buildOrderItemsBySource(userId uint, req models.CreateOrderRequest) ([]models.OrderBuyItem, float64, error) {
	switch req.SourceType {
	case "cart":
		return buildOrderItemsFromCart(userId, req.CartIds)
	case "direct":
		return buildOrderItemsFromDirect(req.SkuId, req.Quantity)
	default:
		return nil, 0, errors.New("非法下单来源")
	}
}

func buildOrderItemsFromCart(userId uint, cartIds []uint) ([]models.OrderBuyItem, float64, error) {
	if len(cartIds) == 0 {
		return nil, 0, errors.New("请选择要下单的购物车商品")
	}

	list, err := dao.GetCartOrderItemsByIds(userId, cartIds)
	if err != nil {
		return nil, 0, err
	}
	if len(list) == 0 {
		return nil, 0, errors.New("未找到可下单购物车商品")
	}

	totalAmount := 0.0
	validItems := make([]models.OrderBuyItem, 0, len(list))

	for _, item := range list {
		if item.ProductId == 0 || item.SkuId == 0 {
			return nil, 0, errors.New("商品不存在")
		}
		if item.ProductStatus != 1 {
			return nil, 0, errors.New("商品已下架")
		}
		if item.SkuStatus != 1 {
			return nil, 0, errors.New("商品规格不可购买")
		}
		if item.Quantity <= 0 {
			return nil, 0, errors.New("购买数量错误")
		}
		if item.Stock < item.Quantity {
			return nil, 0, errors.New("商品库存不足")
		}

		item.TotalAmount = item.Price * float64(item.Quantity)
		totalAmount += item.TotalAmount
		validItems = append(validItems, item)
	}

	if len(validItems) != len(cartIds) {
		return nil, 0, errors.New("部分购物车商品不存在")
	}

	return validItems, totalAmount, nil
}

func buildOrderItemsFromDirect(skuId uint, quantity int) ([]models.OrderBuyItem, float64, error) {
	if skuId == 0 {
		return nil, 0, errors.New("请选择商品规格")
	}
	if quantity <= 0 {
		return nil, 0, errors.New("购买数量错误")
	}

	item, err := dao.GetDirectOrderItemBySkuId(skuId)
	if err != nil {
		return nil, 0, err
	}
	if item == nil || item.SkuId == 0 {
		return nil, 0, errors.New("商品不存在")
	}
	if item.ProductStatus != 1 {
		return nil, 0, errors.New("商品已下架")
	}
	if item.SkuStatus != 1 {
		return nil, 0, errors.New("商品规格不可购买")
	}
	if item.Stock < quantity {
		return nil, 0, errors.New("商品库存不足")
	}

	item.Quantity = quantity
	item.TotalAmount = item.Price * float64(quantity)

	return []models.OrderBuyItem{*item}, item.TotalAmount, nil
}
func generateOrderNo(userId uint) string {
	return fmt.Sprintf("ORD%d%d", userId, time.Now().UnixNano())
}
