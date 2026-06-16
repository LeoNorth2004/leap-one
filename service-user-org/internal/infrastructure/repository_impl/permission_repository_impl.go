package repository_impl

import (
	"context"

	"leap-one/service-user-org/internal/domain/entity"
	"leap-one/service-user-org/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PermissionRepositoryImpl 权限仓库实现
type PermissionRepositoryImpl struct {
	db *gorm.DB
}

// NewPermissionRepository 创建权限仓库实例
func NewPermissionRepository(db *gorm.DB) repository.PermissionRepository {
	return &PermissionRepositoryImpl{db: db}
}

// Create 创建权限
func (r *PermissionRepositoryImpl) Create(ctx context.Context, perm *entity.Permission) error {
	return r.db.WithContext(ctx).Create(perm).Error
}

// GetByID 根据ID获取权限
func (r *PermissionRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error) {
	var perm entity.Permission
	err := r.db.WithContext(ctx).First(&perm, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &perm, nil
}

// GetByCode 根据编码获取权限
func (r *PermissionRepositoryImpl) GetByCode(ctx context.Context, code string) (*entity.Permission, error) {
	var perm entity.Permission
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&perm).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &perm, nil
}

// List 分页查询权限列表（可按资源类型筛选）
func (r *PermissionRepositoryImpl) List(ctx context.Context, page, pageSize int, resource string) ([]*entity.Permission, int64, error) {
	var perms []*entity.Permission
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Permission{})

	if resource != "" {
		query = query.Where("resource = ?", resource)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("resource ASC, action ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&perms).Error

	if err != nil {
		return nil, 0, err
	}

	return perms, total, nil
}

// GetAll 获取所有权限
func (r *PermissionRepositoryImpl) GetAll(ctx context.Context) ([]*entity.Permission, error) {
	var perms []*entity.Permission
	err := r.db.WithContext(ctx).Order("resource ASC, action ASC").Find(&perms).Error
	return perms, err
}

// GetUserPermissions 获取用户的所有权限编码列表
// 通过 用户->角色->权限 的多级关联查询实现
func (r *PermissionRepositoryImpl) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	var codes []string
	err := r.db.WithContext(ctx).
		Distinct().
		Model(&entity.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Pluck("permissions.code", &codes).Error
	return codes, err
}

// CheckPermission 检查用户是否拥有指定权限编码
func (r *PermissionRepositoryImpl) CheckPermission(ctx context.Context, userID uuid.UUID, permissionCode string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.code = ?", userID, permissionCode).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SyncRolePermissions 同步角色权限（全量替换）
func (r *PermissionRepositoryImpl) SyncRolePermissions(ctx context.Context, roleID uuid.UUID, permissionCodes []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 查询权限编码对应的ID
		var permIDs []uuid.UUID
		if len(permissionCodes) > 0 {
			if err := tx.Model(&entity.Permission{}).
				Where("code IN ?", permissionCodes).
				Pluck("id", &permIDs).Error; err != nil {
				return err
			}
		}

		// 清除旧关联
		if err := tx.Where("role_id = ?", roleID).Delete(&entity.RolePermission{}).Error; err != nil {
			return err
		}

		// 插入新关联
		if len(permIDs) > 0 {
			rps := make([]entity.RolePermission, len(permIDs))
			for i, pid := range permIDs {
				rps[i] = entity.RolePermission{
					RoleID:       roleID,
					PermissionID: pid,
				}
			}
			if err := tx.Create(&rps).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetUserRoles 获取用户的角色编码列表
func (r *PermissionRepositoryImpl) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error) {
	var codes []string
	err := r.db.WithContext(ctx).
		Distinct().
		Model(&entity.Role{}).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.status = 1 AND users.deleted_at IS NULL", userID).
		Pluck("roles.code", &codes).Error
	return codes, err
}

// AssignRole 分配角色给用户
func (r *PermissionRepositoryImpl) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	userRole := entity.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.WithContext(ctx).
		Create(&userRole).
		Error
}

// RevokeRole 撤销用户的某个角色
func (r *PermissionRepositoryImpl) RevokeRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&entity.UserRole{}).Error
}

// GetRoleUsers 获取拥有某角色的用户ID列表
func (r *PermissionRepositoryImpl) GetRoleUsers(ctx context.Context, roleID uuid.UUID, page, pageSize int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	baseQuery := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID)

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := baseQuery.
		Order("users.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// SearchUsers 全局搜索用户（供其他服务调用）
func (r *PermissionRepositoryImpl) SearchUsers(ctx context.Context, keyword string, limit int) ([]*entity.User, error) {
	var users []*entity.User

	query := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("status = 1")

	if keyword != "" {
		pattern := "%" + keyword + "%"
		query = query.Where(
			"username LIKE ? OR real_name LIKE ? OR email LIKE ?",
			pattern, pattern, pattern,
		)
	}

	if limit <= 0 || limit > 50 {
		limit = 20
	}

	err := query.Limit(limit).Find(&users).Error
	return users, err
}
