package repository

import (
	"github.com/google/uuid"
	"leap-one/service-requirement/internal/domain/entity"
)

// RequirementRelationRepository 需求关联仓储接�?
type RequirementRelationRepository interface {
	// Create 创建关联关系
	Create(relation *entity.RequirementRelation) error
	// Delete 删除关联关系
	Delete(id uuid.UUID) error
	// ListByRequirementID 根据需求ID查询关联列表
	ListByRequirementID(requirementID uuid.UUID) ([]*entity.RequirementRelation, error)
}
