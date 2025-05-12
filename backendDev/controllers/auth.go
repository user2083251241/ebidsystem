package controllers

import (
	"log"
	"time"

	"github.com/champNoob/ebidsystem/backend/config"
	"github.com/champNoob/ebidsystem/backend/middleware"
	"github.com/champNoob/ebidsystem/backend/services"
	"github.com/champNoob/ebidsystem/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5" //统一使用 v5
)

type AuthController struct {
	userService  *services.UserService
	orderService *services.OrderService
}

func NewAuthController(us *services.UserService) *AuthController {
	return &AuthController{
		userService:  us,
		orderService: services.NewOrderService(config.DB), //确保传递有效的 db 实例
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

// 用户注销（软删除）：
func (ac *AuthController) Logout(c *fiber.Ctx) error {
	log.Println("注销请求开始处理")
	// 从 JWT 中提取用户信息：
	user, err := middleware.GetUserFromJWT(c)
	if err != nil || user == nil {
		log.Printf("用户信息提取失败 | 错误: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "身份验证失败：无法获取用户信息",
		})
	}
	log.Printf("用户信息获取成功 | 用户ID: %d", user.ID)
	// 获取 Token 原始字符串：
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok || token == nil {
		log.Println("Token 格式无效")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "无效的 Token 格式",
		})
	}
	tokenString := token.Raw
	log.Printf("Token 字符串: %s", tokenString)
	// 计算 Token 剩余有效期：
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("Token 声明解析失败 | Token: %s", tokenString)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token 声明解析失败",
		})
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		log.Printf("Token 有效期字段无效 | Exp: %v", claims["exp"])
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token 有效期无效",
		})
	}
	// 计算剩余有效期：
	expiration := time.Until(time.Unix(int64(exp), 0))
	if expiration < 0 {
		log.Printf("Token 已过期 | Exp: %v", time.Unix(int64(exp), 0))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token 已过期",
		})
	}
	log.Printf("Token 剩余有效期: %v", expiration)
	// 将 Token 加入 Redis 黑名单：
	if err := utils.AddToBlacklist(tokenString, expiration); err != nil {
		log.Printf("Token 加入黑名单失败 | 错误: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "注销失败：系统错误",
		})
	}
	log.Println("Token 已加入黑名单")
	// 取消用户未完成订单：
	if err := ac.orderService.CancelUserUnfinishedOrders(user.ID); err != nil {
		log.Printf("取消订单失败 | 用户ID: %d | 错误: %v", user.ID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "注销失败：无法取消未完成订单",
		})
	}
	log.Printf("已取消未完成订单 | 用户ID: %d", user.ID)
	// 标记用户为已删除（软删除）：
	if err := ac.userService.DeleteUser(user.ID); err != nil {
		log.Printf("用户软删除失败 | 用户ID: %d | 错误: %v", user.ID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "注销失败：用户状态更新失败",
		})
	}
	log.Printf("用户标记为已删除 | 用户ID: %d", user.ID)
	// 返回成功响应：
	log.Printf("用户注销成功 | 用户ID: %d", user.ID)
	return c.JSON(fiber.Map{
		"message": "用户已成功注销",
	})
}
