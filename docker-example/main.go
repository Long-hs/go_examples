package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

type Conf struct {
	System struct {
		Name string `yaml:"name"`
	} `yaml:"system"`
}

func initDB() {
	dsn := fmt.Sprintf("root:root@tcp(172.20.0.2:3306)/double_token_example?charset=utf8mb4&parseTime=True&loc=Local")

	var err error
	_, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	log.Println("MySQL connected successfully")
}

func main() {
	r := gin.Default()

	initDB()

	byteData, err := os.ReadFile("settings.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var conf Conf
	err = yaml.Unmarshal(byteData, &conf)
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{"code": 0, "msg": "看到消息就说明部署成功了", "data": gin.H{"name": conf.System.Name}})
	})
	r.Run(":5000")
}
