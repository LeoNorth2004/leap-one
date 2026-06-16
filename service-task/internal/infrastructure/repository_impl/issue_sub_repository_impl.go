package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-task/internal/domain/entity"
	"leap-one/service-task/internal/domain/repository"
	"gorm.io/gorm"
)

// IssueCommentRepositoryImpl 工单评论仓库实现
type IssueCommentRepositoryImpl struct {
	db *gorm.DB
}

func NewIssueCommentRepository(db *gorm.DB) repository.IssueCommentRepository {
	return &IssueCommentRepositoryImpl{db: db}
}

func (r *IssueCommentRepositoryImpl) Create(ctx context.Context, comment *entity.IssueComment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *IssueCommentRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.IssueComment{}, "id = ?", id).Error
}

func (r *IssueCommentRepositoryImpl) ListByIssueID(ctx context.Context, issueID uuid.UUID) ([]*entity.IssueComment, error) {
	var comments []*entity.IssueComment
	err := r.db.WithContext(ctx).Where("issue_id = ?", issueID).Order("created_at ASC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// IssueAttachmentRepositoryImpl 工单附件仓库实现
type IssueAttachmentRepositoryImpl struct {
	db *gorm.DB
}

func NewIssueAttachmentRepository(db *gorm.DB) repository.IssueAttachmentRepository {
	return &IssueAttachmentRepositoryImpl{db: db}
}

func (r *IssueAttachmentRepositoryImpl) Create(ctx context.Context, attachment *entity.IssueAttachment) error {
	return r.db.WithContext(ctx).Create(attachment).Error
}

func (r *IssueAttachmentRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.IssueAttachment{}, "id = ?", id).Error
}

func (r *IssueAttachmentRepositoryImpl) ListByIssueID(ctx context.Context, issueID uuid.UUID) ([]*entity.IssueAttachment, error) {
	var attachments []*entity.IssueAttachment
	err := r.db.WithContext(ctx).Where("issue_id = ?", issueID).Order("created_at DESC").Find(&attachments).Error
	if err != nil {
		return nil, err
	}
	return attachments, nil
}
