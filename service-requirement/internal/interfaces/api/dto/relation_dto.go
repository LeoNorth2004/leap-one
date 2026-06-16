package dto

import (
	"github.com/google/uuid"
)

// CreateRelationRequest 添加关联请求
type CreateRelationRequest struct {
	RelatedType  string    `json:"related_type" binding:"required"` // task/bug/test_case/document
	RelatedID    uuid.UUID `json:"related_id" binding:"required"`
	RelationType string    `json:"relation_type"` // relates_to/depends_on/blocks/duplicates
}

// RelationResponse 关联关系响应
type RelationResponse struct {
	ID            uuid.UUID `json:"id"`
	RequirementID uuid.UUID `json:"requirement_id"`
	RelatedType   string    `json:"related_type"`
	RelatedID     uuid.UUID `json:"related_id"`
	RelationType  string    `json:"relation_type"`
}

// UpdateStatusRequest 更新状态请�?
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
