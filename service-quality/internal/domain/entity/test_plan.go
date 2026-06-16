package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TestPlan 测试计划实体 - 测试执行计划管理
type TestPlan struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name         string         `gorm:"type:varchar(200);not null" json:"name"`   // 计划名称
	Description  string         `gorm:"type:text" json:"description"`             // 计划描述
	ProductID    *uuid.UUID     `gorm:"type:uuid" json:"product_id"`              // 关联产品ID
	ProjectID    *uuid.UUID     `gorm:"type:uuid" json:"project_id"`              // 关联项目ID
	BuildVersion string         `gorm:"type:varchar(100)" json:"build_version"`   // 对应构建版本
	Status       string         `gorm:"size:20;default:'planning'" json:"status"` // planning/executing/completed/cancelled
	StartDate    *time.Time     `json:"start_date"`                               // 计划开始日�?
	EndDate      *time.Time     `json:"end_date"`                                 // 计划结束日期
	CreatorID    uuid.UUID      `gorm:"type:uuid;not null" json:"creator_id"`     // 创建人ID
	ExecutorIDs  string         `gorm:"type:text" json:"executor_ids"`            // JSON数组执行人ID
	Cases        []TestPlanCase `gorm:"foreignKey:PlanID" json:"cases,omitempty"` // 计划中的用例执行记录
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表�?
func (TestPlan) TableName() string {
	return "test_plans"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (p *TestPlan) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// TestPlanCase 测试计划中的用例执行记录实体
type TestPlanCase struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	PlanID       uuid.UUID  `gorm:"type:uuid;index;not null" json:"plan_id"` // 计划ID
	CaseID       uuid.UUID  `gorm:"type:uuid;index;not null" json:"case_id"` // 用例ID
	AssigneeID   *uuid.UUID `gorm:"type:uuid" json:"assignee_id"`            // 执行人ID
	Result       string     `gorm:"size:20" json:"result"`                   // not_run/passed/failed/blocked/skipped
	ExecuteTime  *time.Time `json:"execute_time"`                            // 执行时间
	ActualResult string     `gorm:"type:text" json:"actual_result"`          // 实际结果描述
	BugIDs       string     `gorm:"type:text" json:"bug_ids"`                // 关联Bug ID列表（JSON数组�?
	Comment      string     `gorm:"type:text" json:"comment"`                // 执行备注
	SortOrder    int        `gorm:"default:0" json:"sort_order"`             // 排序顺序
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TableName 指定数据库表�?
func (TestPlanCase) TableName() string {
	return "test_plan_cases"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (pc *TestPlanCase) BeforeCreate(tx *gorm.DB) error {
	if pc.ID == uuid.Nil {
		pc.ID = uuid.New()
	}
	return nil
}
