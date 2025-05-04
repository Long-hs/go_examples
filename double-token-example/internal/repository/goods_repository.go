package repository

import (
	"context"
	"double-token-example/internal/db"
	"double-token-example/internal/model"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GoodsRepository struct {
	collection *mongo.Collection
}

func NewGoodsRepository() *GoodsRepository {
	return &GoodsRepository{
		collection: db.GetMongoDBCollection("goods"),
	}
}

// Create 创建商品
func (r *GoodsRepository) Create(ctx context.Context, goods *model.Goods) error {
	_, err := r.collection.InsertOne(ctx, goods)
	return err
}

// GetList 获取商品列表
func (r *GoodsRepository) GetList(ctx context.Context, req *model.GetGoodsListRequest) ([]*model.Goods, int64, error) {
	// 构建查询条件
	filter := bson.M{}
	if req.Name != "" {
		filter["name"] = bson.M{"$regex": req.Name, "$options": "i"}
	}
	if req.Status == 0 || req.Status == 1 {
		filter["status"] = req.Status
	}

	// 获取总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 设置分页
	opts := options.Find().
		SetSkip(int64((req.Page - 1) * req.Size)).
		SetLimit(int64(req.Size)).
		SetSort(bson.D{{"create_time", -1}})

	// 查询数据
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cursor, ctx)

	var goodsList []*model.Goods
	if err = cursor.All(ctx, &goodsList); err != nil {
		return nil, 0, err
	}

	return goodsList, total, nil
}

// GetByID 根据ID获取商品
func (r *GoodsRepository) GetByID(ctx context.Context, id string) (*model.Goods, error) {
	var goods model.Goods
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&goods)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("商品不存在")
		}
		return nil, err
	}

	return &goods, nil
}

// Update 更新商品
func (r *GoodsRepository) Update(ctx context.Context, goods *model.Goods) error {
	objectID, err := primitive.ObjectIDFromHex(goods.ID)
	if err != nil {
		return errors.New("无效的商品ID")
	}

	update := bson.M{
		"$set": bson.M{
			"name":        goods.Name,
			"description": goods.Description,
			"price":       goods.Price,
			"stock":       goods.Stock,
			"status":      goods.Status,
			"updater_id":  goods.UpdaterID,
			"update_time": goods.UpdateTime,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// Delete 删除商品
func (r *GoodsRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// UpdateStock 更新商品库存（使用乐观锁）
func (r *GoodsRepository) UpdateStock(ctx context.Context, id primitive.ObjectID, stock int64, version int64) error {
	filter := bson.M{
		"_id":     id,
		"version": version,
	}
	update := bson.M{
		"$set": bson.M{
			"stock":       stock,
			"version":     version + 1,
			"update_time": time.Now(),
		},
	}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("update failed, version mismatch")
	}
	return nil
}
