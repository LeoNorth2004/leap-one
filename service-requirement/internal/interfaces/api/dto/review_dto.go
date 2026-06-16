package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateReviewRequest 创建评审请求
type CreateReviewRequest struct {
	Title        string                     `json:"title" binding:"required,max=300"`
	MeetingDate  *time.Time                 `json:"meeting_date"`
	Participants []ReviewParticipantRequest `json:"participants"`
}

// ReviewParticipantRequest 评审参与者请�?
type ReviewParticipantRequest struct {
	UserID  uuid.UUID `json:"user_id" binding:"required"`
	Opinion string    `json:"opinion"` // approve/oppose/abstain
	Comment string    `json:"comment"`
}

// ReviewResponse 评审记录响应
type ReviewResponse struct {
	ID            uuid.UUID             `json:"id"`
	RequirementID uuid.UUID             `json:"requirement_id"`
	Title         string                `json:"title"`
	MeetingDate   *time.Time            `json:"meeting_date"`
	Status        string                `json:"status"`
	Conclusion    string                `json:"conclusion"`
	Decision      string                `json:"decision"`
	CreatorID     uuid.UUID             `json:"creator_id"`
	CreatedAt     time.Time             `json:"created_at"`
	Participants  []ParticipantResponse `json:"participants,omitempty"`
}

// ParticipantResponse 参与者响�?
type ParticipantResponse struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"user_id"`
	Opinion string    `json:"opinion"`
	Comment string    `json:"comment"`
}
