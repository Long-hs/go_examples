package producer

import (
	"testing"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
)

// createMockSyncProducer 创建模拟的同步生产者
func createMockSyncProducer(t *testing.T) *mocks.SyncProducer {
	return mocks.NewSyncProducer(t, nil)
}

// createMockAsyncProducer 创建模拟的异步生产者
func createMockAsyncProducer(t *testing.T) *mocks.AsyncProducer {
	return mocks.NewAsyncProducer(t, nil)
}

// mockMessage 创建测试消息
func mockMessage(topic string, message string) *sarama.ProducerMessage {
	return &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
}
