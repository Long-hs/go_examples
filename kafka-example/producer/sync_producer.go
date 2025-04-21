package producer

import (
	"errors"
	"kafka-example/common"
	"log"
	"time"

	"github.com/IBM/sarama"
)

// SyncProducerService 表示Kafka同步生产者服务
// 该服务提供同步发送消息的功能，确保消息发送成功后才返回
type SyncProducerService struct {
	producer sarama.SyncProducer // Kafka同步生产者实例
	brokers  []string            // Kafka broker地址列表
}

// NewSyncProducerService 创建一个同步生产者服务
// 返回:
//   - *SyncProducerService: 同步生产者服务实例
//   - error: 创建失败时返回错误
func NewSyncProducerService() (*SyncProducerService, error) {
	log.Printf("%s正在创建同步生产者: brokers=%v", common.LogPrefixService, []string{common.Broker})

	// 配置生产者参数
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true          // 要求返回发送成功确认
	config.Producer.Idempotent = true                // 启用幂等性，确保消息不会重复发送
	config.Net.MaxOpenRequests = 1                   // 限制最大并发请求数
	config.Producer.RequiredAcks = sarama.WaitForAll // 等待所有副本确认
	config.Producer.Retry.Max = 5                    // 最大重试次数

	// 创建同步生产者
	producer, err := sarama.NewSyncProducer([]string{common.Broker}, config)
	if err != nil {
		log.Printf("%s创建同步生产者失败: %v", common.LogPrefixService, err)
		return nil, err
	}

	log.Printf("%s同步生产者创建成功", common.LogPrefixService)
	return &SyncProducerService{
		producer: producer,
		brokers:  []string{common.Broker},
	}, nil
}

// SendMessage 同步发送消息到指定的topic
// 参数:
//   - message: 要发送的消息内容
//
// 返回:
//   - error: 发送失败时返回错误
func (s *SyncProducerService) SendMessage(message string) error {
	// 创建生产者消息
	msg := &sarama.ProducerMessage{
		Topic: common.SyncTopic,
		Value: sarama.StringEncoder(message),
	}

	log.Printf("%s开始发送消息: topic=%s, message=%s", common.LogPrefixSync, common.SyncTopic, message)

	// 发送消息并等待结果
	partition, offset, err := s.producer.SendMessage(msg)
	if err != nil {
		log.Printf("%s消息发送失败，准备重试: topic=%s, error=%v", common.LogPrefixSync, common.SyncTopic, err)
		return s.retrySend(msg, 5) // 失败时进行重试
	}

	log.Printf("%s消息发送成功: topic=%s, partition=%d, offset=%d",
		common.LogPrefixSync, common.SyncTopic, partition, offset)
	return nil
}

// retrySend 同步重试发送消息
// 参数:
//   - msg: 要重试发送的消息
//   - maxRetries: 最大重试次数
//
// 返回:
//   - error: 重试失败时返回错误
func (s *SyncProducerService) retrySend(msg *sarama.ProducerMessage, maxRetries int) error {
	for i := 0; i <= maxRetries; i++ {
		// 达到最大重试次数时，不再重试
		if i == maxRetries {
			// 重试失败或达到最大重试次数
			log.Printf("%s重试终止: topic=%s, 重试次数=%d, 错误=%v",
				common.LogPrefixSync, msg.Topic, i+1, errors.New("消息重发失败，超过重试次数上限"))
			return errors.New("消息重发失败，超过重试次数上限")
		}

		log.Printf("%s开始第%d次重试: topic=%s", common.LogPrefixSync, i+1, msg.Topic)

		// 尝试发送消息
		partition, offset, err := s.producer.SendMessage(msg)
		if err == nil {
			log.Printf("%s重试发送成功: topic=%s, 重试次数=%d, partition=%d, offset=%d",
				common.LogPrefixSync, msg.Topic, i+1, partition, offset)
			return nil
		}

		// 如果错误可重试，则等待后继续
		if common.IsRetryableError(err) {
			backoff := time.Duration(1<<i) * 100 * time.Millisecond // 指数退避
			log.Printf("%s重试发送失败: topic=%s, 重试次数=%d, 等待时间=%v, 错误=%v",
				common.LogPrefixSync, msg.Topic, i+1, backoff, err)
			time.Sleep(backoff)
			continue
		}

	}
	return nil
}

// Close 关闭同步生产者服务
// 返回:
//   - error: 关闭失败时返回错误
func (s *SyncProducerService) Close() error {
	log.Printf("%s正在关闭同步生产者服务", common.LogPrefixSync)
	if err := s.producer.Close(); err != nil {
		log.Printf("%s关闭同步生产者失败: %v", common.LogPrefixSync, err)
		return err
	}
	return nil
}
