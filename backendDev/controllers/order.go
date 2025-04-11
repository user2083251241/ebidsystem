package controllers

import (
	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// 创建订单：
func CreateOrder(c *fiber.Ctx) error {
	type OrderRequest struct {
		Symbol    string  `json:"symbol"`
		Quantity  int     `json:"quantity"`
		Direction string  `json:"direction"` // "buy" 或 "sell"
		OrderType string  `json:"type"`      // "market" 或 "limit"
		Price     float64 `json:"price"`     // 限价单必填
	}
	var req OrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	// ============== 输入校验 ==============
	// 1.校验订单类型合法性：
	if req.OrderType != "market" && req.OrderType != "limit" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid order type. Allowed values: market, limit",
		})
	}
	// 2.校验限价单价格：
	if req.OrderType == "limit" && req.Price <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Limit orders require a valid price (> 0)",
		})
	}
	// 3.校验数量：
	if req.Quantity <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Quantity must be greater than 0",
		})
	}
	// ============== 业务逻辑 ==============
	// 从 JWT 中获取用户 ID:
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	// 创建订单：
	order := models.Order{
		UserID:    userID,
		Symbol:    req.Symbol,
		Quantity:  req.Quantity,
		Direction: req.Direction,
		OrderType: req.OrderType,
		Price:     req.Price,
		Status:    "pending",
	}
	// 返回（错误）信息：
	if err := db.Create(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create order"})
	}
	return c.JSON(order)
}

// 查询订单：
func GetOrders(c *fiber.Ctx) error {
	// 从 JWT 中获取用户 ID
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	var orders []models.Order
	if err := db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询订单失败"})
	}

	return c.JSON(orders)
}

// 取消订单（示例）
func CancelOrder(c *fiber.Ctx) error {
	// ...（实现逻辑）
	return c.JSON(fiber.Map{"message": "订单已取消"})
}
