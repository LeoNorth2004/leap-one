package dto

import "github.com/google/uuid"

// ========== 测试套件相关DTO ==========

// CreateTestSuiteRequest 创建测试套件请求
type CreateTestSuiteRequest struct {
	Name        string     `json:"name" binding:"required,max=200"` // 套件名称
	Description string     `json:"description"`                     // 套件描述
	ProductID   *uuid.UUID `json:"product_id"`                      // 产品ID
	ProjectID   *uuid.UUID `json:"project_id"`                      // 项目ID
}

// UpdateTestSuiteRequest 更新测试套件请求
type UpdateTestSuiteRequest struct {
	Name        *string    `json:"name" binding:"omitempty,max=200"`
	Description *string    `json:"description"`
	ProductID   *uuid.UUID `json:"product_id"`
	ProjectID   *uuid.UUID `json:"project_id"`
}

// TestSuiteInfo 测试套件简要信�?
type TestSuiteInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatorID   string `json:"creator_id"`
	CaseCount   int    `json:"case_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// TestSuiteDetailResponse 测试套件详情响应（含用例列表�?
type TestSuiteDetailResponse struct {
	TestSuiteInfo
	Cases []SuiteCaseItem `json:"cases"`
}

// SuiteCaseItem 套件中的用例�?
type SuiteCaseItem struct {
	CaseID    string `json:"case_id"`
	CaseTitle string `json:"case_title"`
	SortOrder int    `json:"sort_order"`
}

// TestSuiteListResponse 测试套件列表响应（分页）
type TestSuiteListResponse struct {
	List  []TestSuiteInfo `json:"list"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
}

// AddCasesToSuiteRequest 添加用例到套件请�?
type AddCasesToSuiteRequest struct {
	CaseIDs []uuid.UUID `json:"case_ids" binding:"required,min=1"` // 要添加的用例ID列表
}

// ========== 测试计划相关DTO ==========

// CreateTestPlanRequest 创建测试计划请求
type CreateTestPlanRequest struct {
	Name         string     `json:"name" binding:"required,max=200"` // 计划名称
	Description  string     `json:"description"`                     // 计划描述
	ProductID    *uuid.UUID `json:"product_id"`                      // 产品ID
	ProjectID    *uuid.UUID `json:"project_id"`                      // 项目ID
	BuildVersion string     `json:"build_version" binding:"max=100"` // 构建版本
	StartDate    string     `json:"start_date"`                      // 开始日期（YYYY-MM-DD�?
	EndDate      string     `json:"end_date"`                        // 结束日期（YYYY-MM-DD�?
	ExecutorIDs  string     `json:"executor_ids"`                    // JSON数组执行人ID
}

// UpdateTestPlanRequest 更新测试计划请求
type UpdateTestPlanRequest struct {
	Name         *string `json:"name" binding:"omitempty,max=200"`
	Description  *string `json:"description"`
	BuildVersion *string `json:"build_version" binding:"omitempty,max=100"`
	StartDate    *string `json:"start_date"`
	EndDate      *string `json:"end_date"`
	ExecutorIDs  *string `json:"executor_ids"`
}

// TestPlanInfo 测试计划简要信�?
type TestPlanInfo struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	BuildVersion string  `json:"build_version"`
	Status       string  `json:"status"`
	StartDate    *string `json:"start_date"`
	EndDate      *string `json:"end_date"`
	CreatorID    string  `json:"creator_id"`
	CaseCount    int     `json:"case_count"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// TestPlanDetailResponse 测试计划详情响应
type TestPlanDetailResponse struct {
	TestPlanInfo
	ExecutorIDs string             `json:"executor_ids"`
	Cases       []TestPlanCaseItem `json:"cases"`
}

// TestPlanCaseItem 计划中的用例执行�?
type TestPlanCaseItem struct {
	PlanCaseID   uuid.UUID `json:"plan_case_id"`
	CaseID       uuid.UUID `json:"case_id"`
	CaseTitle    string    `json:"case_title"`
	AssigneeID   *string   `json:"assignee_id,omitempty"`
	Result       string    `json:"result"`
	ExecuteTime  *string   `json:"execute_time,omitempty"`
	ActualResult string    `json:"actual_result"`
	SortOrder    int       `json:"sort_order"`
}

// TestPlanListResponse 测试计划列表响应（分页）
type TestPlanListResponse struct {
	List  []TestPlanInfo `json:"list"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

// ExecuteTestCaseRequest 执行用例请求
type ExecuteTestCaseRequest struct {
	AssigneeID   *uuid.UUID `json:"assignee_id"`                                                           // 执行人ID
	Result       string     `json:"result" binding:"required,oneof=not_run passed failed blocked skipped"` // 执行结果
	ActualResult string     `json:"actual_result"`                                                         // 实际结果描述
	BugIDs       string     `json:"bug_ids"`                                                               // 关联Bug ID列表（JSON数组�?
	Comment      string     `json:"comment"`                                                               // 执行备注
}

// AddCasesToPlanRequest 添加用例到计划请�?
type AddCasesToPlanRequest struct {
	CaseIDs []uuid.UUID `json:"case_ids" binding:"required,min=1"` // 要添加的用例ID列表
}
