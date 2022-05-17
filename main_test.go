package main

import (
	"testing"

	redis_db "github.com/alirezakargar1380/agar.io-golang/app/service"
	"github.com/go-redis/redis"
)

func TestRedis(t *testing.T) {
	redis_db.Client = &redis_db.RedisDb{
		Client: redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		}),
	}
	_, err := redis_db.Client.Client.Ping().Result()
	if err != nil {
		panic(err)
	}
}
