package controllers

import (
	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// 卖家创建卖出订单
func CreateSellOrder(c *fiber.Ctx) error {
	type SellOrderRequest struct {
		Symbol    string  `json:"symbol"`
		Quantity  int     `json:"quantity"`
		Price     float64 `json:"price"`
		OrderType string  `json:"type"` // market/limit
	}
	var req SellOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	// 从 JWT 中获取卖家 ID:
	userID := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(uint)
	// // 验证买卖方向：
	// if user.Driection == "buy" {
	// 	return c.Status(400).JSON(fiber.Map{"error": "卖方不能买入"})
	// }

	// 创建订单：
	order := models.Order{
		UserID:    userID,
		Symbol:    req.Symbol,
		Quantity:  req.Quantity,
		Price:     req.Price,
		OrderType: req.OrderType,
		Direction: "sell",
		Status:    "pending",
	}

	if err := db.Create(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建订单失败"})
	}

	return c.JSON(order)
}

// 卖家改单
func UpdateOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")
	type UpdateRequest struct {
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	}
	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}

	// 验证订单归属
	userID := uint(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	var order models.Order
	if err := db.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "订单不存在或无权修改"})
	}

	// 更新订单
	if err := db.Model(&order).Updates(models.Order{
		Quantity: req.Quantity,
		Price:    req.Price,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新订单失败"})
	}

	return c.JSON(order)
}

// 卖家查看自己订单
func GetSellerOrders(c *fiber.Ctx) error {
	userID := uint(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	var orders []models.Order
	if err := db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询订单失败"})
	}
	return c.JSON(orders)
}

// 卖家授权销售
func AuthorizeSales(c *fiber.Ctx) error {
	type AuthRequest struct {
		SalesID uint `json:"sales_id"`
	}
	var req AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}

	sellerID := uint(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	auth := models.SellerSalesAuthorization{
		SellerID:      sellerID,
		SalesID:       req.SalesID,
		Authorization: "pending",
	}
	if err := db.Create(&auth).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "授权请求提交失败"})
	}
	return c.JSON(fiber.Map{"message": "授权请求已发送，等待销售确认"})
}
