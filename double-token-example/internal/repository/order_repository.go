package repository

import (
	"context"
	"double-token-example/internal/db"
	"double-token-example/internal/model"
	"gorm.io/gorm"
	"sync"
)

type OrderRepository struct {
	db *gorm.DB
}

var (
	orderRepo *OrderRepository
	orderOnce sync.Once
)

func NewOrderRepository() *OrderRepository {
	orderOnce.Do(func() {
		orderRepo = &OrderRepository{
			db: db.GetMySQL(),
		}
	})
	return orderRepo
}

// Create 创建订单
func (r *OrderRepository) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// Update 更新订单
func (r *OrderRepository) Update(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// Delete 删除订单
func (r *OrderRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Order{}, id).Error
}

// GetByID 根据ID获取订单
func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByOrderNo 根据订单号获取订单
func (r *OrderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByUserID 获取用户的订单列表
func (r *OrderRepository) GetByUserID(ctx context.Context, userID int64, page, size int) ([]*model.Order, int64, error) {
	var orders []*model.Order
	var total int64

	// 获取总数
	err := r.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("create_time DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}
