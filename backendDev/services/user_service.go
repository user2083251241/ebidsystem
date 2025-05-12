package services

import (
	"log"

	"github.com/champNoob/ebidsystem/backend/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
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
	if err := us.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password")
	}
	return &user, nil
}

// 用户注销：
func (us *UserService) Logout(userID uint) error {
	if err := us.db.Model(&models.LiveOrder{}).
		Where("user_id = ? AND status IN ('pending', 'draft')", userID).
		Update("status", "cancelled").Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "注销失败")
	}
	if err := us.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_deleted", true).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "注销失败")
	}
	return nil
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
