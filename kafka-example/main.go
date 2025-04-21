package main

import (
	"context"
	"kafka-example/common"
	"kafka-example/consumer"
	"kafka-example/producer"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 全局服务实例
var (
	syncProducerService        *producer.SyncProducerService        // 同步生产者服务
	asyncProducerService       *producer.AsyncProducerService       // 异步生产者服务
	groupConsumerService       *consumer.GroupConsumerService       // 消费者组服务
	traditionalConsumerService *consumer.TraditionalConsumerService // 传统消费者服务
)

func main() {
	log.Printf("[Main] 正在启动Kafka示例服务...")

	// 初始化生产者服务（使用默认配置）
	var err error
	log.Printf("[Main] 正在初始化同步生产者服务...")
	syncProducerService, err = producer.NewSyncProducerService()
	if err != nil {
		log.Fatalf("[Main] 初始化同步生产者服务失败: %v", err)
	}
	log.Printf("[Main] 同步生产者服务初始化成功")

	log.Printf("[Main] 正在初始化异步生产者服务...")
	asyncProducerService, err = producer.NewAsyncProducerService()
	if err != nil {
		log.Fatalf("[Main] 初始化异步生产者服务失败: %v", err)
	}
	log.Printf("[Main] 异步生产者服务初始化成功")

	// 确保在程序退出时关闭生产者服务
	defer func(sps *producer.SyncProducerService, aps *producer.AsyncProducerService) {
		log.Printf("[Main] 正在关闭生产者服务...")
		if err := sps.Close(); err != nil {
			log.Printf("[Main] 关闭同步生产者服务失败: %v", err)
		}
		if err := aps.Close(); err != nil {
			log.Printf("[Main] 关闭异步生产者服务失败: %v", err)
		}
		log.Printf("[Main] 生产者服务已关闭")
	}(syncProducerService, asyncProducerService)

	// 初始化消费者服务
	log.Printf("[Main] 正在初始化消费者组服务...")
	groupConsumerService, err = consumer.NewGroupConsumerService([]string{common.AsyncTopic})
	if err != nil {
		log.Fatalf("[Main] 初始化消费者组服务失败: %v", err)
	}
	log.Printf("[Main] 消费者组服务初始化成功")

	// 确保在程序退出时关闭消费者组服务
	defer func(g *consumer.GroupConsumerService) {
		log.Printf("[Main] 正在关闭消费者组服务...")
		if err := g.Stop(); err != nil {
			log.Printf("[Main] 关闭消费者组服务失败: %v", err)
		}
		log.Printf("[Main] 消费者组服务已关闭")
	}(groupConsumerService)

	log.Printf("[Main] 正在初始化传统消费者服务...")
	traditionalConsumerService, err = consumer.NewTraditionalConsumerService(common.SyncTopic)
	if err != nil {
		log.Fatalf("[Main] 初始化传统消费者服务失败: %v", err)
	}
	log.Printf("[Main] 传统消费者服务初始化成功")

	// 确保在程序退出时关闭传统消费者服务
	defer func(g *consumer.TraditionalConsumerService) {
		log.Printf("[Main] 正在关闭传统消费者服务...")
		if err := g.Stop(); err != nil {
			log.Printf("[Main] 关闭传统消费者服务失败: %v", err)
		}
		log.Printf("[Main] 传统消费者服务已关闭")
	}(traditionalConsumerService)

	// 启动消费者服务
	log.Printf("[Main] 正在启动消费者服务...")
	if err := groupConsumerService.Start(context.Background()); err != nil {
		log.Fatalf("[Main] 启动消费者组服务失败: %v", err)
	}
	if err := traditionalConsumerService.Start(); err != nil {
		log.Fatalf("[Main] 启动传统消费者服务失败: %v", err)
	}
	log.Printf("[Main] 消费者服务启动成功")

	// 创建 Gin 路由
	log.Printf("[Main] 正在初始化Web服务器...")
	r := gin.Default()

	// 注册路由
	r.GET("/sync", handleSyncSendMessage)
	r.GET("/async", handleAsyncSendMessage)
	log.Printf("[Main] 路由注册完成")

	// 启动服务器
	log.Printf("[Main] 服务器启动在 :8081 端口")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("[Main] 启动服务器失败: %v", err)
	}
}

// handleSyncSendMessage 同步处理发送消息的请求
// 接收GET请求，从查询参数中获取消息内容并同步发送到Kafka
func handleSyncSendMessage(c *gin.Context) {
	log.Printf("[Main] 收到同步发送消息请求")

	// 获取消息参数
	message := c.Query("msg")
	if message == "" {
		log.Printf("[Main] 同步发送消息失败: 消息内容为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "消息不能为空",
		})
		return
	}

	// 发送消息到 Kafka
	log.Printf("[Main] 正在同步发送消息: %s", message)
	err := syncProducerService.SendMessage(message)
	if err != nil {
		log.Printf("[Main] 同步发送消息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "发送消息失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	log.Printf("[Main] 同步发送消息成功: %s", message)
	c.JSON(http.StatusOK, gin.H{
		"message": "消息发送成功",
		"data": gin.H{
			"content": message,
		},
	})
}

// handleAsyncSendMessage 异步处理发送消息的请求
// 接收GET请求，从查询参数中获取消息内容并异步发送到Kafka
func handleAsyncSendMessage(c *gin.Context) {
	log.Printf("[Main] 收到异步发送消息请求")

	// 获取消息参数
	message := c.Query("msg")
	if message == "" {
		log.Printf("[Main] 异步发送消息失败: 消息内容为空")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "消息不能为空",
		})
		return
	}

	// 发送消息到 Kafka
	log.Printf("[Main] 正在异步发送消息: %s", message)
	asyncProducerService.SendMessage(message)

	// 返回成功响应
	log.Printf("[Main] 异步发送消息成功: %s", message)
	c.JSON(http.StatusOK, gin.H{
		"message": "消息发送成功",
		"data": gin.H{
			"content": message,
		},
	})
}
