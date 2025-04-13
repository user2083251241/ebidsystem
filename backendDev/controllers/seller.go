package controllers

import (
	"fmt"

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
	// 从 JWT 中获取卖家 ID:
	userID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "从 JWT 中获取卖家 ID 失败"})
	}
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
	// 从 JWT 中获取卖家 ID:
	userID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	var order models.Order
	// 权限和内容判断：
	if err := db.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "订单不存在或无权修改"})
	}
	if order.Status == "cancelled" {
		return c.Status(400).JSON(fiber.Map{"error": "已取消的订单不可修改"})
	}
	// 定义请求体：
	var req struct {
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	}
	// 判断请求体格式：
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
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
	// 从 JWT 中获取卖家 ID:
	userID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	var order models.Order
	// 查询订单（包括已软删除的记录）：
	if err := db.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "订单不存在或无权操作"})
	}
	// 检查订单状态：
	if order.Status == "cancelled" {
		return c.Status(400).JSON(fiber.Map{"error": "订单已取消，不可重复操作"})
	}
	// 添加事物：
	tx := db.Begin()
	// 执行软删除（设置 DeletedAt）：
	if err := db.Delete(&order).Error; err != nil {
		tx.Rollback() //回滚事务
		return c.Status(500).JSON(fiber.Map{"error": "取消失败"})
	}
	// 更新状态为 "cancelled"（需使用 Unscoped 更新软删除记录）：
	if err := db.Unscoped().Model(&order).Update("status", "cancelled").Error; err != nil {
		tx.Rollback() //回滚事务
		return c.Status(500).JSON(fiber.Map{"error": "状态更新失败"})
	}
	// 提交事务：
	tx.Commit()
	// 返回成功信息：
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
	// 从 JWT 中获取卖家 ID:
	userID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}
	// 开启事务：
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 检查每个订单状态：
	for _, orderID := range req.OrderIDs {
		var order models.Order
		if err := tx.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
			tx.Rollback()
			return c.Status(404).JSON(fiber.Map{"error": fmt.Sprintf("订单 %d 不存在或无权操作", orderID)})
		}
		if order.Status == "cancelled" {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("订单 %d 已取消，不可重复操作", orderID)})
		}
	}
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
