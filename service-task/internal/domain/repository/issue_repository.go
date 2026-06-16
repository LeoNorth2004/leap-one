package repository

import (
	"context"
	"time"

	"leap-one/service-task/internal/domain/entity"

	"github.com/google/uuid"
)

// IssueRepository 工单仓库接口定义
type IssueRepository interface {
	// Create 创建工单
	Create(ctx context.Context, issue *entity.Issue) error

	// GetByID 根据ID获取工单
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Issue, error)

	// Update 更新工单
	Update(ctx context.Context, issue *entity.Issue) error

	// Delete 软删除工单
	Delete(ctx context.Context, id uuid.UUID) error

	// List 分页查询工单列表（支持高级筛选）
	List(ctx context.Context, filter *IssueFilter) ([]*entity.Issue, int64, error)

	// ListByReporterID 查询指定用户的工单（我提报的）
	ListByReporterID(ctx context.Context, reporterID uuid.UUID, page, pageSize int) ([]*entity.Issue, int64, error)

	// ListByAssigneeID 查询指派给指定用户的工单（我的工单）
	ListByAssigneeID(ctx context.Context, assigneeID uuid.UUID, page, pageSize int) ([]*entity.Issue, int64, error)

	// UpdateStatus 更新工单状态
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error

	// UpdateSLADates 更新SLA截止时间
	UpdateSLADates(ctx context.Context, id uuid.UUID, slaDueDate, responseDueDate *time.Time) error
}

// IssueFilter 工单查询筛选条件
type IssueFilter struct {
	Page          int        `json:"page"`
	PageSize      int        `json:"page_size"`
	Keyword       string     `json:"keyword"`     // 标题/描述搜索
	Status        string     `json:"status"`      // 状态筛选
	Type          string     `json:"type"`        // 类型筛选
	Priority      *int       `json:"priority"`    // 优先级筛选
	Severity      *int       `json:"severity"`    // 严重程度筛选
	ProjectID     *uuid.UUID `json:"project_id"`  // 项目筛选
	ProductID     *uuid.UUID `json:"product_id"`  // 产品筛选
	ReporterID    *uuid.UUID `json:"reporter_id"` // 提报人筛选
	AssigneeID    *uuid.UUID `json:"assignee_id"` // 处理人筛选
	Source        string     `json:"source"`      // 来源筛选
	CreatedAtFrom *time.Time `json:"created_at_from"`
	CreatedAtTo   *time.Time `json:"created_at_to"`
	SortBy        string     `json:"sort_by"`    // 排序字段
	SortOrder     string     `json:"sort_order"` // asc/desc
}

// IssueCommentRepository 工单评论仓库接口
type IssueCommentRepository interface {
	// Create 创建评论
	Create(ctx context.Context, comment *entity.IssueComment) error

	// Delete 删除评论
	Delete(ctx context.Context, id uuid.UUID) error

	// ListByIssueID 查询工单评论列表
	ListByIssueID(ctx context.Context, issueID uuid.UUID) ([]*entity.IssueComment, error)
}

// IssueAttachmentRepository 工单附件仓库接口
type IssueAttachmentRepository interface {
	// Create 创建附件记录
	Create(ctx context.Context, attachment *entity.IssueAttachment) error

	// Delete 删除附件
	Delete(ctx context.Context, id uuid.UUID) error

	// ListByIssueID 查询工单附件列表
	ListByIssueID(ctx context.Context, issueID uuid.UUID) ([]*entity.IssueAttachment, error)
}
