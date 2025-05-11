package model

import (
	"time"
)

// RefreshToken 表示一条 Refresh Token 记录
type RefreshToken struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"column:user_id;not null;index:idx_user" json:"user_id"`
	JTI       string    `gorm:"column:jti;size:36;not null;uniqueIndex:idx_jti" json:"jti"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
