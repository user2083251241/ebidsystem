package matching

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	MatchInterval  = 10 * time.Second
	PriceTolerance = 0.01
	MaxRetries     = 3
	OrderBatchSize = 100
)

// OrderPriority 定义订单优先级
type OrderPriority struct {
	TimePriority  int64
	PricePriority float64
	SizePriority  int
}

type MatchingEngine struct {
	db             *gorm.DB
	logger         *logrus.Logger
	priceTolerance float64
	orderCache     sync.Map
	stats          *MatchingStats
}

type MatchingStats struct {
	TotalMatches   int64
	TotalVolume    float64
	LastMatchTime  time.Time
	ProcessingTime time.Duration
	ErrorCount     int64
}

func NewMatchingEngine(db *gorm.DB) *MatchingEngine {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{}) //设置日志格式为 JSON
	logger.SetLevel(logrus.InfoLevel)            //设置日志级别为 Info

	return &MatchingEngine{
		db:             db,
		logger:         logger,
		priceTolerance: PriceTolerance,
		stats:          &MatchingStats{},
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
	startTime := time.Now()
	defer func() {
		e.stats.ProcessingTime = time.Since(startTime)
	}()

	for _, buy := range buyOrders {
		if buy.Quantity <= 0 {
			continue
		}

		// 计算买单优先级
		buyPriority := e.calculatePriority(buy)
		e.logger.WithFields(logrus.Fields{
			"order_id": buy.ID,
			"priority": buyPriority,
		}).Debug("Buy order priority calculated")

		for _, sell := range sellOrders {
			if sell.Quantity <= 0 {
				continue
			}

			// 计算卖单优先级
			sellPriority := e.calculatePriority(sell)
			e.logger.WithFields(logrus.Fields{
				"order_id": sell.ID,
				"priority": sellPriority,
			}).Debug("Sell order priority calculated")

			if !e.isMatchable(buy, sell) {
				continue
			}

			executionPrice := e.determinePrice(buy, sell)
			qty := min(buy.Quantity, sell.Quantity)

			// 记录撮合日志
			e.logMatchAttempt(buy, sell, executionPrice, qty)

			if err := e.executeTrade(buy, sell, qty, executionPrice); err != nil {
				e.logger.WithFields(logrus.Fields{
					"error":         err,
					"buy_order_id":  buy.ID,
					"sell_order_id": sell.ID,
				}).Error("Trade execution failed")
				e.stats.ErrorCount++
				continue
			}

			// 更新统计信息
			e.stats.TotalMatches++
			e.stats.TotalVolume += float64(qty) * executionPrice
			e.stats.LastMatchTime = time.Now()

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

func (e *MatchingEngine) calculatePriority(order models.LiveOrder) OrderPriority {
	return OrderPriority{
		TimePriority:  order.CreatedAt.UnixNano(),
		PricePriority: order.Price,
		SizePriority:  order.Quantity,
	}
}

func (e *MatchingEngine) logMatchAttempt(buy, sell models.LiveOrder, price float64, qty int) {
	logEntry := logrus.WithFields(logrus.Fields{
		"buy_order_id":  buy.ID,
		"sell_order_id": sell.ID,
		"symbol":        buy.Symbol,
		"price":         price,
		"quantity":      qty,
		"buy_type":      buy.OrderType,
		"sell_type":     sell.OrderType,
	})

	logEntry.Info("Match attempt")
}

func (e *MatchingEngine) executeTrade(buy, sell models.LiveOrder, qty int, price float64) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		// 添加乐观锁
		if err := tx.Model(&buy).Where("id = ? AND quantity >= ?", buy.ID, qty).Updates(map[string]interface{}{
			"quantity":   buy.Quantity - qty,
			"status":     e.getOrderStatus(buy.Quantity - qty),
			"updated_at": time.Now(),
		}).Error; err != nil {
			return fmt.Errorf("update buy order failed: %w", err)
		}

		if err := tx.Model(&sell).Where("id = ? AND quantity >= ?", sell.ID, qty).Updates(map[string]interface{}{
			"quantity":   sell.Quantity - qty,
			"status":     e.getOrderStatus(sell.Quantity - qty),
			"updated_at": time.Now(),
		}).Error; err != nil {
			return fmt.Errorf("update sell order failed: %w", err)
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

		if err := tx.Create(&trade).Error; err != nil {
			return fmt.Errorf("create trade record failed: %w", err)
		}

		// 发布撮合结果
		e.publishMatchResult(trade)
		return nil
	})
}

func (e *MatchingEngine) publishMatchResult(trade models.Trade) {
	// 这里可以集成消息队列（如Kafka）来发布撮合结果
	// 示例使用简单的日志输出
	tradeJSON, _ := json.Marshal(trade)
	e.logger.WithField("trade", string(tradeJSON)).Info("Trade executed")
}

// GetStats 返回撮合引擎的统计信息
func (e *MatchingEngine) GetStats() *MatchingStats {
	return e.stats
}

func (e *MatchingEngine) isMatchable(buy, sell models.LiveOrder) bool {
	if buy.Symbol != sell.Symbol {
		return false
	}

	// 获取优先级
	buyPri := e.calculatePriority(buy)
	sellPri := e.calculatePriority(sell)

	// 基本匹配规则
	switch {
	case buy.OrderType == "market":
		return true
	case sell.OrderType == "market":
		return true
	default:
		// 价格匹配检查
		isPriceMatch := (buy.Price - sell.Price) >= -e.priceTolerance
		// 时间优先级检查（较早的订单优先）
		isTimePriorityValid := buyPri.TimePriority <= sellPri.TimePriority
		return isPriceMatch && isTimePriorityValid
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

// 监控接口
func (e *MatchingEngine) Monitor() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		stats := e.GetStats()
		e.logger.WithFields(logrus.Fields{
			"total_matches":   stats.TotalMatches,
			"total_volume":    stats.TotalVolume,
			"error_count":     stats.ErrorCount,
			"processing_time": stats.ProcessingTime,
		}).Info("Matching Engine Status")
	}
}
