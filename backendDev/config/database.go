package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局数据库实例：
var DB *gorm.DB

// 初始化数据库连接：
func InitDB() (*gorm.DB, error) {
	dsn := Get("DB_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("数据库连接字符串未配置，请设置 DB_DSN 环境变量或在配置文件中设置")
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("数据库初始化失败: %v", err)
		return nil, fmt.Errorf("数据库初始化失败: %w", err)
	}
	return DB, nil
}
