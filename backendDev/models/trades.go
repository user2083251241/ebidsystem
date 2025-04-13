package models

import (
	"time"

	"gorm.io/gorm"
)

// Trade 成交记录表
type Trade struct {
	gorm.Model
	BuyOrderID    uint      `gorm:"not null;index"`     //买入订单ID（关联orders表）
	SellOrderID   uint      `gorm:"not null;index"`     //卖出订单ID（关联orders表）
	Symbol        string    `gorm:"not null;size:10"`   //交易标的（如股票代码）
	Price         float64   `gorm:"type:decimal(10,2)"` //实际成交价格
	Quantity      int       `gorm:"not null"`           //成交数量
	ExecutionTime time.Time `gorm:"not null;index"`     //成交时间（撮合完成时间）
}

// 表名自定义（可选）
func (Trade) TableName() string {
	return "trades"
}
