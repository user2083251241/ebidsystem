package models

import (
	"gorm.io/gorm"
)

// 一般用户表：
type User struct {
	gorm.Model
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`      //密码哈希
	Role         string `gorm:"not null"`      //seller 卖方, sales 销售, trader 交易员, client 客户
	IsDeleted    bool   `gorm:"default:false"` //软删除标记
}

// 卖家与销售的授权关联表：
type SellerSalesAuthorization struct {
	gorm.Model
	SellerID      uint   `gorm:"not null"`                                                     //卖家用户ID
	SalesID       uint   `gorm:"not null"`                                                     //销售用户ID
	Authorization string `gorm:"type:enum('pending','approved','rejected');default:'pending'"` //授权状态
}

// 表名自定义（可选）：
func (User) TableName() string {
	return "users"
}
