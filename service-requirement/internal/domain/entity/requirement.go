package entity

import (
	"github.com/google/uuid"
)

// Requirement 需求实体（支持Epic→Feature→Story三级结构�?type Requirement struct {
	BaseModel
	ID             uuid.UUID            `gorm:"primaryKey"`
	Title          string               `gorm:"size:500;not null"`
	Code           string               `gorm:"size:50;uniqueIndex"` // 需求编号如 REQ-001
	Description    string               `gorm:"type:text"`
	Type           string               `gorm:"size:30;default:'story'"` // epic/feature/story
	ParentID       *uuid.UUID           // 父需求（Epic包含Feature，Feature包含Story�?	Level          int                  `gorm:"default:3"` // 1=Epic, 2=Feature, 3=Story
	ProductID      *uuid.UUID
	ProjectID      *uuid.UUID
	Status         string               `gorm:"size:20;default:'draft'"` // draft/reviewing/planning/developing/testing/done/closed/cancelled
	Priority       int                  `gorm:"default:3"` // 1-5
	Source         string               `gorm:"size:30;default:'manual'"` // 来源
	Category       string               `gorm:"size:100"` // 业务需�?用户需�?研发需�?	OwnerID        *uuid.UUID           // 需求负责人(PM/PO)
	ReviewerID     *uuid.UUID           // 评审负责�?	StoryPoints    *float64             // 故事�?	EstimatedHours *float64
	ReleaseVersion string               `gorm:"size:100"` // 目标发布版本
	Stage          string               `gorm:"size:20;default:'requirement'"` // requirement/design/dev/test/release
	SourceURL      string               `gorm:"size:500"` // 原始链接
	Tags           string               `gorm:"type:text"`
	Children       []Requirement        `gorm:"foreignKey:ParentID"`
	Relations      []RequirementRelation `gorm:"foreignKey:RequirementID"`
	Reviews        []RequirementReview  `gorm:"foreignKey:RequirementID"`
	ChangeLogs     []RequirementChangeLog `gorm:"foreignKey:RequirementID"`
}

// TableName 指定表名
func (Requirement) TableName() string {
	return "requirements"
}
