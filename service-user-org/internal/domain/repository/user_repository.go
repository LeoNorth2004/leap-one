package repository

import (
	"context"

	"leap-one/service-user-org/internal/domain/entity"

	"github.com/google/uuid"
)

// UserRepository 用户仓库接口定义
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *entity.User) error

	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// GetByUsername 根据用户名获取用户（用于登录）
	GetByUsername(ctx context.Context, username string) (*entity.User, error)

	// GetByEmail 根据邮箱获取用户
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// Update 更新用户信息
	Update(ctx context.Context, user *entity.User) error

	// Delete 软删除用户
	Delete(ctx context.Context, id uuid.UUID) error

	// List 分页查询用户列表
	List(ctx context.Context, page, pageSize int, keyword string) ([]*entity.User, int64, error)

	// UpdatePassword 更新用户密码
	UpdatePassword(ctx context.Context, id uuid.UUID, hashedPassword string) error

	// UpdateLastLogin 更新最后登录时间和IP
	UpdateLastLogin(ctx context.Context, id uuid.UUID, ip string) error
}
