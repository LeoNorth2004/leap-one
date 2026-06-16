package dto

import "github.com/google/uuid"

// CreateReportRequest 创建报表请求
type CreateReportRequest struct {
	Name      string    `json:"name" binding:"required,max=200"` // 报表名称
	Type      string    `json:"type" binding:"omitempty,max=50"` // 报表类型
	Config    string    `json:"config" binding:"required"`       // JSON查询配置
	ChartType string    `json:"chart_type" binding:"omitempty,oneof=table bar line pie radar funnel scatter heatmap gauge"`
	CreatorID uuid.UUID `json:"creator_id" binding:"required"` // 创建人ID
}

// UpdateReportRequest 更新报表请求
type UpdateReportRequest struct {
	Name      *string `json:"name" binding:"omitempty,max=200"`
	Type      *string `json:"type" binding:"omitempty,max=50"`
	Config    *string `json:"config"`
	ChartType *string `json:"chart_type" binding:"omitempty,oneof=table bar line pie radar funnel scatter heatmap gauge"`
}

// ReportListResponse 报表列表响应
type ReportListResponse struct {
	List  []ReportInfo `json:"list"`
	Total int64        `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
}

// ReportInfo 报表简要信�?
type ReportInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	ChartType string `json:"chart_type"`
	CreatorID string `json:"creator_id"`
	CreatedAt string `json:"created_at"`
}

// ReportDetailResponse 报表详情响应
type ReportDetailResponse struct {
	ReportInfo
	Config    string `json:"config"`
	UpdatedAt string `json:"updated_at"`
}

// DashboardInfo 大屏简要信�?
type DashboardInfo struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	Layout          string `json:"layout"`
	RefreshInterval int    `json:"refresh_interval"`
	IsSystem        bool   `json:"is_system"`
}

// DashboardDetailResponse 大屏详情响应
type DashboardDetailResponse struct {
	DashboardInfo
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// StatsResponse 统计数据响应
type StatsResponse struct {
	MetricType string                 `json:"metric_type"`
	Data       []StatDataPoint        `json:"data"`
	Summary    map[string]interface{} `json:"summary"`
}

// StatDataPoint 统计数据�?
type StatDataPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
	Label string  `json:"label,omitempty"`
}

// ExportRequest 导出请求
type ExportRequest struct {
	Format string `json:"format" binding:"required,oneof=excel csv pdf"`
}
