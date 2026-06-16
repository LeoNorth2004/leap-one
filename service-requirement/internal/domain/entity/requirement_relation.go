package entity

import (
	"github.com/google/uuid"
)

// RequirementRelation 需求关�?type RequirementRelation struct {
	BaseModel
	ID            uuid.UUID `gorm:"primaryKey"`
	RequirementID uuid.UUID `gorm:"index;not null"`
	RelatedType   string    `gorm:"size:20;not null"` // task/bug/test_case/document
	RelatedID     uuid.UUID `gorm:"index;not null"`
	RelationType  string    `gorm:"size:20"` // relates_to/depends_on/blocks/duplicates
}

// TableName 指定表名
func (RequirementRelation) TableName() string {
	return "requirement_relations"
}
