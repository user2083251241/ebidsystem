package container

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/user2083251241/ebidsystem/internal/domain/repository"
	"github.com/user2083251241/ebidsystem/internal/infrastructure/database"
	"github.com/user2083251241/ebidsystem/internal/usecase/order/orderusecase"
	"github.com/user2083251241/ebidsystem/internal/usecase/user/authusecase"
	"github.com/user2083251241/ebidsystem/pkg/utils"
)

type Container struct {
	DB      *gorm.DB
	Redis   *utils.Client
	UserUC  authusecase.UseCase
	OrderUC orderusecase.UseCase
}

func NewContainer() (*Container, error) {
	// 初始化数据库连接
	db, err := database.NewDB()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	// 初始化 Redis 连接
	redisClient := utils.NewRedisClient()

	// 初始化仓库
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// 初始化使用案例
	userUC := authusecase.NewUserService(userRepo, redisClient)
	orderUC := orderusecase.NewOrderService(db, orderRepo)

	return &Container{
		DB:      db,
		Redis:   redisClient,
		UserUC:  userUC,
		OrderUC: orderUC,
	}, nil
}

func (c *Container) Close() error {
	if err := c.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %v", err)
	}
	return nil
}
