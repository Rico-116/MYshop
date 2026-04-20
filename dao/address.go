package dao

import (
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func GetAddressListByUserId(userId uint) ([]models.UserAddress, error) {
	var list []models.UserAddress
	sql := `
		SELECT id, user_id, receiver_name, receiver_phone, province, city, district, detail_address,
		       is_default, status, created_at, updated_at
		FROM user_address
		WHERE user_id = ? AND status = 1
		ORDER BY is_default DESC, id DESC
	`
	err := util.Db.Raw(sql, userId).Scan(&list).Error
	if err != nil {
		logger.Log.Error("查询地址列表失败", zap.Error(err), zap.Uint("user_id", userId))
		return nil, err
	}
	return list, nil
}
func GetAddressByIdAndUserId(id uint, userId uint) (*models.UserAddress, error) {
	var addr models.UserAddress
	sql := `
		SELECT id, user_id, receiver_name, receiver_phone, province, city, district, detail_address,
		       is_default, status, created_at, updated_at
		FROM user_address
		WHERE id = ? AND user_id = ? AND status = 1
		LIMIT 1
	`
	err := util.Db.Raw(sql, id, userId).Scan(&addr).Error
	if err != nil {
		return nil, err
	}
	if addr.Id == 0 {
		return nil, nil
	}
	return &addr, nil
}
func AddAddress(addr *models.UserAddress) error {
	sql := `
		INSERT INTO user_address
		(user_id, receiver_name, receiver_phone, province, city, district, detail_address, is_default, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1, NOW(), NOW())
	`
	return util.Db.Exec(sql,
		addr.UserId,
		addr.ReceiverName,
		addr.ReceiverPhone,
		addr.Province,
		addr.City,
		addr.District,
		addr.DetailAddress,
		addr.IsDefault,
	).Error
}
func UpdateAddress(addr *models.UserAddress) error {
	sql := `
		UPDATE user_address
		SET receiver_name = ?, receiver_phone = ?, province = ?, city = ?, district = ?, detail_address = ?, is_default = ?, updated_at = NOW()
		WHERE id = ? AND user_id = ? AND status = 1
	`
	return util.Db.Exec(sql,
		addr.ReceiverName,
		addr.ReceiverPhone,
		addr.Province,
		addr.City,
		addr.District,
		addr.DetailAddress,
		addr.IsDefault,
		addr.Id,
		addr.UserId,
	).Error
}
func DeleteAddress(id uint, userId uint) error {
	sql := `
		UPDATE user_address
		SET status = 0, updated_at = NOW()
		WHERE id = ? AND user_id = ? AND status = 1
	`
	return util.Db.Exec(sql, id, userId).Error
}
func ClearDefaultAddress(userId uint) error {
	sql := `UPDATE user_address SET is_default = 0, updated_at = NOW() WHERE user_id = ? AND status = 1`
	return util.Db.Exec(sql, userId).Error
}
func SetDefaultAddress(id uint, userId uint) error {
	sql := `UPDATE user_address SET is_default = 1, updated_at = NOW() WHERE id = ? AND user_id = ? AND status = 1`
	return util.Db.Exec(sql, id, userId).Error
}
func GetDefaultAddressByUserId(userId uint) (*models.UserAddress, error) {
	var addr models.UserAddress
	sql := `
		SELECT id, user_id, receiver_name, receiver_phone, province, city, district, detail_address,
		       is_default, status, created_at, updated_at
		FROM user_address
		WHERE user_id = ? AND status = 1 AND is_default = 1
		LIMIT 1
	`
	err := util.Db.Raw(sql, userId).Scan(&addr).Error
	if err != nil {
		return nil, err
	}
	if addr.Id == 0 {
		return nil, nil
	}
	return &addr, nil
}
func GetAddressByIdAndUserIdTx(tx *gorm.DB, id uint, userId uint) (*models.UserAddress, error) {
	var addr models.UserAddress
	sql := `
		SELECT id, user_id, receiver_name, receiver_phone, province, city, district, detail_address,
		       is_default, status, created_at, updated_at
		FROM user_address
		WHERE id = ? AND user_id = ? AND status = 1
		LIMIT 1
	`
	err := tx.Raw(sql, id, userId).Scan(&addr).Error
	if err != nil {
		return nil, err
	}
	if addr.Id == 0 {
		return nil, nil
	}
	return &addr, nil
}
