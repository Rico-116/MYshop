package Service

import (
	"MYshop/dao"
	"MYshop/models"
	"errors"
	"strings"
)

func AddAddress(userId uint, req models.AddAddressRequest) error {
	if strings.TrimSpace(req.ReceiverName) == "" {
		return errors.New("收货人不能为空")
	}
	if strings.TrimSpace(req.ReceiverPhone) == "" {
		return errors.New("手机号不能为空")
	}
	if strings.TrimSpace(req.Province) == "" || strings.TrimSpace(req.City) == "" ||
		strings.TrimSpace(req.District) == "" || strings.TrimSpace(req.DetailAddress) == "" {
		return errors.New("收货地址不能为空")
	}
	if req.IsDefault == 1 {
		if err := dao.ClearDefaultAddress(userId); err != nil {
			return err
		}
	}
	addr := &models.UserAddress{
		UserId:        userId,
		ReceiverName:  req.ReceiverName,
		ReceiverPhone: req.ReceiverPhone,
		Province:      req.Province,
		City:          req.City,
		District:      req.District,
		DetailAddress: req.DetailAddress,
		IsDefault:     req.IsDefault,
	}
	return dao.AddAddress(addr)
}
func GetAddressList(userId uint) ([]models.UserAddress, error) {
	return dao.GetAddressListByUserId(userId)
}
func SetDefaultAddress(userId uint, id uint) error {
	addr, err := dao.GetAddressByIdAndUserId(id, userId)
	if err != nil {
		return err
	}
	if addr == nil {
		return errors.New("地址不存在")
	}
	if err := dao.ClearDefaultAddress(userId); err != nil {
		return err
	}
	return dao.SetDefaultAddress(id, userId)
}
