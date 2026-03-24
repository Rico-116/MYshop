package dao

import (
	"MYshop/models"
	"MYshop/util"
	//"errors"
	"github.com/gin-gonic/gin"
	//"gorm.io/gorm"
	//"github.com/quic-go/quic-go/interop/utils"
)

//	func GetByUsername(username string) (*models.User, error) {
//		var user models.User
//		err := util.Db.Where("username = ?", username).First(&user).Error
//		if err != nil {
//			if errors.Is(err, gorm.ErrRecordNotFound) {
//				return nil, nil
//			}
//			return nil, err
//		}
//		return &user, nil
//	}
func GetByEmail(email string) (*models.User, error) {
	var user models.User
	sql := "SELECT * FROM user WHERE email = ? LIMIT 1"

	err := util.Db.Raw(sql, email).Scan(&user).Error
	if err != nil {
		return nil, err
	}

	// 没查到数据时，主键一般还是 0
	if user.UserId == 0 {
		return nil, nil
	}

	return &user, nil
}
func IsLogin(c *gin.Context) {
	//cookie, err := c.Cookie("login")
}
func CreateUser(user *models.User) error {

	sql := "INSERT INTO `user` (`username`, `password`,`nickname`,`email`,`status`,`phone`,`avatar`,`created_at`) VALUES (?,?,?,?,?,?,?,?)"
	var phone interface{}
	if user.Phone == "" {
		phone = nil
	} else {
		phone = user.Phone
	}
	return util.Db.Exec(sql, user.Username, user.Password, user.Nickname, user.Email, user.Status, phone, user.Avatar, user.CreatedAt).Error
}

//	func SendEmailCode(c *gin.Context, email string) error {
//		exist, err := GetByEmail(email)
//	}
func GetByUsername(username string) (*models.User, error) {
	var user models.User
	sql := "SELECT * FROM user WHERE username = ?"
	err := util.Db.Raw(sql, username).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	if user.UserId == 0 {
		return nil, nil
	}
	return &user, nil
}

func UpdatePassword(email string, password string) error {
	sql := "UPDATE user SET `password` = ? WHERE email = ?"
	return util.Db.Exec(sql, password, email).Error
}
