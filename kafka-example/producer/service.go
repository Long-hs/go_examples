package producer

import (
	"log"
	"time"

	"github.com/IBM/sarama"
)

// ProducerService 表示Kafka生产者服务
type ProducerService struct {
	producer sarama.SyncProducer
	brokers  []string
}

// DefaultProducerService 创建一个默认配置的生产者服务
func DefaultProducerService() (*ProducerService, error) {
	return NewProducerService([]string{"localhost:9092"})
}

// NewProducerService 创建一个新的生产者服务
func NewProducerService(brokers []string) (*ProducerService, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &ProducerService{
		producer: producer,
		brokers:  brokers,
	}, nil
}

// SendMessage 发送消息到指定的topic
func (s *ProducerService) SendMessage(topic string, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := s.producer.SendMessage(msg)
	//如果发送错误，手动进行重试
	if err != nil {
		return retrySend(s.producer, msg, 5)
	}

	log.Printf("消息发送成功: partition=%d, offset=%d\n", partition, offset)
	return nil
}

// Close 关闭生产者服务
func (s *ProducerService) Close() error {
	return s.producer.Close()
}

// retrySend 重试发送消息
func retrySend(producer sarama.SyncProducer, msg *sarama.ProducerMessage, maxRetries int) error {

	for i := 0; i <= maxRetries; i++ {
		partition, offset, err := producer.SendMessage(msg)
		if err == nil {
			log.Printf("消息重发成功，重发次数%d, partition=%d, offset=%d\n", i+1, partition, offset)
			return nil
		}
		if i < maxRetries && isRetryableError(err) {
			time.Sleep(time.Duration(1<<i) * 100 * time.Millisecond) //指数退避
			log.Printf("消息重发失败，重发次数%d, 错误信息: %v\n", i+1, err)
			continue
		}
		//不可重发错误或达到最大重试次数
		return err
	}
	return nil
}

// isRetryableError 判断错误是否可重试
func isRetryableError(err error) bool {
	return err != nil
}
