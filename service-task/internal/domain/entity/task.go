package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Task 任务实体 - 项目任务核心模型
// 支持子任务、多人分配、评论、附件、工作日志、任务关联等完整功能
type Task struct {
	ID             uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title          string           `gorm:"size:500;not null" json:"title"`
	Description    string           `gorm:"type:text" json:"description"`
	Type           string           `gorm:"size:30;default:'task'" json:"type"` // task/sub_task/bug_fix/design/review
	ProjectID      *uuid.UUID       `gorm:"type:uuid;index" json:"project_id,omitempty"`
	IterationID    *uuid.UUID       `gorm:"type:uuid;index" json:"iteration_id,omitempty"`
	RequirementID  *uuid.UUID       `gorm:"type:uuid;index" json:"requirement_id,omitempty"`
	ParentID       *uuid.UUID       `gorm:"type:uuid;index" json:"parent_id,omitempty"` // 父任务（支持子任务）
	AssigneeID     *uuid.UUID       `gorm:"type:uuid;index" json:"assignee_id,omitempty"`
	CreatorID      uuid.UUID        `gorm:"type:uuid;not null;index" json:"creator_id"`
	Status         string           `gorm:"size:20;default:'waiting';index" json:"status"` // waiting/in_progress/done/paused/cancelled/closed
	Priority       int              `gorm:"default:3" json:"priority"`                     // 1-5�?最�?
	Severity       int              `gorm:"default:1" json:"severity"`                     // 严重程度 1-4
	StoryPoints    *float64         `json:"story_points,omitempty"`
	EstimatedHours *float64         `json:"estimated_hours,omitempty"`
	ActualHours    *float64         `json:"actual_hours,omitempty"`
	RemainingHours *float64         `json:"remaining_hours,omitempty"`
	StartDate      *time.Time       `json:"start_date,omitempty"`
	DueDate        *time.Time       `json:"due_date,omitempty"`
	FinishedDate   *time.Time       `json:"finished_date,omitempty"`
	KanbanColumn   string           `gorm:"size:50;default:'todo'" json:"kanban_column"` // 看板列位置
	Tags           string           `gorm:"type:text" json:"tags,omitempty"`             // JSON数组标签
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	DeletedAt      gorm.DeletedAt   `gorm:"index" json:"-"`
	Children       []Task           `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Assignees      []TaskAssignment `gorm:"foreignKey:TaskID" json:"assignees,omitempty"`
	Comments       []TaskComment    `gorm:"foreignKey:TaskID" json:"comments,omitempty"`
	Attachments    []TaskAttachment `gorm:"foreignKey:TaskID" json:"attachments,omitempty"`
}

// TableName 指定数据库表�?
func (Task) TableName() string {
	return "tasks"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// TaskAssignment 任务分配（支持多人分配）
type TaskAssignment struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TaskID     uuid.UUID `gorm:"type:uuid;index;not null" json:"task_id"`
	UserID     uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	Role       string    `gorm:"size:20;default:'assignee'" json:"role"` // assignee/reviewer
	AssignedAt time.Time `json:"assigned_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (TaskAssignment) TableName() string { return "task_assignments" }

func (ta *TaskAssignment) BeforeCreate(tx *gorm.DB) error {
	if ta.ID == uuid.Nil {
		ta.ID = uuid.New()
	}
	return nil
}

// TaskComment 任务评论（支持嵌套回复）
type TaskComment struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TaskID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"task_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	ParentID  *uuid.UUID     `gorm:"type:uuid;index" json:"parent_id,omitempty"` // 支持回复（嵌套评论）
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TaskComment) TableName() string { return "task_comments" }

func (tc *TaskComment) BeforeCreate(tx *gorm.DB) error {
	if tc.ID == uuid.Nil {
		tc.ID = uuid.New()
	}
	return nil
}

// TaskAttachment 任务附件
type TaskAttachment struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TaskID     uuid.UUID      `gorm:"type:uuid;index;not null" json:"task_id"`
	FileName   string         `gorm:"size:255;not null" json:"file_name"`
	FileSize   int64          `json:"file_size"`
	FileType   string         `gorm:"size:100" json:"file_type"`
	FileURL    string         `gorm:"size:500;not null" json:"file_url"`
	UploadedBy uuid.UUID      `gorm:"type:uuid;not null" json:"uploaded_by"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TaskAttachment) TableName() string { return "task_attachments" }

func (ta *TaskAttachment) BeforeCreate(tx *gorm.DB) error {
	if ta.ID == uuid.Nil {
		ta.ID = uuid.New()
	}
	return nil
}

// TaskLink 任务关联（任务间的前置/后续关系）
type TaskLink struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SourceTaskID uuid.UUID `gorm:"type:uuid;index;not null" json:"source_task_id"`
	TargetTaskID uuid.UUID `gorm:"type:uuid;index;not null" json:"target_task_id"`
	LinkType     string    `gorm:"size:20" json:"link_type"` // blocks/is_blocked_by/relates_to/duplicates
	Name         string    `gorm:"size:100" json:"name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (TaskLink) TableName() string { return "task_links" }

func (tl *TaskLink) BeforeCreate(tx *gorm.DB) error {
	if tl.ID == uuid.Nil {
		tl.ID = uuid.New()
	}
	return nil
}

// TaskWorklog 工作日志
type TaskWorklog struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TaskID     uuid.UUID `gorm:"type:uuid;index;not null" json:"task_id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	SpentHours float64   `gorm:"not null" json:"spent_hours"`
	WorkDate   time.Time `gorm:"not null" json:"work_date"`
	Summary    string    `gorm:"size:500" json:"summary,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

func (TaskWorklog) TableName() string { return "task_worklogs" }

func (tw *TaskWorklog) BeforeCreate(tx *gorm.DB) error {
	if tw.ID == uuid.Nil {
		tw.ID = uuid.New()
	}
	return nil
}
