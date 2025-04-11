package routes

import (
	"github.com/champNoob/ebidsystem/backend/config"
	"github.com/champNoob/ebidsystem/backend/controllers"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5" // 添加 JWT 包导入
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
	})

	// 需要认证的路由
	authenticated := app.Group("/api", jwtMiddleware)
	{
		authenticated.Post("/orders", controllers.CreateOrder)
		authenticated.Get("/orders", controllers.GetOrders)
		authenticated.Post("/orders/cancel", TraderOnly, controllers.CancelOrder) // 在此处添加路由
	}
}

// 检查用户是否为交易员
func TraderOnly(c *fiber.Ctx) error {
	// 从 JWT 中提取角色
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	role := claims["role"].(string)

	if role != "trader" {
		return c.Status(403).JSON(fiber.Map{"error": "仅交易员可执行此操作"})
	}
	return c.Next()
}
