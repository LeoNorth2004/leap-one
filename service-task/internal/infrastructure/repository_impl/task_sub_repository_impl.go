package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-task/internal/domain/entity"
	"leap-one/service-task/internal/domain/repository"
	"gorm.io/gorm"
)

// TaskAssignmentRepositoryImpl 任务分配仓库实现
type TaskAssignmentRepositoryImpl struct {
	db *gorm.DB
}

func NewTaskAssignmentRepository(db *gorm.DB) repository.TaskAssignmentRepository {
	return &TaskAssignmentRepositoryImpl{db: db}
}

func (r *TaskAssignmentRepositoryImpl) Create(ctx context.Context, assignment *entity.TaskAssignment) error {
	return r.db.WithContext(ctx).Create(assignment).Error
}

func (r *TaskAssignmentRepositoryImpl) Delete(ctx context.Context, taskID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.TaskAssignment{}, "task_id = ? AND user_id = ?", taskID, userID).Error
}

func (r *TaskAssignmentRepositoryImpl) DeleteByTaskID(ctx context.Context, taskID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.TaskAssignment{}, "task_id = ?", taskID).Error
}

func (r *TaskAssignmentRepositoryImpl) ListByTaskID(ctx context.Context, taskID uuid.UUID) ([]*entity.TaskAssignment, error) {
	var assignments []*entity.TaskAssignment
	err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Find(&assignments).Error
	if err != nil {
		return nil, err
	}
	return assignments, nil
}

// TaskCommentRepositoryImpl 任务评论仓库实现
type TaskCommentRepositoryImpl struct {
	db *gorm.DB
}

func NewTaskCommentRepository(db *gorm.DB) repository.TaskCommentRepository {
	return &TaskCommentRepositoryImpl{db: db}
}

func (r *TaskCommentRepositoryImpl) Create(ctx context.Context, comment *entity.TaskComment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *TaskCommentRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.TaskComment{}, "id = ?", id).Error
}

func (r *TaskCommentRepositoryImpl) ListByTaskID(ctx context.Context, taskID uuid.UUID) ([]*entity.TaskComment, error) {
	var comments []*entity.TaskComment
	err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Order("created_at ASC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// TaskAttachmentRepositoryImpl 任务附件仓库实现
type TaskAttachmentRepositoryImpl struct {
	db *gorm.DB
}

func NewTaskAttachmentRepository(db *gorm.DB) repository.TaskAttachmentRepository {
	return &TaskAttachmentRepositoryImpl{db: db}
}

func (r *TaskAttachmentRepositoryImpl) Create(ctx context.Context, attachment *entity.TaskAttachment) error {
	return r.db.WithContext(ctx).Create(attachment).Error
}

func (r *TaskAttachmentRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.TaskAttachment{}, "id = ?", id).Error
}

func (r *TaskAttachmentRepositoryImpl) ListByTaskID(ctx context.Context, taskID uuid.UUID) ([]*entity.TaskAttachment, error) {
	var attachments []*entity.TaskAttachment
	err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Order("created_at DESC").Find(&attachments).Error
	if err != nil {
		return nil, err
	}
	return attachments, nil
}

// TaskLinkRepositoryImpl 任务关联仓库实现
type TaskLinkRepositoryImpl struct {
	db *gorm.DB
}

func NewTaskLinkRepository(db *gorm.DB) repository.TaskLinkRepository {
	return &TaskLinkRepositoryImpl{db: db}
}

func (r *TaskLinkRepositoryImpl) Create(ctx context.Context, link *entity.TaskLink) error {
	return r.db.WithContext(ctx).Create(link).Error
}

func (r *TaskLinkRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.TaskLink{}, "id = ?", id).Error
}

func (r *TaskLinkRepositoryImpl) ListByTaskID(ctx context.Context, taskID uuid.UUID) ([]*entity.TaskLink, error) {
	var links []*entity.TaskLink
	err := r.db.WithContext(ctx).Where("source_task_id = ? OR target_task_id = ?", taskID, taskID).Find(&links).Error
	if err != nil {
		return nil, err
	}
	return links, nil
}

// TaskWorklogRepositoryImpl 工作日志仓库实现
type TaskWorklogRepositoryImpl struct {
	db *gorm.DB
}

func NewTaskWorklogRepository(db *gorm.DB) repository.TaskWorklogRepository {
	return &TaskWorklogRepositoryImpl{db: db}
}

func (r *TaskWorklogRepositoryImpl) Create(ctx context.Context, worklog *entity.TaskWorklog) error {
	return r.db.WithContext(ctx).Create(worklog).Error
}

func (r *TaskWorklogRepositoryImpl) ListByTaskID(ctx context.Context, taskID uuid.UUID) ([]*entity.TaskWorklog, error) {
	var worklogs []*entity.TaskWorklog
	err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Order("work_date DESC").Find(&worklogs).Error
	if err != nil {
		return nil, err
	}
	return worklogs, nil
}

func (r *TaskWorklogRepositoryImpl) GetTotalSpentHours(ctx context.Context, taskID uuid.UUID) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).Model(&entity.TaskWorklog{}).Where("task_id = ?", taskID).Select("COALESCE(SUM(spent_hours), 0)").Scan(&total).Error
	return total, err
}
