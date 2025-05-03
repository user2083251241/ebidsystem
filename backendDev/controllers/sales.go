package controllers

import (
	"github.com/champNoob/ebidsystem/backend/middleware"
	"github.com/champNoob/ebidsystem/backend/services"
	"github.com/gofiber/fiber/v2"
)

type SalesController struct {
	orderService *services.OrderService
}

func NewSalesController(os *services.OrderService) *SalesController {
	return &SalesController{orderService: os}
}

// SalesCreateDraft 创建草稿订单
func (sc *SalesController) SalesCreateDraft(c *fiber.Ctx) error {
	user, err := middleware.GetUserFromJWT(c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusUnauthorized, "身份验证失败")
	}

	var req services.CreateDraftRequest
	if err := c.BodyParser(&req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "请求格式错误")
	}

	draft, err := sc.orderService.CreateDraftOrder(user.ID, req)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(draft)
}

// SalesUpdateDraft 修改草稿订单
func (sc *SalesController) SalesUpdateDraft(c *fiber.Ctx) error {
	user, err := middleware.GetUserFromJWT(c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusUnauthorized, "身份验证失败")
	}

	draftID, err := c.ParamsInt("id")
	if err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "无效草稿ID")
	}

	var req services.UpdateDraftRequest
	if err := c.BodyParser(&req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "请求格式错误")
	}

	if err := sc.orderService.UpdateDraftOrder(user.ID, uint(draftID), req); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

// SalesSubmitDraft 提交草稿审批
func (sc *SalesController) SalesSubmitDraft(c *fiber.Ctx) error {
	user, err := middleware.GetUserFromJWT(c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusUnauthorized, "身份验证失败")
	}

	draftID, err := c.ParamsInt("id")
	if err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "无效草稿ID")
	}

	if err := sc.orderService.SubmitDraftOrder(user.ID, uint(draftID)); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

// SalesGetAuthorizedDrafts 获取已授权草稿
func (sc *SalesController) SalesGetAuthorizedDrafts(c *fiber.Ctx) error {
	user, err := middleware.GetUserFromJWT(c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusUnauthorized, "身份验证失败")
	}

	drafts, err := sc.orderService.GetAuthorizedDrafts(user.ID)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(drafts)
}

// SalesDeleteDraft 删除草稿
func (sc *SalesController) SalesDeleteDraft(c *fiber.Ctx) error {
	user, err := middleware.GetUserFromJWT(c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusUnauthorized, "身份验证失败")
	}

	draftID, err := c.ParamsInt("id")
	if err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "无效草稿ID")
	}

	if err := sc.orderService.DeleteDraft(user.ID, uint(draftID)); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
