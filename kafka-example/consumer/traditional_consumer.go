package consumer

import (
	"fmt"
	"kafka-example/common"
	"log"
	"time"

	"github.com/IBM/sarama"
)

// TraditionalConsumerService 表示传统的Kafka消费者服务
// 用于处理单个分区的消息消费
type TraditionalConsumerService struct {
	consumer  sarama.Consumer          // Kafka消费者实例
	topic     string                   // 订阅的主题
	broker    string                   // Kafka broker地址
	partition sarama.PartitionConsumer // 分区消费者
	stopChan  chan struct{}            // 停止通道
	config    *sarama.Config           // Kafka配置
}

// NewTraditionalConsumerService 创建一个新的传统消费者服务实例
// topic: 要订阅的主题
// 返回: 传统消费者服务实例和可能的错误
func NewTraditionalConsumerService(topic string) (*TraditionalConsumerService, error) {
	log.Printf("[TraditionalConsumer] 正在初始化传统消费者服务...")

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// 消费者配置
	config.Consumer.Offsets.Initial = sarama.OffsetNewest // 从最新的偏移量开始消费
	config.Consumer.Offsets.AutoCommit.Enable = false     // 禁用自动提交

	// 创建消费者
	consumer, err := sarama.NewConsumer([]string{common.Broker}, config)
	if err != nil {
		log.Printf("[TraditionalConsumer] 创建消费者失败: %v", err)
		return nil, fmt.Errorf("创建消费者失败: %v", err)
	}

	// 创建分区消费者
	partition, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("[TraditionalConsumer] 创建分区消费者失败: %v", err)
		if err := consumer.Close(); err != nil {
			log.Printf("[TraditionalConsumer] 关闭消费者失败: %v", err)
		}
		return nil, fmt.Errorf("创建分区消费者失败: %v", err)
	}

	service := &TraditionalConsumerService{
		consumer:  consumer,
		topic:     common.SyncTopic,
		broker:    common.Broker,
		partition: partition,
		stopChan:  make(chan struct{}),
	}

	log.Printf("[TraditionalConsumer] 传统消费者服务初始化成功，订阅主题: %s", topic)
	return service, nil
}

// processMessage 处理单条消息
// 包含重试机制和错误处理
func (s *TraditionalConsumerService) processMessage(msg *sarama.ConsumerMessage) error {
	maxRetries := 3
	retryInterval := 1 * time.Second

	for i := 0; i < maxRetries; i++ {
		// 处理消息
		log.Printf("[TraditionalConsumer] 处理消息: topic=%s, partition=%d, offset=%d, value=%s",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value))

		// 模拟消息处理
		time.Sleep(100 * time.Millisecond)

		// 如果处理成功，提交偏移量
		if err := s.commitOffset(msg); err != nil {
			log.Printf("[TraditionalConsumer] 提交偏移量失败: %v", err)
			if i < maxRetries-1 {
				log.Printf("[TraditionalConsumer] 将在 %v 后重试", retryInterval)
				time.Sleep(retryInterval)
				continue
			}
			return fmt.Errorf("提交偏移量失败: %v", err)
		}

		return nil
	}

	return fmt.Errorf("消息处理失败，已达到最大重试次数")
}

// commitOffset 提交消息偏移量
func (s *TraditionalConsumerService) commitOffset(msg *sarama.ConsumerMessage) error {
	// 这里可以添加持久化逻辑，确保偏移量不会丢失
	log.Printf("[TraditionalConsumer] 提交偏移量: topic=%s, partition=%d, offset=%d",
		msg.Topic, msg.Partition, msg.Offset)
	return nil
}

// Start 启动消费者服务
// 开始消费消息
func (s *TraditionalConsumerService) Start() error {
	log.Printf("[TraditionalConsumer] 正在启动消费者服务...")

	go func() {
		log.Printf("[TraditionalConsumer] 开始消费消息...")
		for {
			select {
			case msg := <-s.partition.Messages():
				if err := s.processMessage(msg); err != nil {
					log.Printf("[TraditionalConsumer] 处理消息失败: %v", err)
					// 可以在这里添加死信队列逻辑
				}
			case err := <-s.partition.Errors():
				log.Printf("[TraditionalConsumer] 消费错误: %v", err)
				// 可以在这里添加错误恢复逻辑
			case <-s.stopChan:
				log.Printf("[TraditionalConsumer] 收到停止信号")
				return
			}
		}
	}()

	log.Printf("[TraditionalConsumer] 消费者服务启动成功")
	return nil
}

// Stop 停止消费者服务
// 关闭分区消费者和消费者连接
func (s *TraditionalConsumerService) Stop() error {
	log.Printf("[TraditionalConsumer] 正在停止消费者服务...")

	// 发送停止信号
	close(s.stopChan)

	// 等待一段时间确保消息处理完成
	time.Sleep(1 * time.Second)

	// 关闭分区消费者
	if err := s.partition.Close(); err != nil {
		log.Printf("[TraditionalConsumer] 关闭分区消费者失败: %v", err)
		return fmt.Errorf("关闭分区消费者失败: %v", err)
	}
	log.Printf("[TraditionalConsumer] 分区消费者已关闭")

	// 关闭消费者
	if err := s.consumer.Close(); err != nil {
		log.Printf("[TraditionalConsumer] 关闭消费者失败: %v", err)
		return fmt.Errorf("关闭消费者失败: %v", err)
	}
	log.Printf("[TraditionalConsumer] 消费者已关闭")

	log.Printf("[TraditionalConsumer] 消费者服务已成功停止")
	return nil
}
