package service

import (
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-requirement/internal/domain/entity"
	"leap-one/service-requirement/internal/domain/repository"
)

// RelationService 关联关系应用服务
type RelationService struct {
	relationRepo repository.RequirementRelationRepository
	logger       *zap.Logger
}

// NewRelationService 创建关联关系服务实例
func NewRelationService(relationRepo repository.RequirementRelationRepository, logger *zap.Logger) *RelationService {
	return &RelationService{relationRepo: relationRepo, logger: logger}
}

// AddRelation 添加关联关系
func (s *RelationService) AddRelation(requirementID uuid.UUID, relatedType string, relatedID uuid.UUID, relationType string) (*entity.RequirementRelation, error) {
	relation := &entity.RequirementRelation{
		ID:            uuid.New(),
		RequirementID: requirementID,
		RelatedType:   relatedType,
		RelatedID:     relatedID,
		RelationType:  relationType,
	}
	if err := s.relationRepo.Create(relation); err != nil {
		return nil, err
	}
	return relation, nil
}

// GetRelations 获取关联列表
func (s *RelationService) GetRelations(requirementID uuid.UUID) ([]*entity.RequirementRelation, error) {
	return s.relationRepo.ListByRequirementID(requirementID)
}

// RemoveRelation 移除关联关系
func (s *RelationService) RemoveRelation(id uuid.UUID) error {
	return s.relationRepo.Delete(id)
}
