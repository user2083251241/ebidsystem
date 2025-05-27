package mysql

import (
	"context"
	"database/sql"
	"errors"

	"gorm.io/gorm"

	"github.com/user2083251241/ebidsystem/internal/domain/entity"
	"github.com/user2083251241/ebidsystem/internal/domain/repository"
)

type userRepository struct {
	db *gorm.DB
}

// 创建用户仓储实例
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// 创建用户
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	// 检查用户名是否已存在
	exists, err := r.ExistsByUsername(ctx, user.Username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("username already exists")
	}

	// 创建用户
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}

	return nil
}

// 根据用户名查找用户
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ? AND is_deleted = ?", username, false).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &user, nil
}

// 根据ID查找用户
func (r *userRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &user, nil
}

// 更新用户信息
func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	// 检查用户是否存在
	exists, err := r.ExistsByUsername(ctx, user.Username)
	if err != nil {
		return err
	}
	if exists {
		// 检查是否是同一个用户
		existingUser, err := r.FindByUsername(ctx, user.Username)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if existingUser != nil && existingUser.ID != user.ID {
			return errors.New("username already exists")
		}
	}

	// 更新用户信息
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"username":      user.Username,
		"password_hash": user.PasswordHash,
		"role":          user.Role,
	}).Error
}

// 删除用户（软删除）
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// 获取用户列表
func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.WithContext(ctx).
		Where("is_deleted = ?", false).
		Offset(offset).
		Limit(limit).
		Find(&users).Error
	return users, err
}

// 检查用户名是否已存在
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).
		Where("username = ? AND is_deleted = ?", username, false).
		Count(&count).Error
	return count > 0, err
}
