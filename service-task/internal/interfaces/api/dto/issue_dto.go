package dto

import (
	"github.com/google/uuid"
)

// ==================== 工单相关DTO ====================

// CreateIssueRequest 创建工单请求
type CreateIssueRequest struct {
	Title       string     `json:"title" binding:"required,max=500"`
	Description string     `json:"description"`
	Type        string     `json:"type"` // bug/feature/request/incident/question
	ProjectID   *uuid.UUID `json:"project_id"`
	ProductID   *uuid.UUID `json:"product_id"`
	ReporterID  uuid.UUID  `json:"reporter_id"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	Priority    int        `json:"priority" binding:"min=1,max=5"`
	Severity    int        `json:"severity" binding:"min=1,max=4"`
	Source      string     `json:"source"`
	TemplateID  *uuid.UUID `json:"template_id"`
	Tags        []string   `json:"tags"`
}

// UpdateIssueRequest 更新工单请求
type UpdateIssueRequest struct {
	Title       *string    `json:"title" binding:"omitempty,max=500"`
	Description *string    `json:"description"`
	Type        *string    `json:"type"`
	ProjectID   *uuid.UUID `json:"project_id"`
	ProductID   *uuid.UUID `json:"product_id"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	Priority    *int       `json:"priority" binding:"omitempty,min=1,max=5"`
	Severity    *int       `json:"severity" binding:"omitempty,min=1,max=4"`
	Resolution  *string    `json:"resolution"`
	Tags        []string   `json:"tags"`
}

// IssueInfo 工单基本信息
type IssueInfo struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Description     string   `json:"description,omitempty"`
	Type            string   `json:"type"`
	ProjectID       string   `json:"project_id,omitempty"`
	ProductID       string   `json:"product_id,omitempty"`
	ReporterID      string   `json:"reporter_id"`
	AssigneeID      string   `json:"assignee_id,omitempty"`
	Status          string   `json:"status"`
	Priority        int      `json:"priority"`
	Severity        int      `json:"severity"`
	Source          string   `json:"source"`
	TemplateID      string   `json:"template_id,omitempty"`
	SLADueDate      string   `json:"sla_due_date,omitempty"`
	ResponseDueDate string   `json:"response_due_date,omitempty"`
	Satisfaction    *int     `json:"satisfaction,omitempty"`
	Resolution      string   `json:"resolution,omitempty"`
	ResolvedAt      string   `json:"resolved_at,omitempty"`
	ResolvedBy      string   `json:"resolved_by,omitempty"`
	ClosedAt        string   `json:"closed_at,omitempty"`
	ClosedBy        string   `json:"closed_by,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
}

// IssueDetailResponse 工单详情响应
type IssueDetailResponse struct {
	IssueInfo
	Comments    []IssueCommentInfo `json:"comments,omitempty"`
	Attachments []AttachmentInfo   `json:"attachments,omitempty"`
}

// IssueListResponse 工单列表响应（分页）
type IssueListResponse struct {
	List  []IssueInfo `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// TransitionRequest 状态流转请求
type TransitionRequest struct {
	Status     string     `json:"status" binding:"required"`
	Comment    string     `json:"comment"`
	ResolvedBy *uuid.UUID `json:"resolved_by"`
}

// IssueCommentInfo 工单评论信息
type IssueCommentInfo struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	UserName   string `json:"user_name,omitempty"`
	Content    string `json:"content"`
	IsInternal bool   `json:"is_internal"`
	ParentID   string `json:"parent_id,omitempty"`
	CreatedAt  string `json:"created_at"`
}

// CreateIssueCommentRequest 创建工单评论请求
type CreateIssueCommentRequest struct {
	Content    string     `json:"content" binding:"required"`
	IsInternal bool       `json:"is_internal"`
	ParentID   *uuid.UUID `json:"parent_id"`
}

// SatisfactionRequest 满意度评价请�?
type SatisfactionRequest struct {
	Score   int    `json:"score" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

// SLAInfo SLA信息响应
type SLAInfo struct {
	ResponseSLA       int    `json:"response_sla"`        // 响应SLA（分钟）
	ResolveSLA        int    `json:"resolve_sla"`         // 解决SLA（分钟）
	SLADueDate        string `json:"sla_due_date"`        // SLA截止时间
	ResponseDueDate   string `json:"response_due_date"`   // 响应截止时间
	IsOverdue         bool   `json:"is_overdue"`          // 是否已超时
	ResponseOverdue   bool   `json:"response_overdue"`    // 响应是否已超时
	RemainingMinutes  *int64 `json:"remaining_minutes"`   // 剩余分钟数
	BusinessHoursOnly bool   `json:"business_hours_only"` // 仅工作时间
}
