package controllers

import (
	"github.com/champNoob/ebidsystem/backend/services"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// 基础控制器（复用请求解析和校验）
type BaseController struct {
	OrderService *services.OrderService
}

// 通用请求解析与校验
func (bc *BaseController) ValidateRequest(c *fiber.Ctx, req interface{}) error {
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "请求格式错误")
	}
	if err := bc.OrderService.Validate(req); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return fiber.NewError(fiber.StatusInternalServerError, "验证器错误")
		}
		for _, err := range err.(validator.ValidationErrors) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}
	return nil
}
