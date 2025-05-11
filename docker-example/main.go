package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Conf struct {
	System struct {
		Name string `yaml:"name"`
	} `yaml:"system"`
}

func main() {
	r := gin.Default()

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
