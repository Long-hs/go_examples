package producer

import (
	"encoding/binary"
	"kafka-example/common"
	"testing"
	"time"

	"github.com/IBM/sarama"
)

// TestNewAsyncProducerService 测试异步生产者服务的创建
func TestNewAsyncProducerService(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "正常创建异步生产者",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			producer, err := NewAsyncProducerService()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAsyncProducerService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if producer != nil {
				defer func(producer *AsyncProducerService) {
					err := producer.Close()
					if err != nil {
						t.Errorf("关闭异步生产者时发生错误: %v", err)
					}
				}(producer)
			}
			if !tt.wantErr && producer == nil {
				t.Error("NewAsyncProducerService() 返回的生产者为空")
			}
		})
	}
}

// TestAsyncProducerService_SendMessage 测试异步发送消息
func TestAsyncProducerService_SendMessage(t *testing.T) {
	mockProducer := createMockAsyncProducer(t)
	service := &AsyncProducerService{
		producer: mockProducer,
		brokers:  []string{common.Broker},
	}

	// 设置预期
	mockProducer.ExpectInputAndSucceed()

	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "成功发送异步消息",
			message: "test async message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service.SendMessage(tt.message)
			time.Sleep(100 * time.Millisecond) // 等待消息处理
		})
	}
}

// TestAsyncProducerService_retrySend 测试异步重试发送
func TestAsyncProducerService_retrySend(t *testing.T) {
	mockProducer := createMockAsyncProducer(t)
	service := &AsyncProducerService{
		producer: mockProducer,
		brokers:  []string{common.Broker},
	}

	tests := []struct {
		name       string
		maxRetries uint16
		retryCount uint16
		setupMsg   func() *sarama.ProducerMessage
		wantErr    bool
	}{
		{
			name:       "首次重试",
			maxRetries: 5,
			retryCount: 1,
			setupMsg: func() *sarama.ProducerMessage {
				return mockMessage(common.AsyncTopic, "test retry message")
			},
			wantErr: false,
		},
		{
			name:       "第三次重试",
			maxRetries: 5,
			retryCount: 3,
			setupMsg: func() *sarama.ProducerMessage {
				msg := mockMessage(common.AsyncTopic, "test retry message")
				value := make([]byte, 2)
				binary.BigEndian.PutUint16(value, 2)
				msg.Headers = append(msg.Headers, sarama.RecordHeader{
					Key:   []byte("retry_count"),
					Value: value,
				})
				return msg
			},
			wantErr: false,
		},
		{
			name:       "超过最大重试次数",
			maxRetries: 5,
			retryCount: 6,
			setupMsg: func() *sarama.ProducerMessage {
				msg := mockMessage(common.AsyncTopic, "test retry message")
				value := make([]byte, 2)
				binary.BigEndian.PutUint16(value, 5)
				msg.Headers = append(msg.Headers, sarama.RecordHeader{
					Key:   []byte("retry_count"),
					Value: value,
				})
				return msg
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProducer.ExpectInputAndSucceed()
			msg := tt.setupMsg()

			err := service.retrySend(msg, tt.maxRetries)
			if (err != nil) != tt.wantErr {
				t.Errorf("retrySend() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// 验证重试计数
				var foundRetryCount bool
				var actualRetryCount uint16
				for _, header := range msg.Headers {
					if string(header.Key) == "retry_count" {
						foundRetryCount = true
						actualRetryCount = binary.BigEndian.Uint16(header.Value)
						break
					}
				}

				expectedRetryCount := tt.retryCount
				if !foundRetryCount && expectedRetryCount == 1 {
					t.Error("首次重试未添加重试计数头")
				} else if foundRetryCount && actualRetryCount != expectedRetryCount {
					t.Errorf("重试计数不正确，期望 %d，实际 %d", expectedRetryCount, actualRetryCount)
				}
			}
		})
	}
}
