package util

import (
	"MYshop/config"
	"MYshop/package/logger"
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	RDB *redis.Client
	Ctx = context.Background()
)

func InitRedis() error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.Redis.Addr,
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
	})
	res, err := RDB.Ping(Ctx).Result()
	if err != nil {
		return err
	}
	logger.Log.Info("redis connected",
		zap.String("addr", config.AppConfig.Redis.Addr),
		zap.Int("db", config.AppConfig.Redis.DB),
		zap.String("ping", res),
	)
	return nil
}

func CloseRedis() {
	if RDB == nil {
		return
	}
	if err := RDB.Close(); err != nil {
		logger.Log.Warn("redis close failed", zap.Error(err))
	}
}
