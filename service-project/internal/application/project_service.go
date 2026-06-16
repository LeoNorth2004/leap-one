package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"leap-one/service-project/internal/domain/entity"
	"leap-one/service-project/internal/domain/repository"
	"go.uber.org/zap"
)

// 项目服务相关错误定义
var (
	ErrProjectNotFound      = errors.New("项目不存在")
	ErrProjectCodeExists    = errors.New("项目编号已存在")
	ErrInvalidProjectStatus = errors.New("无效的项目状态转换")
)

// 有效的项目状态转换映射（key=当前状态, value=允许的目标状态列表）
var validStatusTransitions = map[string][]string{
	"planning":   {"executing", "cancelled"},
	"executing":  {"paused", "completed", "cancelled"},
	"paused":     {"executing", "completed", "cancelled"},
	"completed":  {"archived"},
	"cancelled":  {},
	"archived":   {},
}

// ProjectService 项目应用服务 - 协调项目相关的业务流程
type ProjectService struct {
	projectRepo repository.ProjectRepository
	memberRepo  repository.ProjectMemberRepository
	logger      *zap.Logger
}

// NewProjectService 创建项目应用服务实例
func NewProjectService(
	projectRepo repository.ProjectRepository,
	memberRepo repository.ProjectMemberRepository,
	logger *zap.Logger,
) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
	memberRepo:  memberRepo,
		logger:      logger,
	}
}

// CreateProject 创建项目用例
func (s *ProjectService) CreateProject(ctx context.Context, input *CreateProjectInput) (*entity.Project, error) {
	// 检查项目编号是否已存在
	if input.Code != "" {
		if existing, _ := s.projectRepo.GetByCode(ctx, input.Code); existing != nil {
			return nil, ErrProjectCodeExists
		}
	}

	now := time.Now()
	project := &entity.Project{
		Name:        input.Name,
		Code:        input.Code,
		Description: input.Description,
		ProgramID:   input.ProgramID,
		PMID:        input.PMID,
		Type:        input.Type,
		Priority:    input.Priority,
		Status:      "planning",
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		Budget:      input.Budget,
		TemplateID:  input.TemplateID,
		CreatedByID: input.CreatedByID,
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if project.Type == "" {
		project.Type = "software"
	}
	if project.Priority == 0 {
		project.Priority = 3
	}
	if project.Status == "" {
		project.Status = "planning"
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		s.logger.Error("创建项目失败", zap.Error(err), zap.String("name", input.Name))
		return nil, errors.New("创建项目失败")
	}

	// 如果指定了项目经理，自动添加为项目成员
	if project.PMID != nil && *project.PMID != uuid.Nil {
		member := &entity.ProjectMember{
			ProjectID: project.ID,
			UserID:    *project.PMID,
			Role:      "pm",
			JoinedAt:  now,
		}
		_ = s.memberRepo.Add(ctx, member)
	}

	s.logger.Info("项目创建成功",
		zap.String("project_id", project.ID.String()),
		zap.String("name", project.Name),
		zap.String("code", project.Code),
	)
	return project, nil
}

// GetProjectDetail 获取项目详情
func (s *ProjectService) GetProjectDetail(ctx context.Context, id uuid.UUID) (*entity.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil || project == nil {
		return nil, ErrProjectNotFound
	}
	return project, nil
}

// UpdateProject 更新项目用例
func (s *ProjectService) UpdateProject(ctx context.Context, id uuid.UUID, input *UpdateProjectInput) (*entity.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil || project == nil {
		return nil, ErrProjectNotFound
	}

	applyProjectUpdate(project, input)
	project.Version++
	project.UpdatedAt = time.Now()

	if err := s.projectRepo.Update(ctx, project); err != nil {
		s.logger.Error("更新项目失败", zap.Error(err), zap.String("project_id", id.String()))
		return nil, errors.New("更新项目失败")
	}

	return project, nil
}

// DeleteProject 删除项目用例
func (s *ProjectService) DeleteProject(ctx context.Context, id uuid.UUID) error {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil || project == nil {
		return ErrProjectNotFound
	}

	if err := s.projectRepo.Delete(ctx, id); err != nil {
		s.logger.Error("删除项目失败", zap.Error(err), zap.String("project_id", id.String()))
		return errors.New("删除项目失败")
	}

	s.logger.Info("项目已删除", zap.String("project_id", id.String()), zap.String("name", project.Name))
	return nil
}

// ListProjects 分页查询项目列表
func (s *ProjectService) ListProjects(ctx context.Context, page, size int, keyword, status, programID, pmID, sortBy, sortOrder string) ([]*entity.Project, int64, error) {
	return s.projectRepo.List(ctx, page, size, keyword, status, programID, pmID, sortBy, sortOrder)
}

// ChangeProjectStatus 项目状态流转
func (s *ProjectService) ChangeProjectStatus(ctx context.Context, id uuid.UUID, newStatus string) error {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil || project == nil {
		return ErrProjectNotFound
	}

	validTargets, ok := validStatusTransitions[project.Status]
	if !ok {
		return ErrInvalidProjectStatus
	}

	allowed := false
	for _, v := range validTargets {
		if v == newStatus {
			allowed = true
			break
		}
	}
	if !allowed {
		return ErrInvalidProjectStatus
	}

	if err := s.projectRepo.UpdateStatus(ctx, id, newStatus); err != nil {
		s.logger.Error("更新项目状态失败", zap.Error(err),
			zap.String("project_id", id.String()),
			zap.String("from", project.Status),
			zap.String("to", newStatus),
		)
		return errors.New("状态更新失败")
	}

	s.logger.Info("项目状态变更",
		zap.String("project_id", id.String()),
		zap.String("old_status", project.Status),
		zap.String("new_status", newStatus),
	)
	return nil
}

// ==================== 输入结构体 ====================

// CreateProjectInput 创建项目输入
type CreateProjectInput struct {
	Name        string
	Code        string
	Description string
	ProgramID   *uuid.UUID
	PMID        *uuid.UUID
	Type        string
	Priority    int
	StartDate   *time.Time
	EndDate     *time.Time
	Budget      *float64
	TemplateID  *uuid.UUID
	CreatedByID uuid.UUID
}

// UpdateProjectInput 更新项目输入
type UpdateProjectInput struct {
	Name        *string
	Description *string
	PMID        *uuid.UUID
	Type        *string
	Priority    *int
	StartDate   *time.Time
	EndDate     *time.Time
	Budget      *float64
	UpdatedByID uuid.UUID
}

// applyProjectUpdate 将更新输入应用到实体
func applyProjectUpdate(project *entity.Project, input *UpdateProjectInput) {
	if input.Name != nil {
		project.Name = *input.Name
	}
	if input.Description != nil {
		project.Description = *input.Description
	}
	if input.PMID != nil {
		project.PMID = input.PMID
	}
	if input.Type != nil {
		project.Type = *input.Type
	}
	if input.Priority != nil {
		project.Priority = *input.Priority
	}
	if input.StartDate != nil {
		project.StartDate = input.StartDate
	}
	if input.EndDate != nil {
		project.EndDate = input.EndDate
	}
	if input.Budget != nil {
		project.Budget = input.Budget
	}
}
