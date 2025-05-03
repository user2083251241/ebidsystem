package controllers

import (
	"github.com/champNoob/ebidsystem/backend/middleware"
	"github.com/champNoob/ebidsystem/backend/services"
	"github.com/gofiber/fiber/v2"
)

type SellerController struct {
	orderService *services.OrderService
}

func NewSellerController(os *services.OrderService) *SellerController {
	return &SellerController{orderService: os}
}

// SellerCreateOrder 卖家创建正式订单
func (sc *SellerController) SellerCreateOrder(c *fiber.Ctx) error {
	user, err := middleware.GetUserFromJWT(c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusUnauthorized, "身份验证失败")
	}

	var req services.CreateSellerOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "请求格式错误")
	}

	order, err := sc.orderService.CreateSellerOrder(user, req)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

// SellerUpdateOrder 卖家修改订单
func (sc *SellerController) SellerUpdateOrder(c *fiber.Ctx) error {
	user, err := middleware.GetUserFromJWT(c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusUnauthorized, "身份验证失败")
	}

	orderID, err := c.ParamsInt("id")
	if err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "无效订单ID")
	}

	var req services.UpdateSellerOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "请求格式错误")
	}

	if err := sc.orderService.UpdateSellerOrder(user.ID, uint(orderID), req); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
func (sc *SalesController) SellerGetOrders(c *fiber.Ctx) error {
	user, err := middleware.GetUserFromJWT(c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusUnauthorized, "身份验证失败")
	}

	orders, err := sc.orderService.GetSellerOrders(user.ID)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(orders)
}

// SellerAuthorizeSales 卖家授权销售
func (sc *SellerController) SellerAuthorizeSales(c *fiber.Ctx) error {
	user, err := middleware.GetUserFromJWT(c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusUnauthorized, "身份验证失败")
	}

	var req services.CreateAuthorizationRequest
	if err := c.BodyParser(&req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "请求格式错误")
	}

	auth, err := sc.orderService.CreateSalesAuthorization(user.ID, req)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(auth)
}

// SellerGetOrders 卖家查看订单
func (sc *SellerController) SellerGetOrders(c *fiber.Ctx) error {
	user, err := middleware.GetUserFromJWT(c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusUnauthorized, "身份验证失败")
	}

	orders, err := sc.orderService.GetSellerOrders(user.ID)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(orders)
}

// SellerCancelOrder 卖家取消订单
func (sc *SellerController) SellerCancelOrder(c *fiber.Ctx) error {
	user, err := middleware.GetUserFromJWT(c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusUnauthorized, "身份验证失败")
	}

	orderID, err := c.ParamsInt("id")
	if err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "无效订单ID")
	}

	if err := sc.orderService.CancelSellerOrder(user.ID, uint(orderID)); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

// SellerBatchCancelOrders 卖家批量取消订单
func (sc *SellerController) SellerBatchCancelOrders(c *fiber.Ctx) error {
	user, err := middleware.GetUserFromJWT(c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusUnauthorized, "身份验证失败")
	}

	var req struct {
		OrderIDs []uint `json:"order_ids"`
	}
	if err := c.BodyParser(&req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "请求格式错误")
	}

	if err := sc.orderService.BatchCancelSellerOrders(user.ID, req.OrderIDs); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
