package application

import (
	"context"
	"time"

	"leap-one/service-project/internal/domain/entity"
	"leap-one/service-project/internal/domain/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ProjectStatisticsService 项目统计分析服务
type ProjectStatisticsService struct {
	projectRepo   repository.ProjectRepository
	memberRepo    repository.ProjectMemberRepository
	milestoneRepo repository.ProjectMilestoneRepository
	riskRepo      repository.ProjectRiskRepository
	iterRepo      repository.IterationRepository
	logger        *zap.Logger
}

// NewProjectStatisticsService 创建项目统计服务实例
func NewProjectStatisticsService(
	projectRepo repository.ProjectRepository,
	memberRepo repository.ProjectMemberRepository,
	milestoneRepo repository.ProjectMilestoneRepository,
	riskRepo repository.ProjectRiskRepository,
	iterRepo repository.IterationRepository,
	logger *zap.Logger,
) *ProjectStatisticsService {
	return &ProjectStatisticsService{
		projectRepo:   projectRepo,
		memberRepo:    memberRepo,
		milestoneRepo: milestoneRepo,
		riskRepo:      riskRepo,
		iterRepo:      iterRepo,
		logger:        logger,
	}
}

// GetProjectStatistics 获取项目的综合统计数据
func (s *ProjectStatisticsService) GetProjectStatistics(ctx context.Context, projectID uuid.UUID) (*ProjectStatisticsOutput, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil || project == nil {
		return nil, ErrProjectNotFound
	}

	output := &ProjectStatisticsOutput{
		ProjectID:   projectID.String(),
		ProjectName: project.Name,
	}

	// 成员统计
	members, _ := s.memberRepo.ListByProjectID(ctx, projectID)
	output.MemberStats.TotalMembers = len(members)
	output.MemberStats.ByRole = make(map[string]int)
	for _, m := range members {
		output.MemberStats.ByRole[m.Role]++
	}

	// 里程碑统计
	milestones, _ := s.milestoneRepo.ListByProjectID(ctx, projectID)
	output.MilestoneStats.Total = len(milestones)
	for _, ms := range milestones {
		switch ms.Status {
		case "completed":
			output.MilestoneStats.Completed++
		case "pending":
			output.MilestoneStats.Pending++
		default:
			output.MilestoneStats.Overdue++
		}
	}
	if output.MilestoneStats.Total > 0 {
		output.MilestoneStats.CompletionRate = float64(output.MilestoneStats.Completed) / float64(output.MilestoneStats.Total) * 100
	}

	// 风险统计
	risks, _ := s.riskRepo.ListByProjectID(ctx, projectID)
	output.RiskStats.Total = len(risks)
	for _, r := range risks {
		switch r.Status {
		case "open":
			output.RiskStats.Open++
		case "mitigating":
			output.RiskStats.Mitigating++
		case "closed":
			output.RiskStats.Closed++
		}
		if r.Severity >= 12 {
			output.RiskStats.HighRisk++
		}
	}

	// 迭代统计
	iterations, _ := s.iterRepo.ListByProjectID(ctx, projectID)
	output.IterationStats.Total = len(iterations)
	for _, iter := range iterations {
		switch iter.Status {
		case "completed":
			output.IterationStats.Completed++
		case "active":
			output.IterationStats.Active++
		}
	}

	// 项目总览信息
	output.Overview.TotalTasks = 0 // TODO: 集成任务服务后填充
	output.Overview.CompletedTasks = 0
	output.Overview.InProgressTasks = 0
	output.Overview.OverdueTasks = 0
	output.Overview.TotalIterations = output.IterationStats.Total
	output.Overview.ActiveIteration = "" // TODO: 从活跃迭代获取名称
	output.Overview.DaysRemaining = calculateDaysRemaining(project.EndDate)

	// 预算相关
	if project.Budget != nil {
		output.Overview.BudgetUsed = 0 // TODO: 从成本服务获取实际支出
		output.Overview.BudgetRate = output.Overview.BudgetUsed / *project.Budget * 100
	}

	return output, nil
}

// GetBurndownData 获取燃尽图数据
func (s *ProjectStatisticsService) GetBurndownData(ctx context.Context, projectID uuid.UUID) (*BurndownOutput, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil || project == nil {
		return nil, ErrProjectNotFound
	}

	activeIter, _ := s.iterRepo.GetActiveIteration(ctx, projectID)

	output := &BurndownOutput{
		SprintName:  activeIter.Name,
		TotalPoints: 0, // TODO: 从任务服务获取实际故事点总数
		IdealLine:   generateIdealLine(activeIter),
		ActualLine:  []BurndownPoint{}, // TODO: 从任务服务获取实际燃尽数据
		Remaining:   0,
	}

	return output, nil
}

// GetGanttData 获取甘特图数据
func (s *ProjectStatisticsService) GetGanttData(ctx context.Context, projectID uuid.UUID) (*GanttOutput, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil || project == nil {
		return nil, ErrProjectNotFound
	}

	output := &GanttOutput{
		ProjectName: project.Name,
		Tasks:       []GanttTaskItem{}, // TODO: 从任务服务获取任务数据
		Milestones:  []GanttMilestoneItem{},
	}

	// 填充里程碑数据到甘特图
	milestones, _ := s.milestoneRepo.ListByProjectID(ctx, projectID)
	for _, ms := range milestones {
		output.Milestones = append(output.Milestones, GanttMilestoneItem{
			ID:      ms.ID.String(),
			Name:    ms.Name,
			DueDate: ms.DueDate.Format("2006-01-02"),
			Status:  ms.Status,
		})
	}

	return output, nil
}

// ==================== 输出结构体 ====================

// ProjectStatisticsOutput 统计输出
type ProjectStatisticsOutput struct {
	ProjectID      string               `json:"project_id"`
	ProjectName    string               `json:"project_name"`
	Overview       OverviewStats        `json:"overview"`
	MemberStats    MemberStatsOutput    `json:"member_stats"`
	MilestoneStats MilestoneStatsOutput `json:"milestone_stats"`
	RiskStats      RiskStatsOutput      `json:"risk_stats"`
	IterationStats IterationStatsOutput `json:"iteration_stats"`
}

// OverviewStats 总览统计输出
type OverviewStats struct {
	TotalTasks      int     `json:"total_tasks"`
	CompletedTasks  int     `json:"completed_tasks"`
	InProgressTasks int     `json:"in_progress_tasks"`
	OverdueTasks    int     `json:"overdue_tasks"`
	CompletionRate  float64 `json:"completion_rate"`
	TotalIterations int     `json:"total_iterations"`
	ActiveIteration string  `json:"active_iteration"`
	DaysRemaining   int     `json:"days_remaining"`
	BudgetUsed      float64 `json:"budget_used"`
	BudgetRate      float64 `json:"budget_rate"`
}

// MemberStatsOutput 成员统计输出
type MemberStatsOutput struct {
	TotalMembers int            `json:"total_members"`
	ByRole       map[string]int `json:"by_role"`
}

// MilestoneStatsOutput 里程碑统计输出
type MilestoneStatsOutput struct {
	Total          int     `json:"total"`
	Completed      int     `json:"completed"`
	Pending        int     `json:"pending"`
	Overdue        int     `json:"overdue"`
	CompletionRate float64 `json:"completion_rate"`
}

// RiskStatsOutput 风险统计输出
type RiskStatsOutput struct {
	Total      int `json:"total"`
	Open       int `json:"open"`
	Mitigating int `json:"mitigating"`
	Closed     int `json:"closed"`
	HighRisk   int `json:"high_risk"`
}

// IterationStatsOutput 迭代统计输出
type IterationStatsOutput struct {
	Total       int     `json:"total"`
	Completed   int     `json:"completed"`
	Active      int     `json:"active"`
	AvgVelocity float64 `json:"avg_velocity"`
}

// BurndownOutput 燃尽图数据输出
type BurndownOutput struct {
	SprintName  string          `json:"sprint_name"`
	TotalPoints float64         `json:"total_points"`
	IdealLine   []BurndownPoint `json:"ideal_line"`
	ActualLine  []BurndownPoint `json:"actual_line"`
	Remaining   float64         `json:"remaining"`
}

// GanttOutput 甘特图数据输出
type GanttOutput struct {
	ProjectName string               `json:"project_name"`
	Tasks       []GanttTaskItem      `json:"tasks"`
	Milestones  []GanttMilestoneItem `json:"milestones"`
}

// GanttTaskItem 甘特图任务项
type GanttTaskItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Progress  int    `json:"progress"`
	Status    string `json:"status"`
	Color     string `json:"color,omitempty"`
}

// GanttMilestoneItem 甘特图里程碑项
type GanttMilestoneItem struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	DueDate string `json:"due_date"`
	Status  string `json:"status"`
}

// BurndownPoint 燃尽图数据点
type BurndownPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// ==================== 辅助函数 ====================

// calculateDaysRemaining 计算项目剩余天数
func calculateDaysRemaining(endDate *time.Time) int {
	if endDate == nil {
		return -1 // 无截止日期
	}
	now := time.Now()
	remaining := endDate.Sub(now).Hours() / 24
	if remaining < 0 {
		return 0
	}
	return int(remaining)
}

// generateIdealLine 生成理想燃尽线数据点
func generateIdealLine(iter *entity.Iteration) []BurndownPoint {
	if iter == nil {
		return []BurndownPoint{}
	}

	totalDays := int(iter.EndDate.Sub(iter.StartDate).Hours()/24) + 1
	if totalDays <= 0 {
		return []BurndownPoint{}
	}

	points := make([]BurndownPoint, totalDays)
	for i := 0; i < totalDays; i++ {
		date := iter.StartDate.AddDate(0, 0, i)
		remainingRatio := float64(totalDays-1-i) / float64(totalDays-1)
		if remainingRatio < 0 {
			remainingRatio = 0
		}
		points[i] = BurndownPoint{
			Date:  date.Format("2006-01-02"),
			Value: remainingRatio, // 归一化值，实际应乘以总故事点
		}
	}
	return points
}

// NowFunc 当前时间函数（与repository_impl共用）
var NowFunc = func() time.Time {
	return time.Now()
}
