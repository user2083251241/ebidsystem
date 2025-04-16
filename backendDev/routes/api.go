package routes

import (
	"github.com/champNoob/ebidsystem/backend/config"
	"github.com/champNoob/ebidsystem/backend/controllers"
	"github.com/champNoob/ebidsystem/backend/middleware"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App) {
	// 添加CORS中间件（必须在路由定义前调用）
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://192.168.1.100:8080", // 替换为前端实际IP和端口
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true, // 允许携带Cookie或Authorization头
	}))
	// 公共路由：
	public := app.Group("/api")
	{
		public.Post("/register", controllers.Register)
		public.Post("/login", controllers.Login)
	}
	// JWT 中间件初始化：
	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(config.Get("JWT_SECRET")),
		},
	})
	// 认证路由组：
	authenticated := app.Group("/api", jwtMiddleware)
	{
		// 所有认证用户均可调用：
		authenticated.Post("/logout", controllers.Logout) // 登出
		// 卖家角色路由组：
		seller := authenticated.Group("/seller", middleware.SellerOnly)
		{
			seller.Post("/orders", controllers.CreateSellOrder)                // 创建卖出订单
			seller.Put("/orders/:id", controllers.UpdateOrder)                 // 修改订单
			seller.Delete("/orders/:id", controllers.CancelOrder)              // 单个撤单
			seller.Post("/orders/batch-cancel", controllers.BatchCancelOrders) // 批量撤单
			seller.Get("/orders", controllers.GetSellerOrders)                 // 查看卖家订单
			seller.Post("/authorize/sales", controllers.AuthorizeSales)        // 授权销售
		}
		// 销售角色路由组：
		sales := authenticated.Group("/sales", middleware.SalesOnly)
		{
			sales.Get("/orders", controllers.GetAuthorizedOrders)     // 查看已授权订单
			sales.Post("/drafts", controllers.CreateDraftOrder)       // 创建订单草稿
			sales.Put("/drafts/:id", controllers.UpdateDraftOrder)    // 修改草稿
			sales.Post("/drafts/:id/submit", controllers.SubmitDraft) // 提交草稿
		}
		// 客户角色路由组：
		client := authenticated.Group("/client", middleware.ClientOnly)
		{
			client.Get("/orders", controllers.GetClientOrders)         //查看匿名处理的卖方订单
			client.Post("/orders/:id/buy", controllers.CreateBuyOrder) //对已有的卖方订单创建自己的买入订单
		}
		// 交易员角色路由组：
		trader := authenticated.Group("/trader", middleware.TraderOnly)
		{
			trader.Get("/orders", controllers.GetAllOrders) // 查看所有订单
		}
	}
	app.Static("/", "./static")
}
