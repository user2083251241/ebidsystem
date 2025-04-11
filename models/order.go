package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID    uint    `gorm:"not null"`           // 外键关联 User.ID
	Symbol    string  `gorm:"not null"`           // 股票代码（如 AAPL）
	Quantity  int     `gorm:"not null"`           // 数量
	Price     float64 `gorm:"type:decimal(10,2)"` // 价格（市价单可为空）
	OrderType string  `gorm:"not null"`           // 订单类型: market/limit
	Direction string  `gorm:"not null"`           // 方向: buy/sell
	Status    string  `gorm:"default:pending"`    // 状态: pending/filled/cancelled
}

// 表名自定义（可选）
func (Order) TableName() string {
	return "orders"
}
