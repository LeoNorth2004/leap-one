package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Issue 工单/事项实体 - 工单管理核心模型
// 支持多种工单类型、SLA管理、工作流状态流转、满意度评价等
type Issue struct {
	ID              uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title           string           `gorm:"size:500;not null" json:"title"`
	Description     string           `gorm:"type:text" json:"description"`
	Type            string           `gorm:"size:30;default:'bug';index" json:"type"` // bug/feature/request/incident/question
	ProjectID       *uuid.UUID       `gorm:"type:uuid;index" json:"project_id,omitempty"`
	ProductID       *uuid.UUID       `gorm:"type:uuid;index" json:"product_id,omitempty"`
	ReporterID      uuid.UUID        `gorm:"type:uuid;not null;index" json:"reporter_id"`
	AssigneeID      *uuid.UUID       `gorm:"type:uuid;index" json:"assignee_id,omitempty"`
	Status          string           `gorm:"size:20;default:'new';index" json:"status"` // new/in_progress/waiting/resolved/closed/cancelled
	Priority        int              `gorm:"default:3;index" json:"priority"`           // 1-5
	Severity        int              `gorm:"default:2" json:"severity"`                 // 1-4
	Source          string           `gorm:"size:30;default:'manual'" json:"source"`     // manual/email/api/webhook
	TemplateID      *uuid.UUID       `gorm:"type:uuid;index" json:"template_id,omitempty"`
	SLADueDate      *time.Time       `json:"sla_due_date,omitempty"`
	ResponseDueDate *time.Time       `json:"response_due_date,omitempty"`
	Satisfaction    *int             `json:"satisfaction,omitempty"` // 满意度评价(1-5)
	Resolution      string           `gorm:"type:text" json:"resolution,omitempty"`
	ResolvedAt      *time.Time       `json:"resolved_at,omitempty"`
	ResolvedBy      *uuid.UUID       `json:"resolved_by,omitempty"`
	ClosedAt        *time.Time       `json:"closed_at,omitempty"`
	ClosedBy        *uuid.UUID       `json:"closed_by,omitempty"`
	Tags            string           `gorm:"type:text" json:"tags,omitempty"` // JSON数组标签
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       gorm.DeletedAt   `gorm:"index" json:"-"`
	Comments        []IssueComment   `gorm:"foreignKey:IssueID" json:"comments,omitempty"`
	Attachments     []IssueAttachment `gorm:"foreignKey:IssueID" json:"attachments,omitempty"`
}

// TableName 指定数据库表名
func (Issue) TableName() string {
	return "issues"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (i *Issue) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// IssueComment 工单评论
type IssueComment struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	IssueID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"issue_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	IsInternal bool       `gorm:"default:false" json:"is_internal"` // 内部备注（客户不可见）
	ParentID  *uuid.UUID `gorm:"type:uuid;index" json:"parent_id,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (IssueComment) TableName() string { return "issue_comments" }

func (ic *IssueComment) BeforeCreate(tx *gorm.DB) error {
	if ic.ID == uuid.Nil {
		ic.ID = uuid.New()
	}
	return nil
}

// IssueAttachment 工单附件
type IssueAttachment struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	IssueID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"issue_id"`
	FileName   string         `gorm:"size:255;not null" json:"file_name"`
	FileSize   int64          `json:"file_size"`
	FileType   string         `gorm:"size:100" json:"file_type"`
	FileURL    string         `gorm:"size:500;not null" json:"file_url"`
	UploadedBy uuid.UUID      `gorm:"type:uuid" json:"uploaded_by"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (IssueAttachment) TableName() string { return "issue_attachments" }

func (ia *IssueAttachment) BeforeCreate(tx *gorm.DB) error {
	if ia.ID == uuid.Nil {
		ia.ID = uuid.New()
	}
	return nil
}

// IssueTemplate 工单模板 - 预定义的工单字段模板
type IssueTemplate struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string     `gorm:"size:200;not null" json:"name"`
	Type      string     `gorm:"size:30;index" json:"type"` // bug/feature/request等
	Fields    string     `gorm:"type:text" json:"fields"`   // JSON模板字段配置
	WorkflowID *uuid.UUID `gorm:"type:uuid;index" json:"workflow_id,omitempty"`
	IsSystem  bool       `gorm:"default:false" json:"is_system"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (IssueTemplate) TableName() string { return "issue_templates" }

func (it *IssueTemplate) BeforeCreate(tx *gorm.DB) error {
	if it.ID == uuid.Nil {
		it.ID = uuid.New()
	}
	return nil
}

// IssueWorkflow 工单工作流定义
type IssueWorkflow struct {
	ID            uuid.UUID                  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name          string                     `gorm:"size:200;not null" json:"name"`
	Type          string                     `gorm:"size:30;index" json:"type"` // 适用类型
	InitialStatus string                    `gorm:"size:20;default:'new'" json:"initial_status"`
	Description   string                    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at"`
	Transitions   []IssueWorkflowTransition `gorm:"foreignKey:WorkflowID" json:"transitions,omitempty"`
}

func (IssueWorkflow) TableName() string { return "issue_workflows" }

func (iw *IssueWorkflow) BeforeCreate(tx *gorm.DB) error {
	if iw.ID == uuid.Nil {
		iw.ID = uuid.New()
	}
	return nil
}

// IssueWorkflowTransition 工作流状态转换规则
type IssueWorkflowTransition struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkflowID uuid.UUID `gorm:"type:uuid;index;not null" json:"workflow_id"`
	FromStatus string    `gorm:"size:20;not null" json:"from_status"`
	ToStatus   string    `gorm:"size:20;not null" json:"to_status"`
	Condition  string    `gorm:"size:200" json:"condition,omitempty"` // 条件表达式
	Name       string    `gorm:"size:100" json:"name"`               // 转换名称
	SortOrder  int       `gorm:"default:0" json:"sort_order"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (IssueWorkflowTransition) TableName() string { return "issue_workflow_transitions" }

func (iwt *IssueWorkflowTransition) BeforeCreate(tx *gorm.DB) error {
	if iwt.ID == uuid.Nil {
		iwt.ID = uuid.New()
	}
	return nil
}

// IssueSLAConfig SLA配置 - 按类型和优先级组合定义响应和解决时限
type IssueSLAConfig struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Type             string    `gorm:"size:30;not null;index" json:"type"` // 工单类型
	Priority         int       `gorm:"not null" json:"priority"`           // 优先级
	ResponseSLA      int       `json:"response_sla"`                       // 响应SLA（分钟）
	ResolveSLA       int       `json:"resolve_sla"`                        // 解决SLA（分钟）
	BusinessHoursOnly bool    `gorm:"default:false" json:"business_hours_only"` // 仅工作时间
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (IssueSLAConfig) TableName() string { return "issue_sla_configs" }

func (isc *IssueSLAConfig) BeforeCreate(tx *gorm.DB) error {
	if isc.ID == uuid.Nil {
		isc.ID = uuid.New()
	}
	return nil
}
