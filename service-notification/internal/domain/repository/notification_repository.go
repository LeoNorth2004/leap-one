package repository

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-notification/internal/domain/entity"
)

// NotificationRepository 通知消息仓库接口
type NotificationRepository interface {
	Create(ctx context.Context, notification *entity.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error)
	ListByReceiver(ctx context.Context, receiverID uuid.UUID, page, pageSize int, unreadOnly bool) ([]*entity.Notification, int64, error)
	CountUnread(ctx context.Context, receiverID uuid.UUID) (int64, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkAllAsRead(ctx context.Context, receiverID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	BatchCreate(ctx context.Context, notifications []*entity.Notification) error
}

// TemplateRepository 消息模板仓库接口
type TemplateRepository interface {
	Create(ctx context.Context, template *entity.NotificationTemplate) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.NotificationTemplate, error)
	GetByCode(ctx context.Context, code string) (*entity.NotificationTemplate, error)
	List(ctx context.Context, page, pageSize int) ([]*entity.NotificationTemplate, int64, error)
	Update(ctx context.Context, template *entity.NotificationTemplate) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// EmailLogRepository 邮件日志仓库接口
type EmailLogRepository interface {
	Create(ctx context.Context, log *entity.EmailLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.EmailLog, error)
	List(ctx context.Context, page, pageSize int, status string) ([]*entity.EmailLog, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status, errorMsg string) error
	IncrementRetry(ctx context.Context, id uuid.UUID) error
}

// WebhookConfigRepository Webhook配置仓库接口
type WebhookConfigRepository interface {
	Create(ctx context.Context, cfg *entity.WebhookConfig) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.WebhookConfig, error)
	ListByCreator(ctx context.Context, creatorID uuid.UUID) ([]*entity.WebhookConfig, error)
	Update(ctx context.Context, cfg *entity.WebhookConfig) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListActive(ctx context.Context) ([]*entity.WebhookConfig, error)
}

// WebhookLogRepository Webhook日志仓库接口
type WebhookLogRepository interface {
	Create(ctx context.Context, log *entity.WebhookLog) error
	ListByWebhookID(ctx context.Context, webhookID uuid.UUID, page, pageSize int) ([]*entity.WebhookLog, int64, error)
}

// SubscriptionRepository 订阅设置仓库接口
type SubscriptionRepository interface {
	GetByUserAndEvent(ctx context.Context, userID uuid.UUID, eventType, channel string) (*entity.NotificationSubscription, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.NotificationSubscription, error)
	Upsert(ctx context.Context, sub *entity.NotificationSubscription) error
	BatchUpsert(ctx context.Context, subs []*entity.NotificationSubscription) error
}
