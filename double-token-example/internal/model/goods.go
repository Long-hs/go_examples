package model

import (
	"time"
)

// Goods 商品结构体
type Goods struct {
	ID          string    `bson:"_id" json:"id"`                  // 商品唯一标识
	Name        string    `bson:"name" json:"name"`               // 商品名称
	Description string    `bson:"description" json:"description"` // 商品描述
	Price       float64   `bson:"price" json:"price"`             // 商品价格
	Stock       int64     `bson:"stock" json:"stock"`             // 商品库存
	Status      int8      `bson:"status" json:"status"`           // 1-上架，0-下架
	CreatorID   int64     `bson:"creator_id" json:"creatorID"`    // 商品创建者ID
	UpdaterID   int64     `bson:"updater_id" json:"updaterID"`    // 商品更新者ID
	CreateTime  time.Time `bson:"create_time" json:"createTime"`  // 商品创建时间
	UpdateTime  time.Time `bson:"update_time" json:"updateTime"`  // 商品更新时间
}

// CreateGoodsRequest 创建商品请求
type CreateGoodsRequest struct {
	Name        string  `json:"name" binding:"required"`        // 商品名称，必填
	Description string  `json:"description" binding:"required"` // 商品描述，必填
	Price       float64 `json:"price" binding:"required,gt=0"`  // 商品价格，必填且大于0
	Stock       int64   `json:"stock" binding:"required,gte=0"` // 商品库存，必填且大于等于0
	CreatorID   int64   `json:"-"`                              // 商品创建者ID，不参与JSON序列化
}

// UpdateGoodsRequest 更新商品请求
type UpdateGoodsRequest struct {
	Name        string  `json:"name"`                       // 商品名称
	Description string  `json:"description"`                // 商品描述
	Price       float64 `json:"price" binding:"gt=0"`       // 商品价格，需大于0
	Stock       int64   `json:"stock" binding:"gte=0"`      // 商品库存，需大于等于0
	Status      int8    `json:"status" binding:"oneof=0 1"` // 商品状态，只能为0或1
	UpdaterID   int64   `json:"-"`                          // 商品更新者ID，不参与JSON序列化
}

// GetGoodsListRequest 获取商品列表请求
type GetGoodsListRequest struct {
	Name   string `form:"name"`                                  // 商品名称，用于筛选
	Status int8   `form:"status"`                                // 商品状态，用于筛选
	Page   int    `form:"page" binding:"required,gte=1"`         // 页码，必填且大于等于1
	Size   int    `form:"size" binding:"required,gte=1,lte=100"` // 每页数量，必填且大于等于1，小于等于100
}
