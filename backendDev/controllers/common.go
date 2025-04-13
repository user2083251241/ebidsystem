package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// 全局数据库实例（供其他控制器使用）
var db *gorm.DB

// 初始化数据库连接
func InitDB(database *gorm.DB) {
	db = database
}

func getCurrentUserID(c *fiber.Ctx) (userID uint, err error) {
	token, ok := c.Locals("user").(*jwt.Token) //获取 token
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "无效的 Token")
	}
	claims, ok := token.Claims.(jwt.MapClaims) //获取 claims
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Token 解析失败")
	}
	userIDFloat, ok := claims["user_id"].(float64) //先转为 float64
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "用户 ID 无效")
	}
	userID = uint(userIDFloat) //再转为 uint
	return
}
