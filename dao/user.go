package dao

import (
	"MYshop/models"
	"MYshop/util"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	//"github.com/quic-go/quic-go/interop/utils"
)

func GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := util.Db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func GetByEmail(email string) (*models.User, error) {
	var user models.User
	sql := "SELECT * FROM user WHERE email = ?"
	err := util.Db.Exec(sql, email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func IsLogin(c *gin.Context) {
	//cookie, err := c.Cookie("login")
}
func CreateUser(user *models.User) error {
	sql := "INSERT INTO `user` (`username`, `password`,`nickname`,`email`,`status`,`phone`,`avatar`,`created_at`) VALUES (?,?,?,?,?,?,?,?)"
	return util.Db.Exec(sql, user.Username, user.Password, user.Nickname, user.Email, user.Status, user.Phone, user.Avatar, user.CreateAt).Error
}

//func SendEmailCode(c *gin.Context, email string) error {
//	exist, err := GetByEmail(email)
//}
