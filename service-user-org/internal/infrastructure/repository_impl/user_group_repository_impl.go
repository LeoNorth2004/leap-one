package repository_impl

import (
	"context"

	"leap-one/service-user-org/internal/domain/entity"
	"leap-one/service-user-org/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserGroupRepositoryImpl 用户组仓库实现
type UserGroupRepositoryImpl struct {
	db *gorm.DB
}

// NewUserGroupRepository 创建用户组仓库实例
func NewUserGroupRepository(db *gorm.DB) repository.UserGroupRepository {
	return &UserGroupRepositoryImpl{db: db}
}

// Create 创建用户组
func (r *UserGroupRepositoryImpl) Create(ctx context.Context, group *entity.UserGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

// GetByID 根据ID获取用户组
func (r *UserGroupRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.UserGroup, error) {
	var group entity.UserGroup
	err := r.db.WithContext(ctx).First(&group, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}

// Update 更新用户组信息
func (r *UserGroupRepositoryImpl) Update(ctx context.Context, group *entity.UserGroup) error {
	return r.db.WithContext(ctx).Save(group).Error
}

// Delete 删除用户组（同时清理成员关联数据）
func (r *UserGroupRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_group_id = ?", id).Delete(&entity.UserGroupMember{}).Error; err != nil {
			return err
		}
		return tx.Delete(&entity.UserGroup{}, "id = ?", id).Error
	})
}

// List 分页查询用户组列表
func (r *UserGroupRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*entity.UserGroup, int64, error) {
	var groups []*entity.UserGroup
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.UserGroup{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&groups).Error

	if err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

// AddMember 添加成员到用户组（幂等操作）
func (r *UserGroupRepositoryImpl) AddMember(ctx context.Context, groupID, userID uuid.UUID) error {
	member := entity.UserGroupMember{
		UserGroupID: groupID,
		UserID:      userID,
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&member).Error
}

// RemoveMember 从用户组移除成员
func (r *UserGroupRepositoryImpl) RemoveMember(ctx context.Context, groupID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_group_id = ? AND user_id = ?", groupID, userID).
		Delete(&entity.UserGroupMember{}).Error
}

// GetMembers 获取用户组成员列表（分页）
func (r *UserGroupRepositoryImpl) GetMembers(ctx context.Context, groupID uuid.UUID, page, pageSize int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	baseQuery := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Joins("JOIN user_group_members ON user_group_members.user_id = users.id").
		Where("user_group_members.user_group_id = ?", groupID)

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := baseQuery.
		Order("user_group_members.joined_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// BatchAddMembers 批量添加成员到用户组
func (r *UserGroupRepositoryImpl) BatchAddMembers(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) error {
	members := make([]entity.UserGroupMember, len(userIDs))
	for i, uid := range userIDs {
		members[i] = entity.UserGroupMember{
			UserGroupID: groupID,
			UserID:      uid,
		}
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&members).Error
}

// BatchRemoveMembers 批量移除用户组成员
func (r *UserGroupRepositoryImpl) BatchRemoveMembers(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_group_id = ? AND user_id IN ?", groupID, userIDs).
		Delete(&entity.UserGroupMember{}).Error
}

// GetUserGroups 获取用户所属的所有用户组
func (r *UserGroupRepositoryImpl) GetUserGroups(ctx context.Context, userID uuid.UUID) ([]*entity.UserGroup, error) {
	var groups []*entity.UserGroup
	err := r.db.WithContext(ctx).
		Distinct().
		Model(&entity.UserGroup{}).
		Joins("JOIN user_group_members ON user_group_members.user_group_id = user_groups.id").
		Where("user_group_members.user_id = ?", userID).
		Order("user_groups.name ASC").
		Find(&groups).Error
	return groups, err
}

// UpdateMemberCount 更新用户组的成员计数
func (r *UserGroupRepositoryImpl) UpdateMemberCount(ctx context.Context, groupID uuid.UUID) error {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.UserGroupMember{}).
		Where("user_group_id = ?", groupID).
		Count(&count).Error
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).
		Model(&entity.UserGroup{}).
		Where("id = ?", groupID).
		Update("member_count", count).Error
}

// GetByCode 根据编码获取用户组
func (r *UserGroupRepositoryImpl) GetByCode(ctx context.Context, code string) (*entity.UserGroup, error) {
	var group entity.UserGroup
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}

// BatchAddMembersRetry 批量添加用户组成员（幂等操作，带重试）
func (r *UserGroupRepositoryImpl) BatchAddMembersRetry(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) error {
	members := make([]entity.UserGroupMember, len(userIDs))
	for i, uid := range userIDs {
		members[i] = entity.UserGroupMember{
			UserGroupID: groupID,
			UserID:      uid,
		}
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&members).Error
}
