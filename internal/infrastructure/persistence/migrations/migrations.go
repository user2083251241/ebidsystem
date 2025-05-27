package migrations

import (
	"gorm.io/gorm"

	"github.com/user2083251241/ebidsystem/internal/domain/entity"
)

// 执行数据库迁移
func RunMigrations(db *gorm.DB) error {
	// 自动迁移所有模型
	err := db.AutoMigrate(
		&entity.User{},
		&entity.SellerSalesAuthorization{},
		&entity.LiveOrder{},
		&entity.DraftOrder{},
		&entity.Trade{},
	)
	if err != nil {
		return err
	}

	// 添加外键约束、索引等【未完成】

	return nil
}
