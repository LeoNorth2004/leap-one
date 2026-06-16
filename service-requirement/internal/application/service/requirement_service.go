package service

import (
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-requirement/internal/domain/entity"
	"leap-one/service-requirement/internal/domain/repository"
)

// RequirementService 需求应用服�?
type RequirementService struct {
	reqRepo       repository.RequirementRepository
	relationRepo  repository.RequirementRelationRepository
	changeLogRepo repository.RequirementChangeLogRepository
	logger        *zap.Logger
}

// NewRequirementService 创建需求服务实�?
func NewRequirementService(
	reqRepo repository.RequirementRepository,
	relationRepo repository.RequirementRelationRepository,
	changeLogRepo repository.RequirementChangeLogRepository,
	logger *zap.Logger,
) *RequirementService {
	return &RequirementService{
		reqRepo:       reqRepo,
		relationRepo:  relationRepo,
		changeLogRepo: changeLogRepo,
		logger:        logger,
	}
}

// CreateRequirement 创建需�?
func (s *RequirementService) CreateRequirement(req *entity.Requirement) (*entity.Requirement, error) {
	// 自动生成需求编�?
	code, err := s.reqRepo.GenerateCode()
	if err != nil {
		return nil, err
	}
	req.Code = code

	// 设置默认�?
	if req.Type == "" {
		req.Type = "story"
	}
	if req.Level == 0 {
		req.Level = 3
	}
	if req.Status == "" {
		req.Status = "draft"
	}
	if req.Priority == 0 {
		req.Priority = 3
	}
	if req.Source == "" {
		req.Source = "manual"
	}
	if req.Stage == "" {
		req.Stage = "requirement"
	}

	// 根据父需求设置层�?
	if req.ParentID != nil {
		parent, err := s.reqRepo.GetByID(*req.ParentID)
		if err != nil {
			return nil, err
		}
		req.Level = parent.Level + 1
	}

	if err := s.reqRepo.Create(req); err != nil {
		s.logger.Error("创建需求失�?, zap.Error(err), zap.String("title", req.Title))
		return nil, err
	}

	// 记录变更日志
	_ = s.recordChangeLog(req.ID, "create", "", "", "创建需�?, uuid.Nil)
	return req, nil
}

// GetRequirement 获取需求详�?
func (s *RequirementService) GetRequirement(id uuid.UUID) (*entity.Requirement, error) {
	return s.reqRepo.GetByID(id)
}

// UpdateRequirement 更新需�?
func (s *RequirementService) UpdateRequirement(id uuid.UUID, updates map[string]interface{}) (*entity.Requirement, error) {
	req, err := s.reqRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 应用更新字段（此处简化处理，实际应逐字段比较）
	for field, value := range updates {
		switch field {
		case "title":
			if v, ok := value.(string); ok {
				req.Title = v
			}
		case "description":
			if v, ok := value.(string); ok {
				req.Description = v
			}
		case "type":
			if v, ok := value.(string); ok {
				req.Type = v
			}
		case "priority":
			if v, ok := value.(int); ok {
				req.Priority = v
			}
		case "status":
			if v, ok := value.(string); ok {
				req.Status = v
			}
		case "stage":
			if v, ok := value.(string); ok {
				req.Stage = v
			}
		case "owner_id":
			if v, ok := value.(uuid.UUID); ok {
				req.OwnerID = &v
			}
		case "release_version":
			if v, ok := value.(string); ok {
				req.ReleaseVersion = v
			}
		}
	}

	if err := s.reqRepo.Update(req); err != nil {
		return nil, err
	}
	return req, nil
}

// DeleteRequirement 删除需�?
func (s *RequirementService) DeleteRequirement(id uuid.UUID) error {
	return s.reqRepo.Delete(id)
}

// ListRequirements 查询需求列�?
func (s *RequirementService) ListRequirements(params *repository.RequirementListParams) ([]*entity.Requirement, int64, error) {
	return s.reqRepo.List(params)
}

// GetRequirementTree 获取需求树
func (s *RequirementService) GetRequirementTree(productID uuid.UUID) ([]*entity.Requirement, error) {
	return s.reqRepo.GetTree(productID)
}

// UpdateStatus 更新需求状�?
func (s *RequirementService) UpdateStatus(id uuid.UUID, status string) error {
	return s.reqRepo.UpdateStatus(id, status)
}

// recordChangeLog 记录变更日志
func (s *RequirementService) recordChangeLog(requirementID uuid.UUID, changeType, fieldName, oldValue, newValue string, changeUserID uuid.UUID) error {
	log := &entity.RequirementChangeLog{
		ID:            uuid.New(),
		RequirementID: requirementID,
		ChangeType:    changeType,
		FieldName:     fieldName,
		OldValue:      oldValue,
		NewValue:      newValue,
		ChangeUserID:  changeUserID,
	}
	return s.changeLogRepo.Create(log)
}
