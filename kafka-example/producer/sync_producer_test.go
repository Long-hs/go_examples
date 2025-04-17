package producer

import (
	"testing"

	"github.com/IBM/sarama"
)

// TestNewSyncProducerService 测试同步生产者服务的创建
func TestNewSyncProducerService(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "正常创建同步生产者",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewSyncProducerService()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSyncProducerService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if service != nil {
				defer func(service *SyncProducerService) {
					err := service.Close()
					if err != nil {
						t.Errorf("关闭同步生产者时发生错误: %v", err)
					}
				}(service)
			}
			if !tt.wantErr && service == nil {
				t.Error("NewSyncProducerService() 返回的生产者为空")
			}
		})
	}
}

// TestSyncProducerService_retrySend 测试同步重试发送
func TestSyncProducerService_retrySend(t *testing.T) {
	mockProducer := createMockSyncProducer(t)
	service := &SyncProducerService{
		producer: mockProducer,
		brokers:  []string{broker},
	}

	msg := mockMessage(syncTopic, "test retry message")

	tests := []struct {
		name       string
		maxRetries int
		setupMock  func()
		wantErr    bool
	}{
		{
			name:       "重试成功",
			maxRetries: 3,
			setupMock: func() {
				mockProducer.ExpectSendMessageAndFail(sarama.ErrLeaderNotAvailable)
				mockProducer.ExpectSendMessageAndFail(sarama.ErrLeaderNotAvailable)
				mockProducer.ExpectSendMessageAndSucceed()
			},
			wantErr: false,
		},
		{
			name:       "达到最大重试次数",
			maxRetries: 2,
			setupMock: func() {
				mockProducer.ExpectSendMessageAndFail(sarama.ErrLeaderNotAvailable)
				mockProducer.ExpectSendMessageAndFail(sarama.ErrLeaderNotAvailable)
				mockProducer.ExpectSendMessageAndFail(sarama.ErrLeaderNotAvailable)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			if err := service.retrySend(msg, tt.maxRetries); (err != nil) != tt.wantErr {
				t.Errorf("retrySend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestSyncProducerService_SendMessage 测试同步消息发送
func TestSyncProducerService_SendMessage(t *testing.T) {
	mockProducer := createMockSyncProducer(t)
	service := &SyncProducerService{
		producer: mockProducer,
		brokers:  []string{broker},
	}

	tests := []struct {
		name      string
		message   string
		setupMock func()
		wantErr   bool
	}{
		{
			name:    "成功发送消息",
			message: "test message",
			setupMock: func() {
				mockProducer.ExpectSendMessageAndSucceed()
			},
			wantErr: false,
		},
		{
			name:    "发送失败后重试成功",
			message: "test retry message",
			setupMock: func() {
				mockProducer.ExpectSendMessageAndFail(sarama.ErrLeaderNotAvailable)
				mockProducer.ExpectSendMessageAndSucceed()
			},
			wantErr: false,
		},
		{
			name:    "发送失败且重试失败",
			message: "test failed message",
			setupMock: func() {
				mockProducer.ExpectSendMessageAndFail(sarama.ErrLeaderNotAvailable)
				mockProducer.ExpectSendMessageAndFail(sarama.ErrLeaderNotAvailable)
				mockProducer.ExpectSendMessageAndFail(sarama.ErrLeaderNotAvailable)
				mockProducer.ExpectSendMessageAndFail(sarama.ErrLeaderNotAvailable)
				mockProducer.ExpectSendMessageAndFail(sarama.ErrLeaderNotAvailable)
				mockProducer.ExpectSendMessageAndFail(sarama.ErrLeaderNotAvailable)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置mock预期
			tt.setupMock()

			// 发送消息
			err := service.SendMessage(tt.message)

			// 验证结果
			if (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}

			// 验证mock生产者的期望是否满足
			if err := mockProducer.Close(); err != nil {
				t.Errorf("关闭mock生产者时发生错误: %v", err)
			}
		})
	}
}
