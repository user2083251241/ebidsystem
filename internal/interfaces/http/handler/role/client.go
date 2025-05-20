package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user2083251241/ebidsystem/middleware"
	"github.com/user2083251241/ebidsystem/services"
)

type ClientController struct {
	orderService *services.OrderService
}

func NewClientController(os *services.OrderService) *ClientController {
	return &ClientController{orderService: os}
}

// 客户创建买入订单：
func (cc *ClientController) ClientCreateBuyOrder(c *fiber.Ctx) error {
	user, err := middleware.GetUserFromJWT(c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusUnauthorized, "身份验证失败")
	}

	var req services.CreateBuyOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "请求格式错误")
	}

	order, err := cc.orderService.CreateClientBuyOrder(user.ID, req)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

// 客户查看可买订单：
func (cc *ClientController) ClientGetOrders(c *fiber.Ctx) error {
	orders, err := cc.orderService.GetAvailableSellOrders()
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(orders)
}
