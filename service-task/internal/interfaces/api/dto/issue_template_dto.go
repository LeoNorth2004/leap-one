package dto

import "github.com/google/uuid"

// ==================== 工单模板相关DTO ====================

// CreateTemplateRequest 创建模板请求
type CreateTemplateRequest struct {
	Name      string     `json:"name" binding:"required,max=200"`
	Type      string     `json:"type" binding:"required,max=30"` // bug/feature/request/incident
	Fields    string     `json:"fields"`                        // JSON模板字段配置
	WorkflowID *uuid.UUID `json:"workflow_id"`
}

// UpdateTemplateRequest 更新模板请求
type UpdateTemplateRequest struct {
	Name       *string     `json:"name" binding:"omitempty,max=200"`
	Type       *string     `json:"type" binding:"omitempty,max=30"`
	Fields     *string     `json:"fields"`
	WorkflowID *uuid.UUID `json:"workflow_id"`
}

// TemplateInfo 模板信息
type TemplateInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Fields    string `json:"fields,omitempty"`
	WorkflowID string `json:"workflow_id,omitempty"`
	IsSystem  bool   `json:"is_system"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// TemplateListResponse 模板列表响应
type TemplateListResponse struct {
	List  []TemplateInfo `json:"list"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}
