package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"leap-one/service-requirement/internal/domain/entity"
	"leap-one/service-requirement/internal/domain/repository"
)

// requirementChangeLogRepository 需求变更日志仓储实�?type requirementChangeLogRepository struct {
	db *gorm.DB
}

// NewRequirementChangeLogRepository 创建需求变更日志仓储实�?func NewRequirementChangeLogRepository(db *gorm.DB) repository.RequirementChangeLogRepository {
	return &requirementChangeLogRepository{db: db}
}

func (r *requirementChangeLogRepository) Create(log *entity.RequirementChangeLog) error {
	return r.db.Create(log).Error
}

func (r *requirementChangeLogRepository) ListByRequirementID(requirementID uuid.UUID) ([]*entity.RequirementChangeLog, error) {
	var logs []*entity.RequirementChangeLog
	err := r.db.Where("requirement_id = ? AND deleted_at IS NULL", requirementID).
		Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *requirementChangeLogRepository) UpdateReviewStatus(id uuid.UUID, status string) error {
	return r.db.Model(&entity.RequirementChangeLog{}).Where("id = ?", id).Update("review_status", status).Error
}
