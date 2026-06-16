package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"leap-one/service-bi/internal/domain/entity"
)

// DashboardConfigRepository 大屏配置仓库接口
type DashboardConfigRepository interface {
	Create(ctx context.Context, config *entity.DashboardConfig) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.DashboardConfig, error)
	GetByType(ctx context.Context, dashboardType string) (*entity.DashboardConfig, error)
	List(ctx context.Context, page, pageSize int) ([]*entity.DashboardConfig, int64, error)
	Update(ctx context.Context, config *entity.DashboardConfig) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListSystemDashboards(ctx context.Context) ([]*entity.DashboardConfig, error)
}

// ReportTemplateRepository 报表模板仓库接口
type ReportTemplateRepository interface {
	Create(ctx context.Context, template *entity.ReportTemplate) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ReportTemplate, error)
	List(ctx context.Context, page, pageSize int, creatorID uuid.UUID, reportType string) ([]*entity.ReportTemplate, int64, error)
	Update(ctx context.Context, template *entity.ReportTemplate) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByCreator(ctx context.Context, creatorID uuid.UUID) ([]*entity.ReportTemplate, error)
}

// DataSnapshotRepository 数据快照仓库接口
type DataSnapshotRepository interface {
	Create(ctx context.Context, snapshot *entity.DataSnapshot) error
	BatchCreate(ctx context.Context, snapshots []*entity.DataSnapshot) error
	GetByMetricAndDate(ctx context.Context, metricType string, date time.Time) ([]*entity.DataSnapshot, error)
	ListByMetricType(ctx context.Context, metricType string, startDate, endDate time.Time) ([]*entity.DataSnapshot, error)
	ListByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.DataSnapshot, error)
	DeleteByDate(ctx context.Context, beforeDate time.Time) (int64, error)
}
