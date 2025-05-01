package middleware

import (
	_ "github.com/champNoob/ebidsystem/backend/config"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/fiber/v2"
	"github.com/user2083251241/ebidsystem/config"
)

// JWT 中间件初始化（已由 routes/api.go 直接调用，此文件可省略）
// 保留此文件以备扩展自定义逻辑

func JWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(config.Get("JWT_SECRET")),
		},
		ContextKey: "user", // 设置 JWT 存储键
	})
}
