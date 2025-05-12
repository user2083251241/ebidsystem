package middleware

import (
	"github.com/champNoob/ebidsystem/backend/config"
	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/champNoob/ebidsystem/backend/utils"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
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

// 从 JWT 中提取用户信息（不查数据库）：
func GetUserFromJWT(c *fiber.Ctx) (*models.User, error) {
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok || token == nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "未找到有效的 JWT 令牌")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "无效的 JWT 声明")
	}
	userIDInterface, exists := claims["user_id"]
	if !exists {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "JWT 声明中缺少 user_id")
	}
	userIDFloat, ok := userIDInterface.(float64)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "user_id 不是有效的浮点数")
	}
	userID := uint(userIDFloat)
	role, ok := claims["role"].(string)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "role 不是有效的字符串")
	}
	ansUser := models.User{
		ID:   userID,
		Role: role,
	}
	return &ansUser, nil
}

// 辅助函数：从 JWT 提取 user_id：
func GetUserIDFromJWT(c *fiber.Ctx) (uint, error) {
	user, err := GetUserFromJWT(c)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func CheckTokenRevoked() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Locals("user").(*jwt.Token)
		if utils.IsTokenRevoked(token.Raw) {
			return c.Status(401).JSON(fiber.Map{"error": "Token 已失效"})
		}
		return c.Next()
	}
}
