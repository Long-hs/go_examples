package repository

import (
	"context"
	"double-token-example/internal/db"
	"double-token-example/internal/model"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GoodsRepository struct {
	collection *mongo.Collection
}

var (
	goodsRepo *GoodsRepository
	goodsOnce sync.Once
)

func NewGoodsRepository() *GoodsRepository {
	goodsOnce.Do(func() {
		goodsRepo = &GoodsRepository{
			collection: db.GetMongoDBCollection("goods"),
		}
	})
	return goodsRepo
}

// Create 创建商品
func (r *GoodsRepository) Create(ctx context.Context, goods *model.Goods) error {
	// 插入商品数据
	_, err := r.collection.InsertOne(ctx, goods)
	if err != nil {
		log.Printf("创建商品失败: %v", err)
		return fmt.Errorf("创建商品失败: %v", err)
	}
	return nil
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

// CreateSeckillGoodsCache 创建秒杀商品缓存
func (r *GoodsRepository) CreateSeckillGoodsCache(ctx context.Context, goods *model.Goods) error {
	seckillKey := fmt.Sprintf("seckill:%s", goods.ID)
	stockKey := fmt.Sprintf("seckill:stock:%s", goods.ID)
	soldKey := fmt.Sprintf("seckill:sold:%s", goods.ID)

	// 使用 Pipeline 批量执行命令
	pipe := db.GetRedisDB().Pipeline()

	// 1. 存储秒杀商品基本信息
	pipe.HSet(ctx, seckillKey, map[string]interface{}{
		"stock":      goods.Stock,
		"start_time": goods.StartTime.Unix(),
		"end_time":   goods.EndTime.Unix(),
		"price":      goods.Price,
		"status":     0, // 初始状态
	})

	// 2. 设置库存计数器
	pipe.Set(ctx, stockKey, goods.Stock, 0)

	// 3. 设置已售数量计数器
	pipe.Set(ctx, soldKey, 0, 0)

	// 4. 设置过期时间
	expiration := goods.EndTime.Unix() - time.Now().Unix()
	if expiration <= 0 {
		expiration = 60 * 60 * 24
	}
	pipe.Expire(ctx, seckillKey, time.Duration(expiration)*time.Second)
	pipe.Expire(ctx, stockKey, time.Duration(expiration)*time.Second)
	pipe.Expire(ctx, soldKey, time.Duration(expiration)*time.Second)
	// 执行 Pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("创建秒杀商品缓存失败: %v", err)
		return fmt.Errorf("创建秒杀商品缓存失败: %v", err)
	}
	return nil
}

// DecreaseSeckillStock 扣减秒杀商品库存
func (r *GoodsRepository) DecreaseSeckillStock(ctx context.Context, goodsID string, quantity int64) (bool, error) {
	stockKey := fmt.Sprintf("seckill:stock:%s", goodsID)
	soldKey := fmt.Sprintf("seckill:sold:%s", goodsID)

	// 使用 Lua 脚本保证原子性
	script := `
		local stock = tonumber(redis.call('GET', KEYS[1]))
		local sold = tonumber(redis.call('GET', KEYS[2]))
		local quantity = tonumber(ARGV[1])
		
		if stock < quantity then
			return 0
		end
		
		redis.call('DECRBY', KEYS[1], quantity)
		redis.call('INCRBY', KEYS[2], quantity)
		return 1
	`

	result, err := db.GetRedisDB().Eval(ctx, script, []string{stockKey, soldKey}, quantity).Result()
	if err != nil {
		log.Printf("扣减库存失败: %v", err)
		return false, fmt.Errorf("扣减库存失败: %v", err)
	}

	success := result.(int64) == 1
	if success {
		log.Printf("扣减库存成功: %s, 数量: %d", goodsID, quantity)
	} else {
		log.Printf("扣减库存失败: %s, 库存不足", goodsID)
	}

	return success, nil
}
