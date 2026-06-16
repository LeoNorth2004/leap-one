package entity

import (
	"github.com/google/uuid"
)

// KanbanBoard 看板
type KanbanBoard struct {
	BaseModel
	ID          uuid.UUID        `gorm:"primaryKey"`
	Name        string           `gorm:"size:200;not null"`
	Type        string           `gorm:"size:30;default:'project'"` // project/product/personal
	RefID       *uuid.UUID       // 关联的项目或产品ID
	OwnerID     uuid.UUID        `gorm:"not null"`
	Description string           `gorm:"type:text"`
	IsDefault   bool             `gorm:"default:false"`
	Columns     []KanbanColumn   `gorm:"foreignKey:BoardID"`
	Cards       []KanbanCard     `gorm:"foreignKey:BoardID"`
	Swimlanes   []KanbanSwimlane `gorm:"foreignKey:BoardID"`
}

// TableName 指定表名
func (KanbanBoard) TableName() string { return "kanban_boards" }
