package databases

import (
	"encoding/json"
	"fmt"

	"github.com/alirezakargar1380/agar.io-golang/app/stars"

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

func (c *RedisDb) AddStar(key string, roomId string) {
	var stars *stars.Star = &stars.Star{
		Star: make(map[string]map[string]bool),
	}
	starsMap, _ := c.Client.Get("stars").Result()
	if starsMap == "" {
		stars.Star[roomId] = make(map[string]bool)
		stars.Star[roomId][key] = true
	} else {
		json.Unmarshal([]byte(starsMap), &stars.Star)
		if stars.Star[roomId] == nil {
			stars.Star[roomId] = make(map[string]bool)
		}
		stars.Star[roomId][key] = true
	}

	pp, err := json.Marshal(stars.Star)
	if err != nil {
		fmt.Println(err)
	}

	err = c.Client.Set("stars", pp, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func (c *RedisDb) CountStars(roomId string) int {
	var stars *stars.Star = &stars.Star{
		Star: make(map[string]map[string]bool),
	}
	vv, _ := c.Client.Get("stars").Result()
	json.Unmarshal([]byte(vv), &stars.Star)

	return len(stars.Star[roomId])
}

func (c *RedisDb) DeleteStart(roomId string, keys []string) {
	var stars *stars.Star = &stars.Star{
		Star: make(map[string]map[string]bool),
	}
	vv, _ := c.Client.Get("stars").Result()
	if vv == "" {
		stars.Star[roomId] = make(map[string]bool)
	}
	json.Unmarshal([]byte(vv), &stars.Star)

	for _, E := range keys {
		delete(stars.Star[roomId], E)
	}

	pp, _ := json.Marshal(stars.Star)

	err := c.Client.Set("stars", pp, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func (c *RedisDb) GetStars(roomId string) map[string]bool {
	var stars *stars.Star = &stars.Star{
		Star: make(map[string]map[string]bool),
	}
	vv, _ := c.Client.Get("stars").Result()
	if vv == "" {
		stars.Star[roomId] = make(map[string]bool)
	}
	json.Unmarshal([]byte(vv), &stars.Star)

	return stars.Star[roomId]
}
