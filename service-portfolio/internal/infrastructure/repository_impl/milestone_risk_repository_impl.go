package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-portfolio/internal/domain/entity"
	"leap-one/service-portfolio/internal/domain/repository"
	"gorm.io/gorm"
)

// MilestoneRepositoryImpl 里程碑仓库实现
type MilestoneRepositoryImpl struct {
	db *gorm.DB
}

// NewMilestoneRepository 创建里程碑仓库实例
func NewMilestoneRepository(db *gorm.DB) repository.MilestoneRepository {
	return &MilestoneRepositoryImpl{db: db}
}

func (r *MilestoneRepositoryImpl) Create(ctx context.Context, milestone *entity.Milestone) error {
	return r.db.WithContext(ctx).Create(milestone).Error
}

func (r *MilestoneRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Milestone, error) {
	var m entity.Milestone
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *MilestoneRepositoryImpl) Update(ctx context.Context, milestone *entity.Milestone) error {
	return r.db.WithContext(ctx).Save(milestone).Error
}

func (r *MilestoneRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Milestone{}, "id = ?", id).Error
}

func (r *MilestoneRepositoryImpl) ListByProgramID(ctx context.Context, programID uuid.UUID) ([]*entity.Milestone, error) {
	var milestones []*entity.Milestone
	err := r.db.WithContext(ctx).
		Where("program_id = ?", programID).
		Order("due_date ASC, created_at ASC").
		Find(&milestones).Error
	return milestones, err
}

// RiskRepositoryImpl 风险仓库实现
type RiskRepositoryImpl struct {
	db *gorm.DB
}

// NewRiskRepository 创建风险仓库实例
func NewRiskRepository(db *gorm.DB) repository.RiskRepository {
	return &RiskRepositoryImpl{db: db}
}

func (r *RiskRepositoryImpl) Create(ctx context.Context, risk *entity.Risk) error {
	return r.db.WithContext(ctx).Create(risk).Error
}

func (r *RiskRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Risk, error) {
	var risk entity.Risk
	err := r.db.WithContext(ctx).First(&risk, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &risk, nil
}

func (r *RiskRepositoryImpl) Update(ctx context.Context, risk *entity.Risk) error {
	return r.db.WithContext(ctx).Save(risk).Error
}

func (r *RiskRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Risk{}, "id = ?", id).Error
}

func (r *RiskRepositoryImpl) ListByProgramID(ctx context.Context, programID uuid.UUID) ([]*entity.Risk, error) {
	var risks []*entity.Risk
	err := r.db.WithContext(ctx).
		Where("program_id = ?", programID).
		Order("created_at DESC").
		Find(&risks).Error
	return risks, err
}
