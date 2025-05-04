package model

import (
	"time"
)

// User 用户结构体
type User struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"ID"`                         // 用户ID
	Username      string     `gorm:"type:varchar(50);not null;unique" json:"username"`           // 用户名
	Password      string     `gorm:"type:varchar(100);not null" json:"-"`                        // 密码（不返回给前端）
	Salt          string     `gorm:"type:varchar(32);not null" json:"-"`                         // 密码盐值
	Phone         string     `gorm:"type:varchar(20);not null;unique" json:"phone"`              // 手机号
	Email         string     `gorm:"type:varchar(100);unique" json:"email"`                      // 邮箱
	Status        int8       `gorm:"type:tinyint;default:1" json:"status"`                       // 状态：1-正常，0-禁用
	LastLoginTime *time.Time `gorm:"type:timestamp" json:"lastLoginTime"`                        // 最后登录时间
	LastLoginIP   string     `gorm:"type:varchar(50)" json:"lastLoginIP"`                        // 最后登录IP
	CreateTime    time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"createTime"` // 创建时间
	UpdateTime    time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updateTime"` // 更新时间
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

// RegisterRequest 用户注册请求结构体
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

// LoginRequest 用户登录请求结构体
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
