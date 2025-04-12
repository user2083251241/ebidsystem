package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// 检查用户是否为卖家
func SellerOnly(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	role := claims["role"].(string)
	if role != "seller" {
		return c.Status(403).JSON(fiber.Map{"error": "仅卖家可执行此操作"})
	}
	return c.Next()
}

// 检查用户是否为销售
func SalesOnly(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	role := claims["role"].(string)
	if role != "sales" {
		return c.Status(403).JSON(fiber.Map{"error": "仅销售可执行此操作"})
	}
	return c.Next()
}

// 检查用户是否为交易员
func TraderOnly(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	role := claims["role"].(string)
	if role != "trader" {
		return c.Status(403).JSON(fiber.Map{"error": "仅交易员可执行此操作"})
	}
	return c.Next()
}
