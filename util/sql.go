package util

import (
	//"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var (
	Db  *gorm.DB
	err error
)

func init() {
	Db, err = gorm.Open(mysql.Open("root:123456@tcp(192.168.0.147)/shop?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		log.Fatal("mysql 连接失败：%v", err)
	}
	log.Println("sql连接成功")
}
