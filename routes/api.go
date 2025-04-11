package routes

import (
	"github.com/champNoob/ebidsystem/backend/config"
	"github.com/champNoob/ebidsystem/backend/controllers"
	jwtware "github.com/gofiber/contrib/jwt" //使用新版中间件，并使用别名 jwtware
	"github.com/gofiber/fiber/v2"
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
		SigningKey: jwtware.SigningKey{
			Key: []byte(config.Get("JWT_SECRET")),
		},
		// 其他配置...
	})

	// 需要认证的路由
	authenticated := app.Group("/api", jwtMiddleware)
	{
		authenticated.Post("/orders", controllers.CreateOrder)
		authenticated.Get("/orders", controllers.GetOrders)
	}
}
