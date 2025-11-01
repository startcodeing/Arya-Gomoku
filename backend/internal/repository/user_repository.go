package repository

import (
	"context"
	"fmt"
	"time"

	"gomoku-backend/internal/model"

	"gorm.io/gorm"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Delete(ctx context.Context, id string) error

	// 查询操作
	List(ctx context.Context, offset, limit int) ([]*model.User, int64, error)
	Search(ctx context.Context, keyword string, offset, limit int) ([]*model.User, int64, error)
	GetActiveUsers(ctx context.Context, offset, limit int) ([]*model.User, int64, error)
	
	// 统计操作
	Count(ctx context.Context) (int64, error)
	CountActiveUsers(ctx context.Context) (int64, error)
	CountByRole(ctx context.Context, role string) (int64, error)
	
	// 业务操作
	UpdateLastActive(ctx context.Context, userID string) error
	SetUserStatus(ctx context.Context, userID string, isActive bool) error
	GetUsersByRole(ctx context.Context, role string, offset, limit int) ([]*model.User, int64, error)
}

// userRepository 用户仓储实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}
	return nil
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}
	return nil
}

// Delete 删除用户（软删除）
func (r *userRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}
	return nil
}

// List 获取用户列表
func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取用户总数失败: %w", err)
	}

	// 获取用户列表
	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("获取用户列表失败: %w", err)
	}

	return users, total, nil
}

// Search 搜索用户
func (r *userRepository) Search(ctx context.Context, keyword string, offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{}).
		Where("username ILIKE ? OR nickname ILIKE ? OR email ILIKE ?", 
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("搜索用户总数失败: %w", err)
	}

	// 获取用户列表
	if err := query.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("搜索用户失败: %w", err)
	}

	return users, total, nil
}

// GetActiveUsers 获取活跃用户列表
func (r *userRepository) GetActiveUsers(ctx context.Context, offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{}).Where("is_active = ?", true)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取活跃用户总数失败: %w", err)
	}

	// 获取用户列表
	if err := query.
		Offset(offset).
		Limit(limit).
		Order("last_active_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("获取活跃用户列表失败: %w", err)
	}

	return users, total, nil
}

// Count 获取用户总数
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("获取用户总数失败: %w", err)
	}
	return count, nil
}

// CountActiveUsers 获取活跃用户总数
func (r *userRepository) CountActiveUsers(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("is_active = ?", true).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("获取活跃用户总数失败: %w", err)
	}
	return count, nil
}

// CountByRole 根据角色获取用户数量
func (r *userRepository) CountByRole(ctx context.Context, role string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("role = ?", role).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("获取角色用户数量失败: %w", err)
	}
	return count, nil
}

// UpdateLastActive 更新最后活跃时间
func (r *userRepository) UpdateLastActive(ctx context.Context, userID string) error {
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Update("last_active_at", time.Now()).Error; err != nil {
		return fmt.Errorf("更新最后活跃时间失败: %w", err)
	}
	return nil
}

// SetUserStatus 设置用户状态
func (r *userRepository) SetUserStatus(ctx context.Context, userID string, isActive bool) error {
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"is_active":  isActive,
			"updated_at": time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("设置用户状态失败: %w", err)
	}
	return nil
}

// GetUsersByRole 根据角色获取用户列表
func (r *userRepository) GetUsersByRole(ctx context.Context, role string, offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{}).Where("role = ?", role)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取角色用户总数失败: %w", err)
	}

	// 获取用户列表
	if err := query.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("获取角色用户列表失败: %w", err)
	}

	return users, total, nil
}