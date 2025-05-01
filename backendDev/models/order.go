package models

import (
	"gorm.io/gorm"
)

// 基类订单（公共字段）：
type BaseOrder struct {
	gorm.Model
	Symbol   string  `gorm:"not null"`           //标的代码
	Quantity int     `gorm:"not null"`           //数量
	Price    float64 `gorm:"type:decimal(10,2)"` //价格
}

// 草稿订单（组合基类）：
type DraftOrder struct {
	BaseOrder           // 嵌套公共字段
	DraftBySales uint   `gorm:"not null"`        //草稿创建者（销售ID）
	Status       string `gorm:"default:'draft'"` //状态：draft/pending_approval
}

// 正式订单（组合基类）：
type LiveOrder struct {
	BaseOrder         // 嵌套公共字段
	Direction  string `gorm:"not null"`          //方向：buy/sell
	Status     string `gorm:"default:'pending'"` //状态：pending/filled/cancelled
	ApprovedBy uint   // 审批人（卖家ID）
}

type OrderDTO struct {
	ID       uint    `json:"id"`
	Symbol   string  `json:"symbol"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity,omitempty"` //对客户隐藏
	Status   string  `json:"status,omitempty"`   //对客户和销售隐藏
	// 隐藏 DraftBySales、SellerID 等字段
}

// 表名自定义（可选）
func (BaseOrder) TableName() string {
	return "orders"
}
