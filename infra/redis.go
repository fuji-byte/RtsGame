package infra

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // パスワードがある場合は設定
		DB:       0,
	})
	return rdb
}
