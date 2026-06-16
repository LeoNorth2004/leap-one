package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-bi/internal/domain/entity"
	"leap-one/service-bi/internal/domain/repository"
	"gorm.io/gorm"
)

// ReportTemplateRepositoryImpl 报表模板仓库实现
type ReportTemplateRepositoryImpl struct {
	db *gorm.DB
}

// NewReportTemplateRepository 创建报表模板仓库实例
func NewReportTemplateRepository(db *gorm.DB) repository.ReportTemplateRepository {
	return &ReportTemplateRepositoryImpl{db: db}
}

func (r *ReportTemplateRepositoryImpl) Create(ctx context.Context, template *entity.ReportTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *ReportTemplateRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.ReportTemplate, error) {
	var tpl entity.ReportTemplate
	err := r.db.WithContext(ctx).First(&tpl, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &tpl, nil
}

func (r *ReportTemplateRepositoryImpl) List(ctx context.Context, page, pageSize int, creatorID uuid.UUID, reportType string) ([]*entity.ReportTemplate, int64, error) {
	var list []*entity.ReportTemplate
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.ReportTemplate{})

	if creatorID != uuid.Nil {
		query = query.Where("creator_id = ?", creatorID)
	}
	if reportType != "" {
		query = query.Where("type = ?", reportType)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *ReportTemplateRepositoryImpl) Update(ctx context.Context, template *entity.ReportTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

func (r *ReportTemplateRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ReportTemplate{}, "id = ?", id).Error
}

func (r *ReportTemplateRepositoryImpl) ListByCreator(ctx context.Context, creatorID uuid.UUID) ([]*entity.ReportTemplate, error) {
	var list []*entity.ReportTemplate
	err := r.db.WithContext(ctx).Where("creator_id = ?", creatorID).Order("created_at DESC").Find(&list).Error
	return list, err
}
