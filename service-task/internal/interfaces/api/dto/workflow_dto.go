package dto

// ==================== 工作流和SLA配置DTO ====================

// CreateWorkflowRequest 创建工作流请求
type CreateWorkflowRequest struct {
	Name          string `json:"name" binding:"required,max=200"`
	Type          string `json:"type" binding:"required,max=30"`
	InitialStatus string `json:"initial_status"`
	Description   string `json:"description"`
}

// UpdateWorkflowRequest 更新工作流请求
type UpdateWorkflowRequest struct {
	Name          *string `json:"name" binding:"omitempty,max=200"`
	Type          *string `json:"type" binding:"omitempty,max=30"`
	InitialStatus *string `json:"initial_status"`
	Description   *string `json:"description"`
}

// WorkflowInfo 工作流信息
type WorkflowInfo struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	Type          string           `json:"type"`
	InitialStatus string           `json:"initial_status"`
	Description   string           `json:"description,omitempty"`
	Transitions   []TransitionInfo `json:"transitions,omitempty"`
	CreatedAt     string           `json:"created_at"`
	UpdatedAt     string           `json:"updated_at"`
}

// WorkflowListResponse 工作流列表响�?
type WorkflowListResponse struct {
	List  []WorkflowInfo `json:"list"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

// CreateTransitionRequest 添加状态转换请求
type CreateTransitionRequest struct {
	FromStatus string `json:"from_status" binding:"required,max=20"`
	ToStatus   string `json:"to_status" binding:"required,max=20"`
	Condition  string `json:"condition" binding:"max=200"`
	Name       string `json:"name" binding:"required,max=100"`
	SortOrder  int    `json:"sort_order"`
}

// TransitionInfo 状态转换信息
type TransitionInfo struct {
	ID         string `json:"id"`
	FromStatus string `json:"from_status"`
	ToStatus   string `json:"to_status"`
	Condition  string `json:"condition,omitempty"`
	Name       string `json:"name"`
	SortOrder  int    `json:"sort_order"`
}

// CreateSLAConfigRequest 创建SLA配置请求
type CreateSLAConfigRequest struct {
	Type              string `json:"type" binding:"required,max=30"`
	Priority          int    `json:"priority" binding:"required,min=1,max=5"`
	ResponseSLA       int    `json:"response_sla"`
	ResolveSLA        int    `json:"resolve_sla"`
	BusinessHoursOnly bool   `json:"business_hours_only"`
}

// UpdateSLAConfigRequest 更新SLA配置请求
type UpdateSLAConfigRequest struct {
	Type              *string `json:"type" binding:"omitempty,max=30"`
	Priority          *int    `json:"priority" binding:"omitempty,min=1,max=5"`
	ResponseSLA       *int    `json:"response_sla"`
	ResolveSLA        *int    `json:"resolve_sla"`
	BusinessHoursOnly *bool   `json:"business_hours_only"`
}

// SLAConfigInfo SLA配置信息
type SLAConfigInfo struct {
	ID                string `json:"id"`
	Type              string `json:"type"`
	Priority          int    `json:"priority"`
	ResponseSLA       int    `json:"response_sla"`
	ResolveSLA        int    `json:"resolve_sla"`
	BusinessHoursOnly bool   `json:"business_hours_only"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

// SLAConfigListResponse SLA配置列表响应
type SLAConfigListResponse struct {
	List  []SLAConfigInfo `json:"list"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
}
