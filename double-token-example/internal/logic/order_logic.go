package logic

import (
	"context"
	"double-token-example/internal/config"
	"double-token-example/internal/kafka"
	"double-token-example/internal/model"
	"double-token-example/internal/repository"
	"double-token-example/pkg/utils"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"time"
)

type OrderLogic struct {
	producer  *kafka.Producer
	orderRepo *repository.OrderRepository
	goodsRepo *repository.GoodsRepository
}

func NewOrderLogic() *OrderLogic {
	return &OrderLogic{
		producer:  kafka.GetProducer(),
		orderRepo: repository.NewOrderRepository(),
		goodsRepo: repository.NewGoodsRepository(),
	}
}

func (l *OrderLogic) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) error {
	uuid := utils.GenerateUUID()
	order := &model.Order{
		CreatorID:     req.CreatorID,
		PaymentNo:     uuid,
		GoodsID:       req.GoodsID,
		Amount:        req.Amount,
		Quantity:      req.Quantity,
		PaymentMethod: req.PaymentMethod,
		Status:        model.OrderStatusPending,
		CreateTime:    time.Time{},
		UpdateTime:    time.Time{},
	}
	stock, err := l.goodsRepo.DecreaseSeckillStock(ctx, req.GoodsID, req.Quantity)
	if err != nil {
		return err
	}
	if !stock {
		return errors.New("库存不足")
	}

	marshal, err := json.Marshal(order)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: config.Cfg.Kafka.Topics.OrderTopic,
		Key:   nil,
		Value: sarama.StringEncoder(marshal),
	}

	err = l.producer.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
