package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/champNoob/ebidsystem/backend/models"
	"gorm.io/gorm"
)

const (
	MatchInterval  = 10 * time.Second
	PriceTolerance = 0.01
	MaxRetries     = 3
	OrderBatchSize = 100
)

type MatchingEngine struct {
	db             *gorm.DB
	logger         *log.Logger
	priceTolerance float64
}

func NewMatchingEngine(db *gorm.DB) *MatchingEngine {
	return &MatchingEngine{
		db:             db,
		logger:         log.New(log.Writer(), "[MATCHING] ", log.LstdFlags|log.Lshortfile),
		priceTolerance: PriceTolerance,
	}
}

func (e *MatchingEngine) Run(ctx context.Context) {
	ticker := time.NewTicker(MatchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			e.logger.Println("Matching engine stopped")
			return
		case <-ticker.C:
			if err := e.matchOrders(); err != nil {
				e.logger.Printf("Matching failed: %v", err)
			}
		}
	}
}

func (e *MatchingEngine) matchOrders() error {
	e.logger.Println("Starting matching cycle")

	// 获取待撮合订单批次
	for offset := 0; ; offset += OrderBatchSize {
		buyOrders, sellOrders, err := e.fetchOrderBatch(offset)
		if err != nil {
			return err
		}
		if len(buyOrders) == 0 && len(sellOrders) == 0 {
			break
		}

		if err := e.processBatch(buyOrders, sellOrders); err != nil {
			return err
		}
	}

	e.logger.Println("Matching cycle completed")
	return nil
}

func (e *MatchingEngine) fetchOrderBatch(offset int) ([]models.LiveOrder, []models.LiveOrder, error) {
	var buyOrders []models.LiveOrder
	if err := e.db.Where("direction = 'buy' AND status = 'pending'").
		Order("order_type DESC, price DESC, created_at ASC").
		Offset(offset).Limit(OrderBatchSize).
		Find(&buyOrders).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch buy orders failed: %w", err)
	}

	var sellOrders []models.LiveOrder
	if err := e.db.Where("direction = 'sell' AND status = 'pending'").
		Order("order_type DESC, price ASC, created_at ASC").
		Offset(offset).Limit(OrderBatchSize).
		Find(&sellOrders).Error; err != nil {
		return nil, nil, fmt.Errorf("fetch sell orders failed: %w", err)
	}

	return buyOrders, sellOrders, nil
}

func (e *MatchingEngine) processBatch(buyOrders, sellOrders []models.LiveOrder) error {
	for _, buy := range buyOrders {
		if buy.Quantity <= 0 {
			continue
		}

		for _, sell := range sellOrders {
			if sell.Quantity <= 0 {
				continue
			}

			if !e.isMatchable(buy, sell) {
				continue
			}

			executionPrice := e.determinePrice(buy, sell)
			qty := min(buy.Quantity, sell.Quantity)

			if err := e.executeTrade(buy, sell, qty, executionPrice); err != nil {
				e.logger.Printf("Trade execution failed: %v", err)
				continue
			}

			// 更新本地副本数量
			buy.Quantity -= qty
			sell.Quantity -= qty

			if buy.Quantity == 0 {
				break
			}
		}
	}
	return nil
}

func (e *MatchingEngine) isMatchable(buy, sell models.LiveOrder) bool {
	if buy.Symbol != sell.Symbol {
		return false
	}

	switch {
	case buy.OrderType == "market":
		return true
	case sell.OrderType == "market":
		return true
	default:
		return (buy.Price - sell.Price) >= -e.priceTolerance
	}
}

func (e *MatchingEngine) determinePrice(buy, sell models.LiveOrder) float64 {
	switch {
	case buy.OrderType == "market":
		return sell.Price
	case sell.OrderType == "market":
		return buy.Price
	default:
		if buy.CreatedAt.Before(sell.CreatedAt) {
			return buy.Price
		}
		return sell.Price
	}
}

func (e *MatchingEngine) executeTrade(buy, sell models.LiveOrder, qty int, price float64) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		// 更新买单
		if err := tx.Model(&buy).Updates(map[string]interface{}{
			"quantity": buy.Quantity - qty,
			"status":   e.getOrderStatus(buy.Quantity - qty),
		}).Error; err != nil {
			return err
		}

		// 更新卖单
		if err := tx.Model(&sell).Updates(map[string]interface{}{
			"quantity": sell.Quantity - qty,
			"status":   e.getOrderStatus(sell.Quantity - qty),
		}).Error; err != nil {
			return err
		}

		// 创建成交记录
		trade := models.Trade{
			BuyOrderID:    buy.ID,
			SellOrderID:   sell.ID,
			ExecutionTime: time.Now(),
			Price:         price,
			Quantity:      qty,
			Symbol:        buy.Symbol,
		}
		return tx.Create(&trade).Error
	})
}

func (e *MatchingEngine) getOrderStatus(remainingQty int) string {
	if remainingQty <= 0 {
		return "filled"
	}
	return "partial_filled"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
