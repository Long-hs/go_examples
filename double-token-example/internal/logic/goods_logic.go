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
	"log"
	"time"

	"github.com/IBM/sarama"
)

type GoodsLogic struct {
	goodsRepo *repository.GoodsRepository
	producer  *kafka.Producer
}

func NewGoodsLogic() *GoodsLogic {
	return &GoodsLogic{
		goodsRepo: repository.NewGoodsRepository(),
		producer:  kafka.GetProducer(),
	}
}

// CreateGoods 创建商品
func (l *GoodsLogic) CreateGoods(ctx context.Context, req *model.CreateGoodsRequest) error {
	uuid := utils.GenerateUUID()
	layout := "2006-01-02 15:04:05"
	startTime, err := time.Parse(layout, req.StartTime)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	endTime, err := time.Parse(layout, req.EndTime)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	goods := &model.Goods{
		ID:          uuid,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Status:      1,
		CreatorID:   req.CreatorID,
		UpdaterID:   req.CreatorID,
		StartTime:   startTime,
		EndTime:     endTime,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	marshal, err := json.Marshal(goods)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: config.Cfg.Kafka.Topics.GoodsTopic,
		Value: sarama.StringEncoder(marshal),
	}
	err = l.producer.Send(msg)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	err = l.goodsRepo.CreateSeckillGoodsCache(ctx, goods)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

// GetGoodsList 获取商品列表
func (l *GoodsLogic) GetGoodsList(ctx context.Context, req *model.GetGoodsListRequest) ([]*model.Goods, int64, error) {
	return l.goodsRepo.GetList(ctx, req)
}

// GetGoodsDetail 获取商品详情
func (l *GoodsLogic) GetGoodsDetail(ctx context.Context, id string) (*model.Goods, error) {
	return l.goodsRepo.GetByID(ctx, id)
}

// UpdateGoods 更新商品
func (l *GoodsLogic) UpdateGoods(ctx context.Context, id string, req *model.UpdateGoodsRequest) error {
	// 检查商品是否存在
	goods, err := l.goodsRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 更新商品信息
	if req.Name != "" {
		goods.Name = req.Name
	}
	if req.Description != "" {
		goods.Description = req.Description
	}
	if req.Price > 0 {
		goods.Price = req.Price
	}
	if req.Stock >= 0 {
		goods.Stock = req.Stock
	}
	if req.Status == 0 || req.Status == 1 {
		goods.Status = req.Status
	}
	goods.UpdaterID = req.UpdaterID
	goods.UpdateTime = time.Now()

	return l.goodsRepo.Update(ctx, goods)
}

// DeleteGoods 删除商品
func (l *GoodsLogic) DeleteGoods(ctx context.Context, id string, userID int64) error {
	// 检查商品是否存在
	goods, err := l.goodsRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查是否有权限删除
	if goods.CreatorID != userID {
		return errors.New("无权限删除该商品")
	}

	return l.goodsRepo.Delete(ctx, id)
}
