package dto

import (
	"github.com/google/uuid"
)

// CreateChangeLogRequest 发起变更请求
type CreateChangeLogRequest struct {
	ChangeType string `json:"change_type" binding:"required"` // create/update/status_change/priority_change/scope_change
	FieldName  string `json:"field_name"`
	OldValue   string `json:"old_value"`
	NewValue   string `json:"new_value"`
	Reason     string `json:"reason"`
	ChangeUserID uuid.UUID `json:"change_user_id" binding:"required"`
}

// ChangeLogResponse 变更日志响应
type ChangeLogResponse struct {
	ID            uuid.UUID `json:"id"`
	RequirementID uuid.UUID `json:"requirement_id"`
	ChangeType    string    `json:"change_type"`
	FieldName     string    `json:"field_name"`
	OldValue      string    `json:"old_value"`
	NewValue      string    `json:"new_value"`
	Reason        string    `json:"reason"`
	ChangeUserID  uuid.UUID `json:"change_user_id"`
	ReviewStatus  string    `json:"review_status"`
	CreatedAt     string    `json:"created_at"`
}

// ChangeLogListResponse 变更日志列表响应
type ChangeLogListResponse struct {
	List []ChangeLogResponse `json:"list"`
	Total int64              `json:"total"`
}
