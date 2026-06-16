package repository

import (
	"github.com/google/uuid"
	"leap-one/service-requirement/internal/domain/entity"
)

// RequirementChangeLogRepository 需求变更日志仓储接�?
type RequirementChangeLogRepository interface {
	// Create 创建变更日志
	Create(log *entity.RequirementChangeLog) error
	// ListByRequirementID 根据需求ID查询变更日志
	ListByRequirementID(requirementID uuid.UUID) ([]*entity.RequirementChangeLog, error)
	// UpdateReviewStatus 更新审核状�?
	UpdateReviewStatus(id uuid.UUID, status string) error
}
