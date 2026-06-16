package repository_impl

import (
	"context"

	"leap-one/service-task/internal/domain/entity"
	"leap-one/service-task/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IssueWorkflowRepositoryImpl 工作流仓库实现
type IssueWorkflowRepositoryImpl struct {
	db *gorm.DB
}

func NewIssueWorkflowRepository(db *gorm.DB) repository.IssueWorkflowRepository {
	return &IssueWorkflowRepositoryImpl{db: db}
}

func (r *IssueWorkflowRepositoryImpl) Create(ctx context.Context, workflow *entity.IssueWorkflow) error {
	return r.db.WithContext(ctx).Create(workflow).Error
}

func (r *IssueWorkflowRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.IssueWorkflow, error) {
	var wf entity.IssueWorkflow
	err := r.db.WithContext(ctx).Preload("Transitions", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).First(&wf, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &wf, nil
}

func (r *IssueWorkflowRepositoryImpl) Update(ctx context.Context, workflow *entity.IssueWorkflow) error {
	return r.db.WithContext(ctx).Save(workflow).Error
}

func (r *IssueWorkflowRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.IssueWorkflow{}, "id = ?", id).Error
}

func (r *IssueWorkflowRepositoryImpl) List(ctx context.Context, page, pageSize int, wfType string) ([]*entity.IssueWorkflow, int64, error) {
	var workflows []*entity.IssueWorkflow
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.IssueWorkflow{})
	if wfType != "" {
		query = query.Where("type = ?", wfType)
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
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&workflows).Error
	if err != nil {
		return nil, 0, err
	}

	return workflows, total, nil
}

func (r *IssueWorkflowRepositoryImpl) AddTransition(ctx context.Context, transition *entity.IssueWorkflowTransition) error {
	return r.db.WithContext(ctx).Create(transition).Error
}

func (r *IssueWorkflowRepositoryImpl) RemoveTransition(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.IssueWorkflowTransition{}, "id = ?", id).Error
}

func (r *IssueWorkflowRepositoryImpl) ListTransitions(ctx context.Context, workflowID uuid.UUID) ([]*entity.IssueWorkflowTransition, error) {
	var transitions []*entity.IssueWorkflowTransition
	err := r.db.WithContext(ctx).Where("workflow_id = ?", workflowID).
		Order("sort_order ASC").Find(&transitions).Error
	if err != nil {
		return nil, err
	}
	return transitions, nil
}

func (r *IssueWorkflowRepositoryImpl) GetByTypeAndStatus(ctx context.Context, wfType, initialStatus string) (*entity.IssueWorkflow, error) {
	var wf entity.IssueWorkflow
	err := r.db.WithContext(ctx).Preload("Transitions", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Where("type = ? AND initial_status = ?", wfType, initialStatus).First(&wf).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &wf, nil
}

// IssueSLAConfigRepositoryImpl SLA配置仓库实现
type IssueSLAConfigRepositoryImpl struct {
	db *gorm.DB
}

func NewIssueSLAConfigRepository(db *gorm.DB) repository.IssueSLAConfigRepository {
	return &IssueSLAConfigRepositoryImpl{db: db}
}

func (r *IssueSLAConfigRepositoryImpl) Create(ctx context.Context, config *entity.IssueSLAConfig) error {
	return r.db.WithContext(ctx).Create(config).Error
}

func (r *IssueSLAConfigRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.IssueSLAConfig, error) {
	var cfg entity.IssueSLAConfig
	err := r.db.WithContext(ctx).First(&cfg, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cfg, nil
}

func (r *IssueSLAConfigRepositoryImpl) Update(ctx context.Context, config *entity.IssueSLAConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}

func (r *IssueSLAConfigRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.IssueSLAConfig{}, "id = ?", id).Error
}

func (r *IssueSLAConfigRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*entity.IssueSLAConfig, int64, error) {
	var configs []*entity.IssueSLAConfig
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.IssueSLAConfig{})

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
	err := query.Order("type ASC, priority ASC").Offset(offset).Limit(pageSize).Find(&configs).Error
	if err != nil {
		return nil, 0, err
	}

	return configs, total, nil
}

func (r *IssueSLAConfigRepositoryImpl) GetByTypeAndPriority(ctx context.Context, slaType string, priority int) (*entity.IssueSLAConfig, error) {
	var cfg entity.IssueSLAConfig
	err := r.db.WithContext(ctx).Where("type = ? AND priority = ?", slaType, priority).First(&cfg).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cfg, nil
}
