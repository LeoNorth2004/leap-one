package repository

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-quality/internal/domain/entity"
)

// BugRepository Bug仓库接口定义
type BugRepository interface {
	// Create 创建Bug
	Create(ctx context.Context, bug *entity.Bug) error

	// GetByID 根据ID获取Bug详情（含评论、附件、历史）
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Bug, error)

	// Update 更新Bug基本信息
	Update(ctx context.Context, bug *entity.Bug) error

	// Delete 删除Bug（软删除�?	Delete(ctx context.Context, id uuid.UUID) error

	// List 分页查询Bug列表（支持高级筛选）
	List(ctx context.Context, page, pageSize int, filter *BugFilter) ([]*entity.Bug, int64, error)

	// ListMyBugs 查询"我的Bug"（我提报�?+ 我负责的�?	ListMyBugs(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*entity.Bug, int64, error)

	// ConfirmBug 确认Bug（new �?confirmed�?	ConfirmBug(ctx context.Context, id uuid.UUID, confirmedBy uuid.UUID) error

	// ResolveBug 解决Bug（in_progress �?resolved�?	ResolveBug(ctx context.Context, id uuid.UUID, resolution string, resolvedBy uuid.UUID) error

	// CloseBug 关闭Bug（resolved �?closed�?	CloseBug(ctx context.Context, id uuid.UUID, closedBy uuid.UUID) error

	// ReopenBug 重新激活Bug（closed/resolved �?reopened/in_progress�?	ReopenBug(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	// AddComment 添加Bug评论
	AddComment(ctx context.Context, comment *entity.BugComment) error

	// ListComments 获取Bug评论列表
	ListComments(ctx context.Context, bugID uuid.UUID) ([]*entity.BugComment, error)

	// AddAttachment 添加Bug附件
	AddAttachment(ctx context.Context, attachment *entity.BugAttachment) error

	// ListAttachments 获取Bug附件列表
	ListAttachments(ctx context.Context, bugID uuid.UUID) ([]*entity.BugAttachment, error)

	// ListHistory 获取Bug变更历史
	ListHistory(ctx context.Context, bugID uuid.UUID) ([]*entity.BugHistory, error)

	// AddHistory 添加变更历史记录
	AddHistory(ctx context.Context, history *entity.BugHistory) error

	// GetStatistics 获取Bug统计数据
	GetStatistics(ctx context.Context, productID, projectID *uuid.UUID) (*BugStatistics, error)
}

// BugFilter Bug高级筛选条�?type BugFilter struct {
	Keyword       string     // 关键词搜索（标题�?	Status        string     // 状�?new/confirmed/in_progress/resolved/closed/reopened/cancelled
	Severity      *int       // 严重程度 1-4
	Priority      *int       // 优先�?1-5
	Type          string     // 类型 code_bug/design_bug/data_bug/config/security/performance/ui
	ProductID     *uuid.UUID // 产品ID
	ProjectID     *uuid.UUID // 项目ID
	ReporterID    *uuid.UUID // 提报人ID
	AssigneeID    *uuid.UUID // 处理人ID
	IterationID   *uuid.UUID // 迭代ID
	Resolution    string     // 解决方案
	StartDate     string     // 开始日期筛�?	EndDate       string     // 结束日期筛�?}

// BugStatistics Bug统计结果
type BugStatistics struct {
	TotalCount   int64            `json:"total_count"`   // 总数
	NewCount     int64            `json:"new_count"`     // 新建�?	ConfirmedCnt int64            `json:"confirmed_cnt"` // 已确认数
	InProgress   int64            `json:"in_progress"`   // 处理中数
	ResolvedCnt  int64            `json:"resolved_cnt"`  // 已解决数
	ClosedCnt    int64            `json:"closed_cnt"`    // 已关闭数
	ReopenedCnt  int64            `json:"reopened_cnt"`  // 重新打开�?	BySeverity   map[int]int64   `json:"by_severity"`   // 按严重程度分�?	ByPriority   map[int]int64   `json:"by_priority"`    // 按优先级分布
	ByType       map[string]int64 `json:"by_type"`       // 按类型分�?}
