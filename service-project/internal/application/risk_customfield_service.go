package application

import (
	"context"
	"errors"

	"leap-one/service-project/internal/domain/entity"
	"leap-one/service-project/internal/domain/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ==================== 风险服务 ====================

// 风险服务相关错误定义
var (
	ErrRiskNotFound = errors.New("风险记录不存在")
)

// RiskService 风险管理应用服务
type RiskService struct {
	riskRepo    repository.ProjectRiskRepository
	projectRepo repository.ProjectRepository
	logger      *zap.Logger
}

// NewRiskService 创建风险管理服务实例
func NewRiskService(
	riskRepo repository.ProjectRiskRepository,
	projectRepo repository.ProjectRepository,
	logger *zap.Logger,
) *RiskService {
	return &RiskService{
		riskRepo:    riskRepo,
		projectRepo: projectRepo,
		logger:      logger,
	}
}

// CreateRisk 创建风险
func (s *RiskService) CreateRisk(ctx context.Context, projectID uuid.UUID, req *CreateRiskInput) (*entity.ProjectRisk, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil || project == nil {
		return nil, ErrProjectNotFound
	}

	// 设置默认概率和影响等级
	probability := 3
	if req.Probability > 0 && req.Probability <= 5 {
		probability = req.Probability
	}
	impact := 3
	if req.Impact > 0 && req.Impact <= 5 {
		impact = req.Impact
	}

	risk := &entity.ProjectRisk{
		ProjectID:   projectID,
		Title:       req.Title,
		Description: req.Description,
		Probability: probability,
		Impact:      impact,
		Severity:    probability * impact, // 自动计算严重程度
		Status:      "open",
		OwnerID:     req.OwnerID,
		Mitigation:  req.Mitigation,
	}

	if err := s.riskRepo.Create(ctx, risk); err != nil {
		s.logger.Error("创建风险失败", zap.Error(err))
		return nil, errors.New("创建风险失败")
	}

	s.logger.Info("风险已创建",
		zap.String("risk_id", risk.ID.String()),
		zap.String("title", req.Title),
		zap.Int("severity", risk.Severity),
	)
	return risk, nil
}

// UpdateRisk 更新风险
func (s *RiskService) UpdateRisk(ctx context.Context, id uuid.UUID, req *UpdateRiskInput) (*entity.ProjectRisk, error) {
	risk, err := s.riskRepo.GetByID(ctx, id)
	if err != nil || risk == nil {
		return nil, ErrRiskNotFound
	}

	if req.Title != nil {
		risk.Title = *req.Title
	}
	if req.Description != nil {
		risk.Description = *req.Description
	}
	if req.Probability != nil {
		risk.Probability = *req.Probability
	}
	if req.Impact != nil {
		risk.Impact = *req.Impact
	}
	if req.OwnerID != nil {
		risk.OwnerID = *req.OwnerID
	}
	if req.Mitigation != nil {
		risk.Mitigation = *req.Mitigation
	}
	if req.Status != nil {
		risk.Status = *req.Status
	}

	// 重新计算严重程度
	risk.Severity = risk.Probability * risk.Impact

	if err := s.riskRepo.Update(ctx, risk); err != nil {
		s.logger.Error("更新风险失败", zap.Error(err))
		return nil, errors.New("更新风险失败")
	}

	return risk, nil
}

// DeleteRisk 删除风险
func (s *RiskService) DeleteRisk(ctx context.Context, id uuid.UUID) error {
	risk, err := s.riskRepo.GetByID(ctx, id)
	if err != nil || risk == nil {
		return ErrRiskNotFound
	}

	if err := s.riskRepo.Delete(ctx, id); err != nil {
		s.logger.Error("删除风险失败", zap.Error(err))
		return errors.New("删除风险失败")
	}

	s.logger.Info("风险已删除", zap.String("risk_id", id.String()))
	return nil
}

// ListRisks 获取项目的风险列表
func (s *RiskService) ListRisks(ctx context.Context, projectID uuid.UUID) ([]*entity.ProjectRisk, error) {
	project, _ := s.projectRepo.GetByID(ctx, projectID)
	if project == nil {
		return nil, ErrProjectNotFound
	}

	return s.riskRepo.ListByProjectID(ctx, projectID)
}

// CreateRiskInput 创建风险输入
type CreateRiskInput struct {
	Title       string
	Description string
	Probability int
	Impact      int
	OwnerID     uuid.UUID
	Mitigation  string
}

// UpdateRiskInput 更新风险输入
type UpdateRiskInput struct {
	Title       *string
	Description *string
	Probability *int
	Impact      *int
	OwnerID     *uuid.UUID
	Mitigation  *string
	Status      *string
}

// ==================== 自定义字段服务 ====================

// CustomFieldService 自定义字段应用服务
type CustomFieldService struct {
	fieldRepo   repository.CustomFieldRepository
	projectRepo repository.ProjectRepository
	logger      *zap.Logger
}

// NewCustomFieldService 创建自定义字段服务实例
func NewCustomFieldService(
	fieldRepo repository.CustomFieldRepository,
	projectRepo repository.ProjectRepository,
	logger *zap.Logger,
) *CustomFieldService {
	return &CustomFieldService{
		fieldRepo:   fieldRepo,
		projectRepo: projectRepo,
		logger:      logger,
	}
}

// AddCustomField 添加自定义字段
func (s *CustomFieldService) AddCustomField(ctx context.Context, projectID uuid.UUID, req *AddCustomFieldInput) (*entity.CustomField, error) {
	project, _ := s.projectRepo.GetByID(ctx, projectID)
	if project == nil {
		return nil, ErrProjectNotFound
	}

	field := &entity.CustomField{
		ProjectID: projectID,
		Name:      req.Name,
		FieldKey:  req.FieldKey,
		FieldType: req.FieldType,
		Options:   req.Options,
		Required:  req.Required,
		SortOrder: req.SortOrder,
	}

	if err := s.fieldRepo.Create(ctx, field); err != nil {
		s.logger.Error("创建自定义字段失败", zap.Error(err))
		return nil, errors.New("创建自定义字段失败")
	}

	return field, nil
}

// UpdateCustomField 更新自定义字段
func (s *CustomFieldService) UpdateCustomField(ctx context.Context, id uuid.UUID, req *UpdateCustomFieldInput) (*entity.CustomField, error) {
	field, err := s.fieldRepo.GetByID(ctx, id)
	if err != nil || field == nil {
		return nil, errors.New("自定义字段不存在")
	}

	if req.Name != nil {
		field.Name = *req.Name
	}
	if req.FieldKey != nil {
		field.FieldKey = *req.FieldKey
	}
	if req.FieldType != nil {
		field.FieldType = *req.FieldType
	}
	if req.Options != nil {
		field.Options = *req.Options
	}
	if req.Required != nil {
		field.Required = *req.Required
	}
	if req.SortOrder != nil {
		field.SortOrder = *req.SortOrder
	}

	if err := s.fieldRepo.Update(ctx, field); err != nil {
		s.logger.Error("更新自定义字段失败", zap.Error(err))
		return nil, errors.New("更新自定义字段失败")
	}

	return field, nil
}

// DeleteCustomField 删除自定义字段
func (s *CustomFieldService) DeleteCustomField(ctx context.Context, id uuid.UUID) error {
	field, err := s.fieldRepo.GetByID(ctx, id)
	if err != nil || field == nil {
		return errors.New("自定义字段不存在")
	}

	if err := s.fieldRepo.Delete(ctx, id); err != nil {
		s.logger.Error("删除自定义字段失败", zap.Error(err))
		return errors.New("删除自定义字段失败")
	}

	return nil
}

// ListCustomFields 获取项目的自定义字段列表
func (s *CustomFieldService) ListCustomFields(ctx context.Context, projectID uuid.UUID) ([]*entity.CustomField, error) {
	project, _ := s.projectRepo.GetByID(ctx, projectID)
	if project == nil {
		return nil, ErrProjectNotFound
	}

	return s.fieldRepo.ListByProjectID(ctx, projectID)
}

// AddCustomFieldInput 添加自定义字段输入
type AddCustomFieldInput struct {
	Name      string
	FieldKey  string
	FieldType string
	Options   string
	Required  bool
	SortOrder int
}

// UpdateCustomFieldInput 更新自定义字段输入
type UpdateCustomFieldInput struct {
	Name      *string
	FieldKey  *string
	FieldType *string
	Options   *string
	Required  *bool
	SortOrder *int
}
