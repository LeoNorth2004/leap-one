package repository_impl

import (
	"context"
	"time"

	"github.com/google/uuid"
	"leap-one/service-bi/internal/domain/entity"
	"leap-one/service-bi/internal/domain/repository"
	"gorm.io/gorm"
)

// DataSnapshotRepositoryImpl 数据快照仓库实现
type DataSnapshotRepositoryImpl struct {
	db *gorm.DB
}

// NewDataSnapshotRepository 创建数据快照仓库实例
func NewDataSnapshotRepository(db *gorm.DB) repository.DataSnapshotRepository {
	return &DataSnapshotRepositoryImpl{db: db}
}

func (r *DataSnapshotRepositoryImpl) Create(ctx context.Context, snapshot *entity.DataSnapshot) error {
	return r.db.WithContext(ctx).Create(snapshot).Error
}

func (r *DataSnapshotRepositoryImpl) BatchCreate(ctx context.Context, snapshots []*entity.DataSnapshot) error {
	if len(snapshots) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&snapshots).Error
}

func (r *DataSnapshotRepositoryImpl) GetByMetricAndDate(ctx context.Context, metricType string, date time.Time) ([]*entity.DataSnapshot, error) {
	var list []*entity.DataSnapshot
	err := r.db.WithContext(ctx).
		Where("metric_type = ? AND metric_date = ?", metricType, date).
		Order("metric_date ASC").
		Find(&list).Error
	return list, err
}

func (r *DataSnapshotRepositoryImpl) ListByMetricType(ctx context.Context, metricType string, startDate, endDate time.Time) ([]*entity.DataSnapshot, error) {
	var list []*entity.DataSnapshot
	err := r.db.WithContext(ctx).
		Where("metric_type = ? AND metric_date BETWEEN ? AND ?", metricType, startDate, endDate).
		Order("metric_date ASC").
		Find(&list).Error
	return list, err
}

func (r *DataSnapshotRepositoryImpl) ListByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.DataSnapshot, error) {
	var list []*entity.DataSnapshot
	err := r.db.WithContext(ctx).
		Where("metric_date BETWEEN ? AND ?", startDate, endDate).
		Order("metric_type ASC, metric_date ASC").
		Find(&list).Error
	return list, err
}

func (r *DataSnapshotRepositoryImpl) DeleteByDate(ctx context.Context, beforeDate time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("metric_date < ?", beforeDate).
		Delete(&entity.DataSnapshot{})
	return result.RowsAffected, result.Error
}
