package dto

// ==================== 模板相关DTO ====================

// CreateTemplateRequest 创建模板请求
type CreateTemplateRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=200"`                               // 模板名称
	Description string `json:"description" binding:"max=2000"`                                      // 模板描述
	Type        string `json:"type" binding:"required,oneof=agile waterfall lightweight lifecycle"` // 模板类型
	Config      string `json:"config" binding:"required"`                                           // JSON配置
}

// UpdateTemplateRequest 更新模板请求
type UpdateTemplateRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=200"`
	Description *string `json:"description" binding:"omitempty,max=2000"`
	Type        *string `json:"type" binding:"omitempty,oneof=agile waterfall lightweight lifecycle"`
	Config      *string `json:"config"`
}

// TemplateListResponse 模板列表响应
type TemplateListResponse struct {
	List  []TemplateInfo `json:"list"`
	Total int64          `json:"total"`
}

// TemplateInfo 模板信息
type TemplateInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	Config      string `json:"config"`
	IsSystem    bool   `json:"is_system"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ==================== 统计相关DTO ====================

// ProjectStatisticsResponse 项目统计响应
type ProjectStatisticsResponse struct {
	ProjectID      string              `json:"project_id"`
	ProjectName    string              `json:"project_name"`
	Overview       StatisticsOverview  `json:"overview"`
	MemberStats    MemberStatistics    `json:"member_stats"`
	MilestoneStats MilestoneStats      `json:"milestone_stats"`
	RiskStats      RiskStatistics      `json:"risk_stats"`
	IterationStats IterationStatistics `json:"iteration_stats"`
}

// StatisticsOverview 总览统计
type StatisticsOverview struct {
	TotalTasks      int     `json:"total_tasks"`       // 总任务数
	CompletedTasks  int     `json:"completed_tasks"`   // 已完成任务数
	InProgressTasks int     `json:"in_progress_tasks"` // 进行中任务数
	OverdueTasks    int     `json:"overdue_tasks"`     // 逾期任务数
	CompletionRate  float64 `json:"completion_rate"`   // 完成率
	TotalIterations int     `json:"total_iterations"`  // 总迭代数
	ActiveIteration string  `json:"active_iteration"`  // 当前活跃迭代名称
	DaysRemaining   int     `json:"days_remaining"`    // 项目剩余天数
	BudgetUsed      float64 `json:"budget_used"`       // 已使用预算
	BudgetRate      float64 `json:"budget_rate"`       // 预算使用率
}

// MemberStatistics 成员统计
type MemberStatistics struct {
	TotalMembers int            `json:"total_members"`
	ByRole       map[string]int `json:"by_role"` // 按角色分组统计
}

// MilestoneStats 里程碑统计
type MilestoneStats struct {
	Total          int     `json:"total"`
	Completed      int     `json:"completed"`
	Pending        int     `json:"pending"`
	Overdue        int     `json:"overdue"`
	CompletionRate float64 `json:"completion_rate"`
}

// RiskStatistics 风险统计
type RiskStatistics struct {
	Total      int `json:"total"`
	Open       int `json:"open"`
	Mitigating int `json:"mitigating"`
	Closed     int `json:"closed"`
	HighRisk   int `json:"high_risk"` // 高风险数量（严重程度>=12）
}

// IterationStatistics 迭代统计
type IterationStatistics struct {
	Total       int     `json:"total"`
	Completed   int     `json:"completed"`
	Active      int     `json:"active"`
	AvgVelocity float64 `json:"avg_velocity"` // 平均速度（故事点/迭代）
}

// ==================== 通用分页参数 ====================

// PaginationQuery 分页查询参数
type PaginationQuery struct {
	Page      int    `form:"page"`       // 页码，默认1
	Size      int    `form:"size"`       // 每页条数，默认20
	Keyword   string `form:"keyword"`    // 搜索关键词
	Status    string `form:"status"`     // 状态筛选
	Type      string `form:"type"`       // 类型筛选
	SortBy    string `form:"sort_by"`    // 排序字段
	SortOrder string `form:"sort_order"` // 排序方向：asc/desc
}

// ListRequest 通用列表请求（用于解析查询参数）
type ListRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Size      int    `form:"size" binding:"omitempty,min=1,max=100"`
	Keyword   string `form:"keyword"`
	Status    string `form:"status"`
	Type      string `form:"type"`
	ProgramID string `form:"program_id"`
	PMID      string `form:"pm_id"`
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order"`
}
