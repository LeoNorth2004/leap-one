package repository

import (
	"context"

	"leap-one/service-user-org/internal/domain/entity"

	"github.com/google/uuid"
)

// UserGroupRepository 用户组仓库接口定义
type UserGroupRepository interface {
	// Create 创建用户组
	Create(ctx context.Context, group *entity.UserGroup) error

	// GetByID 根据ID获取用户组
	GetByID(ctx context.Context, id uuid.UUID) (*entity.UserGroup, error)

	// Update 更新用户组信息
	Update(ctx context.Context, group *entity.UserGroup) error

	// Delete 删除用户组
	Delete(ctx context.Context, id uuid.UUID) error

	// List 分页查询用户组列表
	List(ctx context.Context, page, pageSize int) ([]*entity.UserGroup, int64, error)

	// AddMember 添加成员到用户组
	AddMember(ctx context.Context, groupID, userID uuid.UUID) error

	// RemoveMember 从用户组移除成员
	RemoveMember(ctx context.Context, groupID, userID uuid.UUID) error

	// GetMembers 获取用户组成员列表
	GetMembers(ctx context.Context, groupID uuid.UUID, page, pageSize int) ([]*entity.User, int64, error)

	// GetByCode 根据编码获取用户组
	GetByCode(ctx context.Context, code string) (*entity.UserGroup, error)

	// BatchAddMembers 批量添加用户组成员
	BatchAddMembers(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) error

	// BatchRemoveMembers 批量移除用户组成员
	BatchRemoveMembers(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) error

	// UpdateMemberCount 更新用户组成员计数
	UpdateMemberCount(ctx context.Context, groupID uuid.UUID) error
}
