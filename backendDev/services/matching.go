package services

import (
	"github.com/champNoob/ebidsystem/backend/models"
	"gorm.io/gorm"
)

type MatchingService struct {
	DB *gorm.DB
}

func NewMatchingService(db *gorm.DB) *MatchingService {
	return &MatchingService{DB: db}
}

// 撮合订单
func (s *MatchingService) MatchOrders() {
	// 获取所有未成交买单（按价格降序、时间升序）
	var buyOrders []models.Order
	s.DB.Where("direction = ? AND status = ?", "buy", "pending").
		Order("price DESC, created_at ASC").
		Find(&buyOrders)

	// 获取所有未成交卖单（按价格升序、时间升序）
	var sellOrders []models.Order
	s.DB.Where("direction = ? AND status = ?", "sell", "pending").
		Order("price ASC, created_at ASC").
		Find(&sellOrders)

	// 撮合逻辑
	for _, buy := range buyOrders {
		for _, sell := range sellOrders {
			if buy.Symbol == sell.Symbol && buy.Price >= sell.Price {
				// 更新订单状态为成交
				s.DB.Model(&buy).Update("status", "filled")
				s.DB.Model(&sell).Update("status", "filled")
				break
			}
		}
	}
}
