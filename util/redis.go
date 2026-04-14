package util

import (
	"MYshop/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
	Ctx = context.Background()
)

func init() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.Redis.Addr,
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
	})
	res, err := RDB.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("连接成功：", res)
}
