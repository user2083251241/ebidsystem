package services

import (
	"log"
	"time"

	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/champNoob/ebidsystem/backend/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db           *gorm.DB
	orderService *OrderService
}

func NewUserService(db *gorm.DB, orderService *OrderService) *UserService {
	return &UserService{
		db:           db,
		orderService: orderService,
	}
}

// 用户注册：
func (us *UserService) Register(req RegisterRequest) (*models.User, error) {
	var existingUser models.User
	if err := us.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, fiber.NewError(fiber.StatusConflict, "Username already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}
	user := models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
	}
	if err := us.db.Create(&user).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
	}
	return &user, nil
}

// 用户登录：
func (us *UserService) Login(req LoginRequest) (*models.User, error) {
	var user models.User
	// 验证用户名：
	if err := us.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "用户名或密码错误")
	}
	// 检查用户是否已注销：
	if user.IsDeleted {
		return nil, fiber.NewError(fiber.StatusForbidden, "用户已注销")
	}
	// 验证密码：
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "用户名或密码错误")
	}
	return &user, nil
}

// 用户注销：
func (us *UserService) Logout(userID uint, token string, tokenExp time.Duration) error {
	return us.db.Transaction(func(tx *gorm.DB) error {
		// 1. 将 token 加入黑名单
		if err := utils.AddToBlacklist(token, tokenExp); err != nil {
			log.Printf("Token加入黑名单失败 | 用户ID: %d | 错误: %v", userID, err)
			return fiber.NewError(fiber.StatusInternalServerError, "注销失败：token处理错误")
		}

		// 2. 取消订单和授权
		if err := us.orderService.CancelUserUnfinishedOrders(userID); err != nil {
			log.Printf("取消订单失败 | 用户ID: %d | 错误: %v", userID, err)
			return fiber.NewError(fiber.StatusInternalServerError, "注销失败：订单取消错误")
		}

		// 3. 标记用户为已删除
		if err := tx.Model(&models.User{}).
			Where("id = ?", userID).
			Update("is_deleted", true).Error; err != nil {
			log.Printf("用户标记删除失败 | 用户ID: %d | 错误: %v", userID, err)
			return fiber.NewError(fiber.StatusInternalServerError, "注销失败：用户状态更新错误")
		}

		log.Printf("用户注销成功 | 用户ID: %d", userID)
		return nil
	})
}

// 检查用户名是否已存在：
func (us *UserService) CheckUsernameExists(username string) (bool, error) {
	var count int64
	err := us.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// 创建用户：
func (us *UserService) CreateUser(user *models.User) error {
	return us.db.Create(user).Error
}

func (us *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := us.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 删除用户（标记为已删除）：
func (us *UserService) DeleteUser(userID uint) error {
	result := us.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_deleted", true)

	if result.Error != nil {
		log.Printf("用户软删除失败 | 用户ID: %d | 错误: %v", userID, result.Error)
		return result.Error
	}

	log.Printf("用户标记为已删除 | 用户ID: %d | 影响行数: %d", userID, result.RowsAffected)
	return nil
}
