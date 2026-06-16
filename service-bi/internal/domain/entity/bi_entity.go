package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DashboardConfig BI大屏配置实体
type DashboardConfig struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name            string         `gorm:"size:200;not null" json:"name"`       // 大屏名称
	Type            string         `gorm:"size:50;not null" json:"type"`        // 类型：company_overview/annual_data/ranking/sprint_burndown/annual_summary
	Layout          string         `gorm:"type:text" json:"layout"`             // JSON布局配置
	RefreshInterval int            `gorm:"default:300" json:"refresh_interval"` // 刷新间隔(�?
	IsSystem        bool           `gorm:"default:false" json:"is_system"`      // 是否系统内置
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (DashboardConfig) TableName() string {
	return "dashboard_configs"
}

// BeforeCreate 创建前自动生成UUID
func (d *DashboardConfig) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// ReportTemplate 自定义报表模板实�?
type ReportTemplate struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name      string         `gorm:"size:200;not null" json:"name"`             // 报表名称
	Type      string         `gorm:"size:50" json:"type"`                       // 类型：project_progress/workload/quality/requirement_completion/bug_trend
	Config    string         `gorm:"type:text" json:"config"`                   // JSON查询配置
	ChartType string         `gorm:"size:30;default:'table'" json:"chart_type"` // 图表类型：table/bar/line/pie/radar/funnel/scatter/heatmap/gauge
	CreatorID uuid.UUID      `gorm:"not null" json:"creator_id"`                // 创建人ID
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ReportTemplate) TableName() string {
	return "report_templates"
}

// BeforeCreate 创建前自动生成UUID
func (r *ReportTemplate) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// DataSnapshot 数据快照实体（定时采集的统计数据�?
type DataSnapshot struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	MetricType string    `gorm:"size:50;not null" json:"metric_type"` // 指标类型
	MetricDate time.Time `gorm:"not null" json:"metric_date"`         // 数据日期
	Value      float64   `json:"value"`                               // 数�?
	Dimensions string    `gorm:"type:text" json:"dimensions"`         // JSON维度信息
	CreatedAt  time.Time `json:"created_at"`
}

// TableName 指定表名
func (DataSnapshot) TableName() string {
	return "data_snapshots"
}

// BeforeCreate 创建前自动生成UUID
func (d *DataSnapshot) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}
