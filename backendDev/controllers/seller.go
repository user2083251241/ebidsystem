package controllers

import (
	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// 卖家创建卖出订单：
func CreateSellOrder(c *fiber.Ctx) error {
	type SellOrderRequest struct { //定义请求体
		Symbol    string  `json:"symbol"`
		Quantity  int     `json:"quantity"`
		Price     float64 `json:"price"`
		OrderType string  `json:"type"` //market/limit
	}
	var req SellOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	userID := getCurrentUserID(c) //从 JWT 中获取卖家 ID
	// 校验数量：
	if req.Quantity <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "数量必须大于0"})
	}
	// 校验价格（限价单需价格）：
	if req.OrderType == "limit" && req.Price <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "限价单需指定有效价格"})
	}
	// 创建订单：
	order := models.Order{
		UserID:    userID,
		Symbol:    req.Symbol,
		Quantity:  req.Quantity,
		Price:     req.Price,
		OrderType: req.OrderType,
		Direction: "sell", //强制设置买卖方向为 "sell"
		Status:    "pending",
	}
	if err := db.Create(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建订单失败"})
	}
	// 返回订单信息：
	return c.JSON(order)
}

// 更新订单
func UpdateOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")
	userID := getCurrentUserID(c)
	var req struct {
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	var order models.Order
	// 权限和内容判断：
	if err := db.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "订单不存在或无权修改"})
	}
	// 执行更新：
	if err := db.Model(&order).Updates(models.Order{
		Quantity: req.Quantity,
		Price:    req.Price,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新失败"})
	}

	return c.JSON(order)
}

// 单个撤单（软删除，使用 DELETE 方法）：
func CancelOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")
	userID := getCurrentUserID(c)
	var order models.Order
	if err := db.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "订单不存在或无权操作"})
	}
	if err := db.Model(&order).Update("status", "cancelled").Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "取消失败"})
	}
	return c.JSON(fiber.Map{"message": "订单已取消"})
}

// 批量撤单（软删除，使用 POST 方法）：
func BatchCancelOrders(c *fiber.Ctx) error {
	type BatchCancelRequest struct {
		OrderIDs []uint `json:"order_ids"`
	}
	var req BatchCancelRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	userID := getCurrentUserID(c)
	// 开启事务：
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 批量更新：
	if err := tx.Model(&models.Order{}).
		Where("user_id = ? AND id IN ?", userID, req.OrderIDs).
		Update("status", "cancelled").Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "批量取消失败"})
	}
	tx.Commit()
	return c.JSON(fiber.Map{"message": "批量取消成功"})
}

// 卖家查看自己订单：
func GetSellerOrders(c *fiber.Ctx) error {
	userID := uint(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	var orders []models.Order
	if err := db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询订单失败"})
	}
	return c.JSON(orders)
}

// 卖家授权销售：
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
