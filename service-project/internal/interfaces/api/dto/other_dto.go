package dto

import (
	"time"

	"github.com/google/uuid"
)

// ==================== 里程碑相关DTO ====================

// CreateMilestoneRequest 创建里程碑请求
type CreateMilestoneRequest struct {
	Name        string     `json:"name" binding:"required,min=1,max=200"` // 里程碑名称
	Description string     `json:"description" binding:"max=2000"`        // 描述
	DueDate     time.Time  `json:"due_date" binding:"required"`           // 截止日期
	SortOrder   int        `json:"sort_order"`                            // 排序序号
}

// UpdateMilestoneRequest 更新里程碑请求
type UpdateMilestoneRequest struct {
	Name        *string    `json:"name" binding:"omitempty,min=1,max=200"`
	Description *string    `json:"description" binding:"omitempty,max=2000"`
	DueDate     *time.Time `json:"due_date"`
	SortOrder   *int       `json:"sort_order"`
}

// MilestoneListResponse 里程碑列表响应
type MilestoneListResponse struct {
	List  []MilestoneInfo `json:"list"`
	Total int64           `json:"total"`
}

// MilestoneInfo 里程碑信息
type MilestoneInfo struct {
	ID          string  `json:"id"`
	ProjectID   string  `json:"project_id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	DueDate     string  `json:"due_date"`
	Status      string  `json:"status"`
	CompletedAt string  `json:"completed_at,omitempty"`
	CompletedBy string  `json:"completed_by,omitempty"`
	SortOrder   int     `json:"sort_order"`
	CreatedAt   string  `json:"created_at"`
}

// ==================== 风险相关DTO ====================

// CreateRiskRequest 创建风险请求
type CreateRiskRequest struct {
	Title       string    `json:"title" binding:"required,min=1,max=300"` // 风险标题
	Description string    `json:"description" binding:"max=2000"`         // 风险描述
	Probability int       `json:"probability" binding:"omitempty,min=1,max=5"` // 发生概率1-5
	Impact      int       `json:"impact" binding:"omitempty,min=1,max=5"`     // 影响程度1-5
	OwnerID     uuid.UUID `json:"owner_id" binding:"required"`                 // 负责人ID
	Mitigation  string    `json:"mitigation" binding:"max=2000"`               // 缓解措施
}

// UpdateRiskRequest 更新风险请求
type UpdateRiskRequest struct {
	Title       *string    `json:"title" binding:"omitempty,min=1,max=300"`
	Description *string    `json:"description" binding:"omitempty,max=2000"`
	Probability *int       `json:"probability" binding:"omitempty,min=1,max=5"`
	Impact      *int       `json:"impact" binding:"omitempty,min=1,max=5"`
	OwnerID     *uuid.UUID `json:"owner_id"`
	Mitigation  *string    `json:"mitigation" binding:"omitempty,max=2000"`
	Status      *string    `json:"status" binding:"omitempty,oneof=open mitigating closed"`
}

// RiskListResponse 风险列表响应
type RiskListResponse struct {
	List  []RiskInfo `json:"list"`
	Total int64      `json:"total"`
}

// RiskInfo 风险信息
type RiskInfo struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Probability int    `json:"probability"`
	Impact      int    `json:"impact"`
	Severity    int    `json:"severity"`
	Status      string `json:"status"`
	OwnerID     string `json:"owner_id"`
	Mitigation  string `json:"mitigation,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ==================== 自定义字段相关DTO ====================

// CreateCustomFieldRequest 创建自定义字段请求
type CreateCustomFieldRequest struct {
	Name      string `json:"name" binding:"required,min=1,max=100"`     // 字段显示名称
	FieldKey  string `json:"field_key" binding:"required,min=1,max=50"` // 字段标识符
	FieldType string `json:"field_type" binding:"required,oneof=text number date select user multi_select"` // 字段类型
	Options   string `json:"options"`                                     // JSON选项（select类型用）
	Required  bool   `json:"required"`                                    // 是否必填
	SortOrder int    `json:"sort_order"`                                  // 排序序号
}

// UpdateCustomFieldRequest 更新自定义字段请求
type UpdateCustomFieldRequest struct {
	Name      *string `json:"name" binding:"omitempty,min=1,max=100"`
	FieldKey  *string `json:"field_key" binding:"omitempty,min=1,max=50"`
	FieldType *string `json:"field_type" binding:"omitempty,oneof=text number date select user multi_select"`
	Options   *string `json:"options"`
	Required  *bool   `json:"required"`
	SortOrder *int    `json:"sort_order"`
}

// CustomFieldListResponse 自定义字段列表响应
type CustomFieldListResponse struct {
	List []CustomFieldInfo `json:"list"`
}

// CustomFieldInfo 自定义字段信息
type CustomFieldInfo struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Name      string `json:"name"`
	FieldKey  string `json:"field_key"`
	FieldType string `json:"field_type"`
	Options   string `json:"options,omitempty"`
	Required  bool   `json:"required"`
	SortOrder int    `json:"sort_order"`
	CreatedAt string `json:"created_at"`
}
