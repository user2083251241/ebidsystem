package controllers

import (
	"github.com/champNoob/ebidsystem/backend/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// 用户注册
func Register(c *fiber.Ctx) error {
	type RegisterRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// 密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// 创建用户
	user := models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Username already exists"})
	}

	return c.JSON(fiber.Map{"message": "User created successfully"})
}

// 用户登录（返回 JWT 示例）
func Login(c *fiber.Ctx) error {
	// 实际项目中应实现 JWT 生成逻辑
	return c.JSON(fiber.Map{"token": "dummy_jwt_token"})
}
