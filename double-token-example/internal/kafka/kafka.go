package kafka

import (
	"context"
	"double-token-example/internal/model"
	"double-token-example/internal/repository"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

var kafkaServer *KafkaSever

// KafkaSever Kafka 服务器结构体
type KafkaSever struct {
	GroupConsumer sarama.ConsumerGroup        // Kafka 消费者组
	SyncProducer  sarama.SyncProducer         // Kafka 同步生产者
	Topics        []string                    // 订阅的主题列表
	brokers       []string                    // Kafka 代理地址列表
	config        *sarama.Config              // Kafka 配置
	groupID       string                      // 消费者组 ID
	orderRepo     *repository.OrderRepository // 订单仓储
	goodsRepo     *repository.GoodsRepository // 商品仓储
}

// OrderGroupHandler 实现 sarama.ConsumerGroupHandler 接口
type OrderGroupHandler struct {
	orderRepo *repository.OrderRepository
}

// GoodsGroupHandler 实现 sarama.ConsumerGroupHandler 接口
type GoodsGroupHandler struct {
	goodsRepo *repository.GoodsRepository
}

// Setup 在消费者组会话开始前调用
func (h OrderGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup 在消费者组会话结束后调用
func (h OrderGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 处理从 Kafka 接收到的消息
func (h OrderGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Received order message: topic=%s, partition=%d, offset=%d, value=%s",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value))

		order := model.Order{}
		err := json.Unmarshal(msg.Value, &order)
		if err != nil {
			log.Printf("Failed to unmarshal order message: %v, message content: %s", err, string(msg.Value))
			continue
		}

		// 使用注入的orderRepo处理订单
		if err := h.orderRepo.Create(context.Background(), &order); err != nil {
			log.Printf("Failed to create order: %v", err)
			continue
		}

		// 标记消息为已处理
		sess.MarkMessage(msg, "")
		log.Printf("Marked order message as processed: topic=%s, partition=%d, offset=%d",
			msg.Topic, msg.Partition, msg.Offset)
	}
	return nil
}

// Setup 在消费者组会话开始前调用
func (h GoodsGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup 在消费者组会话结束后调用
func (h GoodsGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 处理从 Kafka 接收到的消息
func (h GoodsGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Received goods message: topic=%s, partition=%d, offset=%d, value=%s",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value))

		goods := model.Goods{}
		err := json.Unmarshal(msg.Value, &goods)
		if err != nil {
			log.Printf("Failed to unmarshal goods message: %v, message content: %s", err, string(msg.Value))
			continue
		}

		// 使用注入的goodsRepo处理商品
		if err := h.goodsRepo.Create(context.Background(), &goods); err != nil {
			log.Printf("Failed to create goods: %v", err)
			continue
		}

		// 标记消息为已处理
		sess.MarkMessage(msg, "")
		log.Printf("Marked goods message as processed: topic=%s, partition=%d, offset=%d",
			msg.Topic, msg.Partition, msg.Offset)
	}
	return nil
}

// InitKafkaServer 创建新的 Kafka 服务器实例
func InitKafkaServer() {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Retry.Max = 5
	cfg.Net.MaxOpenRequests = 1
	cfg.Producer.Idempotent = true
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	// 创建订单消费者组
	orderGroup, err := sarama.NewConsumerGroup([]string{"127.0.0.1:9092"}, "order_consumer_group", cfg)
	if err != nil {
		log.Fatalf("Failed to create order consumer group: %v", err)
	}

	// 创建商品消费者组
	goodsGroup, err := sarama.NewConsumerGroup([]string{"127.0.0.1:9092"}, "goods_consumer_group", cfg)
	if err != nil {
		log.Fatalf("Failed to create goods consumer group: %v", err)
	}

	// 创建生产者
	producer, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, cfg)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}

	// 创建仓储
	orderRepo := repository.NewOrderRepository()
	goodsRepo := repository.NewGoodsRepository()

	server := &KafkaSever{
		SyncProducer: producer,
		Topics:       []string{"goods_topic", "order_topic"},
		brokers:      []string{"127.0.0.1:9092"},
		config:       cfg,
		groupID:      "double_token_example_group",
		orderRepo:    orderRepo,
		goodsRepo:    goodsRepo,
	}

	// 启动订单消费者
	go func() {
		log.Println("Starting order consumer...")
		handler := OrderGroupHandler{
			orderRepo: orderRepo,
		}
		for {
			if err := orderGroup.Consume(context.Background(), []string{"order_topic"}, handler); err != nil {
				log.Printf("Error from order consumer: %v", err)
			}
		}
	}()

	// 启动商品消费者
	go func() {
		log.Println("Starting goods consumer...")
		handler := GoodsGroupHandler{
			goodsRepo: goodsRepo,
		}
		for {
			if err := goodsGroup.Consume(context.Background(), []string{"goods_topic"}, handler); err != nil {
				log.Printf("Error from goods consumer: %v", err)
			}
		}
	}()

	kafkaServer = server
	log.Println("Kafka server initialized successfully")
}

func GetKafkaServer() *KafkaSever {
	return kafkaServer
}
