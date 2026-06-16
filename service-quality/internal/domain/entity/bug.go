package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Bug 缺陷实体 - 质量管理核心跟踪对象
type Bug struct {
	ID            uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	Title         string          `gorm:"type:varchar(500);not null" json:"title"`       // Bug标题
	Description   string          `gorm:"type:text" json:"description"`                  // Bug详细描述
	Steps         string          `gorm:"type:text" json:"steps"`                        // 复现步骤
	Severity      int             `gorm:"default:2" json:"severity"`                     // 严重程度 1致命/2严重/3一�?4提示
	Priority      int             `gorm:"default:3" json:"priority"`                     // 优先�?1-5
	Status        string          `gorm:"size:20;default:'new'" json:"status"`           // new/confirmed/in_progress/resolved/closed/reopened/cancelled
	Type          string          `gorm:"size:30;default:'code_bug'" json:"type"`        // code_bug/design_bug/data_bug/config/security/performance/ui
	ProductID     *uuid.UUID      `gorm:"type:uuid" json:"product_id"`                   // 关联产品ID
	ProjectID     *uuid.UUID      `gorm:"type:uuid" json:"project_id"`                   // 关联项目ID
	IterationID   *uuid.UUID      `gorm:"type:uuid" json:"iteration_id"`                 // 关联迭代ID
	RequirementID *uuid.UUID      `gorm:"type:uuid" json:"requirement_id"`               // 关联需求ID
	TaskID        *uuid.UUID      `gorm:"type:uuid" json:"task_id"`                      // 关联任务ID
	TestCaseID    *uuid.UUID      `gorm:"type:uuid" json:"test_case_id"`                 // 关联用例ID
	ReporterID    uuid.UUID       `gorm:"type:uuid;not null" json:"reporter_id"`         // 提报人ID
	AssigneeID    *uuid.UUID      `gorm:"type:uuid" json:"assignee_id"`                  // 处理人ID
	Resolution    string          `gorm:"size:30" json:"resolution"`                     // 解决方案 fixed/wont_fix/duplicate/by_design/workaround/postponed
	FoundVersion  string          `gorm:"type:varchar(100)" json:"found_version"`        // 发现版本
	FixedVersion  string          `gorm:"type:varchar(100)" json:"fixed_version"`        // 修复版本
	Environment   string          `gorm:"type:varchar(200)" json:"environment"`          // 环境信息
	OS            string          `gorm:"type:varchar(100)" json:"os"`                   // 操作系统
	Browser       string          `gorm:"type:varchar(100)" json:"browser"`              // 浏览�?
	Reproductive  bool            `gorm:"default:true" json:"reproductive"`              // 是否可复�?
	ConfirmedAt   *time.Time      `json:"confirmed_at"`                                  // 确认时间
	ConfirmedBy   *uuid.UUID      `gorm:"type:uuid" json:"confirmed_by"`                 // 确认人ID
	ResolvedAt    *time.Time      `json:"resolved_at"`                                   // 解决时间
	ResolvedBy    *uuid.UUID      `gorm:"type:uuid" json:"resolved_by"`                  // 解决人ID
	ClosedAt      *time.Time      `json:"closed_at"`                                     // 关闭时间
	ClosedBy      *uuid.UUID      `gorm:"type:uuid" json:"closed_by"`                    // 关闭人ID
	Deadline      *time.Time      `json:"deadline"`                                      // 解决期限
	Tags          string          `gorm:"type:text" json:"tags"`                         // 标签
	Comments      []BugComment    `gorm:"foreignKey:BugID" json:"comments,omitempty"`    // 评论列表
	Attachments   []BugAttachment `gorm:"foreignKey:BugID" json:"attachments,omitempty"` // 附件列表
	History       []BugHistory    `gorm:"foreignKey:BugID" json:"history,omitempty"`     // 变更历史
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName 指定数据库表�?
func (Bug) TableName() string {
	return "bugs"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (b *Bug) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// BugComment Bug评论实体
type BugComment struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	BugID     uuid.UUID  `gorm:"type:uuid;index;not null" json:"bug_id"` // 关联Bug ID
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`      // 评论人ID
	Content   string     `gorm:"type:text;not null" json:"content"`      // 评论内容
	ParentID  *uuid.UUID `gorm:"type:uuid" json:"parent_id"`             // 父评论ID（用于回复）
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// TableName 指定数据库表�?
func (BugComment) TableName() string {
	return "bug_comments"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (c *BugComment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// BugAttachment Bug附件实体
type BugAttachment struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	BugID      uuid.UUID `gorm:"type:uuid;index;not null" json:"bug_id"`      // 关联Bug ID
	FileName   string    `gorm:"type:varchar(255);not null" json:"file_name"` // 文件名称
	FileSize   int64     `json:"file_size"`                                   // 文件大小（字节）
	FileType   string    `gorm:"type:varchar(100)" json:"file_type"`          // 文件类型/MIME
	FileURL    string    `gorm:"type:varchar(500);not null" json:"file_url"`  // 文件存储URL
	UploadedBy uuid.UUID `gorm:"type:uuid" json:"uploaded_by"`                // 上传人ID
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName 指定数据库表�?
func (BugAttachment) TableName() string {
	return "bug_attachments"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (a *BugAttachment) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// BugHistory Bug状态变更历史实�?
type BugHistory struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	BugID     uuid.UUID `gorm:"type:uuid;index;not null" json:"bug_id"`      // 关联Bug ID
	FieldName string    `gorm:"type:varchar(50);not null" json:"field_name"` // 变更字段�?
	OldValue  string    `gorm:"type:text" json:"old_value"`                  // 变更前的�?
	NewValue  string    `gorm:"type:text" json:"new_value"`                  // 变更后的�?
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`           // 操作人ID
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定数据库表�?
func (BugHistory) TableName() string {
	return "bug_histories"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (h *BugHistory) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

// BugWorkflow Bug工作流定义实�?
type BugWorkflow struct {
	ID            uuid.UUID               `gorm:"type:uuid;primary_key" json:"id"`
	Name          string                  `gorm:"type:varchar(200);not null" json:"name"`             // 工作流名�?
	InitialStatus string                  `gorm:"size:20;default:'new'" json:"initial_status"`        // 初始状�?
	IsDefault     bool                    `gorm:"default:false" json:"is_default"`                    // 是否默认工作�?
	Transitions   []BugWorkflowTransition `gorm:"foreignKey:WorkflowID" json:"transitions,omitempty"` // 状态转换规�?
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

// TableName 指定数据库表�?
func (BugWorkflow) TableName() string {
	return "bug_workflows"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (w *BugWorkflow) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

// BugWorkflowTransition Bug状态转换规则实�?
type BugWorkflowTransition struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	WorkflowID     uuid.UUID `gorm:"type:uuid;index;not null" json:"workflow_id"` // 所属工作流ID
	FromStatus     string    `gorm:"size:20;not null" json:"from_status"`         // 源状�?
	ToStatus       string    `gorm:"size:20;not null" json:"to_status"`           // 目标状�?
	Name           string    `gorm:"size:varchar(100)" json:"name"`               // 转换名称（如"确认Bug"�?开始处�?�?
	Condition      string    `gorm:"size:200" json:"condition"`                   // 转换条件描述
	RequiredFields string    `gorm:"type:text" json:"required_fields"`            // 必填字段（JSON�?
	SortOrder      int       `gorm:"default:0" json:"sort_order"`                 // 排序顺序
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName 指定数据库表�?
func (BugWorkflowTransition) TableName() string {
	return "bug_workflow_transitions"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (t *BugWorkflowTransition) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
