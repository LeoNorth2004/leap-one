package application

import (
	"context"
	"errors"

	"leap-one/service-task/internal/domain/entity"
	"leap-one/service-task/internal/domain/repository"

	"github.com/google/uuid"
)

// 模板服务相关错误定义
var (
	ErrTemplateNotFound = errors.New("模板不存在")
)

// TemplateService 工单模板应用服务
type TemplateService struct {
	templateRepo repository.IssueTemplateRepository
	logger       interface {
		Info(msg string, fields ...interface{})
	}
}

// NewTemplateService 创建模板服务实例
func NewTemplateService(templateRepo repository.IssueTemplateRepository) *TemplateService {
	return &TemplateService{
		templateRepo: templateRepo,
	}
}

// CreateTemplate 创建模板
func (s *TemplateService) CreateTemplate(ctx context.Context, req interface{}) (*entity.IssueTemplate, error) {
	tmpl := &entity.IssueTemplate{
		Name:     "default",
		Type:     "bug",
		IsSystem: false,
	}
	if err := s.templateRepo.Create(ctx, tmpl); err != nil {
		return nil, errors.New("创建模板失败")
	}
	return tmpl, nil
}

// GetTemplate 获取模板详情
func (s *TemplateService) GetTemplate(ctx context.Context, id uuid.UUID) (*entity.IssueTemplate, error) {
	tmpl, err := s.templateRepo.GetByID(ctx, id)
	if err != nil || tmpl == nil {
		return nil, ErrTemplateNotFound
	}
	return tmpl, nil
}

// UpdateTemplate 更新模板
func (s *TemplateService) UpdateTemplate(ctx context.Context, id uuid.UUID, req interface{}) error {
	tmpl, err := s.templateRepo.GetByID(ctx, id)
	if err != nil || tmpl == nil {
		return ErrTemplateNotFound
	}
	return s.templateRepo.Update(ctx, tmpl)
}

// DeleteTemplate 删除模板
func (s *TemplateService) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	tmpl, err := s.templateRepo.GetByID(ctx, id)
	if err != nil || tmpl == nil {
		return ErrTemplateNotFound
	}
	if tmpl.IsSystem {
		return errors.New("系统预置模板不能删除")
	}
	return s.templateRepo.Delete(ctx, id)
}

// ListTemplates 分页查询模板列表
func (s *TemplateService) ListTemplates(ctx context.Context, page, pageSize int, tmplType string) ([]*entity.IssueTemplate, int64, error) {
	return s.templateRepo.List(ctx, page, pageSize, tmplType)
}

// WorkflowService 工作流应用服务
type WorkflowService struct {
	workflowRepo repository.IssueWorkflowRepository
}

// NewWorkflowService 创建工作流服务实例
func NewWorkflowService(workflowRepo repository.IssueWorkflowRepository) *WorkflowService {
	return &WorkflowService{workflowRepo: workflowRepo}
}

// CreateWorkflow 创建工作流
func (s *WorkflowService) CreateWorkflow(ctx context.Context, req interface{}) (*entity.IssueWorkflow, error) {
	wf := &entity.IssueWorkflow{
		Name:          "default",
		Type:          "bug",
		InitialStatus: "new",
	}
	if err := s.workflowRepo.Create(ctx, wf); err != nil {
		return nil, errors.New("创建工作流失败")
	}
	return wf, nil
}

// GetWorkflow 获取工作流详情（含转换规则）
func (s *WorkflowService) GetWorkflow(ctx context.Context, id uuid.UUID) (*entity.IssueWorkflow, error) {
	wf, err := s.workflowRepo.GetByID(ctx, id)
	if err != nil || wf == nil {
		return nil, ErrWorkflowNotFound
	}
	return wf, nil
}

// UpdateWorkflow 更新工作流
func (s *WorkflowService) UpdateWorkflow(ctx context.Context, id uuid.UUID, req interface{}) error {
	wf, err := s.workflowRepo.GetByID(ctx, id)
	if err != nil || wf == nil {
		return ErrWorkflowNotFound
	}
	return s.workflowRepo.Update(ctx, wf)
}

// DeleteWorkflow 删除工作流
func (s *WorkflowService) DeleteWorkflow(ctx context.Context, id uuid.UUID) error {
	wf, err := s.workflowRepo.GetByID(ctx, id)
	if err != nil || wf == nil {
		return ErrWorkflowNotFound
	}
	return s.workflowRepo.Delete(ctx, id)
}

// ListWorkflows 分页查询工作流列表
func (s *WorkflowService) ListWorkflows(ctx context.Context, page, pageSize int, wfType string) ([]*entity.IssueWorkflow, int64, error) {
	return s.workflowRepo.List(ctx, page, pageSize, wfType)
}

// AddTransition 添加状态转换规则
func (s *WorkflowService) AddTransition(ctx context.Context, workflowID uuid.UUID, req interface{}) (*entity.IssueWorkflowTransition, error) {
	wf, err := s.workflowRepo.GetByID(ctx, workflowID)
	if err != nil || wf == nil {
		return nil, ErrWorkflowNotFound
	}

	transition := &entity.IssueWorkflowTransition{
		WorkflowID: workflowID,
		FromStatus: "new",
		ToStatus:   "in_progress",
		Name:       "开始处理",
		SortOrder:  0,
	}

	if err := s.workflowRepo.AddTransition(ctx, transition); err != nil {
		return nil, errors.New("添加状态转换规则失败")
	}
	return transition, nil
}

// SLAConfigService SLA配置应用服务
type SLAConfigService struct {
	slaConfigRepo repository.IssueSLAConfigRepository
}

// NewSLAConfigService 创建SLA配置服务实例
func NewSLAConfigService(slaConfigRepo repository.IssueSLAConfigRepository) *SLAConfigService {
	return &SLAConfigService{slaConfigRepo: slaConfigRepo}
}

// CreateSLAConfig 创建SLA配置
func (s *SLAConfigService) CreateSLAConfig(ctx context.Context, req interface{}) (*entity.IssueSLAConfig, error) {
	cfg := &entity.IssueSLAConfig{
		Type:        "bug",
		Priority:    3,
		ResponseSLA: 30,
		ResolveSLA:  480,
	}
	if err := s.slaConfigRepo.Create(ctx, cfg); err != nil {
		return nil, errors.New("创建SLA配置失败")
	}
	return cfg, nil
}

// GetSLAConfig 获取SLA配置详情
func (s *SLAConfigService) GetSLAConfig(ctx context.Context, id uuid.UUID) (*entity.IssueSLAConfig, error) {
	cfg, err := s.slaConfigRepo.GetByID(ctx, id)
	if err != nil || cfg == nil {
		return nil, errors.New("SLA配置不存在")
	}
	return cfg, nil
}

// UpdateSLAConfig 更新SLA配置
func (s *SLAConfigService) UpdateSLAConfig(ctx context.Context, id uuid.UUID, req interface{}) error {
	cfg, err := s.slaConfigRepo.GetByID(ctx, id)
	if err != nil || cfg == nil {
		return errors.New("SLA配置不存在")
	}
	return s.slaConfigRepo.Update(ctx, cfg)
}

// DeleteSLAConfig 删除SLA配置
func (s *SLAConfigService) DeleteSLAConfig(ctx context.Context, id uuid.UUID) error {
	return s.slaConfigRepo.Delete(ctx, id)
}

// ListSLAConfigs 分页查询SLA配置列表
func (s *SLAConfigService) ListSLAConfigs(ctx context.Context, page, pageSize int) ([]*entity.IssueSLAConfig, int64, error) {
	return s.slaConfigRepo.List(ctx, page, pageSize)
}
