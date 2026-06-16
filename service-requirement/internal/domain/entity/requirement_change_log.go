package entity

import (
	"github.com/google/uuid"
)

// RequirementChangeLog 需求变更日�?type RequirementChangeLog struct {
	BaseModel
ID            uuid.UUID `gorm:"primaryKey"`
	RequirementID uuid.UUID `gorm:"index;not null"`
	ChangeType    string    `gorm:"size:30;not null"` // create/update/status_change/priority_change/scope_change
	FieldName     string    `gorm:"size:50"`
	OldValue      string    `gorm:"type:text"`
	NewValue      string    `gorm:"type:text"`
	Reason        string    `gorm:"type:text"` // 变更原因
	ChangeUserID  uuid.UUID `gorm:"not null"`
	ReviewStatus  string    `gorm:"size:20;default:'pending'"` // pending/approved/rejected
}

// TableName 指定表名
func (RequirementChangeLog) TableName() string {
	return "requirement_change_logs"
}
