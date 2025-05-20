package authusecase

import (
	"context"
	"errors"
	"time"

	"github.com/user2083251241/ebidsystem/internal/domain/entity"
	"github.com/user2083251241/ebidsystem/internal/domain/repository"
	"github.com/user2083251241/ebidsystem/pkg/utils"
)

type userService struct {
	userRepo repository.UserRepository
	redis    *utils.RedisClient
}

func NewUserService(repo repository.UserRepository, redis *utils.RedisClient) UseCase {
	return &userService{
		userRepo: repo,
		redis:    redis,
	}
}

func (s *userService) Register(req RegisterRequest) (*entity.User, error) {
	// 验证请求
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(context.Background(), req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// 创建用户实体
	user := &entity.User{
		Username:  req.Username,
		Password:  req.Password, // 实际使用时需要加密
		Role:      req.Role,
		CreatedAt: time.Now(),
	}

	// 保存用户
	if err := s.userRepo.Create(context.Background(), user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(req LoginRequest) (*entity.User, error) {
	// 验证请求
	if err := validateLoginRequest(req); err != nil {
		return nil, err
	}

	// 查找用户
	user, err := s.userRepo.FindByUsername(context.Background(), req.Username)
	if err != nil {
		return nil, err
	}

	// 验证密码
	if !checkPassword(user.Password, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *userService) Logout(userID uint, token string, expiration time.Duration) error {
	// 将 token 加入 Redis 黑名单
	if err := s.redis.AddToBlacklist(context.Background(), token, expiration); err != nil {
		return err
	}
	return nil
}

func validateRegisterRequest(req RegisterRequest) error {
	if req.Username == "" || req.Password == "" || req.Role == "" {
		return errors.New("username, password, and role are required")
	}
}
