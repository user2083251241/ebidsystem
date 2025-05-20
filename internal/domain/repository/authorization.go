package repository

import (
	"context"

	"github.com/user2083251241/ebidsystem/internal/domain/entity"
)

// Authorization 仓库接口
type AuthorizationRepository interface {
	// 创建授权
	Create(ctx context.Context, auth *entity.SellerSalesAuthorization) error

	// 根据订单ID获取授权
	FindByOrderID(ctx context.Context, orderID uint) (*entity.SellerSalesAuthorization, error)

	// 根据销售ID获取授权列表
	ListBySalesID(ctx context.Context, salesID uint, offset, limit int) ([]*entity.SellerSalesAuthorization, error)

	// 更新授权状态
	UpdateStatus(ctx context.Context, id uint, status entity.SellerSalesAuthorization) error

	// 删除授权
	Delete(ctx context.Context, id uint) error
}
