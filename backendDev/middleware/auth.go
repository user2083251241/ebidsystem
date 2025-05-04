package middleware

import (
	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	// "github.com/user2083251241/ebidsystem/models"
	"gorm.io/gorm"
)

func RoleRequired(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		userRole := claims["role"].(string)
		if userRole != role {
			errStr := "权限不足！仅" + userRole + "可访问"
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": errStr})
		}
		return c.Next()
	}
}

func AttachUserToContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := GetUserFromJWT(c)
		if err != nil {
			return err
		}
		c.Locals("curUser", user) // 存入上下文
		return c.Next()
	}
}

// 完整用户信息（需查数据库）
func GetFullUserInfoFromDB(c *fiber.Ctx, db *gorm.DB) (*models.User, error) {
	userID, err := GetUserIDFromJWT(c)
	if err != nil {
		return nil, err
	}
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "用户不存在")
	}
	return &user, nil
}

func DraftAuthorization(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		draftID := c.Params("id")
		var draft models.DraftOrder
		if err := db.First(&draft, draftID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "草稿不存在"})
		}
		user, _ := GetUserFromJWT(c)
		if draft.CreatorID != user.ID {
			return c.Status(403).JSON(fiber.Map{"error": "无权操作此草稿"})
		}
		return c.Next()
	}
}
