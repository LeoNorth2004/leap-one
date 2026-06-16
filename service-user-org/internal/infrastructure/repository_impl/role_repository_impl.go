package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-user-org/internal/domain/entity"
	"leap-one/service-user-org/internal/domain/repository"
	"gorm.io/gorm"
)

// RoleRepositoryImpl 角色仓库实现
type RoleRepositoryImpl struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓库实例
func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &RoleRepositoryImpl{db: db}
}

// Create 创建角色
func (r *RoleRepositoryImpl) Create(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

// GetByID 根据ID获取角色
func (r *RoleRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).First(&role, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// GetByCode 根据编码获取角色（如"admin", "pm", "dev"等预置角色）
func (r *RoleRepositoryImpl) GetByCode(ctx context.Context, code string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// Update 更新角色信息
func (r *RoleRepositoryImpl) Update(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

// Delete 软删除角色
func (r *RoleRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Role{}, "id = ?", id).Error
}

// List 分页查询角色列表
func (r *RoleRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*entity.Role, int64, error) {
	var roles []*entity.Role
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Role{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&roles).Error

	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// GetAll 获取所有角色（用于下拉选择框等场景）
func (r *RoleRepositoryImpl) GetAll(ctx context.Context) ([]*entity.Role, error) {
	var roles []*entity.Role
	err := r.db.WithContext(ctx).Where("status = 1").Order("created_at ASC").Find(&roles).Error
	return roles, err
}

// AssignPermissions 为角色分配权限（先清除旧权限再批量添加新权限）
func (r *RoleRepositoryImpl) AssignPermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除该角色的所有现有权限关联
		if err := tx.Where("role_id = ?", roleID).Delete(&entity.RolePermission{}).Error; err != nil {
			return err
		}

		// 批量插入新的权限关联
		if len(permissionIDs) > 0 {
			rps := make([]entity.RolePermission, len(permissionIDs))
			for i, permID := range permissionIDs {
				rps[i] = entity.RolePermission{
					RoleID:       roleID,
					PermissionID: permID,
				}
			}
			if err := tx.Create(&rps).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetRolePermissions 获取角色的权限列表
func (r *RoleRepositoryImpl) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*entity.Permission, error) {
	var permissions []*entity.Permission
	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}
