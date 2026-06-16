package repository

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-user-org/internal/domain/entity"
)

// RoleRepository 角色仓库接口定义
type RoleRepository interface {
	// Create 创建角色
	Create(ctx context.Context, role *entity.Role) error

	// GetByID 根据ID获取角色
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)

	// GetByCode 根据编码获取角色
	GetByCode(ctx context.Context, code string) (*entity.Role, error)

	// Update 更新角色信息
	Update(ctx context.Context, role *entity.Role) error

	// Delete 软删除角色
	Delete(ctx context.Context, id uuid.UUID) error

	// List 分页查询角色列表
	List(ctx context.Context, page, pageSize int) ([]*entity.Role, int64, error)

	// GetAll 获取所有角色（用于下拉选择）
	GetAll(ctx context.Context) ([]*entity.Role, error)

	// AssignPermissions 为角色分配权限
	AssignPermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error

	// GetRolePermissions 获取角色的权限列表
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*entity.Permission, error)
}
