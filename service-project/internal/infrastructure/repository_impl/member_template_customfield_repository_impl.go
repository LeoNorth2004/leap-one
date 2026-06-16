package repository_impl

import (
	"context"

	"leap-one/service-project/internal/domain/entity"
	"leap-one/service-project/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProjectMemberRepositoryImpl 项目成员仓库实现
type ProjectMemberRepositoryImpl struct {
	db *gorm.DB
}

// NewProjectMemberRepository 创建项目成员仓库实例
func NewProjectMemberRepository(db *gorm.DB) repository.ProjectMemberRepository {
	return &ProjectMemberRepositoryImpl{db: db}
}

// Add 添加项目成员
func (r *ProjectMemberRepositoryImpl) Add(ctx context.Context, member *entity.ProjectMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

// Remove 移除项目成员
func (r *ProjectMemberRepositoryImpl) Remove(ctx context.Context, projectID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&entity.ProjectMember{}, "project_id = ? AND user_id = ?", projectID, userID).Error
}

// Get 获取单个成员
func (r *ProjectMemberRepositoryImpl) Get(ctx context.Context, projectID, userID uuid.UUID) (*entity.ProjectMember, error) {
	var member entity.ProjectMember
	err := r.db.WithContext(ctx).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

// ListByProjectID 获取项目的所有成员
func (r *ProjectMemberRepositoryImpl) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entity.ProjectMember, error) {
	var members []*entity.ProjectMember
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("join_time ASC").
		Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

// UpdateRole 更新成员角色
func (r *ProjectMemberRepositoryImpl) UpdateRole(ctx context.Context, projectID, userID uuid.UUID, role string) error {
	return r.db.WithContext(ctx).
		Model(&entity.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Update("role", role).Error
}

// CountByProjectID 统计项目成员�?
func (r *ProjectMemberRepositoryImpl) CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.ProjectMember{}).
		Where("project_id = ?", projectID).
		Count(&count).Error
	return count, err
}

// ==================== 模板仓库实现 ====================

// ProjectTemplateRepositoryImpl 项目模板仓库实现
type ProjectTemplateRepositoryImpl struct {
	db *gorm.DB
}

// NewProjectTemplateRepository 创建模板仓库实例
func NewProjectTemplateRepository(db *gorm.DB) repository.ProjectTemplateRepository {
	return &ProjectTemplateRepositoryImpl{db: db}
}

// Create 创建模板
func (r *ProjectTemplateRepositoryImpl) Create(ctx context.Context, template *entity.ProjectTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

// GetByID 根据ID获取模板
func (r *ProjectTemplateRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.ProjectTemplate, error) {
	var template entity.ProjectTemplate
	err := r.db.WithContext(ctx).First(&template, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}

// List 列出所有模板（分页�?
func (r *ProjectTemplateRepositoryImpl) List(ctx context.Context, page, pageSize int, templateType string) ([]*entity.ProjectTemplate, int64, error) {
	var templates []*entity.ProjectTemplate
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.ProjectTemplate{})

	if templateType != "" {
		query = query.Where("type = ?", templateType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("is_system ASC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&templates).Error
	if err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

// ListSystemTemplates 列出系统预置模板
func (r *ProjectTemplateRepositoryImpl) ListSystemTemplates(ctx context.Context) ([]*entity.ProjectTemplate, error) {
	var templates []*entity.ProjectTemplate
	err := r.db.WithContext(ctx).
		Where("is_system = ?", true).
		Order("created_at ASC").
		Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// Update 更新模板
func (r *ProjectTemplateRepositoryImpl) Update(ctx context.Context, template *entity.ProjectTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

// Delete 删除模板（仅非系统预置）
func (r *ProjectTemplateRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	// 先检查是否为系统模板
	var template entity.ProjectTemplate
	if err := r.db.WithContext(ctx).First(&template, "id = ?", id).Error; err != nil {
		return err
	}
	if template.IsSystem {
		return ErrSystemTemplateCannotDelete
	}
	return r.db.WithContext(ctx).Delete(&entity.ProjectTemplate{}, "id = ?", id).Error
}

// ==================== 自定义字段仓库实现 ====================

// CustomFieldRepositoryImpl 自定义字段仓库实�?
type CustomFieldRepositoryImpl struct {
	db *gorm.DB
}

// NewCustomFieldRepository 创建自定义字段仓库实�?
func NewCustomFieldRepository(db *gorm.DB) repository.CustomFieldRepository {
	return &CustomFieldRepositoryImpl{db: db}
}

// Create 创建自定义字段
func (r *CustomFieldRepositoryImpl) Create(ctx context.Context, field *entity.CustomField) error {
	return r.db.WithContext(ctx).Create(field).Error
}

// GetByID 根据ID获取字段
func (r *CustomFieldRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.CustomField, error) {
	var field entity.CustomField
	err := r.db.WithContext(ctx).First(&field, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &field, nil
}

// ListByProjectID 获取项目的所有自定义字段
func (r *CustomFieldRepositoryImpl) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entity.CustomField, error) {
	var fields []*entity.CustomField
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("sort_order ASC, created_at ASC").
		Find(&fields).Error
	if err != nil {
		return nil, err
	}
	return fields, nil
}

// Update 更新字段
func (r *CustomFieldRepositoryImpl) Update(ctx context.Context, field *entity.CustomField) error {
	return r.db.WithContext(ctx).Save(field).Error
}

// Delete 删除字段
func (r *CustomFieldRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.CustomField{}, "id = ?", id).Error
}
