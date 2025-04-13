package controllers

import (
	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/gofiber/fiber/v2"
)

// 客户竞拍（买入订单）：
func CreateBuyOrder(c *fiber.Ctx) error {
	type BuyRequest struct {
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"` // 仅限限价单需要
	}
	var req BuyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	// 获取目标卖家订单 ID:
	sellerOrderID := c.Params("id")
	var sellerOrder models.Order
	if err := db.Where("id = ? AND direction = 'sell' AND status = 'pending'", sellerOrderID).
		First(&sellerOrder).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "卖家订单无效或已关闭"})
	}
	// 创建买入订单：
	userID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}
	buyOrder := models.Order{
		UserID:    userID,
		Symbol:    sellerOrder.Symbol,
		Quantity:  req.Quantity,
		Price:     req.Price,
		OrderType: "limit", // 或根据需求调整
		Direction: "buy",
		Status:    "pending",
	}
	if err := db.Create(&buyOrder).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建买入订单失败"})
	}
	// 更新卖家订单状态：
	return c.JSON(buyOrder)
}

// 客户查看自己订单：
func GetClientOrders(c *fiber.Ctx) error {
	//查询所有未取消的卖家订单（隐藏卖家ID）：
	var orders []models.Order
	if err := db.Where("direction = 'sell' AND status NOT IN ('cancelled', 'filled')").
		Select("id, symbol, quantity, price, order_type, created_at"). //排除敏感字段
		Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询订单失败"})
	}
	return c.JSON(orders)
}
