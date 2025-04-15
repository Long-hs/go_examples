package main

import (
	"context"
	"kafka-example/consumer"
	"kafka-example/producer"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	producerService *producer.ProducerService
	consumerService *consumer.ConsumerService
)

func main() {
	// 初始化生产者服务（使用默认配置）
	var err error
	producerService, err = producer.DefaultProducerService()
	if err != nil {
		log.Fatal(err)
	}
	defer func(producerService *producer.ProducerService) {
		err := producerService.Close()
		if err != nil {

		}
	}(producerService)

	// 初始化消费者服务（使用默认配置）
	consumerService, err = consumer.DefaultConsumerService()
	if err != nil {
		log.Fatal(err)
	}
	defer func(consumerService *consumer.ConsumerService) {
		err := consumerService.Stop()
		if err != nil {
			log.Fatal(err)
		}
	}(consumerService)

	// 启动消费者服务
	if err := consumerService.Start(context.Background()); err != nil {
		log.Fatal(err)
	}

	// 创建 Gin 路由
	r := gin.Default()

	// 注册路由
	r.GET("/send", handleSendMessage)

	// 启动服务器
	log.Println("服务器启动在 :8081 端口")
	if err := r.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}

// handleSendMessage 处理发送消息的请求
func handleSendMessage(c *gin.Context) {
	// 获取消息参数
	message := c.Query("msg")
	if message == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "消息不能为空",
		})
		return
	}

	// 发送消息到 Kafka
	err := producerService.SendMessage("kafka-example-sync", message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "发送消息失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "消息发送成功",
		"data": gin.H{
			"content": message,
		},
	})
}
