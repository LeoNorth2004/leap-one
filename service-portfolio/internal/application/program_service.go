package application

import (
	"context"
	"errors"
	"time"

	"leap-one/service-portfolio/internal/domain/entity"
	"leap-one/service-portfolio/internal/domain/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// 项目集服务相关错误定义
var (
	ErrProgramNotFound     = errors.New("项目集不存在")
	ErrProgramCodeExists   = errors.New("项目集编号已存在")
	ErrProductNotFound     = errors.New("产品不存在")
	ErrProductCodeExists   = errors.New("产品编码已存在")
	ErrProductLineNotFound = errors.New("产品线不存在")
	ErrVersionNotFound     = errors.New("版本不存在")
	ErrRoadmapItemNotFound = errors.New("路线图项不存在")
	ErrPlanNotFound        = errors.New("计划不存在")
	ErrMilestoneNotFound   = errors.New("里程碑不存在")
	ErrRiskNotFound        = errors.New("风险项不存在")
)

// ProgramService 项目集应用服务 - 协调项目集相关的业务流程
type ProgramService struct {
	programRepo   repository.ProgramRepository
	milestoneRepo repository.MilestoneRepository
	riskRepo      repository.RiskRepository
	productRepo   repository.ProductRepository
	logger        *zap.Logger
}

// NewProgramService 创建项目集应用服务实例
func NewProgramService(
	programRepo repository.ProgramRepository,
	milestoneRepo repository.MilestoneRepository,
	riskRepo repository.RiskRepository,
	productRepo repository.ProductRepository,
	logger *zap.Logger,
) *ProgramService {
	return &ProgramService{
		programRepo:   programRepo,
		milestoneRepo: milestoneRepo,
		riskRepo:      riskRepo,
		productRepo:   productRepo,
		logger:        logger,
	}
}

// CreateProgram 创建项目集用例
func (s *ProgramService) CreateProgram(ctx context.Context, req *entity.Program) (*entity.Program, error) {
	// 校验编号唯一性
	if existing, _ := s.programRepo.GetByCode(ctx, req.Code); existing != nil {
		return nil, ErrProgramCodeExists
	}

	if err := s.programRepo.Create(ctx, req); err != nil {
		s.logger.Error("创建项目集失败", zap.Error(err), zap.String("code", req.Code))
		return nil, errors.New("创建项目集失败")
	}

	s.logger.Info("项目集创建成功",
		zap.String("program_id", req.ID.String()),
		zap.String("code", req.Code),
	)
	return req, nil
}

// GetProgramDetail 获取项目集详情
func (s *ProgramService) GetProgramDetail(ctx context.Context, id uuid.UUID) (*entity.Program, error) {
	program, err := s.programRepo.GetByID(ctx, id)
	if err != nil || program == nil {
		return nil, ErrProgramNotFound
	}
	return program, nil
}

// UpdateProgram 更新项目集信息
func (s *ProgramService) UpdateProgram(ctx context.Context, program *entity.Program) error {
	// 检查是否存在
	existing, err := s.programRepo.GetByID(ctx, program.ID)
	if err != nil || existing == nil {
		return ErrProgramNotFound
	}

	// 如果修改了编号，检查唯一性
	if program.Code != existing.Code {
		if dup, _ := s.programRepo.GetByCode(ctx, program.Code); dup != nil && dup.ID != program.ID {
			return ErrProgramCodeExists
		}
	}

	if err := s.programRepo.Update(ctx, program); err != nil {
		s.logger.Error("更新项目集失败", zap.Error(err))
		return errors.New("更新项目集失败")
	}
	return nil
}

// DeleteProgram 删除项目集（软删除）
func (s *ProgramService) DeleteProgram(ctx context.Context, id uuid.UUID) error {
	program, err := s.programRepo.GetByID(ctx, id)
	if err != nil || program == nil {
		return ErrProgramNotFound
	}

	if err := s.programRepo.Delete(ctx, id); err != nil {
		s.logger.Error("删除项目集失败", zap.Error(err))
		return errors.New("删除项目集失败")
	}
	return nil
}

// ListPrograms 分页查询项目集列表
func (s *ProgramService) ListPrograms(ctx context.Context, page, pageSize int, keyword, status string) ([]*entity.Program, int64, error) {
	return s.programRepo.List(ctx, page, pageSize, keyword, status)
}

// GetProgramTree 获取项目集树形结构
func (s *ProgramService) GetProgramTree(ctx context.Context) ([]*entity.Program, error) {
	return s.programRepo.GetTree(ctx)
}

// CreateMilestone 创建里程碑
func (s *ProgramService) CreateMilestone(ctx context.Context, milestone *entity.Milestone) error {
	// 检查项目集是否存在
	if _, err := s.programRepo.GetByID(ctx, milestone.ProgramID); err != nil {
		return ErrProgramNotFound
	}

	if err := s.milestoneRepo.Create(ctx, milestone); err != nil {
		s.logger.Error("创建里程碑失败", zap.Error(err))
		return errors.New("创建里程碑失败")
	}
	return nil
}

// ListMilestones 获取项目集的里程碑列表
func (s *ProgramService) ListMilestones(ctx context.Context, programID uuid.UUID) ([]*entity.Milestone, error) {
	// 检查项目集是否存在
	if _, err := s.programRepo.GetByID(ctx, programID); err != nil {
		return nil, ErrProgramNotFound
	}
	return s.milestoneRepo.ListByProgramID(ctx, programID)
}

// CreateRisk 创建风险项
func (s *ProgramService) CreateRisk(ctx context.Context, risk *entity.Risk) error {
	// 检查项目集是否存在
	if _, err := s.programRepo.GetByID(ctx, risk.ProgramID); err != nil {
		return ErrProgramNotFound
	}

	if err := s.riskRepo.Create(ctx, risk); err != nil {
		s.logger.Error("创建风险项失败", zap.Error(err))
		return errors.New("创建风险项失败")
	}
	return nil
}

// ListRisks 获取项目集的风险列表
func (s *ProgramService) ListRisks(ctx context.Context, programID uuid.UUID) ([]*entity.Risk, error) {
	// 检查项目集是否存在
	if _, err := s.programRepo.GetByID(ctx, programID); err != nil {
		return nil, ErrProgramNotFound
	}
	return s.riskRepo.ListByProgramID(ctx, programID)
}

// GetProgramStatistics 获取项目集统计信息
func (s *ProgramService) GetProgramStatistics(ctx context.Context, programID uuid.UUID) (map[string]interface{}, error) {
	program, err := s.programRepo.GetByID(ctx, programID)
	if err != nil || program == nil {
		return nil, ErrProgramNotFound
	}

	stats := make(map[string]interface{})
	stats["program_id"] = programID.String()

	// 统计关联产品数和活跃产品数
	products, _ := s.productRepo.ListByProgramID(ctx, programID)
	stats["total_products"] = int64(len(products))
	activeCount := int64(0)
	for _, p := range products {
		if p.Status == "active" {
			activeCount++
		}
	}
	stats["active_products"] = activeCount

	// 统计里程碑
	milestones, _ := s.milestoneRepo.ListByProgramID(ctx, programID)
	stats["total_milestones"] = int64(len(milestones))
	doneCount := int64(0)
	for _, m := range milestones {
		if m.Status == "completed" {
			doneCount++
		}
	}
	stats["done_milestones"] = doneCount

	// 统计风险
	risks, _ := s.riskRepo.ListByProgramID(ctx, programID)
	stats["total_risks"] = int64(len(risks))
	openCount := int64(0)
	for _, r := range risks {
		if r.Status == "open" || r.Status == "mitigating" {
			openCount++
		}
	}
	stats["open_risks"] = openCount

	// 统计子项目集数量
	children, _ := s.programRepo.GetChildren(ctx, programID)
	stats["child_count"] = int64(len(children))

	return stats, nil
}

// parseDate 解析日期字符串为time.Time指针
func parseDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
