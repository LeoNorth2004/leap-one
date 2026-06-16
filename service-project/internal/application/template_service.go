package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"leap-one/service-project/internal/domain/entity"
	"leap-one/service-project/internal/domain/repository"
	"go.uber.org/zap"
)

// 模板服务相关错误定义
var (
	ErrTemplateNotFound = errors.New("模板不存�?)
)

// TemplateService 项目模板应用服务
type TemplateService struct {
	templateRepo repository.ProjectTemplateRepository
	logger       *zap.Logger
}

// NewTemplateService 创建模板服务实例
func NewTemplateService(
	templateRepo repository.ProjectTemplateRepository,
	logger *zap.Logger,
) *TemplateService {
	return &TemplateService{
		templateRepo: templateRepo,
		logger:       logger,
	}
}

// CreateTemplate 创建自定义模板（非系统预置）
func (s *TemplateService) CreateTemplate(ctx context.Context, req *CreateTemplateInput) (*entity.ProjectTemplate, error) {
	template := &entity.ProjectTemplate{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Config:      req.Config,
		IsSystem:    false, // 用户创建的模板标记为非系统预�?	}

	if err := s.templateRepo.Create(ctx, template); err != nil {
		s.logger.Error("创建模板失败", zap.Error(err))
		return nil, errors.New("创建模板失败")
	}

	s.logger.Info("项目模板已创�?,
		zap.String("template_id", template.ID.String()),
		zap.String("name", req.Name),
	)
	return template, nil
}

// GetTemplateByID 根据ID获取模板详情
func (s *TemplateService) GetTemplateByID(ctx context.Context, id uuid.UUID) (*entity.ProjectTemplate, error) {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil || template == nil {
		return nil, ErrTemplateNotFound
	}
	return template, nil
}

// ListTemplates 分页查询模板列表
func (s *TemplateService) ListTemplates(
	ctx context.Context,
	page, pageSize int,
	templateType string,
) ([]*entity.ProjectTemplate, int64, error) {
	return s.templateRepo.List(ctx, page, pageSize, templateType)
}

// UpdateTemplate 更新模板
func (s *TemplateService) UpdateTemplate(ctx context.Context, id uuid.UUID, req *UpdateTemplateInput) (*entity.ProjectTemplate, error) {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil || template == nil {
		return nil, ErrTemplateNotFound
	}

	if req.Name != nil {
		template.Name = *req.Name
	}
	if req.Description != nil {
		template.Description = *req.Description
	}
	if req.Type != nil {
		template.Type = *req.Type
	}
	if req.Config != nil {
		template.Config = *req.Config
	}

	if err := s.templateRepo.Update(ctx, template); err != nil {
		s.logger.Error("更新模板失败", zap.Error(err))
		return nil, errors.New("更新模板失败")
	}

	return template, nil
}

// DeleteTemplate 删除模板（仅允许删除非系统预置模板）
func (s *TemplateService) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	if err := s.templateRepo.Delete(ctx, id); err != nil {
		if err.Error() == "系统预置模板不允许删�? {
			return errors.New("系统预置模板不允许删�?)
		}
		s.logger.Error("删除模板失败", zap.Error(err))
		return errors.New("删除模板失败")
	}

	s.logger.Info("模板已删�?, zap.String("template_id", id.String()))
	return nil
}

// ListSystemTemplates 获取所有系统预置模�?func (s *TemplateService) ListSystemTemplates(ctx context.Context) ([]*entity.ProjectTemplate, error) {
	return s.templateRepo.ListSystemTemplates(ctx)
}

// CreateTemplateInput 创建模板输入
type CreateTemplateInput struct {
	Name        string
	Description string
	Type        string
	Config      string
}

// UpdateTemplateInput 更新模板输入
type UpdateTemplateInput struct {
	Name        *string
	Description *string
	Type        *string
	Config      *string
}
