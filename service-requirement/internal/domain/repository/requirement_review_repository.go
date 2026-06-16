package repository

import (
	"github.com/google/uuid"
	"leap-one/service-requirement/internal/domain/entity"
)

// RequirementReviewRepository 需求评审仓储接�?type RequirementReviewRepository interface {
	// Create 创建评审记录
	Create(review *entity.RequirementReview) error
	// GetByID 根据ID获取评审记录
	GetByID(id uuid.UUID) (*entity.RequirementReview, error)
	// ListByRequirementID 根据需求ID查询评审列表
	ListByRequirementID(requirementID uuid.UUID) ([]*entity.RequirementReview, error)
	// AddParticipant 添加评审参与�?	AddParticipant(participant *entity.RequirementReviewParticipant) error
}
