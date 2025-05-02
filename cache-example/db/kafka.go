package db

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

// KafkaServer 全局 Kafka 服务器实例
var KafkaServer *KafkaSever

// KafkaSever Kafka 服务器结构体
type KafkaSever struct {
	GroupConsumer sarama.ConsumerGroup // Kafka 消费者组
	SyncProducer  sarama.SyncProducer  // Kafka 同步生产者
	Topics        []string             // 订阅的主题列表
	brokers       []string             // Kafka 代理地址列表
	config        *sarama.Config       // Kafka 配置
	groupID       string               // 消费者组 ID
}

// ConsumerGroupHandler 实现 sarama.ConsumerGroupHandler 接口
type ConsumerGroupHandler struct{}

// Setup 在消费者组会话开始前调用
func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {

	return nil
}

// Cleanup 在消费者组会话结束后调用
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {

	return nil
}

// ConsumeClaim 处理从 Kafka 接收到的消息
func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Received Kafka message: topic=%s, partition=%d, offset=%d, value=%s",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value))

		// 解析消息为 Info 对象
		info := &Info{}
		err := json.Unmarshal(msg.Value, info)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v, message content: %s", err, string(msg.Value))
			continue
		}

		// 更新数据库
		err = DB.Model(&Info{}).Where("id = ?", info.ID).Updates(info).Error
		if err != nil {
			log.Printf("Failed to update database: %v, info: %+v", err, info)
			continue
		}
		log.Printf("Successfully updated database for info: %+v", info)

		// 标记消息为已处理
		sess.MarkMessage(msg, "")
		log.Printf("Marked message as processed: topic=%s, partition=%d, offset=%d",
			msg.Topic, msg.Partition, msg.Offset)
	}
	return nil
}

// NewKafkaServer 创建新的 Kafka 服务器实例
func NewKafkaServer(brokers []string, topics []string, groupID string) *KafkaSever {
	log.Printf("Initializing Kafka server with brokers: %v, topics: %v, groupID: %s",
		brokers, topics, groupID)

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Retry.Max = 5
	config.Net.MaxOpenRequests = 1
	config.Producer.Idempotent = true
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// 创建消费者组
	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}

	// 创建生产者
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}

	server := &KafkaSever{
		GroupConsumer: group,
		SyncProducer:  producer,
		Topics:        topics,
		brokers:       brokers,
		config:        config,
		groupID:       groupID,
	}

	// 启动消费者
	go func() {
		log.Println("Starting Kafka consumer...")
		handler := ConsumerGroupHandler{}
		for {
			if err := group.Consume(context.Background(), topics, handler); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
		}
	}()

	log.Println("Kafka server initialized successfully")
	return server
}
