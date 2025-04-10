package controllers

import (
	"github.com/champNoob/ebidsystem/backend/models"

	"github.com/gofiber/fiber/v2"
)

// 创建订单
func CreateOrder(c *fiber.Ctx) error {
	type OrderRequest struct {
		Symbol    string `json:"symbol"`
		Quantity  int    `json:"quantity"`
		Direction string `json:"direction"` // "buy" or "sell"
	}

	var req OrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// 创建订单（简化的逻辑）
	order := models.Order{
		UserID:    1, // 示例用户ID，实际应从 JWT 中获取
		Symbol:    req.Symbol,
		Quantity:  req.Quantity,
		Direction: req.Direction,
		Status:    "pending",
	}

	if err := db.Create(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create order"})
	}

	return c.JSON(order)
}

// 查询订单
func GetOrders(c *fiber.Ctx) error {
	var orders []models.Order
	if err := db.Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch orders"})
	}
	return c.JSON(orders)
}
