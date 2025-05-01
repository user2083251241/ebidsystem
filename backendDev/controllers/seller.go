package controllers

import (
	"time"

	"github.com/champNoob/ebidsystem/backend/middleware"
	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/user2083251241/ebidsystem/services"
	"github.com/user2083251241/ebidsystem/utils"
)

type SellerController struct {
	orderService *services.OrderService
}

func NewSellerController(os *services.OrderService) *SellerController {
	return &SellerController{orderService: os}
}

// 卖家创建卖出订单：
func (sc *SellerController) SellerCreateOrder(c *fiber.Ctx) error {
	user, _ := middleware.GetCurrentUserFromJWT(c)
	var req services.CreateSellerOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	order, err := sc.orderService.CreateSellerOrder(user, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(order)
}

// 卖家修改订单：
func (sc *SellerController) SellerUpdateOrder(c *fiber.Ctx) error {
	user, _ := middleware.GetCurrentUserFromJWT(c)
	orderID, _ := c.ParamsInt("id")
	var req services.UpdateSellerOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	if err := sc.orderService.UpdateSellerOrder(user.ID, uint(orderID), req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(200)
}

// 向指定销售授权：
func AuthorizeSales(c *fiber.Ctx) error {
	type AuthRequest struct {
		SalesID   uint      `json:"sales_id"`
		ExpiresAt time.Time `json:"expires_at"` // 授权有效期
	}
	var req AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	sellerID, _ := middleware.GetCurrentUserIDFromJWT(c)
	auth := models.SellerSalesAuthorization{
		SellerID:      sellerID,
		SalesID:       req.SalesID,
		Authorization: "pending",
		ExpiresAt:     req.ExpiresAt,
	}
	if err := db.Create(&auth).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "授权请求提交失败"})
	}
	return c.JSON(fiber.Map{"message": "授权请求已发送，等待审批"})
}

func (sc *SellerController) GetOrders(c *fiber.Ctx) error {
	// 从 JWT 中获取用户信息：
	user, err := middleware.GetCurrentUserFromJWT(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "身份验证失败"})
	}
	// 构造查询条件：
	conditions := services.QueryCondition{
		UserID:        &user.ID,          // 卖家只能看自己的订单
		Direction:     utils.Ptr("sell"), // 限制卖出方向
		HideSensitive: false,
	}
	// 调用服务层查询订单：
	orders, err := sc.orderService.GetOrders(conditions)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询失败"})
	}
	// 返回响应：
	return c.JSON(orders)
}
