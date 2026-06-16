package service

import (
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-requirement/internal/domain/entity"
	"leap-one/service-requirement/internal/domain/repository"
)

// ChangeLogService 变更日志应用服务
type ChangeLogService struct {
	changeLogRepo repository.RequirementChangeLogRepository
	logger        *zap.Logger
}

// NewChangeLogService 创建变更日志服务实例
func NewChangeLogService(changeLogRepo repository.RequirementChangeLogRepository, logger *zap.Logger) *ChangeLogService {
	return &ChangeLogService{changeLogRepo: changeLogRepo, logger: logger}
}

// CreateChangeLog 发起变更记录
func (s *ChangeLogService) CreateChangeLog(log *entity.RequirementChangeLog) error {
	log.ID = uuid.New()
	log.ReviewStatus = "pending"
	return s.changeLogRepo.Create(log)
}

// GetChangeLogs 获取变更日志列表
func (s *ChangeLogService) GetChangeLogs(requirementID uuid.UUID) ([]*entity.RequirementChangeLog, error) {
	return s.changeLogRepo.ListByRequirementID(requirementID)
}

// ApproveChange 审批通过变更
func (s *ChangeLogService) ApproveChange(id uuid.UUID) error {
	return s.changeLogRepo.UpdateReviewStatus(id, "approved")
}

// RejectChange 审批拒绝变更
func (s *ChangeLogService) RejectChange(id uuid.UUID) error {
	return s.changeLogRepo.UpdateReviewStatus(id, "rejected")
}
