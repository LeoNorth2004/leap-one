package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-task/internal/domain/entity"
	"leap-one/service-task/internal/domain/repository"
	"gorm.io/gorm"
)

// IssueTemplateRepositoryImpl 工单模板仓库实现
type IssueTemplateRepositoryImpl struct {
	db *gorm.DB
}

func NewIssueTemplateRepository(db *gorm.DB) repository.IssueTemplateRepository {
	return &IssueTemplateRepositoryImpl{db: db}
}

func (r *IssueTemplateRepositoryImpl) Create(ctx context.Context, template *entity.IssueTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *IssueTemplateRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.IssueTemplate, error) {
	var tmpl entity.IssueTemplate
	err := r.db.WithContext(ctx).First(&tmpl, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &tmpl, nil
}

func (r *IssueTemplateRepositoryImpl) Update(ctx context.Context, template *entity.IssueTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

func (r *IssueTemplateRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.IssueTemplate{}, "id = ?", id).Error
}

func (r *IssueTemplateRepositoryImpl) List(ctx context.Context, page, pageSize int, tmplType string) ([]*entity.IssueTemplate, int64, error) {
	var templates []*entity.IssueTemplate
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.IssueTemplate{})
	if tmplType != "" {
		query = query.Where("type = ?", tmplType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&templates).Error
	if err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

func (r *IssueTemplateRepositoryImpl) ListByType(ctx context.Context, tmplType string) ([]*entity.IssueTemplate, error) {
	var templates []*entity.IssueTemplate
	err := r.db.WithContext(ctx).Where("type = ?", tmplType).Order("created_at ASC").Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}
