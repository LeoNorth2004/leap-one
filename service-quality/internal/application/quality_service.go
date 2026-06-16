package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"leap-one/service-quality/internal/domain/entity"
	"leap-one/service-quality/internal/domain/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 质量管理服务相关错误定义
var (
	ErrBugNotFound      = errors.New("Bug不存�?)
	ErrInvalidStatus    = errors.New("无效的状态转�?)
	ErrTestCaseNotFound = errors.New("测试用例不存�?)
)

// Bug状态机：合法的状态转换映�?
var validBugTransitions = map[string][]string{
	"new":         {"confirmed", "in_progress", "resolved", "closed"},
	"confirmed":   {"in_progress", "resolved", "closed"},
	"in_progress": {"resolved", "closed"},
	"resolved":    {"closed", "reopened"},
	"closed":      {"reopened"},
	"reopened":    {"in_progress", "resolved", "closed"},
	"cancelled":   {},
}

// QualityService 质量管理应用服务 - 协调Bug状态机、统计计算等核心业务逻辑
type QualityService struct {
	bugRepo  repository.BugRepository
	caseRepo repository.TestCaseRepository
	logger   *zap.Logger
}

// NewQualityService 创建质量管理应用服务实例
func NewQualityService(
	bugRepo repository.BugRepository,
	caseRepo repository.TestCaseRepository,
	logger *zap.Logger,
) *QualityService {
	return &QualityService{
		bugRepo:  bugRepo,
		caseRepo: caseRepo,
		logger:   logger,
	}
}

// ValidateTransition 验证Bug状态转换是否合�?
// 根据预定义的状态机规则检查当前状态是否可以转换到目标状�?
func (s *QualityService) ValidateTransition(currentStatus, targetStatus string) error {
	if currentStatus == targetStatus {
		return nil // 同状态不报错
	}

	allowedTargets, exists := validBugTransitions[currentStatus]
	if !exists {
		return ErrInvalidStatus
	}

	for _, allowed := range allowedTargets {
		if allowed == targetStatus {
			return nil
		}
	}

	return ErrInvalidStatus
}

// GetSeverityName 获取严重程度名称
func (s *QualityService) GetSeverityName(severity int) string {
	names := map[int]string{
		1: "致命",
		2: "严重",
		3: "一�?,
		4: "提示",
	}
	if name, ok := names[severity]; ok {
		return name
	}
	return "未知"
}

// GetPriorityName 获取优先级名�?
func (s *QualityService) GetPriorityName(priority int) string {
	names := map[int]string{
		1: "最�?,
		2: "�?,
		3: "�?,
		4: "�?,
		5: "最�?,
	}
	if name, ok := names[priority]; ok {
		return name
	}
	return "未知"
}

// GetBugTypeName 获取Bug类型名称
func (s *QualityService) GetBugTypeName(bugType string) string {
	names := map[string]string{
		"code_bug":    "代码缺陷",
		"design_bug":  "设计缺陷",
		"data_bug":    "数据问题",
		"config":      "配置问题",
		"security":    "安全漏洞",
		"performance": "性能问题",
		"ui":          "界面问题",
	}
	if name, ok := names[bugType]; ok {
		return name
	}
	return bugType
}

// GetResolutionName 获取解决方案名称
func (s *QualityService) GetResolutionName(resolution string) string {
	names := map[string]string{
		"fixed":      "已修�?,
		"wont_fix":   "不予修复",
		"duplicate":  "重复问题",
		"by_design":  "设计如此",
		"workaround": "临时方案",
		"postponed":  "延期处理",
	}
	if name, ok := names[resolution]; ok {
		return name
	}
	return resolution
}

// CalculatePassRate 计算通过率（百分比）
func (s *QualityService) CalculatePassRate(passed, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(passed) / float64(total) * 100
}

// GetBugWorkflowDefault 获取默认Bug工作流配置信�?
// 返回标准工作流的状态转换规则描�?
func (s *QualityService) GetBugWorkflowDefault() []entity.BugWorkflowTransition {
	return []entity.BugWorkflowTransition{
		{FromStatus: "new", ToStatus: "confirmed", Name: "确认Bug", Condition: "确认Bug有效且可复现", SortOrder: 1},
		{FromStatus: "new", ToStatus: "in_progress", Name: "开始处�?, Condition: "直接分配并开始处�?, SortOrder: 2},
		{FromStatus: "new", ToStatus: "closed", Name: "关闭Bug", Condition: "无效或无法复现的Bug", SortOrder: 3},
		{FromStatus: "confirmed", ToStatus: "in_progress", Name: "开始处�?, Condition: "分配给开发人�?, SortOrder: 1},
		{FromStatus: "confirmed", ToStatus: "resolved", Name: "解决Bug", Condition: "修复完成并验证通过", SortOrder: 2},
		{FromStatus: "confirmed", ToStatus: "closed", Name: "关闭Bug", Condition: "确认为无效Bug或设计如�?, SortOrder: 3},
		{FromStatus: "in_progress", ToStatus: "resolved", Name: "解决Bug", Condition: "代码已修复并通过验证", SortOrder: 1},
		{FromStatus: "in_progress", ToStatus: "closed", Name: "关闭Bug", Condition: "不再需要修�?, SortOrder: 2},
		{FromStatus: "resolved", ToStatus: "closed", Name: "关闭Bug", Condition: "验证通过，正式关�?, SortOrder: 1},
		{FromStatus: "resolved", ToStatus: "reopened", Name: "重新打开", Condition: "回归测试未通过", SortOrder: 2},
		{FromStatus: "closed", ToStatus: "reopened", Name: "重新激�?, Condition: "需要重新处�?, SortOrder: 1},
		{FromStatus: "reopened", ToStatus: "in_progress", Name: "重新处理", Condition: "再次分配处理", SortOrder: 1},
		{FromStatus: "reopened", ToStatus: "resolved", Name: "再次解决", Condition: "修复后再次验证通过", SortOrder: 2},
		{FromStatus: "reopened", ToStatus: "closed", Name: "关闭Bug", Condition: "确认无需再处�?, SortOrder: 3},
	}
}

// InitDefaultWorkflow 初始化默认Bug工作流数�?
// 在首次启动时创建标准工作流及转换规则
func (s *QualityService) InitDefaultWorkflow(ctx context.Context, db *gorm.DB) error {
	var existing entity.BugWorkflow
	err := db.First(&existing, "is_default = ?", true).Error
	if err == nil && existing.ID != uuid.Nil {
		s.logger.Debug("默认Bug工作流已存在，跳过初始化")
		return nil
	}

	workflow := &entity.BugWorkflow{
		Name:          "标准Bug工作�?,
		InitialStatus: "new",
		IsDefault:     true,
		Transitions:   s.GetBugWorkflowDefault(),
	}

	if err := db.Create(workflow).Error; err != nil {
		return err
	}

	s.logger.Info("默认Bug工作流初始化成功",
		zap.String("workflow_id", workflow.ID.String()),
		zap.Int("transition_count", len(workflow.Transitions)),
	)
	return nil
}

// InitDefaultEnvironments 初始化默认测试环�?
// 在首次启动时创建开发、测试、预发布环境
func (s *QualityService) InitDefaultEnvironments(ctx context.Context, db *gorm.DB) error {
	defaultEnvs := []struct {
		Name        string
		URL         string
		Type        string
		Description string
	}{
		{"开发环�?, "http://localhost:8080", "dev", "本地开发调试环�?},
		{"测试环境", "http://test.example.com", "test", "功能测试与集成测试环�?},
		{"预发布环�?, "http://staging.example.com", "staging", "上线前验收测试环�?},
	}

	for _, env := range defaultEnvs {
		var count int64
		db.Model(&entity.TestEnvironment{}).Where("name = ?", env.Name).Count(&count)
		if count > 0 {
			continue
		}

		e := &entity.TestEnvironment{
			Name:        env.Name,
			URL:         env.URL,
			Type:        env.Type,
			Description: env.Description,
			IsActive:    true,
		}

		if createErr := db.Create(e).Error; createErr != nil {
			s.logger.Warn("创建默认测试环境失败",
				zap.String("name", env.Name),
				zap.Error(createErr),
			)
		} else {
			s.logger.Debug("默认测试环境已创�?,
				zap.String("name", env.Name),
				zap.String("id", e.ID.String()),
			)
		}
	}

	return nil
}
