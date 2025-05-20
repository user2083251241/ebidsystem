package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/user2083251241/ebidsystem/internal/app/config"
	"github.com/user2083251241/ebidsystem/internal/interfaces/http/dto"
	"github.com/user2083251241/ebidsystem/internal/interfaces/http/middleware"
)

type AuthController struct {
	userService *user.UseCase
}

func NewAuthController(us *user.UseCase) *AuthController {
	return &AuthController{
		userService: us,
	}
}

// 注册：
func (ac *AuthController) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}
	// 输入校验
	if req.Username == "" || req.Password == "" || req.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, password, and role are required",
		})
	}
	// 检查角色合法性
	if req.Role != "seller" && req.Role != "client" && req.Role != "sales" && req.Role != "trader" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role. Allowed values: client, sales, trader",
		})
	}
	user, err := ac.userService.Register(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"user_id": user.ID,
	})
}

// 登录：
func (ac *AuthController) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}
	user, err := ac.userService.Login(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 过期时间 3 天
	})
	tokenString, err := token.SignedString([]byte(config.Get("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}
	return c.JSON(fiber.Map{
		"token": tokenString,
		"role":  user.Role, // 返回用户角色
	})
}

// 用户注销：
func (ac *AuthController) Logout(c *fiber.Ctx) error {
	// 1. 获取用户信息
	user, err := middleware.GetUserFromJWT(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "身份验证失败",
		})
	}

	// 2. 获取token信息
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	exp, ok := claims["exp"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token格式无效",
		})
	}
	expiration := time.Until(time.Unix(int64(exp), 0))

	// 3. 执行注销
	if err := ac.userService.Logout(user.ID, token.Raw, expiration); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "用户已成功注销",
	})
}
