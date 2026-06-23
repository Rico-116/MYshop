package Service

import (
	"MYshop/dao"
	"MYshop/models"
	"MYshop/util"
	"errors"
	"strings"
)

func GetAdminOrderList(status int, orderNo string, page, pageSize int) (*models.OrderListResult, error) {
	orderNo = strings.TrimSpace(orderNo)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if !isValidOrderStatusForQuery(status) {
		return nil, errors.New("订单状态参数错误")
	}

	list, total, err := dao.AdminGetMergedOrderList(status, orderNo, page, pageSize)
	if err != nil {
		return nil, err
	}
	orderNos := make([]string, 0, len(list))
	for _, item := range list {
		orderNos = append(orderNos, item.OrderNo)
	}
	items, err := dao.GetOrderItemsByOrderNos(orderNos)
	if err != nil {
		return nil, err
	}
	itemMap := make(map[string][]models.OrderItemVO)
	for _, item := range items {
		itemMap[item.OrderNo] = append(itemMap[item.OrderNo], item)
	}
	for i := range list {
		list[i].StatusText = models.GetOrderStatusText(list[i].Status)
		list[i].Items = itemMap[list[i].OrderNo]
	}

	return &models.OrderListResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Tabs:     models.GetOrderTabs(),
	}, nil
}

func GetAdminOrderDetail(orderNo string) (*models.OrderDetailResult, error) {
	orderNo = strings.TrimSpace(orderNo)
	if orderNo == "" {
		return nil, errors.New("订单号不能为空")
	}

	order, err := dao.AdminGetOrderByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}
	if order == nil || order.Id == 0 {
		return nil, errors.New("订单不存在")
	}
	items, err := dao.GetOrderItemsByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}
	seckillOrder, err := dao.GetSeckillOrderByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}
	isSeckill := seckillOrder != nil
	var seckillActivityId uint
	if isSeckill {
		seckillActivityId = seckillOrder.ActivityId
	}

	return &models.OrderDetailResult{
		OrderId:               order.Id,
		OrderNo:               order.OrderNo,
		UserId:                order.UserId,
		Status:                order.Status,
		StatusText:            models.GetOrderStatusText(order.Status),
		TotalAmount:           order.TotalAmount,
		PayAmount:             order.PayAmount,
		FreightAmount:         order.FreightAmount,
		CouponAmount:          order.CouponAmount,
		ReceiverName:          order.ReceiverName,
		ReceiverPhone:         order.ReceiverPhone,
		ReceiveProvince:       order.ReceiverProvince,
		ReceiverCity:          order.ReceiverCity,
		ReceiverDistrict:      order.ReceiverDistrict,
		ReceiverDetailAddress: order.ReceiverDetailAddress,
		Remark:                order.Remark,
		PayTime:               order.PayTime,
		DeliverTime:           order.DeliveryTime,
		FinishTime:            order.FinishTime,
		CloseTime:             order.CloseTime,
		CreateTime:            order.CreatedAt,
		UpdateTime:            order.UpdatedAt,
		IsSeckill:             isSeckill,
		SeckillActivityId:     seckillActivityId,
		Items:                 items,
	}, nil
}

func ShipAdminOrder(req models.AdminShipOrderRequest) (*models.OrderDetailResult, error) {
	orderNo := strings.TrimSpace(req.OrderNo)
	if orderNo == "" {
		return nil, errors.New("订单号不能为空")
	}
	order, err := dao.AdminGetOrderByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}
	if order == nil || order.Id == 0 {
		return nil, errors.New("订单不存在")
	}
	if order.Status != models.OrderStatusPaid {
		return nil, errors.New("只有待发货订单才能发货")
	}
	rows, err := dao.AdminUpdateOrderToShipped(orderNo)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, errors.New("发货失败，订单状态已变化")
	}
	return GetAdminOrderDetail(orderNo)
}

func CancelAdminOrder(req models.AdminCancelOrderRequest) (*models.OrderDetailResult, error) {
	orderNo := strings.TrimSpace(req.OrderNo)
	if orderNo == "" {
		return nil, errors.New("订单号不能为空")
	}
	order, err := dao.AdminGetOrderByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}
	if order == nil || order.Id == 0 {
		return nil, errors.New("订单不存在")
	}
	if order.Status == models.OrderStatusCanceled {
		return nil, errors.New("订单已取消，请勿重复操作")
	}
	if order.Status != models.OrderStatusUnpaid && order.Status != models.OrderStatusPaid {
		return nil, errors.New("只有待支付或待发货订单可以取消")
	}

	tx := util.Db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	rows, err := dao.AdminUpdateOrderToCanceledTx(tx, orderNo)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, errors.New("取消失败，订单状态已变化")
	}

	items, err := dao.GetOrderItemsByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if err := dao.RestoreSkuStockTx(tx, item.SkuId, item.Quantity); err != nil {
			return nil, err
		}
	}

	seckillOrder, err := dao.CancelSeckillOrderAndRestoreActivityStockTx(tx, orderNo)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	committed = true
	if seckillOrder != nil {
		rollbackRedisSeckill(seckillOrder.ActivityId, seckillOrder.UserId)
	}
	return GetAdminOrderDetail(orderNo)
}

func GetAdminUserList(keyword string, status *int, page, pageSize int) ([]models.User, int64, error) {
	keyword = strings.TrimSpace(keyword)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if status != nil && *status != 0 && *status != 1 {
		return nil, 0, errors.New("用户状态只能是0或1")
	}
	return dao.AdminGetUserList(keyword, status, page, pageSize)
}

func UpdateAdminUserStatus(req models.AdminUpdateUserStatusRequest) (*models.User, error) {
	if req.UserId == 0 {
		return nil, errors.New("用户id不能为空")
	}
	if req.Status != 0 && req.Status != 1 {
		return nil, errors.New("用户状态只能是0或1")
	}

	user, err := dao.AdminGetUserById(req.UserId)
	if err != nil {
		return nil, err
	}
	if user == nil || user.UserId == 0 {
		return nil, errors.New("用户不存在")
	}
	if err := dao.AdminUpdateUserStatus(req.UserId, req.Status); err != nil {
		return nil, err
	}
	return dao.AdminGetUserById(req.UserId)
}
