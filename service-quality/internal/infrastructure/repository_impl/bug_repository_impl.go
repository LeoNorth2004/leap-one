package repository_impl

import (
	"context"
	"time"

	"github.com/google/uuid"
	"leap-one/service-quality/internal/domain/entity"
	"leap-one/service-quality/internal/domain/repository"
	"gorm.io/gorm"
)

// BugRepositoryImpl Bug仓库实现
type BugRepositoryImpl struct {
	db *gorm.DB
}

// NewBugRepository 创建Bug仓库实例
func NewBugRepository(db *gorm.DB) repository.BugRepository {
	return &BugRepositoryImpl{db: db}
}

// Create 创建Bug
func (r *BugRepositoryImpl) Create(ctx context.Context, bug *entity.Bug) error {
	return r.db.WithContext(ctx).Create(bug).Error
}

// GetByID 根据ID获取Bug详情（预加载评论、附件、历史）
func (r *BugRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Bug, error) {
	var bug entity.Bug
	err := r.db.WithContext(ctx).
		Preload("Comments").
		Preload("Attachments").
		Preload("History", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		First(&bug, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &bug, nil
}

// Update 更新Bug基本信息
func (r *BugRepositoryImpl) Update(ctx context.Context, bug *entity.Bug) error {
	return r.db.WithContext(ctx).Save(bug).Error
}

// Delete 删除Bug（软删除�?func (r *BugRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Bug{}, "id = ?", id).Error
}

// List 分页查询Bug列表（支持高级筛选）
func (r *BugRepositoryImpl) List(ctx context.Context, page, pageSize int, filter *repository.BugFilter) ([]*entity.Bug, int64, error) {
	var bugs []*entity.Bug
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Bug{})

	query = r.applyBugFilter(query, filter)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&bugs).Error

	if err != nil {
		return nil, 0, err
	}

	return bugs, total, nil
}

// ListMyBugs 查询"我的Bug"（我提报�?+ 我负责的�?func (r *BugRepositoryImpl) ListMyBugs(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*entity.Bug, int64, error) {
	var bugs []*entity.Bug
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Bug{}).
		Where("reporter_id = ? OR assignee_id = ?", userID, userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("updated_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&bugs).Error

	if err != nil {
		return nil, 0, err
	}

	return bugs, total, nil
}

// ConfirmBug 确认Bug（new �?confirmed�?func (r *BugRepositoryImpl) ConfirmBug(ctx context.Context, id uuid.UUID, confirmedBy uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.Bug{}).
		Where("id = ? AND status = 'new'", id).
		Updates(map[string]interface{}{
			"status":      "confirmed",
			"confirmed_at": now,
			"confirmed_by": confirmedBy,
		}).Error
}

// ResolveBug 解决Bug（in_progress �?resolved�?func (r *BugRepositoryImpl) ResolveBug(ctx context.Context, id uuid.UUID, resolution string, resolvedBy uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.Bug{}).
		Where("id = ? AND status IN ('in_progress','confirmed')", id).
		Updates(map[string]interface{}{
			"status":     "resolved",
			"resolution": resolution,
			"resolved_at": now,
			"resolved_by": resolvedBy,
		}).Error
}

// CloseBug 关闭Bug（resolved �?closed�?func (r *BugRepositoryImpl) CloseBug(ctx context.Context, id uuid.UUID, closedBy uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.Bug{}).
		Where("id = ? AND status = 'resolved'", id).
		Updates(map[string]interface{}{
			"status":    "closed",
			"closed_at": now,
			"closed_by": closedBy,
		}).Error
}

// ReopenBug 重新激活Bug
func (r *BugRepositoryImpl) ReopenBug(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.Bug{}).
		Where("id = ? AND status IN ('resolved','closed')", id).
		Update("status", "reopened").Error
}

// AddComment 添加Bug评论
func (r *BugRepositoryImpl) AddComment(ctx context.Context, comment *entity.BugComment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

// ListComments 获取Bug评论列表
func (r *BugRepositoryImpl) ListComments(ctx context.Context, bugID uuid.UUID) ([]*entity.BugComment, error) {
	var comments []*entity.BugComment
	err := r.db.WithContext(ctx).
		Where("bug_id = ?", bugID).
		Order("created_at ASC").
		Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// AddAttachment 添加Bug附件
func (r *BugRepositoryImpl) AddAttachment(ctx context.Context, attachment *entity.BugAttachment) error {
	return r.db.WithContext(ctx).Create(attachment).Error
}

// ListAttachments 获取Bug附件列表
func (r *BugRepositoryImpl) ListAttachments(ctx context.Context, bugID uuid.UUID) ([]*entity.BugAttachment, error) {
	var attachments []*entity.BugAttachment
	err := r.db.WithContext(ctx).
		Where("bug_id = ?", bugID).
		Order("created_at ASC").
		Find(&attachments).Error
	if err != nil {
		return nil, err
	}
	return attachments, nil
}

// ListHistory 获取Bug变更历史
func (r *BugRepositoryImpl) ListHistory(ctx context.Context, bugID uuid.UUID) ([]*entity.BugHistory, error) {
	var histories []*entity.BugHistory
	err := r.db.WithContext(ctx).
		Where("bug_id = ?", bugID).
		Order("created_at ASC").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

// AddHistory 添加变更历史记录
func (r *BugRepositoryImpl) AddHistory(ctx context.Context, history *entity.BugHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetStatistics 获取Bug统计数据
func (r *BugRepositoryImpl) GetStatistics(ctx context.Context, productID, projectID *uuid.UUID) (*repository.BugStatistics, error) {
	stats := &repository.BugStatistics{
		BySeverity: make(map[int]int64),
		ByPriority: make(map[int]int64),
		ByType:     make(map[string]int64),
	}

	baseQuery := r.db.WithContext(ctx).Model(&entity.Bug{})

	if productID != nil {
		baseQuery = baseQuery.Where("product_id = ?", *productID)
	}
	if projectID != nil {
		baseQuery = baseQuery.Where("project_id = ?", *projectID)
	}

	// 总数统计
	baseQuery.Count(&stats.TotalCount)

	// 各状态数量统�?	statuses := []struct {
		Field *int64
		Value string
	}{
		{&stats.NewCount, "new"},
		{&stats.ConfirmedCnt, "confirmed"},
		{&stats.InProgress, "in_progress"},
		{&stats.ResolvedCnt, "resolved"},
		{&stats.ClosedCnt, "closed"},
		{&stats.ReopenedCnt, "reopened"},
	}
	for _, s := range statuses {
		q := baseQuery.Session(&gorm.Session{})
		if productID != nil {
			q = q.Where("product_id = ?", *productID)
		}
		if projectID != nil {
			q = q.Where("project_id = ?", *projectID)
		}
		q.Where("status = ?", s.Value).Count(s.Field)
	}

	// 按严重程度分�?	var severityStats []struct {
		Severity int
		Count    int64
	}
	sevQ := baseQuery.Session(&gorm.Session{}).Select("severity, COUNT(*) as count").Group("severity")
	if productID != nil {
		sevQ = sevQ.Where("product_id = ?", *productID)
	}
	if projectID != nil {
		sevQ = sevQ.Where("project_id = ?", *projectID)
	}
	sevQ.Scan(&severityStats)
	for _, s := range severityStats {
		stats.BySeverity[s.Severity] = s.Count
	}

	// 按优先级分布
	var priorityStats []struct {
		Priority int
		Count    int64
	}
	priQ := baseQuery.Session(&gorm.Session{}).Select("priority, COUNT(*) as count").Group("priority")
	if productID != nil {
		priQ = priQ.Where("product_id = ?", *productID)
	}
	if projectID != nil {
		priQ = priQ.Where("project_id = ?", *projectID)
	}
	priQ.Scan(&priorityStats)
	for _, p := range priorityStats {
		stats.ByPriority[p.Priority] = p.Count
	}

	// 按类型分�?	var typeStats []struct {
		BugType string
		Count   int64
	}
	typQ := baseQuery.Session(&gorm.Session{}).Select("type, COUNT(*) as count").Group("type")
	if productID != nil {
		typQ = typQ.Where("product_id = ?", *productID)
	}
	if projectID != nil {
		typQ = typQ.Where("project_id = ?", *projectID)
	}
	typQ.Scan(&typeStats)
	for _, t := range typeStats {
		stats.ByType[t.BugType] = t.Count
	}

	return stats, nil
}

// applyBugFilter 应用Bug高级筛选条件到查询构建�?func (r *BugRepositoryImpl) applyBugFilter(query *gorm.DB, filter *repository.BugFilter) *gorm.DB {
	if filter == nil {
		return query
	}

	if filter.Keyword != "" {
		searchPattern := "%" + filter.Keyword + "%"
		query = query.Where("title LIKE ?", searchPattern)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Severity != nil {
		query = query.Where("severity = ?", *filter.Severity)
	}
	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.ProductID != nil {
		query = query.Where("product_id = ?", *filter.ProductID)
	}
	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.ReporterID != nil {
		query = query.Where("reporter_id = ?", *filter.ReporterID)
	}
	if filter.AssigneeID != nil {
		query = query.Where("assignee_id = ?", *filter.AssigneeID)
	}
	if filter.IterationID != nil {
		query = query.Where("iteration_id = ?", *filter.IterationID)
	}
	if filter.Resolution != "" {
		query = query.Where("resolution = ?", filter.Resolution)
	}
	if filter.StartDate != "" {
		query = query.Where("created_at >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("created_at <= ?", filter.EndDate+" 23:59:59")
	}

	return query
}
