package service

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-requirement/internal/domain/entity"
	"leap-one/service-requirement/internal/domain/repository"
)

// ReviewService 评审应用服务
type ReviewService struct {
	reviewRepo repository.RequirementReviewRepository
	logger     *zap.Logger
}

// NewReviewService 创建评审服务实例
func NewReviewService(reviewRepo repository.RequirementReviewRepository, logger *zap.Logger) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo, logger: logger}
}

// SubmitReview 提交评审
func (s *ReviewService) SubmitReview(requirementID uuid.UUID, title string, meetingDate *time.Time, creatorID uuid.UUID, participants []entity.RequirementReviewParticipant) (*entity.RequirementReview, error) {
	review := &entity.RequirementReview{
		ID:           uuid.New(),
		RequirementID: requirementID,
		Title:        title,
		MeetingDate:  meetingDate,
		Status:       "planning",
		CreatorID:    creatorID,
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, err
	}

	// 添加参与�?	for _, p := range participants {
		p.ID = uuid.New()
		p.ReviewID = review.ID
		if err := s.reviewRepo.AddParticipant(&p); err != nil {
			s.logger.Error("添加评审参与者失�?, zap.Error(err))
		}
	}

	return review, nil
}

// GetReviews 获取评审记录列表
func (s *ReviewService) GetReviews(requirementID uuid.UUID) ([]*entity.RequirementReview, error) {
	return s.reviewRepo.ListByRequirementID(requirementID)
}
