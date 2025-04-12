package middleware

import (
	_ "github.com/gofiber/fiber/v2"
	// _ jwtware "github.com/gofiber/contrib/jwt"
	_ "github.com/champNoob/ebidsystem/backend/config"
)

// JWT 中间件初始化（已由 routes/api.go 直接调用，此文件可省略）
// 保留此文件以备扩展自定义逻辑
