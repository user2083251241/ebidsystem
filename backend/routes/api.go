package routes

import (
	"github.com/champNoob/ebidsystem/backend/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func SetupRoutes(app *fiber.App) {
	// 公共路由（无需认证）
	public := app.Group("/api")
	{
		public.Post("/register", controllers.Register) // 用户注册
		public.Post("/login", controllers.Login)       // 用户登录
	}

	// 需要认证的路由
	authenticated := app.Group("/api", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "admin", // 示例基础认证（实际项目应使用 JWT）
		},
	}))
	{
		authenticated.Post("/orders", controllers.CreateOrder) // 创建订单
		authenticated.Get("/orders", controllers.GetOrders)    // 查询订单
	}
}
