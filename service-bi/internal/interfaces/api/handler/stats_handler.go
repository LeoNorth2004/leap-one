package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"leap-one/service-bi/internal/interfaces/api/dto"
)

// StatsHandler 统计数据Handler
type StatsHandler struct {
	logger *zap.Logger
}

// NewStatsHandler 创建统计Handler实例
func NewStatsHandler(logger *zap.Logger) *StatsHandler {
	return &StatsHandler{logger: logger}
}

// ProjectProgress 项目进度统计 (GET /api/v1/stats/project-progress)
func (h *StatsHandler) ProjectProgress(c *gin.Context) {
	c.JSON(http.StatusOK, dto.StatsResponse{
		MetricType: "project_progress",
		Summary: map[string]interface{}{
			"total_projects":    25,
			"completed":         18,
			"in_progress":       5,
			"not_started":       2,
			"completion_rate":   72.0,
			"avg_duration_days": 45.5,
		},
		Data: []dto.StatDataPoint{
			{Date: "Q1", Value: 85.0, Label: "完成�?},
			{Date: "Q2", Value: 72.0, Label: "完成�?},
		},
	})
}

// Workload 人员工作量统�?(GET /api/v1/stats/workload)
func (h *StatsHandler) Workload(c *gin.Context) {
	c.JSON(http.StatusOK, dto.StatsResponse{
		MetricType: "workload",
		Summary: map[string]interface{}{
			"total_hours":     12500,
			"completed_hours": 9800,
			"pending_hours":   2700,
			"team_size":       80,
			"avg_per_person":  156.3,
			"overtime_rate":   12.5,
		},
		Data: []dto.StatDataPoint{
			{Date: "开发组", Value: 5200.0, Label: "已完成工�?},
			{Date: "测试�?, Value: 2800.0, Label: "已完成工�?},
			{Date: "产品�?, Value: 1800.0, Label: "已完成工�?},
		},
	})
}

// Quality 产品质量统计 (GET /api/v1/stats/quality)
func (h *StatsHandler) Quality(c *gin.Context) {
	c.JSON(http.StatusOK, dto.StatsResponse{
		MetricType: "quality",
		Summary: map[string]interface{}{
			"total_bugs":        45,
			"resolved":          38,
			"critical":          2,
			"high":              8,
			"medium":            15,
			"low":               20,
			"resolution_rate":   84.4,
			"avg_resolve_hours": 24.5,
		},
		Data: []dto.StatDataPoint{
			{Date: "严重", Value: 2.0, Label: "Bug�?},
			{Date: "�?, Value: 8.0, Label: "Bug�?},
			{Date: "�?, Value: 15.0, Label: "Bug�?},
			{Date: "�?, Value: 20.0, Label: "Bug�?},
		},
	})
}

// RequirementCompletion 需求完成率 (GET /api/v1/stats/requirement-completion)
func (h *StatsHandler) RequirementCompletion(c *gin.Context) {
	c.JSON(http.StatusOK, dto.StatsResponse{
		MetricType: "requirement_completion",
		Summary: map[string]interface{}{
			"total_requirements": 120,
			"completed":          95,
			"in_progress":        15,
			"not_started":        10,
			"completion_rate":    79.2,
			"on_time_rate":       68.3,
		},
		Data: []dto.StatDataPoint{
			{Date: "Q1", Value: 78.0, Label: "完成�?},
			{Date: "Q2", Value: 82.0, Label: "完成�?},
			{Date: "Q3", Value: 75.0, Label: "完成�?},
			{Date: "Q4", Value: 81.5, Label: "完成�?},
		},
	})
}

// BugTrends Bug趋势分析 (GET /api/v1/stats/bug-trends)
func (h *StatsHandler) BugTrends(c *gin.Context) {
	c.JSON(http.StatusOK, dto.StatsResponse{
		MetricType: "bug_trends",
		Summary: map[string]interface{}{
			"trend_direction":  "下降",
			"month_over_month": -15.5,
			"peak_month":       "2026-03",
		},
		Data: []dto.StatDataPoint{
			{Date: "2026-01", Value: 22.0, Label: "新增Bug"},
			{Date: "2026-02", Value: 19.0, Label: "新增Bug"},
			{Date: "2026-03", Value: 28.0, Label: "新增Bug"},
			{Date: "2026-04", Value: 16.0, Label: "新增Bug"},
			{Date: "2026-05", Value: 14.0, Label: "新增Bug"},
			{Date: "2026-06", Value: 11.0, Label: "新增Bug"},
		},
	})
}
