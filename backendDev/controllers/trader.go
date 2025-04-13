package controllers

import (
	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/gofiber/fiber/v2"
)

// 交易员查看所有订单
func GetAllOrders(c *fiber.Ctx) error {
	var orders []models.Order
	if err := db.Unscoped().Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询订单失败"})
	}
	return c.JSON(orders)
}
