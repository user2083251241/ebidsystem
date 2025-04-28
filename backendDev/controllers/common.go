package controllers

import (
	"github.com/champNoob/ebidsystem/backend/models"
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

// 获取当前用户完整信息（包含角色和其他字段）
func GetCurrentUser(c *fiber.Ctx) (*models.User, error) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		return nil, err
	}
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "用户不存在")
	}
	return &user, nil
}

// 统一错误响应：
func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"error": message})
}
