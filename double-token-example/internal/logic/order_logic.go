package logic

import (
	"context"
	"double-token-example/internal/model"
	"double-token-example/internal/repository"
	"double-token-example/pkg/utils"
	"time"
)

type OrderLogic struct {
	orderRepo *repository.OrderRepository
}

func NewOrderLogic() *OrderLogic {
	return &OrderLogic{
		orderRepo: repository.NewOrderRepository(),
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
	return l.orderRepo.Create(ctx, order)
}
