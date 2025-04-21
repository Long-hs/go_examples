package consumer

import (
	"fmt"
	"kafka-example/common"
	"log"

	"github.com/IBM/sarama"
)

// TraditionalConsumerService 表示传统的Kafka消费者服务
// 用于处理单个分区的消息消费
type TraditionalConsumerService struct {
	consumer  sarama.Consumer          // Kafka消费者实例
	topic     string                   // 订阅的主题
	broker    string                   // Kafka broker地址
	partition sarama.PartitionConsumer // 分区消费者
}

// NewTraditionalConsumerService 创建一个新的传统消费者服务实例
// topic: 要订阅的主题
// 返回: 传统消费者服务实例和可能的错误
func NewTraditionalConsumerService(topic string) (*TraditionalConsumerService, error) {
	log.Printf("[TraditionalConsumer] 正在初始化传统消费者服务...")

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest // 从最新的偏移量开始消费

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
		err := consumer.Close()
		if err != nil {
			return nil, err
		} // 关闭消费者
		return nil, fmt.Errorf("创建分区消费者失败: %v", err)
	}

	service := &TraditionalConsumerService{
		consumer:  consumer,
		topic:     common.SyncTopic,
		broker:    common.Broker,
		partition: partition,
	}

	log.Printf("[TraditionalConsumer] 传统消费者服务初始化成功，订阅主题: %s", topic)
	return service, nil
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
				log.Printf("[TraditionalConsumer] 收到消息: topic=%s, partition=%d, offset=%d, value=%s",
					msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
			case err := <-s.partition.Errors():
				log.Printf("[TraditionalConsumer] 消费错误: %v", err)
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
