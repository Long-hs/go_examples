package db

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mysqlDB *gorm.DB

// InitMySQL 初始化MySQL连接
func InitMySQL() {
	dsn := fmt.Sprintf("root:root@tcp(127.0.0.1:8806)/double_token_example?charset=utf8mb4&parseTime=True&loc=Local")

	var err error
	mysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	log.Println("MySQL connected successfully")
}

// GetMySQL 获取MySQL连接
func GetMySQL() *gorm.DB {
	return mysqlDB
}
