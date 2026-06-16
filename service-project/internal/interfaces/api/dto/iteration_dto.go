package dto

import (
	"time"

	"github.com/google/uuid"
)

// ==================== 迭代相关DTO ====================

// CreateIterationRequest 创建迭代请求
type CreateIterationRequest struct {
	ProjectID   uuid.UUID `json:"project_id" binding:"required"`         // 所属项目ID
	Name        string    `json:"name" binding:"required,min=1,max=200"` // 迭代名称
	Description string    `json:"description" binding:"max=2000"`        // 描述
	StartDate   time.Time `json:"start_date" binding:"required"`         // 开始日期
	EndDate     time.Time `json:"end_date" binding:"required"`           // 结束日期
	Capacity    *float64  `json:"capacity"`                              // 容量（故事点/工时）
	Goal        string    `json:"goal" binding:"max=500"`                // 迭代目标
	SortOrder   int       `json:"sort_order"`                            // 排序序号
}

// UpdateIterationRequest 更新迭代请求
type UpdateIterationRequest struct {
	Name        *string    `json:"name" binding:"omitempty,min=1,max=200"`
	Description *string    `json:"description" binding:"omitempty,max=2000"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Capacity    *float64   `json:"capacity"`
	Goal        *string    `json:"goal" binding:"omitempty,max=500"`
	SortOrder   *int       `json:"sort_order"`
}

// IterationListResponse 迭代列表响应
type IterationListResponse struct {
	List  []IterationInfo `json:"list"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
}

// IterationInfo 迭代简要信息
type IterationInfo struct {
	ID          string  `json:"id"`
	ProjectID   string  `json:"project_id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Status      string  `json:"status"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	Capacity    float64 `json:"capacity,omitempty"`
	Goal        string  `json:"goal,omitempty"`
	SortOrder   int     `json:"sort_order"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// IterationDetailResponse 迭代详情响应
type IterationDetailResponse struct {
	IterationInfo
	BoardData interface{}  `json:"board_data"` // 看板数据
	Burndown  BurndownData `json:"burndown"`   // 燃尽图数据
}

// ==================== 燃尽图数据 ====================

// BurndownData 燃尽图数据结构
type BurndownData struct {
	SprintName  string          `json:"sprint_name"`  // 迭代名称
	TotalPoints float64         `json:"total_points"` // 总故事点
	IdealLine   []BurndownPoint `json:"ideal_line"`   // 理想燃尽线
	ActualLine  []BurndownPoint `json:"actual_line"`  // 实际燃尽线
	Remaining   float64         `json:"remaining"`    // 剩余工作量
}

// BurndownPoint 燃尽图数据点
type BurndownPoint struct {
	Date  string  `json:"date"`  // 日期
	Value float64 `json:"value"` // 工作量值
}

// ==================== 甘特图数据 ====================

// GanttData 甘特图数据结构
type GanttData struct {
	ProjectName string           `json:"project_name"`
	Tasks       []GanttTask      `json:"tasks"`
	Milestones  []GanttMilestone `json:"milestones"`
}

// GanttTask 甘特图任务项
type GanttTask struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Progress  int    `json:"progress"` // 进度百分比
	Status    string `json:"status"`
	DependsOn string `json:"depends_on,omitempty"` // 依赖任务ID
	Color     string `json:"color,omitempty"`
}

// GanttMilestone 甘特图里程碑
type GanttMilestone struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	DueDate string `json:"due_date"`
	Status  string `json:"status"`
}
