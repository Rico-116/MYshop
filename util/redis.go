package util

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func init() {
	Rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.103:6379",
		Password: "123456",
		DB:       0,
	})
	Ctx := context.Background()
	res, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}
	RDB = Rdb
	fmt.Println("连接成功：", res)
}
