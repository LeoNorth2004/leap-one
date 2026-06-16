package entity

import (
	"time"

	"github.com/google/uuid"
)

// RequirementReview 需求评审记�?
type RequirementReview struct {
	BaseModel
	ID            uuid.UUID `gorm:"primaryKey"`
	RequirementID uuid.UUID `gorm:"index;not null"`
	Title         string    `gorm:"size:300;not null"`
	MeetingDate   *time.Time
	Status        string                         `gorm:"size:20;default:'planning'"` // planned/completed/cancelled
	Conclusion    string                         `gorm:"type:text"`                  // 评审结论
	Decision      string                         `gorm:"size:20;default:'pending'"`  // approved/rejected/conditional
	CreatorID     uuid.UUID                      `gorm:"not null"`
	Participants  []RequirementReviewParticipant `gorm:"foreignKey:ReviewID"`
}

// TableName 指定表名
func (RequirementReview) TableName() string {
	return "requirement_reviews"
}

// RequirementReviewParticipant 评审参与�?
type RequirementReviewParticipant struct {
	BaseModel
	ID       uuid.UUID `gorm:"primaryKey"`
	ReviewID uuid.UUID `gorm:"index;not null"`
	UserID   uuid.UUID `gorm:"not null"`
	Opinion  string    `gorm:"size:20"` // approve/oppose/abstain
	Comment  string    `gorm:"type:text"`
}

// TableName 指定表名
func (RequirementReviewParticipant) TableName() string {
	return "requirement_review_participants"
}
