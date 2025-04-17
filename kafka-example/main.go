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
	syncProducerService  *producer.SyncProducerService
	asyncProducerService *producer.AsyncProducerService
	consumerService      *consumer.ConsumerService
)

func main() {
	// 初始化生产者服务（使用默认配置）
	var err error
	syncProducerService, err = producer.NewSyncProducerService()
	if err != nil {
		log.Fatal(err)
	}
	asyncProducerService, err = producer.NewAsyncProducerService()
	if err != nil {
		log.Fatal(err)
	}
	defer func(sps *producer.SyncProducerService, aps *producer.AsyncProducerService) {
		err := sps.Close()
		if err != nil {
			log.Fatal(err)
		}
		err = aps.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(syncProducerService, asyncProducerService)

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

	err = consumerService.SwitchMode(consumer.ModeStandalone)
	if err != nil {
		log.Fatal(err)
	}

	// 启动消费者服务
	if err := consumerService.Start(context.Background()); err != nil {
		log.Fatal(err)
	}

	// 创建 Gin 路由
	r := gin.Default()

	// 注册路由
	r.GET("/sync", handleSyncSendMessage)
	r.GET("/async", handleAsyncSendMessage)

	// 启动服务器
	log.Println("服务器启动在 :8081 端口")
	if err := r.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}

// handleSyncSendMessage 同步处理发送消息的请求
func handleSyncSendMessage(c *gin.Context) {
	// 获取消息参数
	message := c.Query("msg")
	if message == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "消息不能为空",
		})
		return
	}

	// 发送消息到 Kafka
	err := syncProducerService.SendMessage(message)
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

// handleAsyncSendMessage 异步处理发送消息的请求
func handleAsyncSendMessage(c *gin.Context) {
	// 获取消息参数
	message := c.Query("msg")
	if message == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "消息不能为空",
		})
		return
	}

	// 发送消息到 Kafka
	asyncProducerService.SendMessage(message)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "消息发送成功",
		"data": gin.H{
			"content": message,
		},
	})
}
