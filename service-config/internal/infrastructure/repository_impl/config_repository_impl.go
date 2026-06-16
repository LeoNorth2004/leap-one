package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-config/internal/domain/entity"
	"leap-one/service-config/internal/domain/repository"
	"gorm.io/gorm"
)

type SystemConfigRepositoryImpl struct {
	db *gorm.DB
}

func NewSystemConfigRepository(db *gorm.DB) repository.SystemConfigRepository {
	return &SystemConfigRepositoryImpl{db: db}
}

func (r *SystemConfigRepositoryImpl) Create(ctx context.Context, cfg *entity.SystemConfig) error {
	return r.db.WithContext(ctx).Create(cfg).Error
}

func (r *SystemConfigRepositoryImpl) GetByCategoryAndKey(ctx context.Context, category, key string) (*entity.SystemConfig, error) {
	var c entity.SystemConfig
	err := r.db.WithContext(ctx).Where("category=? AND key=?", category, key).First(&c).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *SystemConfigRepositoryImpl) ListByCategory(ctx context.Context, category string) ([]*entity.SystemConfig, error) {
	var list []*entity.SystemConfig
	err := r.db.WithContext(ctx).Where("category=?", category).Order("sort_order ASC").Find(&list).Error
	return list, err
}

func (r *SystemConfigRepositoryImpl) ListAll(ctx context.Context) ([]*entity.SystemConfig, error) {
	var list []*entity.SystemConfig
	err := r.db.WithContext(ctx).Order("category ASC, sort_order ASC").Find(&list).Error
	return list, err
}

func (r *SystemConfigRepositoryImpl) Update(ctx context.Context, cfg *entity.SystemConfig) error {
	return r.db.WithContext(ctx).Save(cfg).Error
}

func (r *SystemConfigRepositoryImpl) BatchUpdate(ctx context.Context, configs []*entity.SystemConfig) error {
	for _, cfg := range configs {
		if e := r.db.WithContext(ctx).Save(cfg).Error; e != nil {
			return e
		}
	}
	return nil
}

func (r *SystemConfigRepositoryImpl) GetGroups(ctx context.Context) (map[string][]*entity.SystemConfig, error) {
	var all []*entity.SystemConfig
	if err := r.db.WithContext(ctx).Order("category ASC, sort_order ASC").Find(&all).Error; err != nil {
		return nil, err
	}
	groups := make(map[string][]*entity.SystemConfig)
	for _, cfg := range all {
		groups[cfg.Category] = append(groups[cfg.Category], cfg)
	}
	return groups, nil
}

type FeatureFlagRepositoryImpl struct {
	db *gorm.DB
}

func NewFeatureFlagRepository(db *gorm.DB) repository.FeatureFlagRepository {
	return &FeatureFlagRepositoryImpl{db: db}
}

func (r *FeatureFlagRepositoryImpl) Create(ctx context.Context, f *entity.FeatureFlag) error {
	return r.db.WithContext(ctx).Create(f).Error
}

func (r *FeatureFlagRepositoryImpl) GetByKey(ctx context.Context, key string) (*entity.FeatureFlag, error) {
	var f entity.FeatureFlag
	err := r.db.WithContext(ctx).Where("key=?", key).First(&f).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FeatureFlagRepositoryImpl) List(ctx context.Context) ([]*entity.FeatureFlag, error) {
	var list []*entity.FeatureFlag
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *FeatureFlagRepositoryImpl) Update(ctx context.Context, f *entity.FeatureFlag) error {
	return r.db.WithContext(ctx).Save(f).Error
}

func (r *FeatureFlagRepositoryImpl) Delete(ctx context.Context, key string) error {
	return r.db.WithContext(ctx).Where("key=?", key).Delete(&entity.FeatureFlag{}).Error
}

func (r *FeatureFlagRepositoryImpl) IsEnabled(ctx context.Context, key string) (bool, error) {
	var f entity.FeatureFlag
	err := r.db.WithContext(ctx).Select("enabled").Where("key=?", key).First(&f).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return f.Enabled, nil
}

type AuditLogRepositoryImpl struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) repository.AuditLogRepository {
	return &AuditLogRepositoryImpl{db: db}
}

func (r *AuditLogRepositoryImpl) Create(ctx context.Context, log *entity.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *AuditLogRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.AuditLog, error) {
	var l entity.AuditLog
	err := r.db.WithContext(ctx).First(&l, "id=?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *AuditLogRepositoryImpl) List(ctx context.Context, page, pageSize int, userID uuid.UUID, action, resource string) ([]*entity.AuditLog, int64, error) {
	var list []*entity.AuditLog
	var total int64
	q := r.db.WithContext(ctx).Model(&entity.AuditLog{})
	if userID != uuid.Nil {
		q = q.Where("user_id=?", userID)
	}
	if action != "" {
		q = q.Where("action=?", action)
	}
	if resource != "" {
		q = q.Where("resource=?", resource)
	}
	q.Count(&total)
	offset := (page - 1) * pageSize
	if err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
