package services

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/champNoob/ebidsystem/backend/config"
	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	// 设置测试环境
	setupTestEnv()

	// 运行测试
	code := m.Run()

	// 清理测试环境
	cleanupTestEnv()

	os.Exit(code)
}

func setupTestEnv() {
	// 设置测试数据库连接
	os.Setenv("DB_DSN", "root:jiongs@tcp(localhost:3306)/ebidsystem?charset=utf8mb4&parseTime=True&loc=Local")
}

func cleanupTestEnv() {
	// 清理环境变量
	os.Unsetenv("DB_DSN")
}

// 初始化测试数据库
func setupTestDB(t testing.TB) *gorm.DB {
	// 使用已有的数据库连接，但切换到测试数据库
	db := config.DB
	if db == nil {
		var err error
		db, err = config.InitDB()
		if err != nil {
			t.Fatalf("failed to connect database: %v", err)
		}
	}

	// 切换到测试数据库
	err := db.Exec("CREATE DATABASE IF NOT EXISTS ebidsystem_test").Error
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	err = db.Exec("USE ebidsystem_test").Error
	if err != nil {
		t.Fatalf("failed to switch to test database: %v", err)
	}

	// 清理测试数据
	db.Exec("DROP TABLE IF EXISTS trades")
	db.Exec("DROP TABLE IF EXISTS live_orders")

	// 自动迁移模式
	err = db.AutoMigrate(&models.LiveOrder{}, &models.Trade{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

// 清理测试数据库
func cleanupTestDB(t testing.TB, db *gorm.DB) {
	if db != nil {
		// 清理测试数据
		db.Exec("DROP TABLE IF EXISTS trades")
		db.Exec("DROP TABLE IF EXISTS live_orders")
	}
}

// 生成测试订单
func generateTestOrders(count int, direction string) []models.LiveOrder {
	orders := make([]models.LiveOrder, count)
	for i := 0; i < count; i++ {
		price := 50000.0 + rand.Float64()*1000 // 随机价格
		baseOrder := models.BaseOrder{
			Symbol:    "BTC/USDT",
			OrderType: "limit",
			Quantity:  rand.Intn(10) + 1, // 1-10的随机数量
			Price:     price,
			CreatorID: uint(1), // 测试用户ID
		}

		orders[i] = models.LiveOrder{
			BaseOrder: baseOrder,
			Direction: direction,
			Status:    "pending",
		}
	}
	return orders
}

// 基础功能测试
func TestMatchingEngine(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	t.Run("Basic Matching", func(t *testing.T) {
		engine := NewMatchingEngine(db)

		// 创建基础订单信息
		baseOrderBuy := models.BaseOrder{
			Symbol:    "BTC/USDT",
			OrderType: "limit",
			Quantity:  1,
			Price:     50000,
			CreatorID: uint(1),
		}

		// 创建并保存买单
		buyOrder := models.LiveOrder{
			BaseOrder: baseOrderBuy,
			Direction: "buy",
			Status:    "pending",
		}
		if err := db.Create(&buyOrder).Error; err != nil {
			t.Fatalf("Failed to create buy order: %v", err)
		}

		// 创建并保存卖单
		baseOrderSell := models.BaseOrder{
			Symbol:    "BTC/USDT",
			OrderType: "limit",
			Quantity:  1,
			Price:     50000,
			CreatorID: uint(2),
		}
		sellOrder := models.LiveOrder{
			BaseOrder: baseOrderSell,
			Direction: "sell",
			Status:    "pending",
		}
		if err := db.Create(&sellOrder).Error; err != nil {
			t.Fatalf("Failed to create sell order: %v", err)
		}

		// 执行撮合
		err := engine.processBatch([]models.LiveOrder{buyOrder}, []models.LiveOrder{sellOrder})
		assert.NoError(t, err)

		// 验证结果
		var trade models.Trade
		err = db.First(&trade).Error
		assert.NoError(t, err)
		assert.Equal(t, buyOrder.ID, trade.BuyOrderID)
		assert.Equal(t, sellOrder.ID, trade.SellOrderID)
	})
}

func BenchmarkMatchingEngine(b *testing.B) {
	db := setupTestDB(b)
	defer cleanupTestDB(b, db)

	engine := NewMatchingEngine(db)

	// 准备测试数据
	buyOrders := generateTestOrders(100, "buy")
	sellOrders := generateTestOrders(100, "sell")

	// 保存订单到数据库
	for i := range buyOrders {
		if err := db.Create(&buyOrders[i]).Error; err != nil {
			b.Fatalf("Failed to create buy order: %v", err)
		}
	}
	for i := range sellOrders {
		if err := db.Create(&sellOrders[i]).Error; err != nil {
			b.Fatalf("Failed to create sell order: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.processBatch(buyOrders, sellOrders)
	}
}

func TestConcurrentMatching(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	engine := NewMatchingEngine(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 启动撮合引擎
	go engine.Run(ctx)

	// 模拟并发订单提交
	for i := 0; i < 5; i++ {
		go func() {
			buyOrders := generateTestOrders(5, "buy")
			sellOrders := generateTestOrders(5, "sell")

			// 保存订单到数据库
			for j := range buyOrders {
				if err := db.Create(&buyOrders[j]).Error; err != nil {
					t.Errorf("Failed to create buy order: %v", err)
					return
				}
			}
			for j := range sellOrders {
				if err := db.Create(&sellOrders[j]).Error; err != nil {
					t.Errorf("Failed to create sell order: %v", err)
					return
				}
			}

			err := engine.processBatch(buyOrders, sellOrders)
			assert.NoError(t, err)
		}()
	}

	// 等待一段时间让并发操作完成
	time.Sleep(3 * time.Second)

	// 验证是否有成交记录
	var tradeCount int64
	db.Model(&models.Trade{}).Count(&tradeCount)
	assert.Greater(t, tradeCount, int64(0))
}
