package repository_impl

import (
	"context"
	"time"

	"github.com/google/uuid"
	"leap-one/service-task/internal/domain/entity"
	"leap-one/service-task/internal/domain/repository"
	"gorm.io/gorm"
)

// IssueRepositoryImpl 工单仓库实现
type IssueRepositoryImpl struct {
	db *gorm.DB
}

func NewIssueRepository(db *gorm.DB) repository.IssueRepository {
	return &IssueRepositoryImpl{db: db}
}

func (r *IssueRepositoryImpl) Create(ctx context.Context, issue *entity.Issue) error {
	return r.db.WithContext(ctx).Create(issue).Error
}

func (r *IssueRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Issue, error) {
	var issue entity.Issue
	err := r.db.WithContext(ctx).Preload("Comments").Preload("Attachments").
		First(&issue, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &issue, nil
}

func (r *IssueRepositoryImpl) Update(ctx context.Context, issue *entity.Issue) error {
	return r.db.WithContext(ctx).Save(issue).Error
}

func (r *IssueRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Issue{}, "id = ?", id).Error
}

func (r *IssueRepositoryImpl) List(ctx context.Context, filter *repository.IssueFilter) ([]*entity.Issue, int64, error) {
	var issues []*entity.Issue
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Issue{})
	query = buildIssueQuery(query, filter)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := filter.SortOrder
	if sortOrder == "" {
		sortOrder = "DESC"
	}

	offset := (page - 1) * pageSize
	err := query.Order(sortBy + " " + sortOrder).
		Offset(offset).
		Limit(pageSize).
		Find(&issues).Error
	if err != nil {
		return nil, 0, err
	}

	return issues, total, nil
}

func (r *IssueRepositoryImpl) ListByReporterID(ctx context.Context, reporterID uuid.UUID, page, pageSize int) ([]*entity.Issue, int64, error) {
	var issues []*entity.Issue
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Issue{}).Where("reporter_id = ?", reporterID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&issues).Error
	if err != nil {
		return nil, 0, err
	}
	return issues, total, nil
}

func (r *IssueRepositoryImpl) ListByAssigneeID(ctx context.Context, assigneeID uuid.UUID, page, pageSize int) ([]*entity.Issue, int64, error) {
	var issues []*entity.Issue
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Issue{}).Where("assignee_id = ?", assigneeID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&issues).Error
	if err != nil {
		return nil, 0, err
	}
	return issues, total, nil
}

func (r *IssueRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": now,
	}
	switch status {
	case "resolved":
		updates["resolved_at"] = now
	case "closed":
		updates["closed_at"] = now
	}

	return r.db.WithContext(ctx).Model(&entity.Issue{}).Where("id = ?", id).Updates(updates).Error
}

func (r *IssueRepositoryImpl) UpdateSLADates(ctx context.Context, id uuid.UUID, slaDueDate, responseDueDate *time.Time) error {
	updates := map[string]interface{}{}
	if slaDueDate != nil {
		updates["sla_due_date"] = *slaDueDate
	}
	if responseDueDate != nil {
		updates["response_due_date"] = *responseDueDate
	}
	if len(updates) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&entity.Issue{}).Where("id = ?", id).Updates(updates).Error
}

// buildIssueQuery 构建工单查询条件
func buildIssueQuery(query *gorm.DB, filter *repository.IssueFilter) *gorm.DB {
	if filter == nil {
		return query
	}

	if filter.Keyword != "" {
		pattern := "%" + filter.Keyword + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", pattern, pattern)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}
	if filter.Severity != nil {
		query = query.Where("severity = ?", *filter.Severity)
	}
	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.ProductID != nil {
		query = query.Where("product_id = ?", *filter.ProductID)
	}
	if filter.ReporterID != nil {
		query = query.Where("reporter_id = ?", *filter.ReporterID)
	}
	if filter.AssigneeID != nil {
		query = query.Where("assignee_id = ?", *filter.AssigneeID)
	}
	if filter.Source != "" {
		query = query.Where("source = ?", filter.Source)
	}
	if filter.CreatedAtFrom != nil {
		query = query.Where("created_at >= ?", *filter.CreatedAtFrom)
	}
	if filter.CreatedAtTo != nil {
		query = query.Where("created_at <= ?", *filter.CreatedAtTo)
	}

	return query
}
