package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-bi/internal/domain/entity"
	"leap-one/service-bi/internal/domain/repository"
	"gorm.io/gorm"
)

// DashboardConfigRepositoryImpl 大屏配置仓库实现
type DashboardConfigRepositoryImpl struct {
	db *gorm.DB
}

// NewDashboardConfigRepository 创建大屏配置仓库实例
func NewDashboardConfigRepository(db *gorm.DB) repository.DashboardConfigRepository {
	return &DashboardConfigRepositoryImpl{db: db}
}

func (r *DashboardConfigRepositoryImpl) Create(ctx context.Context, config *entity.DashboardConfig) error {
	return r.db.WithContext(ctx).Create(config).Error
}

func (r *DashboardConfigRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.DashboardConfig, error) {
	var cfg entity.DashboardConfig
	err := r.db.WithContext(ctx).First(&cfg, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cfg, nil
}

func (r *DashboardConfigRepositoryImpl) GetByType(ctx context.Context, dashboardType string) (*entity.DashboardConfig, error) {
	var cfg entity.DashboardConfig
	err := r.db.WithContext(ctx).Where("type = ?", dashboardType).First(&cfg).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cfg, nil
}

func (r *DashboardConfigRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*entity.DashboardConfig, int64, error) {
	var list []*entity.DashboardConfig
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.DashboardConfig{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *DashboardConfigRepositoryImpl) Update(ctx context.Context, config *entity.DashboardConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}

func (r *DashboardConfigRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.DashboardConfig{}, "id = ?", id).Error
}

func (r *DashboardConfigRepositoryImpl) ListSystemDashboards(ctx context.Context) ([]*entity.DashboardConfig, error) {
	var list []*entity.DashboardConfig
	err := r.db.WithContext(ctx).Where("is_system = true").Order("type ASC").Find(&list).Error
	return list, err
}
