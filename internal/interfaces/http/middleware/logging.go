package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func LoggingMiddleware(c *fiber.Ctx) error {
	// 请求前日志：
	fmt.Printf("[%s] %s - Request\n", c.Method(), c.Path())
	// 继续处理请求：
	err := c.Next()
	// 请求后日志：
	fmt.Printf("[%s] %s - Response Status: %d\n", c.Method(), c.Path(), c.Response().StatusCode())
	return err
}
