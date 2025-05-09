package kafka

import (
	"context"
	"double-token-example/internal/config"
	"double-token-example/internal/model"
	"double-token-example/internal/repository"
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
)

type Consumer struct {
	orderConsumer sarama.ConsumerGroup
	goodsConsumer sarama.ConsumerGroup
}

var consumer *Consumer

func init() {
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	orderGroup, err := sarama.NewConsumerGroup(config.Cfg.Kafka.Brokers, config.Cfg.Kafka.Groups.OrderGroup, cfg)
	if err != nil {
		log.Fatalf("Failed to create order consumer group: %v", err)
	}

	goodsGroup, err := sarama.NewConsumerGroup(config.Cfg.Kafka.Brokers, config.Cfg.Kafka.Groups.GoodsGroup, cfg)
	if err != nil {
		log.Fatalf("Failed to create goods consumer group: %v", err)
	}
	go func() {
		for {
			if err := orderGroup.Consume(context.Background(), []string{config.Cfg.Kafka.Topics.OrderTopic}, OrderGroupHandler{repository.NewOrderRepository()}); err != nil {
				log.Fatalf("Error from order consumer: %v", err)
			}
		}
	}()
	go func() {
		for {
			if err := goodsGroup.Consume(context.Background(), []string{config.Cfg.Kafka.Topics.GoodsTopic}, GoodsGroupHandler{repository.NewGoodsRepository()}); err != nil {
				log.Fatalf("Error from goods consumer: %v", err)
			}
		}
	}()
	consumer = &Consumer{
		orderConsumer: orderGroup,
		goodsConsumer: goodsGroup,
	}
}

func GetConsumer() *Consumer {
	return consumer
}

// OrderGroupHandler 实现 sarama.ConsumerGroupHandler 接口
type OrderGroupHandler struct {
	orderRepo *repository.OrderRepository
}

// GoodsGroupHandler 实现 sarama.ConsumerGroupHandler 接口
type GoodsGroupHandler struct {
	goodsRepo *repository.GoodsRepository
}

// Setup 在消费者组会话开始前调用
func (h OrderGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup 在消费者组会话结束后调用
func (h OrderGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 处理从 Kafka 接收到的消息
func (h OrderGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		order := model.Order{}
		err := json.Unmarshal(msg.Value, &order)
		if err != nil {
			log.Printf("Failed to unmarshal order message: %v, message content: %s", err, string(msg.Value))
			continue
		}

		// 使用注入的orderRepo处理订单
		if err := h.orderRepo.Create(context.Background(), &order); err != nil {
			log.Printf("Failed to create order: %v", err)
			continue
		}

		// 标记消息为已处理
		sess.MarkMessage(msg, "")
	}
	return nil
}

// Setup 在消费者组会话开始前调用
func (h GoodsGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup 在消费者组会话结束后调用
func (h GoodsGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 处理从 Kafka 接收到的消息
func (h GoodsGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		goods := model.Goods{}
		err := json.Unmarshal(msg.Value, &goods)
		if err != nil {
			log.Printf("Failed to unmarshal goods message: %v, message content: %s", err, string(msg.Value))
			continue
		}

		// 使用注入的goodsRepo处理商品
		if err := h.goodsRepo.Create(context.Background(), &goods); err != nil {
			log.Printf("Failed to create goods: %v", err)
			continue
		}

		// 标记消息为已处理
		sess.MarkMessage(msg, "")
	}
	return nil
}
