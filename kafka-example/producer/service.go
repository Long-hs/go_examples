package producer

import (
	"log"

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
	if err != nil {
		return err
	}

	log.Printf("消息发送成功: partition=%d, offset=%d\n", partition, offset)
	return nil
}

// Close 关闭生产者服务
func (s *ProducerService) Close() error {
	return s.producer.Close()
}
