package repository

import (
	"context"
	"time"

	"github.com/user2083251241/ebidsystem/internal/domain/entity"
)

// AuditLog 仓库接口
type AuditLogRepository interface {
	// 创建审计日志
	Create(ctx context.Context, log *entity.AuditLog) error

	// 获取日志列表
	List(ctx context.Context, offset, limit int) ([]*entity.AuditLog, error)

	// 根据用户ID获取日志
	ListByUserID(ctx context.Context, userID uint, offset, limit int) ([]*entity.AuditLog, error)

	// 根据时间范围获取日志
	ListByTimeRange(ctx context.Context, startTime, endTime time.Time, offset, limit int) ([]*entity.AuditLog, error)
}
