package handler

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/user2083251241/ebidsystem/services"
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
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if validationErrors, ok := err.(*validator.ValidationErrors); ok {
			return fiber.NewError(fiber.StatusBadRequest, validationErrors.Error()) //取第一个验证错误信息
		}
	}
	return nil
}
