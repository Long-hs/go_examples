package repository

import (
	"cache-example/db"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// InfoRepository 信息仓储接口
type InfoRepository interface {
	GetFromMysql(id int64) (*db.Info, error)
	GetFromCache(id int64, ctx context.Context) (*db.Info, error)
	SaveToCache(info *db.Info, ctx context.Context) error
	UpdateToMysql(info *db.Info) error
	DeleteFromCache(id int64, ctx context.Context) error
}

// infoRepository 信息仓储实现
type infoRepository struct{}

func (r *infoRepository) DeleteFromCache(id int64, ctx context.Context) error {
	key := fmt.Sprintf("info:%d", id)
	// 删除缓存（设置过期时间避免大key问题）
	if err := db.RedisDB.PExpire(ctx, key, time.Millisecond*1).Err(); err != nil {
		log.Printf("[Cache] 设置过期时间失败: %v", err)
		return fmt.Errorf("设置过期时间失败: %v", err)
	}
	log.Printf("[Cache] 设置过期时间成功: key=%s", key)
	return nil
}

// NewInfoRepository 创建信息仓储实例
func NewInfoRepository() InfoRepository {
	return &infoRepository{}
}

// GetFromMysql 从MySQL获取信息
func (r *infoRepository) GetFromMysql(id int64) (*db.Info, error) {
	info := &db.Info{}
	if err := db.DB.Table(info.TableName()).Where("id = ?", id).First(info).Error; err != nil {
		log.Printf("[DB] 查询失败: %v", err)
		return nil, fmt.Errorf("查询失败: %v", err)
	}
	return info, nil
}

// GetFromCache 从缓存获取信息
func (r *infoRepository) GetFromCache(id int64, ctx context.Context) (*db.Info, error) {
	info := &db.Info{}
	key := fmt.Sprintf("info:%d", id)

	// 获取缓存
	result, err := db.RedisDB.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			log.Printf("[Cache] 缓存未命中: key=%s", key)
			return nil, fmt.Errorf("缓存未命中")
		}
		log.Printf("[Cache] 获取缓存失败: %v", err)
		return nil, fmt.Errorf("获取缓存失败: %v", err)
	}

	// 解析JSON
	if err := json.Unmarshal([]byte(result), info); err != nil {
		log.Printf("[Cache] 解析缓存数据失败: %v", err)
		return nil, fmt.Errorf("解析缓存数据失败: %v", err)
	}

	log.Printf("[Cache] 缓存命中: key=%s, value=%+v", key, info)
	return info, nil
}

// SaveToCache 保存信息到缓存
func (r *infoRepository) SaveToCache(info *db.Info, ctx context.Context) error {
	key := fmt.Sprintf("info:%d", info.ID)
	data, err := json.Marshal(info)
	if err != nil {
		log.Printf("[Cache] 序列化数据失败: %v", err)
		return fmt.Errorf("序列化数据失败: %v", err)
	}

	// 设置缓存，过期时间可配置
	if err := db.RedisDB.Set(ctx, key, data, time.Minute*5).Err(); err != nil {
		log.Printf("[Cache] 保存缓存失败: %v", err)
		return fmt.Errorf("保存缓存失败: %v", err)
	}

	return nil
}

func (r *infoRepository) UpdateToMysql(info *db.Info) error {
	if err := db.DB.Table(info.TableName()).Where("id = ?", info.ID).Updates(info).First(info).Error; err != nil {
		log.Printf("[DB] 更新失败: %v", err)
		return fmt.Errorf("更新失败: %v", err)
	}
	return nil
}
