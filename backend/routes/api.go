package routes

import (
	"github.com/champNoob/ebidsystem/backend/config"
	"github.com/champNoob/ebidsystem/backend/controllers"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func SetupRoutes(app *fiber.App) {
	// 公开路由
	public := app.Group("/api")
	{
		public.Post("/register", controllers.Register)
		public.Post("/login", controllers.Login)
	}

	// JWT 中间件
	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey: []byte(config.Get("JWT_SECRET")),
	})

	// 需要认证的路由
	authenticated := app.Group("/api", jwtMiddleware)
	{
		authenticated.Post("/orders", controllers.CreateOrder)
		authenticated.Get("/orders", controllers.GetOrders)
	}
}
