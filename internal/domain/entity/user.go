package entity

import (
	"gorm.io/gorm"
)

// 一般用户表：
type User struct {
	gorm.Model
	ID           uint   `gorm:"primarykey"`
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`      //密码哈希
	Role         string `gorm:"not null"`      //seller 卖方, sales 销售, trader 交易员, client 客户
	IsDeleted    bool   `gorm:"default:false"` //软删除标记
}

// // 表名自定义（可选）：
// func (User) TableName() string {
// 	return "users"
// }
