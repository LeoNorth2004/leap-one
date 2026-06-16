package entity

import (
	"time"

	"github.com/google/uuid"
)

// KanbanCard 看板卡片
type KanbanCard struct {
	BaseModel
	ID          uuid.UUID  `gorm:"primaryKey"`
	BoardID     uuid.UUID  `gorm:"index;not null"`
	ColumnID    uuid.UUID  `gorm:"index;not null"`
	SwimlaneID  *uuid.UUID // 泳道
	CardType    string     `gorm:"size:20;default:'task'"` // task/requirement/bug
	RefID       uuid.UUID  `gorm:"not null"` // 关联的任�?需�?Bug ID
	Title       string     `gorm:"size:500;not null"`
	Priority    int        `gorm:"default:3"` // 1-5
	AssigneeID  *uuid.UUID
	DueDate     *time.Time
	Tags        string     `gorm:"type:text"`
	BlockReason string     `gorm:"type:text"` // 阻塞原因
	SortOrder   int        `gorm:"default:0"`
	MovedAt     time.Time  // 最后移动时�?	MovedBy     *uuid.UUID
}

// TableName 指定表名
func (KanbanCard) TableName() string { return "kanban_cards" }
