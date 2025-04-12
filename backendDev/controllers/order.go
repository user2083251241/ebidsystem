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
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// ============== 权限校验 ==============
	// 从 JWT 中提取角色：
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	role := claims["role"].(string)
	// 客户只能创建买入订单，交易员可创建卖出订单：
	if req.Direction == "sell" && role != "trader" {
		return c.Status(403).JSON(fiber.Map{"error": "仅交易员可创建卖出订单"})
	}

	// ============== 输入校验 ==============
	// 1. 校验订单类型合法性：
	if req.OrderType != "market" && req.OrderType != "limit" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid order type. Allowed values: market, limit",
		})
	}
	// 2. 校验限价单价格：
	if req.OrderType == "limit" && req.Price <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Limit orders require a valid price (> 0)",
		})
	}
	// 3. 校验数量：
	if req.Quantity <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Quantity must be greater than 0",
		})
	}

	// ============== 业务逻辑 ==============
	userID := uint(claims["user_id"].(float64))
	order := models.Order{
		UserID:    userID,
		Symbol:    req.Symbol,
		Quantity:  req.Quantity,
		Direction: req.Direction,
		OrderType: req.OrderType,
		Price:     req.Price,
		Status:    "pending",
	}
	// 创建订单：
	if err := db.Create(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create order"})
	}
	// 返回订单信息：
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

// 取消订单：
func CancelOrder(c *fiber.Ctx) error {
	// 从请求体中获取订单 ID:
	type CancelRequest struct {
		OrderID uint `json:"order_id"`
	}
	var req CancelRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}
	// 从 JWT 中获取用户 ID 和角色：
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	role := claims["role"].(string)
	// 查询订单：
	var order models.Order
	if err := db.First(&order, req.OrderID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "订单不存在"})
	}
	// 权限校验：仅交易员可取消他人订单，客户只能取消自己的订单：
	if role != "trader" && order.UserID != userID {
		return c.Status(403).JSON(fiber.Map{"error": "无权操作此订单"})
	}
	// 更新订单状态为 "cancelled"：
	if err := db.Model(&order).Update("status", "cancelled").Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "取消失败"})
	}
	// 返回成功信息：
	return c.JSON(fiber.Map{"message": "订单已取消"})
}
