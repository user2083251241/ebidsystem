package mysql

import (
	"fmt"
	"time"

	"github.com/user2083251241/ebidsystem/internal/app/config"
	"github.com/user2083251241/ebidsystem/internal/infrastructure/persistence/migrations"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBClient struct {
	DB *gorm.DB
}

// 创建数据库客户端
func NewDBClient(cfg *config.Database) (*DBClient, error) {
	// 创建数据库连接
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 执行数据库迁移
	if err := migrations.RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &DBClient{DB: db}, nil
}

// 关闭数据库连接
func (c *DBClient) Close() error {
	if c.DB == nil {
		return nil
	}

	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
