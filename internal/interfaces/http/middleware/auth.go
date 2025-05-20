package middleware

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/user2083251241/ebidsystem/internal/app/config"
	"github.com/user2083251241/ebidsystem/internal/domain/entity"
	"github.com/user2083251241/ebidsystem/pkg/utils"
	"gorm.io/gorm"
)

// GetUserIDFromJWT 从JWT中获取用户ID
func GetUserIDFromJWT(c *fiber.Ctx) (uint, error) {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(*JWTClaims)
	return claims.UserID, nil
}

// GetRoleFromJWT 从JWT中获取用户角色
func GetRoleFromJWT(c *fiber.Ctx) string {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(*JWTClaims)
	return claims.Role
}

// GetFullUserInfoFromDB 从数据库获取完整用户信息
func GetFullUserInfoFromDB(c *fiber.Ctx, db *gorm.DB) (*entity.User, error) {
	userID, err := GetUserIDFromJWT(c)
	if err != nil {
		return nil, err
	}

	var user entity.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "用户不存在")
	}
	return &user, nil
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "未提供授权令牌",
			})
		}

		token = token[7:] // 去掉 "Bearer " 前缀
		claims := &JWTClaims{}

		// 验证token
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Get("JWT_SECRET")), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "无效的令牌",
			})
		}

		// 验证token是否在黑名单中
		isBlacklisted, err := checkTokenBlacklist(c.Context(), token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "检查黑名单失败",
			})
		}
		if isBlacklisted {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "令牌已被撤销",
			})
		}

		c.Locals("user", claims)
		return c.Next()
	}
}

// RoleMiddleware 角色验证中间件
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := GetRoleFromJWT(c)
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "无权访问",
		})
	}
}

// DraftAuthorization 草稿授权中间件
func DraftAuthorization(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		draftID := c.Params("id")
		var draft entity.DraftOrder
		if err := db.First(&draft, draftID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "草稿不存在"})
		}

		user, err := GetFullUserInfoFromDB(c, db)
		if err != nil {
			return err
		}

		if draft.CreatorID != user.ID {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "无权访问此草稿",
			})
		}

		c.Locals("draft", draft)
		return c.Next()
	}
}

// checkTokenBlacklist 检查token是否在黑名单中
func checkTokenBlacklist(ctx context.Context, token string) (bool, error) {
	redisClient := utils.RedisClient
	val, err := redisClient.Get(ctx, "jwt_blacklist:"+token).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "1", nil
}
