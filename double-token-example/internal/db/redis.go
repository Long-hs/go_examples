package db

import (
	"context"
	"double-token-example/internal/config"
	"double-token-example/internal/model"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

var (
	redisDB *redis.Client
	ctx     = context.Background()
)

func init() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	redisDB = client
	log.Println("Redis connection successfully")
	err := CreateBloomFilter()
	if err != nil {
		log.Printf("Failed to create bloom filter:%v", err)
	}
}

// GetRedisDB 获取Redis连接
func GetRedisDB() *redis.Client {
	return redisDB
}

// CreateBloomFilter 创建布隆过滤器
func CreateBloomFilter() error {
	name := config.Cfg.Redis.Bloom.Name
	exists, err := redisDB.Exists(ctx, name).Result()
	if err != nil {
		return fmt.Errorf("检查布隆过滤器是否存在失败: %v", err)
	}

	// 2. 若存在则删除
	if exists > 0 {
		log.Printf("布隆过滤器 %s 已存在，准备删除", name)
		if _, err := redisDB.Del(ctx, name).Result(); err != nil {
			return fmt.Errorf("删除布隆过滤器失败: %v", err)
		}
		log.Printf("布隆过滤器 %s 已删除", name)
	}
	_, err = redisDB.BFReserve(ctx, name, config.Cfg.Redis.Bloom.ErrorRate, config.Cfg.Redis.Bloom.ExpectedItems).Result()
	if err != nil {
		return err
	}
	mySQL := GetMySQL()
	var usernames []interface{}
	err = mySQL.Model(&model.User{}).Select("username").Find(&usernames).Error
	if err != nil {
		return err
	}
	log.Println("usernames:", usernames)
	redisDB.BFMAdd(ctx, name, usernames...)
	return nil
}
