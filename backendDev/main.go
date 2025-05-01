package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/champNoob/ebidsystem/backend/config"
	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/champNoob/ebidsystem/backend/routes"
	"github.com/user2083251241/ebidsystem/middleware"
	"github.com/user2083251241/ebidsystem/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	/* 加载环境变量 */
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	/* 初始化数据库连接 */
	dsn := config.Get("DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	/* 自动迁移数据库表 */
	if err := db.AutoMigrate(
		&models.BaseOrder{},
		&models.DraftOrder{},
		&models.LiveOrder{},
		&models.SellerSalesAuthorization{},
		&models.Trade{},
		// &models.Stock{},
	); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	/* 初始化Fiber应用 */
	app := fiber.New()
	/* 注册全局中间件 */
	// 异常恢复：
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			log.Printf("[PANIC] %v", e)
		},
	}))
	app.Use(middleware.LoggingMiddleware) //自定义日志中间件
	// 跨域请求：
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", //允许所有来源
		AllowMethods: "GET,POST,PUT,DELETE",
	}))
	serviceLog, _ := os.OpenFile("bin/logs/service.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	errorLog, _ := os.OpenFile("bin/logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// 配置请求日志中间件：
	app.Use(logger.New(logger.Config{
		Output: serviceLog, // HTTP 请求日志输出到 service.log
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))
	// 配置全局错误处理中间件：
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			errorLog.WriteString(fmt.Sprintf("[PANIC] %v\n", e))
		},
	}))
	// controllers.InitDB(db)
	/* 注册路由（依赖注入） */
	routes.SetupRoutes(app, db)
	// 启动撮合引擎：
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			services.MatchOrders(db, 10*time.Minute, 0.01)
			<-ticker.C
			if err := services.MatchOrders(db, 10*time.Minute, 0.0001); err != nil {
				log.Printf("撮合引擎错误: %v", err)
			}
		}
	}()
	/* 启动服务器 */
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	// 输出启动信息：
	fmt.Printf("🚀 Server started on port %s\n", port)
	if err := app.Listen("0.0.0.0:" + port); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
