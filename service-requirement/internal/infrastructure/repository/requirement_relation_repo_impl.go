package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"leap-one/service-requirement/internal/domain/entity"
	"leap-one/service-requirement/internal/domain/repository"
)

// requirementRelationRepository 需求关联仓储实�?type requirementRelationRepository struct {
	db *gorm.DB
}

// NewRequirementRelationRepository 创建需求关联仓储实�?func NewRequirementRelationRepository(db *gorm.DB) repository.RequirementRelationRepository {
	return &requirementRelationRepository{db: db}
}

func (r *requirementRelationRepository) Create(relation *entity.RequirementRelation) error {
	return r.db.Create(relation).Error
}

func (r *requirementRelationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.RequirementRelation{}, "id = ?", id).Error
}

func (r *requirementRelationRepository) ListByRequirementID(requirementID uuid.UUID) ([]*entity.RequirementRelation, error) {
	var relations []*entity.RequirementRelation
	err := r.db.Where("requirement_id = ? AND deleted_at IS NULL", requirementID).
		Find(&relations).Error
	return relations, err
}
