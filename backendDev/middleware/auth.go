package middleware

import (
	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

// 轻量级用户信息（仅从 JWT 提取）：
func GetCurrentUserFromJWT(c *fiber.Ctx) (*models.User, error) {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	role := claims["role"].(string)
	return &models.User{ID: userID, Role: role}, nil
}

// 辅助函数：从 JWT 提取 user_id
func GetCurrentUserIDFromJWT(c *fiber.Ctx) (uint, error) {
	user, err := GetCurrentUserFromJWT(c)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func AttachUserToContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := GetCurrentUserFromJWT(c)
		if err != nil {
			return err
		}
		c.Locals("currentUser", user) // 存入上下文
		return c.Next()
	}
}

// 完整用户信息（需查数据库）
func GetFullUserInfoFromDB(c *fiber.Ctx, db *gorm.DB) (*models.User, error) {
	userID, err := GetCurrentUserIDFromJWT(c)
	if err != nil {
		return nil, err
	}
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "用户不存在")
	}
	return &user, nil
}

func DraftAuthorization() fiber.Handler {
	return func(c *fiber.Ctx) error {
		draftID := c.Params("id")
		var draft models.DraftOrder
		if err := db.First(&draft, draftID).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "草稿不存在"})
		}
		user, _ := GetCurrentUserFromJWT(c)
		if draft.CreatorID != user.ID {
			return c.Status(403).JSON(fiber.Map{"error": "无权操作此草稿"})
		}
		return c.Next()
	}
}
