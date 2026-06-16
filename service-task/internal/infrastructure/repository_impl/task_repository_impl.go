package repository_impl

import (
	"context"
	"time"

	"github.com/google/uuid"
	"leap-one/service-task/internal/domain/entity"
	"leap-one/service-task/internal/domain/repository"
	"gorm.io/gorm"
)

// TaskRepositoryImpl 任务仓库实现
type TaskRepositoryImpl struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) repository.TaskRepository {
	return &TaskRepositoryImpl{db: db}
}

func (r *TaskRepositoryImpl) Create(ctx context.Context, task *entity.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *TaskRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error) {
	var task entity.Task
	err := r.db.WithContext(ctx).Preload("Assignees").Preload("Comments").Preload("Attachments").
		First(&task, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepositoryImpl) Update(ctx context.Context, task *entity.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *TaskRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Task{}, "id = ?", id).Error
}

func (r *TaskRepositoryImpl) List(ctx context.Context, filter *repository.TaskFilter) ([]*entity.Task, int64, error) {
	var tasks []*entity.Task
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Task{})

	query = buildTaskQuery(query, filter)

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
		Find(&tasks).Error

	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *TaskRepositoryImpl) ListByAssigneeID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*entity.Task, int64, error) {
	var tasks []*entity.Task
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Task{}).
		Where("assignee_id = ?", userID)

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
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *TaskRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": now,
	}
	switch status {
	case "in_progress":
		updates["start_date"] = now
		updates["kanban_column"] = "doing"
	case "done":
		updates["finished_date"] = now
		updates["kanban_column"] = "done"
	case "paused":
		updates["kanban_column"] = "todo"
	}

	return r.db.WithContext(ctx).Model(&entity.Task{}).Where("id = ?", id).Updates(updates).Error
}

// buildTaskQuery 构建任务查询条件
func buildTaskQuery(query *gorm.DB, filter *repository.TaskFilter) *gorm.DB {
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
	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.IterationID != nil {
		query = query.Where("iteration_id = ?", *filter.IterationID)
	}
	if filter.AssigneeID != nil {
		query = query.Where("assignee_id = ?", *filter.AssigneeID)
	}
	if filter.CreatorID != nil {
		query = query.Where("creator_id = ?", *filter.CreatorID)
	}
	if !filter.IsSubTask && filter.ParentID == nil {
		query = query.Where("parent_id IS NULL")
	} else if filter.ParentID != nil {
		query = query.Where("parent_id = ?", *filter.ParentID)
	}
	if filter.StartDateFrom != nil {
		query = query.Where("start_date >= ?", *filter.StartDateFrom)
	}
	if filter.StartDateTo != nil {
		query = query.Where("start_date <= ?", *filter.StartDateTo)
	}
	if filter.DueDateFrom != nil {
		query = query.Where("due_date >= ?", *filter.DueDateFrom)
	}
	if filter.DueDateTo != nil {
		query = query.Where("due_date <= ?", *filter.DueDateTo)
	}

	return query
}
