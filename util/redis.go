package util

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func init() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	Ctx := context.Background()
	res, err := RDB.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("连接成功：", res)
}
