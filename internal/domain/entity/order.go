package entity

import (
	"time"

	"gorm.io/gorm"
)

// 基类订单（公共字段）：
type BaseOrder struct {
	gorm.Model
	Symbol    string  `gorm:"not null"`           //标的代码
	OrderType string  `gorm:"not null"`           //订单类型：limit/market
	Quantity  int     `gorm:"not null"`           //数量
	Price     float64 `gorm:"type:decimal(10,2)"` //价格
	CreatorID uint    `gorm:"not null"`           //创建者 ID
}

// 草稿订单（组合基类）：
type DraftOrder struct {
	BaseOrder         //嵌套公共字段
	RefOrderID *uint  `gorm:"index"`           //关联原订单 ID
	Status     string `gorm:"default:'draft'"` //状态：draft/pending_approval
}

// 正式订单（组合基类）：
type LiveOrder struct {
	BaseOrder        //嵌套公共字段
	Direction string `gorm:"not null"`          //方向：buy/sell
	Status    string `gorm:"default:'pending'"` //状态：pending/filled/cancelled
	// ApprovedBy uint   //审批人（卖家ID）
}

// 订单状态枚举
type OrderStatus string

const (
	StatusPending         OrderStatus = "pending"
	StatusFilled          OrderStatus = "filled"
	StatusCancelled       OrderStatus = "cancelled"
	StatusDraft           OrderStatus = "draft"
	StatusPendingApproval OrderStatus = "pending_approval"
)

// 订单统计信息
type OrderStats struct {
	TotalOrders      int            `gorm:"-"` // 总订单数
	AveragePrice     float64        `gorm:"-"` // 平均价格
	TotalQuantity    int            `gorm:"-"` // 总数量
	LastUpdated      time.Time      `gorm:"-"` // 最后更新时间
	TotalValue       float64        `gorm:"-"` // 总金额
	AverageQuantity  int            `gorm:"-"` // 平均数量
	OrderCountByType map[string]int `gorm:"-"` // 按类型统计
}

// 表名定义
func (LiveOrder) TableName() string {
	return "live_orders"
}

func (DraftOrder) TableName() string {
	return "draft_orders"
}
