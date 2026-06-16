package repository

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-config/internal/domain/entity"
)

type SystemConfigRepository interface {
	Create(ctx context.Context, cfg *entity.SystemConfig) error
	GetByCategoryAndKey(ctx context.Context, category, key string) (*entity.SystemConfig, error)
	ListByCategory(ctx context.Context, category string) ([]*entity.SystemConfig, error)
	ListAll(ctx context.Context) ([]*entity.SystemConfig, error)
	Update(ctx context.Context, cfg *entity.SystemConfig) error
	BatchUpdate(ctx context.Context, configs []*entity.SystemConfig) error
	GetGroups(ctx context.Context) (map[string][]*entity.SystemConfig, error)
}

type FeatureFlagRepository interface {
	Create(ctx context.Context, f *entity.FeatureFlag) error
	GetByKey(ctx context.Context, key string) (*entity.FeatureFlag, error)
	List(ctx context.Context) ([]*entity.FeatureFlag, error)
	Update(ctx context.Context, f *entity.FeatureFlag) error
	Delete(ctx context.Context, key string) error
	IsEnabled(ctx context.Context, key string) (bool, error)
}

type AuditLogRepository interface {
	Create(ctx context.Context, log *entity.AuditLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.AuditLog, error)
	List(ctx context.Context, page, pageSize int, userID uuid.UUID, action, resource string) ([]*entity.AuditLog, int64, error)
}
