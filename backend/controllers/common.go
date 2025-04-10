package controllers

import "gorm.io/gorm"

// 全局数据库实例（供其他控制器使用）
var db *gorm.DB

// 初始化数据库连接
func InitDB(database *gorm.DB) {
	db = database
}
