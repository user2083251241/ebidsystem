package middleware

import (
	"github.com/champNoob/ebidsystem/backend/config"
	_ "github.com/champNoob/ebidsystem/backend/config"
	"github.com/champNoob/ebidsystem/backend/models"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

// 从 JWT 中提取用户信息（不查数据库）
func GetUserFromJWT(c *fiber.Ctx) (*models.User, error) {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	role := claims["role"].(string)
	return &models.User{ID: userID, Role: role}, nil
}
