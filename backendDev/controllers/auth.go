package controllers

import (
	"time"

	"github.com/champNoob/ebidsystem/backend/config"
	"github.com/champNoob/ebidsystem/backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5" //统一使用 v5

	// "github.com/gofiber/contrib/jwt" //使用支持 v5 的 Fiber 中间件
	"golang.org/x/crypto/bcrypt"
)

// ========================== 用户注册 ==========================
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"` // 角色: client/sales/trader
}

func Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}
	// 输入校验：
	if req.Username == "" || req.Password == "" || req.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, password, and role are required",
		})
	}
	// 检查角色合法性：
	if req.Role != "seller" && req.Role != "client" && req.Role != "sales" && req.Role != "trader" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role. Allowed values: client, sales, trader",
		})
	}
	// 检查用户名是否已存在：
	var existingUser models.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Username already exists",
		})
	}
	// 哈希加密：
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}
	// 创建用户：
	user := models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
	}
	// 将用户保存到数据库：
	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}
	// 返回成功响应：
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"user_id": user.ID,
	})
}

// ========================== 用户登录 ==========================
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}
	// 查询用户是否存在：
	var user models.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}
	// 验证密码：
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}
	// 生成 JWT:
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 过期时间 3 天
	})

	// 签名令牌：
	tokenString, err := token.SignedString([]byte(config.Get("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ //500 服务器错误
			"error": "Failed to generate token",
		})
	}
	// 返回 JWT:
	return c.JSON(fiber.Map{
		"token": tokenString,
		"role":  user.Role, // 返回用户角色
	})
}

// ========================== 用户注销 ==========================
// 注销逻辑（软删除）：
/*当前代码仅标记用户为已删除，但已签发的 JWT 仍可正常使用（需结合黑名单等机制彻底禁用）！
若需立即失效 Token，必须实现 Token 吊销逻辑（如基于 Redis 的短期 Token 有效期）！*/
func Logout(c *fiber.Ctx) error {
	// 获取用户信息：
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	// 取消所有未完成订单：
	if err := db.Model(&models.Order{}).
		Where("user_id = ? AND status IN ('pending', 'draft')", userID).
		Update("status", "cancelled").Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "注销失败"})
	}
	// 标记用户为已删除：
	if err := db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_deleted", true).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "注销失败"})
	}
	// 返回成功信息：
	return c.JSON(fiber.Map{"message": "用户已注销"})
}
