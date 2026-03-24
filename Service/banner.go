package Service

import (
	"MYshop/dao"
	"MYshop/models"
)

func GetBannerList() ([]models.Banner, error) {
	return dao.GetBannerList()
}
