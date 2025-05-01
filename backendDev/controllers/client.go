package controllers

import (
	// "github.com/champNoob/ebidsystem/backend/services"
	"github.com/champNoob/ebidsystem/backend/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/user2083251241/ebidsystem/services"
)

// controllers/client.go
type ClientController struct {
	orderService *services.OrderService
}

func NewClientController(os *services.OrderService) *ClientController {
	return &ClientController{orderService: os}
}

// ClientCreateBuyOrder 客户创建买入订单
func (cc *ClientController) ClientCreateBuyOrder(c *fiber.Ctx) error {
	user, _ := middleware.GetCurrentUserFromJWT(c)
	var req services.CreateBuyOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "请求格式错误"})
	}
	order, err := cc.orderService.CreateClientBuyOrder(user.ID, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(order)
}
