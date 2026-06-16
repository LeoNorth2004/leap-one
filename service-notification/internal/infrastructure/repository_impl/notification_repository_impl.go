package repository_impl

import (
	"context"
	"time"

	"github.com/google/uuid"
	"leap-one/service-notification/internal/domain/entity"
	"leap-one/service-notification/internal/domain/repository"
	"gorm.io/gorm"
)

// NotificationRepositoryImpl 通知消息仓库实现
type NotificationRepositoryImpl struct{ db *gorm.DB }

func NewNotificationRepository(db *gorm.DB) repository.NotificationRepository {
	return &NotificationRepositoryImpl{db: db}
}

func (r *NotificationRepositoryImpl) Create(ctx context.Context, n *entity.Notification) error {
	return r.db.WithContext(ctx).Create(n).Error
}
func (r *NotificationRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error) {
	var n entity.Notification
	err := r.db.WithContext(ctx).First(&n, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}
func (r *NotificationRepositoryImpl) ListByReceiver(ctx context.Context, receiverID uuid.UUID, page, pageSize int, unreadOnly bool) ([]*entity.Notification, int64, error) {
	var list []*entity.Notification
	var total int64
	query := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("receiver_id = ?", receiverID)
	if unreadOnly {
		query = query.Where("is_read = false")
	}
	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Order("sent_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
func (r *NotificationRepositoryImpl) CountUnread(ctx context.Context, receiverID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("receiver_id = ? AND is_read = false", receiverID).Count(&count).Error
	return count, err
}
func (r *NotificationRepositoryImpl) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.Notification{}).Where("id = ?", id).Updates(map[string]interface{}{"is_read": true, "read_at": now}).Error
}
func (r *NotificationRepositoryImpl) MarkAllAsRead(ctx context.Context, receiverID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.Notification{}).Where("receiver_id = ? AND is_read = false", receiverID).Updates(map[string]interface{}{"is_read": true, "read_at": now}).Error
}
func (r *NotificationRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Notification{}, "id = ?", id).Error
}
func (r *NotificationRepositoryImpl) BatchCreate(ctx context.Context, ns []*entity.Notification) error {
	if len(ns) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&ns).Error
}

// TemplateRepositoryImpl 消息模板仓库实现
type TemplateRepositoryImpl struct{ db *gorm.DB }

func NewTemplateRepository(db *gorm.DB) repository.TemplateRepository {
	return &TemplateRepositoryImpl{db: db}
}

func (r *TemplateRepositoryImpl) Create(ctx context.Context, t *entity.NotificationTemplate) error {
	return r.db.WithContext(ctx).Create(t).Error
}
func (r *TemplateRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.NotificationTemplate, error) {
	var t entity.NotificationTemplate
	err := r.db.WithContext(ctx).First(&t, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}
func (r *TemplateRepositoryImpl) GetByCode(ctx context.Context, code string) (*entity.NotificationTemplate, error) {
	var t entity.NotificationTemplate
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&t).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}
func (r *TemplateRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*entity.NotificationTemplate, int64, error) {
	var list []*entity.NotificationTemplate
	var total int64
	query := r.db.WithContext(ctx).Model(&entity.NotificationTemplate{})
	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
func (r *TemplateRepositoryImpl) Update(ctx context.Context, t *entity.NotificationTemplate) error {
	return r.db.WithContext(ctx).Save(t).Error
}
func (r *TemplateRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.NotificationTemplate{}, "id = ?", id).Error
}
