package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateRequirementRequest 创建需求请�?
type CreateRequirementRequest struct {
	Title          string     `json:"title" binding:"required,max=500"`
	Description    string     `json:"description"`
	Type           string     `json:"type"` // epic/feature/story
	ParentID       *uuid.UUID `json:"parent_id"`
	ProductID      *uuid.UUID `json:"product_id"`
	ProjectID      *uuid.UUID `json:"project_id"`
	Priority       int        `json:"priority"`
	Source         string     `json:"source"`
	Category       string     `json:"category"`
	OwnerID        *uuid.UUID `json:"owner_id"`
	ReviewerID     *uuid.UUID `json:"reviewer_id"`
	StoryPoints    *float64   `json:"story_points"`
	EstimatedHours *float64   `json:"estimated_hours"`
	ReleaseVersion string     `json:"release_version"`
	Stage          string     `json:"stage"`
	SourceURL      string     `json:"source_url"`
	Tags           string     `json:"tags"`
}

// UpdateRequirementRequest 更新需求请�?
type UpdateRequirementRequest struct {
	Title          *string    `json:"title" binding:"omitempty,max=500"`
	Description    *string    `json:"description"`
	Type           *string    `json:"type"`
	ParentID       *uuid.UUID `json:"parent_id"`
	ProductID      *uuid.UUID `json:"product_id"`
	ProjectID      *uuid.UUID `json:"project_id"`
	Status         *string    `json:"status"`
	Priority       *int       `json:"priority"`
	Source         *string    `json:"source"`
	Category       *string    `json:"category"`
	OwnerID        *uuid.UUID `json:"owner_id"`
	ReviewerID     *uuid.UUID `json:"reviewer_id"`
	StoryPoints    *float64   `json:"story_points"`
	EstimatedHours *float64   `json:"estimated_hours"`
	ReleaseVersion *string    `json:"release_version"`
	Stage          *string    `json:"stage"`
	SourceURL      *string    `json:"source_url"`
	Tags           *string    `json:"tags"`
}

// RequirementResponse 需求响�?
type RequirementResponse struct {
	ID             uuid.UUID             `json:"id"`
	Code           string                `json:"code"`
	Title          string                `json:"title"`
	Description    string                `json:"description"`
	Type           string                `json:"type"`
	ParentID       *uuid.UUID            `json:"parent_id"`
	Level          int                   `json:"level"`
	ProductID      *uuid.UUID            `json:"product_id"`
	ProjectID      *uuid.UUID            `json:"project_id"`
	Status         string                `json:"status"`
	Priority       int                   `json:"priority"`
	Source         string                `json:"source"`
	Category       string                `json:"category"`
	OwnerID        *uuid.UUID            `json:"owner_id"`
	ReviewerID     *uuid.UUID            `json:"reviewer_id"`
	StoryPoints    *float64              `json:"story_points"`
	EstimatedHours *float64              `json:"estimated_hours"`
	ReleaseVersion string                `json:"release_version"`
	Stage          string                `json:"stage"`
	SourceURL      string                `json:"source_url"`
	Tags           string                `json:"tags"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	Children       []RequirementResponse `json:"children,omitempty"`
}

// RequirementListResponse 需求列表响�?
type RequirementListResponse struct {
	List  []RequirementResponse `json:"list"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
}
