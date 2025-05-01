package controllers

import (
	// "github.com/champNoob/ebidsystem/backend/services"

	"github.com/champNoob/ebidsystem/backend/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/user2083251241/ebidsystem/services"
)

// controllers/trader.go
type TraderController struct {
	orderService *services.OrderService
}

func NewTraderController(os *services.OrderService) *TraderController {
	return &TraderController{orderService: os}
}

// TraderGetAllOrders 交易员查看所有订单
func (tc *TraderController) TraderGetAllOrders(c *fiber.Ctx) error {
	orders, err := tc.orderService.GetAllOrders()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(orders)
}

// TraderCancelOrder 交易员紧急撤单
func (tc *TraderController) TraderCancelOrder(c *fiber.Ctx) error {
	orderID, _ := c.ParamsInt("id")
	user, _ := middleware.GetCurrentUserFromJWT(c)
	if err := tc.orderService.EmergencyCancelOrder(uint(orderID), user.ID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(200)
}

// 交易员查看所有订单（无过滤）：

func (tc *TraderController) GetAllOrders(c *fiber.Ctx) error {
	orders, err := tc.orderService.GetOrders(services.QueryCondition{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "查询失败"})
	}
	return c.JSON(orders)
}
