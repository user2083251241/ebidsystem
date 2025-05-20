package entity

import (
	"time"

	"gorm.io/gorm"
)

// 审计日志结构体
type AuditLog struct {
	gorm.Model
	UserID    uint      `gorm:"not null;index"` // 用户 ID
	Action    string    `gorm:"not null"`       // 操作类型，如 emergency_cancel
	OrderID   uint      `gorm:"not null;index"` // 订单 ID
	Timestamp time.Time `gorm:"not null"`       // 操作时间
}

// 表名自定义（可选）
func (AuditLog) TableName() string {
	return "audit_logs"
}
