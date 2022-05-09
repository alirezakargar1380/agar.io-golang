package redis_db

import (
	"fmt"

	"github.com/go-redis/redis"
)

var Client *RedisDb

type RedisDb struct {
	Client *redis.Client
}

func (c *RedisDb) Test() {
	// c.Client.Set("test", "test", 0).Err()
	fmt.Println("hello im redis")
}
