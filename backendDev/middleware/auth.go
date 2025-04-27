package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
