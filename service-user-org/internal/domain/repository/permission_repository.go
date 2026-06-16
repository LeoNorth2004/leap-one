package repository

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-user-org/internal/domain/entity"
)

// PermissionRepository 权限仓库接口定义
type PermissionRepository interface {
	// Create 创建权限
	Create(ctx context.Context, perm *entity.Permission) error

	// GetByID 根据ID获取权限
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error)

	// GetByCode 根据编码获取权限
	GetByCode(ctx context.Context, code string) (*entity.Permission, error)

	// List 分页查询权限列表
	List(ctx context.Context, page, pageSize int, resource string) ([]*entity.Permission, int64, error)

	// GetAll 获取所有权限
	GetAll(ctx context.Context) ([]*entity.Permission, error)

	// GetUserPermissions 获取用户的所有权限编码列表
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error)

	// CheckPermission 检查用户是否拥有指定权限
	CheckPermission(ctx context.Context, userID uuid.UUID, permissionCode string) (bool, error)

	// GetRoleUsers 获取拥有某角色的用户列表（分页）
	GetRoleUsers(ctx context.Context, roleID uuid.UUID, page, pageSize int) ([]*entity.User, int64, error)

	// SearchUsers 全局搜索用户
	SearchUsers(ctx context.Context, keyword string, limit int) ([]*entity.User, error)
}
