package routes

import (
	"github.com/champNoob/ebidsystem/backend/config"
	"github.com/champNoob/ebidsystem/backend/controllers"
	"github.com/champNoob/ebidsystem/backend/middleware"
	"github.com/champNoob/ebidsystem/backend/services"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// 初始化服务层：
	orderService := services.NewOrderService(db)
	userService := services.NewUserService(db)
	// 初始化控制器：
	authController := controllers.NewAuthController(userService)
	sellerController := controllers.NewSellerController(orderService)
	salesController := controllers.NewSalesController(orderService)
	traderController := controllers.NewTraderController(orderService)
	clientController := controllers.NewClientController(orderService)
	// 添加 CORS 中间件（必须在路由定义前调用）：
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://192.168.1.100:8080", //替换为前端实际IP和端口
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true, //允许携带 Cookie 或 Authorization 头
	}))
	// 公共路由：
	public := app.Group("/api")
	{
		public.Post("/register", authController.Register)
		public.Post("/login", authController.Login)
	}
	// 初始化 JWT 中间件：
	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(config.Get("JWT_SECRET")),
		},
	})
	// 认证路由组：
	authenticated := app.Group("/api", jwtMiddleware, middleware.CheckTokenRevoked())
	{
		// 所有认证用户均可调用：
		authenticated.Post("/logout", authController.Logout) //登出
		// 卖家路由组：
		seller := authenticated.Group("/seller", middleware.RoleRequired("seller"))
		{
			seller.Post("/orders", sellerController.SellerCreateOrder)                    //创建卖出订单
			seller.Put("/orders/:id", sellerController.SellerUpdateOrder)                 //修改订单
			seller.Delete("/orders/:id", sellerController.SellerCancelOrder)              //单个撤单
			seller.Post("/orders/batch-cancel", sellerController.SellerBatchCancelOrders) //批量撤单
			seller.Get("/orders", sellerController.SellerGetOrders)                       //查看卖家订单
			seller.Post("/authorize/sales", sellerController.SellerAuthorizeSales)        //授权销售
		}
		// 销售路由组：
		sales := authenticated.Group("/sales", middleware.RoleRequired("sales"))
		{
			sales.Get("/orders", salesController.SalesGetAuthorizedDrafts)     //查看已授权草稿
			sales.Post("/drafts", salesController.SalesCreateDraft)            //创建草稿
			sales.Put("/drafts/:id", salesController.SalesUpdateDraft)         //修改草稿
			sales.Post("/drafts/:id/submit", salesController.SalesSubmitDraft) //提交草稿
			sales.Delete("/drafts/:id", salesController.SalesDeleteDraft)      //删除草稿
		}
		// 客户路由组：
		client := authenticated.Group("/client", middleware.RoleRequired("client"))
		{
			client.Get("/orders", clientController.ClientGetOrders)               //查看匿名处理的卖方订单
			client.Post("/orders/:id/buy", clientController.ClientCreateBuyOrder) //对已有的卖方订单创建自己的买入订单
		}
		// 交易员路由组：
		trader := authenticated.Group("/trader", middleware.RoleRequired("trader"))
		{
			trader.Get("/orders", traderController.TraderGetAllOrders) // 查看所有订单
			// trader.Post("/orders/:id/cancel", traderController.EmergencyCancel) //手动操作，暂不实现
		}
	}
	app.Static("/", "./static")
	app.Static("/assets", "./static/assets")
}
