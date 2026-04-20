package dao

import (
	"MYshop/models"
	"MYshop/util"
)

func GetAdminByUsername(username string) (*models.AdminUser, error) {
	var admin models.AdminUser
	sql := "SELECT * FROM admin_user WHERE username = ? LIMIT 1"

	err := util.Db.Raw(sql, username).Scan(&admin).Error
	if err != nil {
		return nil, err
	}

	if admin.Id == 0 {
		return nil, nil
	}

	return &admin, nil
}
