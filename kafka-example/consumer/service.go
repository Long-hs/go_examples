package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

// 消费模式
const (
	ModeStandalone = "standalone" // 独立模式
	ModeGroup      = "group"      // 消费者组模式
)

// 日志前缀
const logPrefixService = "消费者服务"

// ConsumerService 表示Kafka消费者服务
type ConsumerService struct {
	group   sarama.ConsumerGroup
	topics  []string
	groupID string
	brokers []string
	mode    string // 消费模式
}

// DefaultConsumerService 创建一个默认配置的消费者服务
func DefaultConsumerService() (*ConsumerService, error) {
	return NewConsumerService(
		[]string{"localhost:9092"},     // 默认broker地址
		"my-consumer-group",            // 默认消费者组ID
		[]string{"kafka-example-sync"}, // 默认topic
	)
}

// NewConsumerService 创建一个新的消费者服务
func NewConsumerService(brokers []string, groupID string, topics []string) (*ConsumerService, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &ConsumerService{
		group:   group,
		topics:  topics,
		groupID: groupID,
		brokers: brokers,
	}, nil
}

// Start 启动消费者服务
func (s *ConsumerService) Start(ctx context.Context) error {
	// 处理消息的函数
	messageHandler := func(msg *sarama.ConsumerMessage) error {
		log.Printf("收到消息: topic=%s, partition=%d, offset=%d, value=%s\n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		return nil
	}

	// 实现 ConsumerGroupHandler 接口
	handler := &ConsumerGroupHandler{
		handler: messageHandler,
	}

	// 启动消费组
	go func() {
		for {
			if err := s.group.Consume(ctx, s.topics, handler); err != nil {
				log.Printf("消费组错误: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	return nil
}

// ConsumerGroupHandler 实现 sarama.ConsumerGroupHandler 接口
type ConsumerGroupHandler struct {
	handler func(message *sarama.ConsumerMessage) error
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if err := h.handler(message); err != nil {
			log.Printf("处理消息失败: %v", err)
		}
		session.MarkMessage(message, "")
	}
	return nil
}

// Stop 停止消费者服务
func (s *ConsumerService) Stop() error {
	return s.group.Close()
}

// SwitchMode 切换消费模式
func (s *ConsumerService) SwitchMode(mode string) error {
	if mode != ModeStandalone && mode != ModeGroup {
		return fmt.Errorf("不支持的消费模式: %s", mode)
	}

	s.mode = mode
	log.Printf("%s切换到%s模式", logPrefixService, mode)
	return nil
}
