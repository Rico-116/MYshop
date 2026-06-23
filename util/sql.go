package util

import (
	"MYshop/config"
	"MYshop/package/logger"
	//"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	Db  *gorm.DB
	err error
)

func InitMySQL() error {
	Db, err = gorm.Open(mysql.Open(config.AppConfig.MySQL.DSN))
	if err != nil {
		return err
	}
	sqlDB, err := Db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	logger.Log.Info("mysql connected",
		zap.Int("max_open_conns", 100),
		zap.Int("max_idle_conns", 20),
	)
	return nil
}

func CloseMySQL() {
	if Db == nil {
		return
	}
	sqlDB, err := Db.DB()
	if err != nil {
		logger.Log.Warn("mysql raw db get failed", zap.Error(err))
		return
	}
	if err := sqlDB.Close(); err != nil {
		logger.Log.Warn("mysql close failed", zap.Error(err))
	}
}
