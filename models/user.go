package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"unique;not null"` // 用户名（唯一）
	PasswordHash string `gorm:"not null"`        // 密码哈希
	Role         string `gorm:"not null"`        // 角色: client/sales/trader
}

// 表名自定义（可选）
func (User) TableName() string {
	return "users"
}
