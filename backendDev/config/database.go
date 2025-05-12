package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局数据库实例：
var DB *gorm.DB

// 初始化数据库连接：
func InitDB() (*gorm.DB, error) {
	dsn := Get("DB_DSN")
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库初始化失败")
		return nil, err
	}
	return DB, nil
}
