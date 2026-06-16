package dto

import "github.com/google/uuid"

// ========== 测试用例相关DTO ==========

// CreateTestCaseRequest 创建测试用例请求
type CreateTestCaseRequest struct {
	Title          string     `json:"title" binding:"required,max=500"`                                 // 用例标题
	Module         string     `json:"module" binding:"max=200"`                                         // 所属模�?
	Precondition   string     `json:"precondition"`                                                     // 前置条件
	Steps          string     `json:"steps"`                                                            // JSON数组测试步骤
	ExpectedResult string     `json:"expected_result"`                                                  // 预期结果
	Priority       int        `json:"priority" binding:"min=1,max=5"`                                   // 优先�?1-5
	Type           string     `json:"type" binding:"omitempty,oneof=manual automated smoke regression"` // 用例类型
	Automation     bool       `json:"automation"`                                                       // 是否自动�?
	ProductID      *uuid.UUID `json:"product_id"`                                                       // 产品ID
	ProjectID      *uuid.UUID `json:"project_id"`                                                       // 项目ID
	RequirementID  *uuid.UUID `json:"requirement_id"`                                                   // 需求ID
	Tags           string     `json:"tags"`                                                             // 标签
}

// UpdateTestCaseRequest 更新测试用例请求
type UpdateTestCaseRequest struct {
	Title          *string    `json:"title" binding:"omitempty,max=500"`
	Module         *string    `json:"module" binding:"omitempty,max=200"`
	Precondition   *string    `json:"precondition"`
	Steps          *string    `json:"steps"`
	ExpectedResult *string    `json:"expected_result"`
	Priority       *int       `json:"priority" binding:"omitempty,min=1,max=5"`
	Type           *string    `json:"type" binding:"omitempty,oneof=manual automated smoke regression"`
	Status         *string    `json:"status" binding:"omitempty,oneof=draft active archived"`
	Automation     *bool      `json:"automation"`
	ProductID      *uuid.UUID `json:"product_id"`
	ProjectID      *uuid.UUID `json:"project_id"`
	RequirementID  *uuid.UUID `json:"requirement_id"`
	Tags           *string    `json:"tags"`
}

// TestCaseInfo 测试用例简要信�?
type TestCaseInfo struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	Module     string  `json:"module"`
	Priority   int     `json:"priority"`
	Type       string  `json:"type"`
	Status     string  `json:"status"`
	Automation bool    `json:"automation"`
	ProductID  *string `json:"product_id,omitempty"`
	ProjectID  *string `json:"project_id,omitempty"`
	CreatorID  string  `json:"creator_id"`
	Version    int     `json:"version"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

// TestCaseDetailResponse 测试用例详情响应
type TestCaseDetailResponse struct {
	TestCaseInfo
	Precondition   string  `json:"precondition"`
	Steps          string  `json:"steps"`
	ExpectedResult string  `json:"expected_result"`
	RequirementID  *string `json:"requirement_id,omitempty"`
	ReviewerID     *string `json:"reviewer_id,omitempty"`
	ReviewedAt     *string `json:"reviewed_at,omitempty"`
	Tags           string  `json:"tags"`
}

// TestCaseListResponse 测试用例列表响应（分页）
type TestCaseListResponse struct {
	List  []TestCaseInfo `json:"list"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

// ImportTestCasesRequest 导入用例请求
type ImportTestCasesRequest struct {
	Cases []CreateTestCaseRequest `json:"cases" binding:"required,dive"` // 用例列表
}

// ReviewTestCaseRequest 评审用例请求
type ReviewTestCaseRequest struct {
	ReviewerID uuid.UUID `json:"reviewer_id" binding:"required"` // 评审人ID
}
