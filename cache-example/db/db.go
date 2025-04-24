package db

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	MysqlHost = "127.0.0.1:8806"
	MysqlUser = "root"
	MysqlPass = "root"
	MysqlDb   = "cache_example"
	RedisHost = "localhost:6379"
	RedisPass = ""
)

var (
	DB      *gorm.DB
	RedisDB *redis.Client
)

func NewMysqlDB() {
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", MysqlUser, MysqlPass, MysqlHost, MysqlDb)
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Mysql connected")
	}
	DB = db
}

func NewRedisDB() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     RedisHost,
		Password: RedisPass,
		DB:       0,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Redis connected")
	}
	RedisDB = rdb
}
