package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/champNoob/ebidsystem/backend/models"
	"gorm.io/gorm"
)

var (
	matchLogger = log.New(os.Stdout, "", 0) // 控制台输出
	fileLogger  *log.Logger                 // 文件日志
)

func init() {
	// 初始化日志输出
	matchLogger = log.New(os.Stdout, "[MATCH] ", log.Ltime|log.Lshortfile)

	// 获取可执行文件所在目录
	exeDir, err := os.Getwd()
	if err != nil {
		log.Println("获取工作目录失败:", err)
		return
	}

	// 构建跨平台日志路径
	logPath := filepath.Join(exeDir, "bin", "matchLog", "matchLog.txt")
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		log.Println("创建日志目录失败:", err)
		return
	}

	logFile, err := os.OpenFile(logPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("打开日志文件失败:", err)
		return
	}

	fileLogger = log.New(logFile, "", log.LstdFlags)
}

// MatchOrders 订单撮合引擎核心逻辑
/* 参数说明：
* db: 数据库连接实例
* matchInterval: 撮合间隔（用于限价单时间优先规则）
* priceTolerance: 价格浮动容忍度（用于处理浮点数精度问题）
 */
func MatchOrders(db *gorm.DB, matchInterval time.Duration, priceTolerance float64) error {
	output := func(format string, v ...interface{}) {
		msg := fmt.Sprintf(format, v...)
		if matchLogger != nil {
			matchLogger.Print(msg) // 控制台输出
		}
		if fileLogger != nil {
			fileLogger.Print(msg) // 文件输出
		}
	}
	// output("=== 撮合引擎启动 ===")
	// defer output("=== 撮合引擎结束 ===")
	// ==================== 阶段1：查询待撮合订单 ====================
	// 获取当前时间戳，用于时间优先规则
	now := time.Now()

	// 查询所有未完全成交的买入订单（按价格降序、时间升序排序）：
	var buyOrders []models.Order
	if err := db.Where(
		"direction = 'buy' AND status = 'pending' AND "+
			"(order_type = 'market' OR (order_type = 'limit' AND created_at <= ?))", // 限价单需在时间窗口内
		now.Add(-matchInterval),
	).
		Order("CASE WHEN order_type = 'market' THEN 0 ELSE 1 END, price DESC, created_at ASC").
		Find(&buyOrders).Error; err != nil {
		return fmt.Errorf("查询买入订单失败: %v", err)
	}
	// 查询所有未完全成交的卖出订单（按价格升序、时间升序排序）：
	var sellOrders []models.Order
	if err := db.Where(
		"direction = 'sell' AND status = 'pending' AND "+
			"(order_type = 'market' OR (order_type = 'limit' AND created_at <= ?))",
		now.Add(-matchInterval),
	).
		Order("CASE WHEN order_type = 'market' THEN 0 ELSE 1 END, price ASC, created_at ASC").
		Find(&sellOrders).Error; err != nil {
		return fmt.Errorf("查询卖出订单失败: %v", err)
	}
	// ==================== 阶段2：订单撮合处理 ====================
	for i := 0; i < len(buyOrders); i++ {
		buy := &buyOrders[i] // 使用指针以便直接修改

		// 跳过已完全成交的订单
		if buy.Quantity <= 0 {
			continue
		}

		// 遍历卖出订单寻找匹配
		for j := 0; j < len(sellOrders); j++ {
			sell := &sellOrders[j]

			// 跳过已完全成交的订单
			if sell.Quantity <= 0 {
				continue
			}
			// 输出订单数量信息：
			output("有效买入订单数量: %d, 有效卖出订单数量: %d", len(buyOrders), len(sellOrders))
			// // 输出匹配尝试信息：
			// log.Printf("尝试撮合: 买单ID=%d (价格%.2f) vs 卖单ID=%d (价格%.2f)",
			// 	buy.ID, buy.Price, sell.ID, sell.Price)
			// 检查基础匹配条件
			if !isMatchable(buy, sell, priceTolerance) {
				continue
			}
			// ================ 计算成交细节 ================
			/* 确定成交价（价格优先规则）：
			* 当两者都是限价单时，取先到达的价格
			* 市价单接受对方限价
			 */
			executionPrice := determineExecutionPrice(buy, sell)

			// 计算成交量（取两者剩余数量的最小值）
			executionQty := min(buy.Quantity, sell.Quantity)

			// ================ 数据库事务处理 ================
			tx := db.Begin()
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()

			// 更新买入订单状态
			if err := updateOrder(tx, buy, executionQty); err != nil {
				tx.Rollback()
				break // 跳出当前卖出订单循环
			}

			// 更新卖出订单状态
			if err := updateOrder(tx, sell, executionQty); err != nil {
				tx.Rollback()
				break
			}

			// 生成成交记录
			trade := models.Trade{
				BuyOrderID:    buy.ID,
				SellOrderID:   sell.ID,
				ExecutionTime: now,
				Price:         executionPrice,
				Quantity:      executionQty,
				Symbol:        buy.Symbol,
			}
			if err := tx.Create(&trade).Error; err != nil {
				tx.Rollback()
				break
			}
			log.Printf("撮合成功: 买单%d@%.2f -> 卖单%d@%.2f (数量%d)",
				buy.ID, buy.Price, sell.ID, sell.Price, executionQty)
			// 提交事务
			if err := tx.Commit().Error; err != nil {
				return fmt.Errorf("事务提交失败: %v", err)
			}

			// 更新内存中的数量状态
			buy.Quantity -= executionQty
			sell.Quantity -= executionQty

			// 如果当前买入订单已完全成交，跳出卖出订单循环
			if buy.Quantity == 0 {
				break
			}
		}
	}
	return nil
}

// ==================== 工具函数 ====================

// isMatchable 检查订单是否可撮合
func isMatchable(buy *models.Order, sell *models.Order, tolerance float64) bool {
	// 标的证券代码必须一致
	if buy.Symbol != sell.Symbol {
		return false
	}

	// 市价买单可以匹配任何价格的卖单
	if buy.OrderType == "market" {
		return true
	}

	// 市价卖单可以匹配任何价格的买单
	if sell.OrderType == "market" {
		return true
	}

	// 限价单匹配条件：买入价 >= 卖出价（考虑浮点数精度）
	return (buy.Price - sell.Price) >= -tolerance
}

// determineExecutionPrice 确定成交价格
func determineExecutionPrice(buy *models.Order, sell *models.Order) float64 {
	switch {
	case buy.OrderType == "market":
		return sell.Price // 市价买单以卖出价成交
	case sell.OrderType == "market":
		return buy.Price // 市价卖单以买入价成交
	default:
		// 都是限价单时，取先到达市场的订单价格
		if buy.CreatedAt.Before(sell.CreatedAt) {
			return buy.Price
		}
		return sell.Price
	}
}

// updateOrder 更新订单状态（事务安全）
func updateOrder(tx *gorm.DB, order *models.Order, executedQty int) error {
	// 计算剩余数量
	remainingQty := order.Quantity - executedQty

	updateData := map[string]interface{}{
		"quantity": remainingQty,
	}

	// 如果完全成交则更新状态
	if remainingQty <= 0 {
		updateData["status"] = "filled"
	}

	// 使用 Select 明确更新字段以避免 GORM 的自动更新干扰
	return tx.Model(order).
		Select("quantity", "status", "updated_at").
		Updates(updateData).
		Error
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
