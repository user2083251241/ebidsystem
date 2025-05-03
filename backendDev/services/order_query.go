package services

import (
	"github.com/champNoob/ebidsystem/backend/models"
	"gorm.io/gorm"
)

// QueryStrategy 订单查询策略接口
type QueryStrategy interface {
	Apply(query *gorm.DB) *gorm.DB
	GetDTOConverter() func(models.LiveOrder) models.OrderDTO
}

// ---------------------- 具体策略实现 ----------------------

// SellerOrdersStrategy 卖家订单策略
type SellerOrdersStrategy struct {
	UserID uint
}

func (s *SellerOrdersStrategy) Apply(query *gorm.DB) *gorm.DB {
	return query.Model(&models.LiveOrder{}).Where("user_id = ? AND direction = 'sell'", s.UserID) //明确指定模型
}

func (s *SellerOrdersStrategy) GetDTOConverter() func(models.LiveOrder) models.OrderDTO {
	return func(o models.LiveOrder) models.OrderDTO {
		return models.OrderDTO{
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
	return query.Model(&models.LiveOrder{}).Where("direction = 'sell' AND status = 'pending'")
}

func (c *ClientOrdersStrategy) GetDTOConverter() func(models.LiveOrder) models.OrderDTO {
	return func(o models.LiveOrder) models.OrderDTO {
		return models.OrderDTO{
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

func (t *TraderOrdersStrategy) GetDTOConverter() func(models.LiveOrder) models.OrderDTO {
	return func(o models.LiveOrder) models.OrderDTO {
		return models.OrderDTO{
			ID:       o.ID,
			Symbol:   o.Symbol,
			Quantity: o.Quantity,
			Price:    o.Price,
			Status:   o.Status,
		}
	}
}
