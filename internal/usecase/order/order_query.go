package orderusecase

import (
	"github.com/user2083251241/ebidsystem/internal/domain/entity"
	"github.com/user2083251241/ebidsystem/internal/interfaces/http/dto"

	"gorm.io/gorm"
)

// QueryStrategy 订单查询策略接口
type QueryStrategy interface {
	Apply(query *gorm.DB) *gorm.DB
	GetDTOConverter() func(entity.LiveOrder) dto.OrderDTO
}

// ---------------------- 具体策略实现 ----------------------

// SellerOrdersStrategy 卖家订单策略
type SellerOrdersStrategy struct {
	UserID uint
}

func (s *SellerOrdersStrategy) Apply(query *gorm.DB) *gorm.DB {
	return query.Model(&entity.LiveOrder{}).Where("creator_id = ? AND direction = 'sell'", s.UserID) //明确指定模型
}

func (s *SellerOrdersStrategy) GetDTOConverter() func(entity.LiveOrder) entity.OrderDTO {
	return func(o entity.LiveOrder) dto.OrderDTO {
		return dto.OrderDTO{
			ID:       o.ID,
			Symbol:   o.Symbol,
			Quantity: o.Quantity,
			Price:    o.Price,
			Status:   o.Status,
		}
	}
}

// ClientOrdersStrategy 客户可见订单策略
type ClientOrdersStrategy struct{}

func (c *ClientOrdersStrategy) Apply(query *gorm.DB) *gorm.DB {
	return query.Model(&entity.LiveOrder{}).Where("direction = 'sell' AND status = 'pending'")
}

func (c *ClientOrdersStrategy) GetDTOConverter() func(entity.LiveOrder) entity.OrderDTO {
	return func(o entity.LiveOrder) dto.OrderDTO {
		return dto.OrderDTO{
			ID:     o.ID,
			Symbol: o.Symbol,
			Price:  o.Price,
		}
	}
}

// TraderOrdersStrategy 交易员订单策略
type TraderOrdersStrategy struct{}

func (t *TraderOrdersStrategy) Apply(query *gorm.DB) *gorm.DB {
	return query.Unscoped() // 查看所有订单包括软删除
}

func (t *TraderOrdersStrategy) GetDTOConverter() func(entity.LiveOrder) entity.OrderDTO {
	return func(o entity.LiveOrder) entity.OrderDTO {
		return entity.OrderDTO{
			ID:       o.ID,
			Symbol:   o.Symbol,
			Quantity: o.Quantity,
			Price:    o.Price,
			Status:   o.Status,
		}
	}
}
