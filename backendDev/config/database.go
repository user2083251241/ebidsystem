package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局数据库实例：
var DB *gorm.DB

// 初始化数据库连接：
func InitDB() (DB *gorm.DB, err error) {
	dsn := Get("DB_DSN")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return
}
