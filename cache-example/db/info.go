package db

import (
	"time"
)

// Info 信息表模型
type Info struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`    // 主键ID
	Name       string    `gorm:"type:varchar(50);not null" json:"name"` // 名称
	CreateTime time.Time `gorm:"type:timestamp" json:"create_time"`     // 创建时间
	UpdateTime time.Time `gorm:"type:timestamp" json:"update_time"`     // 更新时间
}

// TableName 指定表名
func (Info) TableName() string {
	return "info"
}
