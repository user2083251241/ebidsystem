package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID       uint    `gorm:"not null"`           //订单创建者（卖家或客户），外键关联 User.ID
	Symbol       string  `gorm:"not null"`           //股票代码（如 AAPL）
	Quantity     int     `gorm:"not null"`           //股票数量
	Price        float64 `gorm:"type:decimal(10,2)"` //股票价格（市价单可为空）
	OrderType    string  `gorm:"not null"`           //market, limit
	Direction    string  `gorm:"not null"`           //方向: buy/sell
	Status       string  `gorm:"default:'pending'"`  //状态：draft/pending/filled/cancelled
	DraftBySales uint    //销售草稿的销售用户ID（0表示非草稿）
	ApprovedBy   uint    //卖家批准用户ID
}

// 表名自定义（可选）
func (Order) TableName() string {
	return "orders"
}
