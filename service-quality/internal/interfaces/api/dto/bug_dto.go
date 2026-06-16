package dto

import "github.com/google/uuid"

// ========== Bug管理相关DTO ==========

// CreateBugRequest 创建Bug请求
type CreateBugRequest struct {
	Title         string     `json:"title" binding:"required,max=500"`                                                           // Bug标题
	Description   string     `json:"description"`                                                                                // 详细描述
	Steps         string     `json:"steps"`                                                                                      // 复现步骤
	Severity      int        `json:"severity" binding:"min=1,max=4"`                                                             // 严重程度 1-4
	Priority      int        `json:"priority" binding:"min=1,max=5"`                                                             // 优先�?1-5
	Type          string     `json:"type" binding:"omitempty,oneof=code_bug design_bug data_bug config security performance ui"` // Bug类型
	ProductID     *uuid.UUID `json:"product_id"`                                                                                 // 产品ID
	ProjectID     *uuid.UUID `json:"project_id"`                                                                                 // 项目ID
	IterationID   *uuid.UUID `json:"iteration_id"`                                                                               // 迭代ID
	RequirementID *uuid.UUID `json:"requirement_id"`                                                                             // 需求ID
	TaskID        *uuid.UUID `json:"task_id"`                                                                                    // 任务ID
	TestCaseID    *uuid.UUID `json:"test_case_id"`                                                                               // 关联用例ID
	AssigneeID    *uuid.UUID `json:"assignee_id"`                                                                                // 处理人ID
	FoundVersion  string     `json:"found_version" binding:"max=100"`                                                            // 发现版本
	FixedVersion  string     `json:"fixed_version" binding:"max=100"`                                                            // 修复版本
	Environment   string     `json:"environment" binding:"max=200"`                                                              // 环境信息
	OS            string     `json:"os" binding:"max=100"`                                                                       // 操作系统
	Browser       string     `json:"browser" binding:"max=100"`                                                                  // 浏览�?
	Reproductive  bool       `json:"reproductive"`                                                                               // 是否可复�?
	Deadline      string     `json:"deadline"`                                                                                   // 解决期限
	Tags          string     `json:"tags"`                                                                                       // 标签
}

// UpdateBugRequest 更新Bug请求
type UpdateBugRequest struct {
	Title        *string    `json:"title" binding:"omitempty,max=500"`
	Description  *string    `json:"description"`
	Steps        *string    `json:"steps"`
	Severity     *int       `json:"severity" binding:"omitempty,min=1,max=4"`
	Priority     *int       `json:"priority" binding:"omitempty,min=1,max=5"`
	Type         *string    `json:"type" binding:"omitempty,oneof=code_bug design_bug data_bug config security performance ui"`
	AssigneeID   *uuid.UUID `json:"assignee_id"`
	FoundVersion *string    `json:"found_version" binding:"omitempty,max=100"`
	FixedVersion *string    `json:"fixed_version" binding:"omitempty,max=100"`
	Environment  *string    `json:"environment" binding:"omitempty,max=200"`
	OS           *string    `json:"os" binding:"omitempty,max=100"`
	Browser      *string    `json:"browser" binding:"omitempty,max=100"`
	Reproductive *bool      `json:"reproductive"`
	Deadline     *string    `json:"deadline"`
	Tags         *string    `json:"tags"`
}

// ResolveBugRequest 解决Bug请求
type ResolveBugRequest struct {
	Resolution string `json:"resolution" binding:"required,oneof=fixed wont_fix duplicate by_design workaround postponed"` // 解决方案
}

// BugInfo Bug简要信�?
type BugInfo struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	Severity     int     `json:"severity"`
	Priority     int     `json:"priority"`
	Status       string  `json:"status"`
	Type         string  `json:"type"`
	ReporterID   string  `json:"reporter_id"`
	AssigneeID   *string `json:"assignee_id,omitempty"`
	Resolution   string  `json:"resolution"`
	FoundVersion string  `json:"found_version"`
	FixedVersion string  `json:"fixed_version"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// BugDetailResponse Bug详情响应（含评论、附件、历史）
type BugDetailResponse struct {
	BugInfo
	Description   string              `json:"description"`
	Steps         string              `json:"steps"`
	ProductID     *string             `json:"product_id,omitempty"`
	ProjectID     *string             `json:"project_id,omitempty"`
	IterationID   *string             `json:"iteration_id,omitempty"`
	RequirementID *string             `json:"requirement_id,omitempty"`
	TaskID        *string             `json:"task_id,omitempty"`
	TestCaseID    *string             `json:"test_case_id,omitempty"`
	Environment   string              `json:"environment"`
	OS            string              `json:"os"`
	Browser       string              `json:"browser"`
	Reproductive  bool                `json:"reproductive"`
	ConfirmedAt   *string             `json:"confirmed_at,omitempty"`
	ResolvedAt    *string             `json:"resolved_at,omitempty"`
	ClosedAt      *string             `json:"closed_at,omitempty"`
	Deadline      *string             `json:"deadline,omitempty"`
	Tags          string              `json:"tags"`
	Comments      []BugCommentInfo    `json:"comments"`
	Attachments   []BugAttachmentInfo `json:"attachments"`
	History       []BugHistoryInfo    `json:"history"`
}

// BugCommentInfo Bug评论信息
type BugCommentInfo struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Content   string  `json:"content"`
	ParentID  *string `json:"parent_id,omitempty"`
	CreatedAt string  `json:"created_at"`
}

// BugAttachmentInfo Bug附件信息
type BugAttachmentInfo struct {
	ID         string `json:"id"`
	FileName   string `json:"file_name"`
	FileSize   int64  `json:"file_size"`
	FileType   string `json:"file_type"`
	FileURL    string `json:"file_url"`
	UploadedBy string `json:"uploaded_by"`
	CreatedAt  string `json:"created_at"`
}

// BugHistoryInfo Bug变更历史信息
type BugHistoryInfo struct {
	ID        string `json:"id"`
	FieldName string `json:"field_name"`
	OldValue  string `json:"old_value"`
	NewValue  string `json:"new_value"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

// BugListResponse Bug列表响应（分页）
type BugListResponse struct {
	List  []BugInfo `json:"list"`
	Total int64     `json:"total"`
	Page  int       `json:"page"`
	Size  int       `json:"size"`
}

// AddBugCommentRequest 添加Bug评论请求
type AddBugCommentRequest struct {
	Content  string     `json:"content" binding:"required"` // 评论内容
	ParentID *uuid.UUID `json:"parent_id"`                  // 父评论ID（回复时使用�?
}

// UploadAttachmentRequest 上传附件请求
type UploadAttachmentRequest struct {
	FileName string `json:"file_name" binding:"required,max=255"` // 文件名称
	FileSize int64  `json:"file_size"`                            // 文件大小
	FileType string `json:"file_type" binding:"max=100"`          // 文件类型/MIME
	FileURL  string `json:"file_url" binding:"required,max=500"`  // 文件存储URL
}
