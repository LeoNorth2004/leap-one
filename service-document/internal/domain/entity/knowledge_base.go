package entity

import (
	"github.com/google/uuid"
)

// KnowledgeBase 知识�?type KnowledgeBase struct {
	BaseModel
	ID          uuid.UUID `gorm:"primaryKey"`
	Name        string    `gorm:"size:200;not null"`
	Description string    `gorm:"type:text"`
	OwnerID     uuid.UUID `gorm:"not null"`
	IsPublic    bool      `gorm:"default:false"`
}

// TableName 指定表名
func (KnowledgeBase) TableName() string { return "knowledge_bases" }
