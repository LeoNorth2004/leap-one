package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"leap-one/service-bi/internal/domain/entity"
	"leap-one/service-bi/internal/domain/repository"
	"go.uber.org/zap"
)

// BIStatService BI统计应用服务 - 协调统计相关的业务逻辑
type BIStatService struct {
	dashboardRepo repository.DashboardConfigRepository
	reportRepo    repository.ReportTemplateRepository
	snapshotRepo  repository.DataSnapshotRepository
	logger        *zap.Logger
}

// NewBIStatService 创建BI统计应用服务实例
func NewBIStatService(
	dashboardRepo repository.DashboardConfigRepository,
	reportRepo repository.ReportTemplateRepository,
	snapshotRepo repository.DataSnapshotRepository,
	logger *zap.Logger,
) *BIStatService {
	return &BIStatService{
		dashboardRepo: dashboardRepo,
		reportRepo:    reportRepo,
		snapshotRepo:  snapshotRepo,
		logger:        logger,
	}
}

// CollectMetricData 用例：采集指定类型的指标数据并保存快�?func (s *BIStatService) CollectMetricData(ctx context.Context, metricType string, value float64, dimensions string) error {
	snapshot := &entity.DataSnapshot{
		MetricType: metricType,
		MetricDate: time.Now().Truncate(24 * time.Hour),
		Value:      value,
		Dimensions: dimensions,
	}

	if err := s.snapshotRepo.Create(ctx, snapshot); err != nil {
		s.logger.Error("采集指标数据失败",
			zap.String("metric_type", metricType),
			zap.Error(err),
		)
		return err
	}

	s.logger.Info("指标数据采集成功",
		zap.String("metric_type", metricType),
		zap.Float64("value", value),
	)
	return nil
}

// GetDashboardConfig 获取大屏配置（含缓存策略�?func (s *BIStatService) GetDashboardConfig(ctx context.Context, dashboardType string) (*entity.DashboardConfig, error) {
	config, err := s.dashboardRepo.GetByType(ctx, dashboardType)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// GetReportData 获取报表数据（聚合计算）
func (s *BIStatService) GetReportData(ctx context.Context, reportID uuid.UUID) (interface{}, error) {
	tpl, err := s.reportRepo.GetByID(ctx, reportID)
	if err != nil || tpl == nil {
		return nil, err
	}

	// 根据报表类型查询对应的数据快�?	now := time.Now()
	monthAgo := now.AddDate(0, -1, 0)
	snapshots, err := s.snapshotRepo.ListByDateRange(ctx, monthAgo, now)
	if err != nil {
		s.logger.Warn("查询数据快照失败，返回空数据", zap.Error(err))
		return map[string]interface{}{"data": []interface{}{}}, nil
	}

	result := map[string]interface{}{
		"report_id":   reportID.String(),
		"report_name": tpl.Name,
		"type":        tpl.Type,
		"snapshots":   snapshots,
	}
	return result, nil
}

// CleanupOldSnapshots 清理过期的历史快照数�?func (s *BIStatService) CleanupOldSnapshots(ctx context.Context, retentionDays int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	count, err := s.snapshotRepo.DeleteByDate(ctx, cutoffDate)
	if err != nil {
		return 0, err
	}
	s.logger.Info("清理过期快照完成",
		zap.Int64("deleted_count", count),
		zap.Int("retention_days", retentionDays),
	)
	return count, nil
}
