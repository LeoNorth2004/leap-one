package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"leap-one/service-requirement/internal/domain/entity"
	"leap-one/service-requirement/internal/domain/repository"
)

// requirementReviewRepository 需求评审仓储实�?
type requirementReviewRepository struct {
	db *gorm.DB
}

// NewRequirementReviewRepository 创建需求评审仓储实�?
func NewRequirementReviewRepository(db *gorm.DB) repository.RequirementReviewRepository {
	return &requirementReviewRepository{db: db}
}

func (r *requirementReviewRepository) Create(review *entity.RequirementReview) error {
	return r.db.Create(review).Error
}

func (r *requirementReviewRepository) GetByID(id uuid.UUID) (*entity.RequirementReview, error) {
	var review entity.RequirementReview
	err := r.db.Preload("Participants").
		Where("id = ? AND deleted_at IS NULL", id).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *requirementReviewRepository) ListByRequirementID(requirementID uuid.UUID) ([]*entity.RequirementReview, error) {
	var reviews []*entity.RequirementReview
	err := r.db.Preload("Participants").
		Where("requirement_id = ? AND deleted_at IS NULL", requirementID).
		Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

func (r *requirementReviewRepository) AddParticipant(participant *entity.RequirementReviewParticipant) error {
	return r.db.Create(participant).Error
}
