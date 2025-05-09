package db

import (
	"github.com/redis/go-redis/v9"
	"log"
)

var RedisDB *redis.Client

func init() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	RedisDB = client
	log.Println("Redis connection successfully")
}

// GetRedisDB 获取Redis连接
func GetRedisDB() *redis.Client {
	return RedisDB
}
