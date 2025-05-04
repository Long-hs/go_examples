package model

import (
	"time"
)

// OrderStatus 订单状态常量
const (
	OrderStatusPending   int8 = 1 // 待支付
	OrderStatusPaid      int8 = 2 // 已支付
	OrderStatusCancelled int8 = 3 // 已取消
	OrderStatusRefunded  int8 = 4 // 已退款
)

// Order 订单支付信息结构体
type Order struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"ID"`                         // 支付ID
	CreatorID     int64      `gorm:"not null" json:"creatorID"`                                  // 创建者ID
	GoodsID       string     `gorm:"type:varchar(48);not null" json:"goodsID"`                   // 商品ID                // 商品ID
	PaymentNo     string     `gorm:"type:varchar(36);not null;unique" json:"paymentNo"`          // 支付流水号
	Quantity      int8       `gorm:"type:tinyint;not null" json:"quantity"`                      // 购买数量（至少1件）
	Amount        float64    `gorm:"type:decimal(10,2);not null" json:"amount"`                  // 支付金额
	PaymentMethod int8       `gorm:"type:tinyint;not null" json:"paymentMethod"`                 // 支付方式：1-支付宝，2-微信
	Status        int8       `gorm:"type:tinyint;not null" json:"status"`                        // 支付状态：1-待支付，2-支付成功，3-支付失败
	PayTime       *time.Time `gorm:"type:timestamp" json:"payTime"`                              // 支付时间
	CreateTime    time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"createTime"` // 创建时间
	UpdateTime    time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updateTime"` // 更新时间
}

// TableName 指定表名
func (Order) TableName() string {
	return "order"
}

type CreateOrderRequest struct {
	CreatorID     int64   `json:"creatorId"`
	GoodsID       string  `json:"goodsId"`
	Quantity      int8    `json:"quantity"`
	Amount        float64 `json:"amount"`
	PaymentMethod int8    `json:"paymentMethod"`
}

type UpdateOrderRequest struct {
	CreatorID     int64   `json:"creatorID"`
	PaymentNo     string  `json:"paymentNo"`
	Quantity      int8    `json:"quantity"`
	Amount        float64 `json:"amount"`
	PaymentMethod int8    `json:"paymentMethod"`
	Status        int8    `json:"status"`
	PayTime       string  `json:"payTime"`
}
