package dto

// ========== 测试环境相关DTO ==========

// CreateEnvironmentRequest 创建测试环境请求
type CreateEnvironmentRequest struct {
	Name        string `json:"name" binding:"required,max=200"`                      // 环境名称
	URL         string `json:"url" binding:"max=500"`                                // 环境访问地址
	Type        string `json:"type" binding:"omitempty,oneof=dev test staging prod"` // 环境类型
	OS          string `json:"os" binding:"max=100"`                                 // 操作系统
	Browser     string `json:"browser" binding:"max=100"`                            // 默认浏览�?
	Description string `json:"description"`                                          // 环境描述
	IsActive    *bool  `json:"is_active"`                                            // 是否启用
}

// UpdateEnvironmentRequest 更新测试环境请求
type UpdateEnvironmentRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=200"`
	URL         *string `json:"url" binding:"omitempty,max=500"`
	Type        *string `json:"type" binding:"omitempty,oneof=dev test staging prod"`
	OS          *string `json:"os" binding:"omitempty,max=100"`
	Browser     *string `json:"browser" binding:"omitempty,max=100"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

// EnvironmentInfo 测试环境信息
type EnvironmentInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	OS          string `json:"os"`
	Browser     string `json:"browser"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// EnvironmentListResponse 测试环境列表响应
type EnvironmentListResponse struct {
	List []EnvironmentInfo `json:"list"`
}

// ========== 报表统计相关DTO ==========

// QualityStatisticsResponse 质量统计概览响应
type QualityStatisticsResponse struct {
	TotalBugs      int64              `json:"total_bugs"`       // Bug总数
	OpenBugs       int64              `json:"open_bugs"`        // 未关闭Bug�?
	ResolvedRate   float64            `json:"resolved_rate"`    // 解决�?
	AvgResolveDays float64            `json:"avg_resolve_days"` // 平均解决天数
	ByStatus       map[string]int64   `json:"by_status"`        // 按状态分�?
	BySeverity     map[int]int64      `json:"by_severity"`      // 按严重程度分�?
	ByPriority     map[int]int64      `json:"by_priority"`      // 按优先级分布
	ByType         map[string]int64   `json:"by_type"`          // 按类型分�?
	TestCaseStats  TestCaseStatistics `json:"test_case_stats"`  // 用例统计
	TestPlanStats  TestPlanStatistics `json:"test_plan_stats"`  // 计划统计
}

// TestCaseStatistics 用例统计数据
type TestCaseStatistics struct {
	TotalCount    int64 `json:"total_count"`    // 总数
	ActiveCount   int64 `json:"active_count"`   // 活跃�?
	DraftCount    int64 `json:"draft_count"`    // 草稿�?
	ArchivedCount int64 `json:"archived_count"` // 归档�?
	AutoCount     int64 `json:"auto_count"`     // 自动化数�?
}

// TestPlanStatistics 计划统计数据
type TestPlanStatistics struct {
	TotalCount    int64 `json:"total_count"`    // 总数
	PlanningCount int64 `json:"planning_count"` // 规划�?
	ExecutingCnt  int64 `json:"executing_cnt"`  // 执行�?
	CompletedCnt  int64 `json:"completed_cnt"`  // 已完�?
	CancelledCnt  int64 `json:"cancelled_cnt"`  // 已取�?
}

// BugTrendItem Bug趋势数据�?
type BugTrendItem struct {
	Date     string `json:"date"`     // 日期
	Created  int64  `json:"created"`  // 新增�?
	Resolved int64  `json:"resolved"` // 解决�?
	Reopened int64  `json:"reopened"` // 重开�?
}

// BugTrendResponse Bug趋势分析响应
type BugTrendResponse struct {
	Trends []BugTrendItem `json:"trends"`
}

// PassRateItem 通过率数据项
type PassRateItem struct {
	PlanName   string  `json:"plan_name"`   // 计划名称
	TotalCases int     `json:"total_cases"` // 总用例数
	Passed     int     `json:"passed"`      // 通过�?
	Failed     int     `json:"failed"`      // 失败�?
	Blocked    int     `json:"blocked"`     // 阻塞�?
	Skipped    int     `json:"skipped"`     // 跳过�?
	NotRun     int     `json:"not_run"`     // 未执�?
	PassRate   float64 `json:"pass_rate"`   // 通过率百分比
}

// PassRateResponse 通过率统计响�?
type PassRateResponse struct {
	OverallPassRate float64        `json:"overall_pass_rate"` // 整体通过�?
	Plans           []PassRateItem `json:"plans"`             // 各计划通过�?
}
