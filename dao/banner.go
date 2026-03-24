package dao

import (
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
)

func GetBannerList() ([]models.Banner, error) {
	var banners []models.Banner
	sql := "SELECT id,title,image_url,link,sort,status,created_at,updated_at FROM banner WHERE status = 1 ORDER BY sort ASC, id DESC"
	err := util.Db.Raw(sql).Scan(&banners).Error
	if err != nil {
		logger.Log.Error(err.Error())
	}
	return banners, err
}
