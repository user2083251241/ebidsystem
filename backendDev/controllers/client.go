package controllers

import (
	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// 客户竞拍（买入订单）
func CreateBuyOrder(c *fiber.Ctx) error {
	type BuyOrderRequest struct {
		Symbol    string  `json:"symbol"`
		Quantity  int     `json:"quantity"`
		Price     float64 `json:"price"`
		OrderType string  `json:"type"` // market/limit
	}
	var req BuyOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}

	userID := uint(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	order := models.Order{
		UserID:    userID,
		Symbol:    req.Symbol,
		Quantity:  req.Quantity,
		Price:     req.Price,
		OrderType: req.OrderType,
		Direction: "buy",
		Status:    "pending",
	}

	if err := db.Create(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建订单失败"})
	}

	return c.JSON(order)
}

// 客户查看自己订单
func GetClientOrders(c *fiber.Ctx) error {
	userID := uint(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	var orders []models.Order
	if err := db.Where("user_id = ? AND direction = 'buy'", userID).Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询订单失败"})
	}
	return c.JSON(orders)
}
