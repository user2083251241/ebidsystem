package controllers

import (
	"time"

	"github.com/champNoob/ebidsystem/backend/config"
	"github.com/champNoob/ebidsystem/backend/middleware"
	"github.com/champNoob/ebidsystem/backend/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5" //统一使用 v5
)

type AuthController struct {
	userService  *services.UserService
	orderService *services.OrderService
}

func NewAuthController(us *services.UserService) *AuthController {
	return &AuthController{
		userService: us,
	}
}

// 注册：
func (ac *AuthController) Register(c *fiber.Ctx) error {
	var req services.RegisterRequest
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
	var req services.LoginRequest
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

// 注销：
/* 注销逻辑（软删除）*/
/*当前代码仅标记用户为已删除，但已签发的 JWT 仍可正常使用（需结合黑名单等机制彻底禁用）！
若需立即失效 Token，必须实现 Token 吊销逻辑（如基于 Redis 的短期 Token 有效期）！*/
func (ac *AuthController) Logout(c *fiber.Ctx) error {
	user, err := middleware.GetUserFromJWT(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to get user information",
		})
	}
	userID := user.ID
	// 取消所有未完成订单：
	if err := ac.orderService.CancelUserUnfinishedOrders(userID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "注销失败"})
	}
	// 删除用户（标记为已删除）：
	if err := ac.userService.DeleteUser(userID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "注销失败"})
	}
	// 返回成功信息：
	return c.JSON(fiber.Map{"message": "用户已注销"})
}
