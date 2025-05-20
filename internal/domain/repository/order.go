package repository

import (
	"context"

	"github.com/user2083251241/ebidsystem/internal/domain/entity"
)

// Order 仓库接口
type OrderRepository interface {
	// 创建草稿订单
	CreateDraft(ctx context.Context, order *entity.DraftOrder) error

	// 创建正式订单
	CreateLive(ctx context.Context, order *entity.LiveOrder) error

	// 根据ID查找草稿订单
	FindDraftByID(ctx context.Context, id uint) (*entity.DraftOrder, error)

	// 根据ID查找正式订单
	FindLiveByID(ctx context.Context, id uint) (*entity.LiveOrder, error)

	// 根据用户ID查找草稿订单
	ListDraftsByUserID(ctx context.Context, userID uint, offset, limit int) ([]*entity.DraftOrder, error)

	// 根据用户ID查找正式订单
	ListLiveOrdersByUserID(ctx context.Context, userID uint, offset, limit int) ([]*entity.LiveOrder, error)

	// 更新草稿订单状态
	UpdateDraftStatus(ctx context.Context, id uint, status string) error

	// 更新正式订单状态
	UpdateLiveStatus(ctx context.Context, id uint, status entity.OrderStatus) error

	// 更新正式订单价格
	UpdateLivePrice(ctx context.Context, id uint, price float64) error

	// 更新正式订单数量
	UpdateLiveQuantity(ctx context.Context, id uint, quantity int) error

	// 删除草稿订单
	DeleteDraft(ctx context.Context, id uint) error

	// 删除正式订单
	DeleteLive(ctx context.Context, id uint) error

	// 获取订单统计信息
	GetLiveStats(ctx context.Context, userID uint) (*entity.OrderStats, error)
}
