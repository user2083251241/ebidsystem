package entity

import (
	"time"

	"gorm.io/gorm"
)

// 卖家-销售授权关联表
type SellerSalesAuthorization struct {
	gorm.Model
	SellerID      uint      `gorm:"not null;index"`
	SalesID       uint      `gorm:"not null;index"`
	Authorization string    `gorm:"type:enum('pending','approved','rejected');default:'pending'"`
	ExpiresAt     time.Time `gorm:"not null"`
}

func (SellerSalesAuthorization) TableName() string {
	return "seller_sales_authorizations"
}
