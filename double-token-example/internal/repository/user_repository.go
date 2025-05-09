package repository

import (
	"context"
	"double-token-example/internal/config"
	"double-token-example/internal/db"
	"double-token-example/internal/model"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"sync"
)

type UserRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

var (
	userRepo *UserRepository
	userOnce sync.Once
)

func NewUserRepository() *UserRepository {
	userOnce.Do(func() {
		userRepo = &UserRepository{
			db:    db.GetMySQL(),
			redis: db.GetRedisDB(),
		}
	})
	return userRepo
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	// 使用数据库事务确保数据一致性
	err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.WithContext(ctx).Create(user).Error; err != nil {
			log.Printf("创建用户失败: %v", err)
			return err
		}
		err := r.CreateToRedis(ctx, user)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// CreateToRedis 将写入 Redis 缓存
func (r *UserRepository) CreateToRedis(ctx context.Context, user *model.User) error {
	result, err := r.redis.Exists(ctx, "user:"+user.Username).Result()
	if result == 1 {
		return nil
	}
	bloomFilterName := config.Cfg.Redis.Bloom.Name
	luaScript := `
			-- 添加到布隆过滤器
			redis.call('BF.ADD', KEYS[1], ARGV[1])
			
			-- 设置用户哈希（使用 HSET 替代 HMSET）
			local userKey = 'user:' .. ARGV[1]
			redis.call('HSET', userKey, 'username', ARGV[1])
			redis.call('HSET', userKey, 'phone', ARGV[2])
			redis.call('HSET', userKey, 'email', ARGV[3])
			
			-- 设置过期时间
			redis.call('EXPIRE', userKey, ARGV[4])
			
			return 1
		`
	_, err = r.redis.Eval(ctx, luaScript,
		[]string{bloomFilterName},   // KEYS[1]
		user.Username,               // ARGV[1]
		user.Phone,                  // ARGV[2]
		user.Email,                  // ARGV[3]
		config.Cfg.Redis.UserExpiry, // ARGV[4]
	).Result()

	if err != nil {
		log.Printf("执行 Redis Lua 脚本失败: %v", err)
		return fmt.Errorf("更新缓存失败: %w", err)
	}
	return nil
}

// Update 更新用户
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}
	err = r.CreateToRedis(ctx, &user)
	if err != nil {
		log.Printf("更新 Redis 缓存失败: %v", err)
	}
	return &user, nil
}

// GetByPhone 根据手机号获取用户
func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// ExistByUsernameFromBloomFilter 判断用户名是否存在于布隆过滤器中
func (r *UserRepository) ExistByUsernameFromBloomFilter(ctx context.Context, username string) (bool, error) {
	result, err := r.redis.BFExists(ctx, config.Cfg.Redis.Bloom.Name, username).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

// ExistByUsernameFromRedis 判断用户名是否存在于 Redis 中的缓存
func (r *UserRepository) ExistByUsernameFromRedis(ctx context.Context, username string) (bool, error) {
	key := "user:" + username
	result, err := r.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
