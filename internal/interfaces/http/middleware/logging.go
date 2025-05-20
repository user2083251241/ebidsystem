package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/user2083251241/ebidsystem/internal/app/config"
)

// 日志中间件
func LoggerMiddleware() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} ${status} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   config.Get("TIMEZONE"),
	})
}
