package repository

import (
	"context"

	"github.com/user2083251241/ebidsystem/internal/domain/entity"
)

// User 仓库接口
type UserRepository interface {
	// 创建用户
	Create(ctx context.Context, user *entity.User) error

	// 根据用户名查找用户
	FindByUsername(ctx context.Context, username string) (*entity.User, error)

	// 根据ID查找用户
	FindByID(ctx context.Context, id uint) (*entity.User, error)

	// 更新用户信息
	Update(ctx context.Context, user *entity.User) error

	// 删除用户
	Delete(ctx context.Context, id uint) error

	// 获取用户列表
	List(ctx context.Context, offset, limit int) ([]*entity.User, error)

	// 检查用户名是否已存在
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}
