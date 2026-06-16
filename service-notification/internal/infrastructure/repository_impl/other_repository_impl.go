package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-notification/internal/domain/entity"
	"leap-one/service-notification/internal/domain/repository"
	"gorm.io/gorm"
)

// EmailLogRepositoryImpl 邮件日志仓库实现
type EmailLogRepositoryImpl struct{ db *gorm.DB }

func NewEmailLogRepository(db *gorm.DB) repository.EmailLogRepository {
	return &EmailLogRepositoryImpl{db: db}
}

func (r *EmailLogRepositoryImpl) Create(ctx context.Context, log *entity.EmailLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
func (r *EmailLogRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.EmailLog, error) {
	var l entity.EmailLog
	err := r.db.WithContext(ctx).First(&l, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}
func (r *EmailLogRepositoryImpl) List(ctx context.Context, page, pageSize int, status string) ([]*entity.EmailLog, int64, error) {
	var list []*entity.EmailLog
	var total int64
	query := r.db.WithContext(ctx).Model(&entity.EmailLog{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
func (r *EmailLogRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status, errorMsg string) error {
	updates := map[string]interface{}{"status": status}
	if status == "sent" {
		now := entity.NowTime()
		updates["sent_at"] = &now
	}
	if errorMsg != "" {
		updates["error_msg"] = errorMsg
	}
	return r.db.WithContext(ctx).Model(&entity.EmailLog{}).Where("id = ?", id).Updates(updates).Error
}
func (r *EmailLogRepositoryImpl) IncrementRetry(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entity.EmailLog{}).Where("id = ?", id).UpdateColumn("retry_count", gorm.Expr("retry_count + 1")).Error
}

// WebhookConfigRepositoryImpl Webhook配置仓库实现
type WebhookConfigRepositoryImpl struct{ db *gorm.DB }

func NewWebhookConfigRepository(db *gorm.DB) repository.WebhookConfigRepository {
	return &WebhookConfigRepositoryImpl{db: db}
}

func (r *WebhookConfigRepositoryImpl) Create(ctx context.Context, cfg *entity.WebhookConfig) error {
	return r.db.WithContext(ctx).Create(cfg).Error
}
func (r *WebhookConfigRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.WebhookConfig, error) {
	var c entity.WebhookConfig
	err := r.db.WithContext(ctx).First(&c, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
func (r *WebhookConfigRepositoryImpl) ListByCreator(ctx context.Context, creatorID uuid.UUID) ([]*entity.WebhookConfig, error) {
	var list []*entity.WebhookConfig
	err := r.db.WithContext(ctx).Where("creator_id = ?", creatorID).Order("created_at DESC").Find(&list).Error
	return list, err
}
func (r *WebhookConfigRepositoryImpl) Update(ctx context.Context, cfg *entity.WebhookConfig) error {
	return r.db.WithContext(ctx).Save(cfg).Error
}
func (r *WebhookConfigRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.WebhookConfig{}, "id = ?", id).Error
}
func (r *WebhookConfigRepositoryImpl) ListActive(ctx context.Context) ([]*entity.WebhookConfig, error) {
	var list []*entity.WebhookConfig
	err := r.db.WithContext(ctx).Where("is_active = true").Find(&list).Error
	return list, err
}

// WebhookLogRepositoryImpl Webhook日志仓库实现
type WebhookLogRepositoryImpl struct{ db *gorm.DB }

func NewWebhookLogRepository(db *gorm.DB) repository.WebhookLogRepository {
	return &WebhookLogRepositoryImpl{db: db}
}

func (r *WebhookLogRepositoryImpl) Create(ctx context.Context, log *entity.WebhookLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
func (r *WebhookLogRepositoryImpl) ListByWebhookID(ctx context.Context, webhookID uuid.UUID, page, pageSize int) ([]*entity.WebhookLog, int64, error) {
	var list []*entity.WebhookLog
	var total int64
	query := r.db.WithContext(ctx).Model(&entity.WebhookLog{}).Where("webhook_id = ?", webhookID)
	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// SubscriptionRepositoryImpl 订阅设置仓库实现
type SubscriptionRepositoryImpl struct{ db *gorm.DB }

func NewSubscriptionRepository(db *gorm.DB) repository.SubscriptionRepository {
	return &SubscriptionRepositoryImpl{db: db}
}

func (r *SubscriptionRepositoryImpl) GetByUserAndEvent(ctx context.Context, userID uuid.UUID, eventType, channel string) (*entity.NotificationSubscription, error) {
	var s entity.NotificationSubscription
	err := r.db.WithContext(ctx).Where("user_id = ? AND event_type = ? AND channel = ?", userID, eventType, channel).First(&s).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}
func (r *SubscriptionRepositoryImpl) ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.NotificationSubscription, error) {
	var list []*entity.NotificationSubscription
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("event_type ASC").Find(&list).Error
	return list, err
}
func (r *SubscriptionRepositoryImpl) Upsert(ctx context.Context, sub *entity.NotificationSubscription) error {
	existing, _ := r.GetByUserAndEvent(ctx, sub.UserID, sub.EventType, sub.Channel)
	if existing != nil {
		sub.ID = existing.ID
		return r.db.WithContext(ctx).Save(sub).Error
	}
	return r.db.WithContext(ctx).Create(sub).Error
}
func (r *SubscriptionRepositoryImpl) BatchUpsert(ctx context.Context, subs []*entity.NotificationSubscription) error {
	for _, sub := range subs {
		if err := r.Upsert(ctx, sub); err != nil {
			return err
		}
	}
	return nil
}
