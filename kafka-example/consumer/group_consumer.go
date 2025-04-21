package consumer

import (
	"context"
	"fmt"
	"kafka-example/common"
	"log"

	"github.com/IBM/sarama"
)

// GroupConsumerService 表示Kafka消费者组服务
// 用于处理消费者组模式下的消息消费
type GroupConsumerService struct {
	group   sarama.ConsumerGroup // Kafka消费者组实例
	topics  []string             // 订阅的主题列表
	brokers []string             // Kafka broker地址列表
	handler consumerGroupHandler // 消费者组处理器
}

// consumerGroupHandler 实现 sarama.ConsumerGroupHandler 接口
// 用于处理消费者组的生命周期事件和消息消费
type consumerGroupHandler struct{}

// NewGroupConsumerService 创建一个新的消费者组服务实例
// topics: 要订阅的主题列表
// 返回: 消费者组服务实例和可能的错误
func NewGroupConsumerService(topics []string) (*GroupConsumerService, error) {
	log.Printf("[GroupConsumer] 正在初始化消费者组服务...")

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest // 从最新的偏移量开始消费

	consumerGroup, err := sarama.NewConsumerGroup([]string{common.Broker}, "group_consumer", config)
	if err != nil {
		log.Printf("[GroupConsumer] 创建消费者组失败: %v", err)
		return nil, fmt.Errorf("创建消费者组失败: %v", err)
	}

	consumerService := &GroupConsumerService{
		group:   consumerGroup,
		topics:  topics,
		brokers: []string{common.Broker},
		handler: consumerGroupHandler{},
	}

	log.Printf("[GroupConsumer] 消费者组服务初始化成功，订阅主题: %v", topics)
	return consumerService, nil
}

// Setup 在消费者组会话开始前调用
// 用于准备消费者组会话
func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	log.Printf("[GroupConsumer] 消费者组会话开始")
	return nil
}

// Cleanup 在消费者组会话结束后调用
// 用于清理消费者组会话
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	log.Printf("[GroupConsumer] 消费者组会话结束")
	return nil
}

// ConsumeClaim 处理分配给消费者的消息
// 这是实际处理消息的地方
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log.Printf("[GroupConsumer] 开始消费分区 %d 的消息", claim.Partition())

	for msg := range claim.Messages() {
		log.Printf("[GroupConsumer] 收到消息: topic=%s, partition=%d, offset=%d, value=%s",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		sess.MarkMessage(msg, "")
	}

	log.Printf("[GroupConsumer] 分区 %d 的消息消费完成", claim.Partition())
	return nil
}

// Start 启动消费者组服务
// ctx: 上下文，用于控制服务的生命周期
func (g *GroupConsumerService) Start(ctx context.Context) error {
	log.Printf("[GroupConsumer] 正在启动消费者组服务...")

	go func() {
		for {
			if err := g.group.Consume(ctx, g.topics, g.handler); err != nil {
				log.Printf("[GroupConsumer] 消费错误: %v", err)
				continue
			}
			if ctx.Err() != nil {
				log.Printf("[GroupConsumer] 消费者组服务已停止")
				return
			}
		}
	}()

	log.Printf("[GroupConsumer] 消费者组服务启动成功")
	return nil
}

// Stop 停止消费者组服务
// 关闭消费者组连接
func (g *GroupConsumerService) Stop() error {
	log.Printf("[GroupConsumer] 正在停止消费者组服务...")

	if err := g.group.Close(); err != nil {
		log.Printf("[GroupConsumer] 停止服务失败: %v", err)
		return fmt.Errorf("停止消费者组服务失败: %v", err)
	}

	log.Printf("[GroupConsumer] 消费者组服务已成功停止")
	return nil
}
