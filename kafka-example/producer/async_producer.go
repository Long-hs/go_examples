package producer

import (
	"encoding/binary"
	"errors"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

// 全局变量，用于确保错误处理协程只启动一次
var (
	errorHandler   sync.Once // 错误处理协程的Once对象
	successHandler sync.Once // 成功处理协程的Once对象
)

// AsyncProducerService 表示Kafka异步生产者服务
// 该服务提供异步发送消息的功能，不等待消息发送完成就返回
type AsyncProducerService struct {
	producer sarama.AsyncProducer // Kafka异步生产者实例
	brokers  []string             // Kafka broker地址列表
}

// NewAsyncProducerService 创建一个异步生产者服务
// 返回:
//   - *AsyncProducerService: 异步生产者服务实例
//   - error: 创建失败时返回错误
func NewAsyncProducerService() (*AsyncProducerService, error) {
	log.Printf("%s正在创建异步生产者: brokers=%v", logPrefixService, []string{broker})

	// 配置生产者参数
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // 等待所有副本确认
	config.Producer.Idempotent = true                // 启用幂等性，确保消息不会重复发送
	config.Producer.Retry.Max = 5                    // 最大重试次数
	config.Net.MaxOpenRequests = 1                   // 最大并发请求数
	config.Producer.Return.Errors = true             // 返回错误信息

	// 创建异步生产者
	producer, err := sarama.NewAsyncProducer([]string{broker}, config)
	if err != nil {
		log.Printf("%s创建异步生产者失败: %v", logPrefixService, err)
		return nil, err
	}

	s := &AsyncProducerService{
		producer: producer,
		brokers:  []string{broker},
	}

	// 启动错误处理协程
	errorHandler.Do(s.errorHanding)
	log.Printf("%s异步生产者创建成功", logPrefixService)
	return s, nil
}

// errorHanding 处理异步发送过程中的错误
// 启动一个协程监听错误通道，对可重试的错误进行重试
func (s *AsyncProducerService) errorHanding() {
	log.Printf("%s启动错误处理协程", logPrefixAsync)
	go func() {
		// 从错误通道中读取错误
		for result := range s.producer.Errors() {
			if isRetryableError(result.Err) {
				log.Printf("%s消息发送失败，准备重试: topic=%s, partition=%d, error=%v",
					logPrefixAsync, result.Msg.Topic, result.Msg.Partition, result.Err)
				err := s.retrySend(result.Msg, 5)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Printf("%s消息发送失败(不可重试): topic=%s, partition=%d, error=%v",
					logPrefixAsync, result.Msg.Topic, result.Msg.Partition, result.Err)
			}
		}
	}()
}

// SendMessage 异步发送消息
// 参数:
//   - message: 要发送的消息内容
func (s *AsyncProducerService) SendMessage(message string) {
	// 创建生产者消息
	msg := &sarama.ProducerMessage{
		Topic: asyncTopic,
		Value: sarama.StringEncoder(message),
	}
	log.Printf("%s发送消息: topic=%s, message=%s", logPrefixAsync, asyncTopic, message)
	// 将消息发送到输入通道
	s.producer.Input() <- msg
}

// retrySend 异步重试发送消息
// 参数:
//   - msg: 要重试发送的消息
//   - maxRetries: 最大重试次数
//
// 返回:
//   - error: 重试失败时返回错误
func (s *AsyncProducerService) retrySend(msg *sarama.ProducerMessage, maxRetries uint16) error {
	value := make([]byte, 2)
	var retryCount uint16 = 0

	// 检查消息头中是否已有重试计数
	for i := range msg.Headers {
		if string(msg.Headers[i].Key) == "retry_count" {
			retryCount = binary.BigEndian.Uint16(msg.Headers[i].Value)
			value := make([]byte, 2)
			binary.BigEndian.PutUint16(value, retryCount+1)
			msg.Headers[i].Value = value
			break
		}
	}

	// 检查是否超过最大重试次数
	if retryCount >= maxRetries {
		log.Printf("%s消息重发失败: topic=%s, 重试次数=%d, 达到重试次数上限=%d",
			logPrefixAsync, msg.Topic, retryCount, maxRetries)
		return errors.New("消息重发失败，超过重试次数上限")
	}

	// 如果是首次重试，添加重试计数头
	if retryCount == 0 {
		binary.BigEndian.PutUint16(value, retryCount+1)
		msg.Headers = append(msg.Headers, sarama.RecordHeader{
			Key:   []byte("retry_count"),
			Value: value,
		})
	}

	log.Printf("%s开始第%d次重试: topic=%s", logPrefixAsync, retryCount+1, msg.Topic)
	// 将消息重新发送到输入通道
	s.producer.Input() <- msg
	return nil
}

// Close 关闭异步生产者服务
// 返回:
//   - error: 关闭失败时返回错误
func (s *AsyncProducerService) Close() error {
	log.Printf("%s正在关闭异步生产者服务", logPrefixAsync)
	if err := s.producer.Close(); err != nil {
		log.Printf("%s关闭异步生产者失败: %v", logPrefixAsync, err)
		return err
	}
	return nil
}
