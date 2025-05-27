package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user2083251241/ebidsystem/internal/app/config"
	"github.com/user2083251241/ebidsystem/internal/domain/matching"
	"github.com/user2083251241/ebidsystem/internal/infrastructure/persistence/mysql"
	"github.com/user2083251241/ebidsystem/internal/interfaces/http/middleware"
	"github.com/user2083251241/ebidsystem/internal/interfaces/http/routes"
	"github.com/user2083251241/ebidsystem/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	/* 加载配置 */
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Database Config: %+v", cfg.Database)

	/* 数据库初始化并连接 */

	// 初始化数据库：
	dbClient, err := mysql.NewDBClient(&config.Database{
		DSN:         cfg.Database.DSN,
		MaxIdleConn: cfg.Database.MaxIdleConn,
		MaxOpenConn: cfg.Database.MaxOpenConn,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbClient.Close()

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
	app.Use(middleware.LoggerMiddleware) //自定义日志中间件
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
	me := matching.NewMatchingEngine(db)
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
