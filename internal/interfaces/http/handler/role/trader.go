package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user2083251241/ebidsystem/middleware"
	"github.com/user2083251241/ebidsystem/services"
)

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
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(orders)
}

// TraderEmergencyCancel 交易员紧急撤单
func (tc *TraderController) TraderEmergencyCancel(c *fiber.Ctx) error {
	user, err := middleware.GetUserFromJWT(c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusUnauthorized, "身份验证失败")
	}

	orderID, err := c.ParamsInt("id")
	if err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "无效订单ID")
	}

	if err := tc.orderService.EmergencyCancelOrder(uint(orderID), user.ID); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
