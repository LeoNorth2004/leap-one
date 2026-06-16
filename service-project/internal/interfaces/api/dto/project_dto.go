package dto

import (
	"time"

	"github.com/google/uuid"
)

// ==================== 项目相关DTO ====================

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	Name        string     `json:"name" binding:"required,min=1,max=200"` // 项目名称
	Code        string     `json:"code" binding:"required,min=1,max=50"`  // 项目编号
	ProgramID   *uuid.UUID `json:"program_id"`                             // 关联项目集ID
	Description string     `json:"description" binding:"max=2000"`        // 项目描述
	PMID        uuid.UUID  `json:"pm_id" binding:"required"`              // 项目经理ID
	Type        string     `json:"type" binding:"omitempty,oneof=agile waterfall lightweight lifecycle"` // 项目类型
	Priority    int        `json:"priority" binding:"omitempty,min=1,max=5"` // 优先级 1-5
	StartDate   *time.Time `json:"start_date"`                            // 开始日期
	EndDate     *time.Time `json:"end_date"`                              // 结束日期
	Budget      *float64   `json:"budget"`                                // 预算
	TemplateID  *uuid.UUID `json:"template_id"`                           // 模板ID
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	Name        *string    `json:"name" binding:"omitempty,min=1,max=200"`
	Description *string    `json:"description" binding:"omitempty,max=2000"`
	PMID        *uuid.UUID `json:"pm_id"`
	Type        *string    `json:"type" binding:"omitempty,oneof=agile waterfall lightweight lifecycle"`
	Priority    *int       `json:"priority" binding:"omitempty,min=1,max=5"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Budget      *float64   `json:"budget"`
}

// ProjectListResponse 项目列表响应（分页）
type ProjectListResponse struct {
	List  []ProjectInfo `json:"list"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
}

// ProjectInfo 项目简要信息
type ProjectInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	ProgramID   string  `json:"program_id,omitempty"`
	Description string  `json:"description,omitempty"`
	PMID        string  `json:"pm_id"`
	Status      string  `json:"status"`
	Type        string  `json:"type"`
	Priority    int     `json:"priority"`
	StartDate   string  `json:"start_date,omitempty"`
	EndDate     string  `json:"end_date,omitempty"`
	Budget      float64 `json:"budget,omitempty"`
	MemberCount int     `json:"member_count"`
	CreatedAt   string  `json:"created_at"`
}

// ProjectDetailResponse 项目详情响应
type ProjectDetailResponse struct {
	ProjectInfo
	TemplateID    string            `json:"template_id,omitempty"`
	CreatedByID   string            `json:"created_by_id"`
	UpdatedByID   string            `json:"updated_by_id,omitempty"`
	UpdatedAt     string            `json:"updated_at"`
	Version       int               `json:"version"`
	Members       []MemberInfo      `json:"members"`
	Milestones    []MilestoneInfo   `json:"milestones"`
	Risks         []RiskInfo        `json:"risks"`
	CustomFields  []CustomFieldInfo `json:"custom_fields"`
}

// ==================== 项目成员相关DTO ====================

// AddMemberRequest 添加成员请求
type AddMemberRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"` // 用户ID
	Role   string    `json:"role" binding:"required,oneof=pm po dev qa viewer ba sm"` // 角色
}

// UpdateMemberRequest 更新成员角色请求
type UpdateMemberRequest struct {
	Role string `json:"role" binding:"required,oneof=pm po dev qa viewer ba sm"` // 新角色
}

// MemberListResponse 成员列表响应
type MemberListResponse struct {
	List  []MemberInfo `json:"list"`
	Total int64        `json:"total"`
}

// MemberInfo 成员信息
type MemberInfo struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Role     string `json:"role"`
	JoinTime string `json:"join_time"`
}
