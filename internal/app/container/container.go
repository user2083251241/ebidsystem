package container

import (
	"context"
	"fmt"
	"log"
	"sync"

	"gorm.io/gorm"

	"github.com/redis/go-redis/v8"
	"github.com/user2083251241/ebidsystem/internal/domain/repository"
	mysqlrepo "github.com/user2083251241/ebidsystem/internal/infrastructure/persistence/mysql"
	"github.com/user2083251241/ebidsystem/internal/usecase/order/orderusecase"
	"github.com/user2083251241/ebidsystem/internal/usecase/user/authusecase"
	"github.com/user2083251241/ebidsystem/pkg/logger"
)

type Container struct {
	config *Config
	db     *gorm.DB
	redis  *redis.Client
	mu     sync.RWMutex

	// 仓库
	userRepo  repository.UserRepository
	orderRepo repository.OrderRepository

	// 用例
	userUC  authusecase.UseCase
	orderUC orderusecase.UseCase

	// 其他服务
	logger *logger.Logger
}

type Config struct {
	DB struct {
		DSN         string `mapstructure:"dsn"`
		MaxIdleConn int    `mapstructure:"max_idle_conn"`
		MaxOpenConn int    `mapstructure:"max_open_conn"`
	} `mapstructure:"database"`
	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
	JWT struct {
		Secret          string `mapstructure:"secret"`
		AccessTokenTTL  int    `mapstructure:"access_token_ttl"`
		RefreshTokenTTL int    `mapstructure:"refresh_token_ttl"`
	} `mapstructure:"jwt"`
	AppEnv string `mapstructure:"app_env"`
}

func NewContainer(cfg *Config) (*Container, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	c := &Container{
		config: cfg,
		logger: logger.NewLogger(cfg.AppEnv == "production"),
	}

	// 初始化数据库
	if err := c.initDB(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// 初始化 Redis
	if err := c.initRedis(); err != nil {
		return nil, fmt.Errorf("failed to initialize redis: %w", err)
	}

	// 初始化仓库
	c.initRepositories()

	// 初始化用例
	c.initUseCases()

	return c, nil
}

// 初始化数据库
func (c *Container) initDB() error {
	dbClient, err := mysqlrepo.NewDBClient(&mysqlrepo.Config{
		DSN:         c.config.DB.DSN,
		MaxIdleConn: c.config.DB.MaxIdleConn,
		MaxOpenConn: c.config.DB.MaxOpenConn,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	c.db = dbClient.DB
	return nil
}

// 初始化 Redis
func (c *Container) initRedis() error {
	c.redis = redis.NewClient(&redis.Options{
		Addr:     c.config.Redis.Addr,
		Password: c.config.Redis.Password,
		DB:       c.config.Redis.DB,
	})

	// 测试连接
	_, err := c.redis.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}
	return nil
}

// initRepositories 初始化仓库
func (c *Container) initRepositories() {
	c.userRepo = mysqlrepo.NewUserRepository(c.db)
	// 初始化其他仓储...
}

// initUseCases 初始化用例
func (c *Container) initUseCases() {
	secret, accessTTL, refreshTTL := c.GetJWTConfig()
	c.userUC = authusecase.NewUserService(c.userRepo, c.redis, secret, accessTTL, refreshTTL)
	c.orderUC = orderusecase.NewOrderService(c.db, c.orderRepo)
}

// GetDB 获取数据库实例
func (c *Container) GetDB() *gorm.DB {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.db
}

// GetRedis 获取 Redis 客户端
func (c *Container) GetRedis() *redis.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.redis
}

// GetJWTConfig 获取 JWT 配置
func (c *Container) GetJWTConfig() (string, int, int) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config.JWT.Secret, c.config.JWT.AccessTokenTTL, c.config.JWT.RefreshTokenTTL
}

// GetUserUseCase 获取用户用例
func (c *Container) GetUserUseCase() authusecase.UseCase {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.userUC
}

// GetOrderUseCase 获取订单用例
func (c *Container) GetOrderUseCase() orderusecase.UseCase {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.orderUC
}

// GetLogger 获取日志记录器
func (c *Container) GetLogger() *logger.Logger {
	return c.logger
}

// Close 关闭容器，释放资源
func (c *Container) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var errs []error

	// 关闭数据库连接
	if c.db != nil {
		if sqlDB, err := c.db.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errs = append(errs, fmt.Errorf("error closing database: %w", err))
			}
		}
	}

	// 关闭 Redis 连接
	if c.redis != nil {
		if err := c.redis.Close(); err != nil {
			errs = append(errs, fmt.Errorf("error closing redis: %w", err))
		}
	}

	// 如果有错误，返回第一个错误
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

// Shutdown 优雅关闭
func (c *Container) Shutdown() {
	if err := c.Close(); err != nil {
		log.Printf("Error during container shutdown: %v", err)
	}
}
