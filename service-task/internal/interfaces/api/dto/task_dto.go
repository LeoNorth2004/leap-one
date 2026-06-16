package dto

import (
	"time"

	"github.com/google/uuid"
)

// ==================== 任务相关DTO ====================

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Title          string     `json:"title" binding:"required,max=500"`
	Description    string     `json:"description"`
	Type           string     `json:"type"` // task/sub_task/bug_fix/design/review
	ProjectID      *uuid.UUID `json:"project_id"`
	IterationID    *uuid.UUID `json:"iteration_id"`
	RequirementID  *uuid.UUID `json:"requirement_id"`
	ParentID       *uuid.UUID `json:"parent_id"`
	AssigneeID     *uuid.UUID `json:"assignee_id"`
	CreatorID      uuid.UUID  `json:"creator_id"`
	Priority       int        `json:"priority" binding:"min=1,max=5"`
	Severity       int        `json:"severity" binding:"min=1,max=4"`
	StoryPoints    *float64   `json:"story_points"`
	EstimatedHours *float64   `json:"estimated_hours"`
	StartDate      *time.Time `json:"start_date"`
	DueDate        *time.Time `json:"due_date"`
	KanbanColumn   string     `json:"kanban_column"`
	Tags           []string   `json:"tags"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	Title          *string     `json:"title" binding:"omitempty,max=500"`
	Description    *string     `json:"description"`
	Type           *string     `json:"type"`
	ProjectID      *uuid.UUID `json:"project_id"`
	IterationID    *uuid.UUID `json:"iteration_id"`
	RequirementID  *uuid.UUID `json:"requirement_id"`
	AssigneeID     *uuid.UUID `json:"assignee_id"`
	Priority       *int       `json:"priority" binding:"omitempty,min=1,max=5"`
	Severity       *int       `json:"severity" binding:"omitempty,min=1,max=4"`
	StoryPoints    *float64   `json:"story_points"`
	EstimatedHours *float64   `json:"estimated_hours"`
	ActualHours    *float64   `json:"actual_hours"`
	RemainingHours *float64   `json:"remaining_hours"`
	StartDate      *time.Time `json:"start_date"`
	DueDate        *time.Time `json:"due_date"`
	KanbanColumn   *string    `json:"kanban_column"`
	Tags           []string   `json:"tags"`
}

// TaskInfo 任务基本信息
type TaskInfo struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	Description    string     `json:"description,omitempty"`
	Type           string     `json:"type"`
	ProjectID      string     `json:"project_id,omitempty"`
	IterationID    string     `json:"iteration_id,omitempty"`
	RequirementID  string     `json:"requirement_id,omitempty"`
	ParentID       string     `json:"parent_id,omitempty"`
	AssigneeID     string     `json:"assignee_id,omitempty"`
	CreatorID      string     `json:"creator_id"`
	Status         string     `json:"status"`
	Priority       int        `json:"priority"`
	Severity       int        `json:"severity"`
	StoryPoints    *float64   `json:"story_points,omitempty"`
	EstimatedHours *float64   `json:"estimated_hours,omitempty"`
	ActualHours    *float64   `json:"actual_hours,omitempty"`
	RemainingHours *float64   `json:"remaining_hours,omitempty"`
	StartDate      string     `json:"start_date,omitempty"`
	DueDate        string     `json:"due_date,omitempty"`
	FinishedDate   string     `json:"finished_date,omitempty"`
	KanbanColumn   string     `json:"kanban_column"`
	Tags           []string   `json:"tags,omitempty"`
	CreatedAt      string     `json:"created_at"`
	UpdatedAt      string     `json:"updated_at"`
}

// TaskDetailResponse 任务详情响应
type TaskDetailResponse struct {
	TaskInfo
	Assignees   []AssignmentInfo `json:"assignees,omitempty"`
	Comments    []CommentInfo    `json:"comments,omitempty"`
	Attachments []AttachmentInfo `json:"attachments,omitempty"`
}

// TaskListResponse 任务列表响应（分页）
type TaskListResponse struct {
	List  []TaskInfo `json:"list"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

// AssignTaskRequest 分配任务请求
type AssignTaskRequest struct {
	UserIDs []string `json:"user_ids" binding:"required,min=1"`
	Role    string   `json:"role"` // assignee/reviewer
}

// AssignmentInfo 分配信息
type AssignmentInfo struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	UserName   string `json:"user_name,omitempty"`
	Role       string `json:"role"`
	AssignedAt string `json:"assigned_at"`
}

// CreateCommentRequest 创建评论请求（任务/工单通用）
type CreateCommentRequest struct {
	Content  string     `json:"content" binding:"required"`
	ParentID *uuid.UUID `json:"parent_id"`
}

// CommentInfo 评论信息
type CommentInfo struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name,omitempty"`
	Content   string `json:"content"`
	ParentID  string `json:"parent_id,omitempty"`
	CreatedAt string `json:"created_at"`
}

// UploadAttachmentRequest 上传附件请求
type UploadAttachmentRequest struct {
	FileName string `json:"file_name" binding:"required,max=255"`
	FileSize int64  `json:"file_size"`
	FileType string `json:"file_type"`
	FileURL  string `json:"file_url" binding:"required,url,max=500"`
}

// AttachmentInfo 附件信息
type AttachmentInfo struct {
	ID         string `json:"id"`
	FileName   string `json:"file_name"`
	FileSize   int64  `json:"file_size"`
	FileType   string `json:"file_type"`
	FileURL    string `json:"file_url"`
	UploadedBy string `json:"uploaded_by"`
	CreatedAt  string `json:"created_at"`
}

// CreateWorklogRequest 添加工作日志请求
type CreateWorklogRequest struct {
	SpentHours float64 `json:"spent_hours" binding:"required,gt=0"`
	WorkDate   string  `json:"work_date" binding:"required"` // RFC3339格式
	Summary    string  `json:"summary" binding:"max=500"`
}

// WorklogInfo 工作日志信息
type WorklogInfo struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	UserName   string  `json:"user_name,omitempty"`
	SpentHours float64 `json:"spent_hours"`
	WorkDate   string  `json:"work_date"`
	Summary    string  `json:"summary,omitempty"`
	CreatedAt  string  `json:"created_at"`
}

// CreateSubTaskRequest 创建子任务请求
type CreateSubTaskRequest struct {
	Title       string  `json:"title" binding:"required,max=500"`
	Description string  `json:"description"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	Priority    int     `json:"priority" binding:"min=1,max=5"`
}

// CreateTaskLinkRequest 创建任务关联请求
type CreateTaskLinkRequest struct {
	TargetTaskID uuid.UUID `json:"target_task_id" binding:"required"`
	LinkType     string    `json:"link_type" binding:"required,oneof=blocks is_blocked_by relates_to duplicates"`
	Name         string    `json:"name"`
}

// TaskLinkInfo 任务关联信息
type TaskLinkInfo struct {
	ID           string `json:"id"`
	SourceTaskID string `json:"source_task_id"`
	TargetTaskID string `json:"target_task_id"`
	LinkType     string `json:"link_type"`
	Name         string `json:"name,omitempty"`
	CreatedAt    string `json:"created_at"`
}

// UpdateStatusRequest 更改状态请求
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
