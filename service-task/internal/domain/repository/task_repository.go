package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"leap-one/service-task/internal/domain/entity"
)

// TaskRepository 任务仓库接口定义
type TaskRepository interface {
	// Create 创建任务
	Create(ctx context.Context, task *entity.Task) error

	// GetByID 根据ID获取任务
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error)

	// Update 更新任务
	Update(ctx context.Context, task *entity.Task) error

	// Delete 软删除任务
	Delete(ctx context.Context, id uuid.UUID) error

	// List 分页查询任务列表（支持高级筛选）
	List(ctx context.Context, filter *TaskFilter) ([]*entity.Task, int64, error)

	// ListByAssigneeID 查询指定用户的任务
	ListByAssigneeID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*entity.Task, int64, error)

	// UpdateStatus 更新任务状态
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

// TaskFilter 任务查询筛选条件
type TaskFilter struct {
	Page         int        `json:"page"`
	PageSize     int        `json:"page_size"`
	Keyword      string     `json:"keyword"`       // 标题/描述搜索
	Status       string     `json:"status"`        // 状态筛选
	Type         string     `json:"type"`          // 类型筛选
	Priority     *int       `json:"priority"`      // 优先级筛选
	ProjectID    *uuid.UUID `json:"project_id"`    // 项目筛选
	IterationID  *uuid.UUID `json:"iteration_id"`  // 迭代筛选
	AssigneeID   *uuid.UUID `json:"assignee_id"`   // 执行人筛选
	CreatorID    *uuid.UUID `json:"creator_id"`    // 创建人筛选
	ParentID     *uuid.UUID `json:"parent_id"`     // 父任务筛选（为nil且IsSubTask=false时查顶级）
	IsSubTask    bool       `json:"is_sub_task"`   // 是否只查子任务
	StartDateFrom *time.Time `json:"start_date_from"`
	StartDateTo   *time.Time `json:"start_date_to"`
	DueDateFrom   *time.Time `json:"due_date_from"`
	DueDateTo     *time.Time `json:"due_date_to"`
	SortBy        string    `json:"sort_by"`        // 排序字段
	SortOrder     string    `json:"sort_order"`     // asc/desc
}

// TaskAssignmentRepository 任务分配仓库接口
type TaskAssignmentRepository interface {
	Create(ctx context.Context, assignment *entity.TaskAssignment) error
	Delete(ctx context.Context, taskID, userID uuid.UUID) error
	DeleteByTaskID(ctx context.Context, taskID uuid.UUID) error
	ListByTaskID(ctx context.Context, taskID uuid.UUID) ([]*entity.TaskAssignment, error)
}

// TaskCommentRepository 任务评论仓库接口
type TaskCommentRepository interface {
	Create(ctx context.Context, comment *entity.TaskComment) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByTaskID(ctx context.Context, taskID uuid.UUID) ([]*entity.TaskComment, error)
}

// TaskAttachmentRepository 任务附件仓库接口
type TaskAttachmentRepository interface {
	Create(ctx context.Context, attachment *entity.TaskAttachment) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByTaskID(ctx context.Context, taskID uuid.UUID) ([]*entity.TaskAttachment, error)
}

// TaskLinkRepository 任务关联仓库接口
type TaskLinkRepository interface {
	Create(ctx context.Context, link *entity.TaskLink) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByTaskID(ctx context.Context, taskID uuid.UUID) ([]*entity.TaskLink, error)
}

// TaskWorklogRepository 工作日志仓库接口
type TaskWorklogRepository interface {
	Create(ctx context.Context, worklog *entity.TaskWorklog) error
	ListByTaskID(ctx context.Context, taskID uuid.UUID) ([]*entity.TaskWorklog, error)
	GetTotalSpentHours(ctx context.Context, taskID uuid.UUID) (float64, error)
}
