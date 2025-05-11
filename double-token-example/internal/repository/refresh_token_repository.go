package repository

import (
	"context"
	"double-token-example/internal/db"
	"double-token-example/internal/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type RefreshTokenRepository struct {
	db      *gorm.DB
	redisDB *redis.Client
}

var (
	refreshTokenRepo *RefreshTokenRepository
	refreshTokenOnce sync.Once
)

func NewRefreshTokenRepository() *RefreshTokenRepository {
	refreshTokenOnce.Do(func() {
		refreshTokenRepo = &RefreshTokenRepository{
			db:      db.GetMySQL(),
			redisDB: db.GetRedisDB(),
		}
	})
	return refreshTokenRepo
}

func (r *RefreshTokenRepository) GetTokenByJTI(jti string) (string, error) {
	return "", nil
}

// Create 持久化RefreshToken
func (r *RefreshTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	err := r.db.WithContext(ctx).Create(token).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RefreshTokenRepository) Refresh(ctx context.Context, token *model.RefreshToken) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := r.db.WithContext(ctx).Delete(&model.RefreshToken{}, "user_id = ?", token.UserID).Error
		if err != nil {
			return err
		}
		err = r.db.WithContext(ctx).Create(token).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (r *RefreshTokenRepository) DeleteByUserID(ctx context.Context, jti string, expiresAt time.Time, userID int64) error {
	log.Printf("expiresAt: %v，now: %v", expiresAt, time.Now().Unix())
	r.db.Transaction(func(tx *gorm.DB) error {
		err := r.db.WithContext(ctx).Delete(&model.RefreshToken{}, "user_id = ?", userID).Error
		if err != nil {
			return err
		}
		_, err = r.redisDB.Set(ctx, "jwt_blacklist:"+jti, 1, time.Until(expiresAt)).Result()
		if err != nil {
			return err
		}
		return nil
	})

	return nil
}
