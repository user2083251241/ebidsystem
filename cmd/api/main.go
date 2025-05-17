package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/champNoob/ebidsystem/backend/config"
	"github.com/champNoob/ebidsystem/backend/middleware"
	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/champNoob/ebidsystem/backend/routes"
	"github.com/champNoob/ebidsystem/backend/services"
	"github.com/champNoob/ebidsystem/backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	/* 加载环境变量 */
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	/* 数据库初始化并连接 */

	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	} else {
		log.Println("Successfully connected to database")
	}
	config.DB = db
	//log.Printf("LiveOrder struct: %+v", reflect.TypeOf(models.LiveOrder{})) //打印 LiveOrder 结构体定义
	// 自动迁移数据库表：
	if err := db.AutoMigrate(
		&models.BaseOrder{},
		&models.DraftOrder{},
		&models.LiveOrder{},
		&models.SellerSalesAuthorization{},
		&models.Trade{},
		&models.User{},
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

	/* 初始化 Redis */
	// 在路由注册前调用：
	utils.InitRedis()
	// 测试 Redis 连接：
	if err := utils.RedisClient.Ping(utils.Ctx).Err(); err != nil {
		log.Fatalf("Redis 连接测试失败: %v", err)
	}
	log.Println("Redis 连接测试成功")

	/* 注册路由（依赖注入） */

	routes.SetupRoutes(app, db)

	/* 启动撮合引擎 */

	// 启动撮合引擎：
	me := services.NewMatchingEngine(db)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		me.Run(ctx)
	}()
	// 捕获系统信号，用于优雅关闭：
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel() //取消上下文，停止匹配引擎
		if err := app.Shutdown(); err != nil {
			log.Printf("Server shutdown failed: %v", err)
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
	} else {
		log.Printf("端口 %s 已成功绑定", port)
	}

	/* 测试 Redis 写入 */
	testErr := utils.AddToBlacklist("test_token", 10*time.Minute)
	if testErr != nil {
		log.Fatalf("Redis 写入测试失败: %v", testErr)
	} else {
		log.Println("Redis 写入测试成功")
	}
}
