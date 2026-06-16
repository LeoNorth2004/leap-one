package entity

import (
	"time"

	"github.com/google/uuid"
)

// KanbanSwimlane 泳道
type KanbanSwimlane struct {
	BaseModel
	ID        uuid.UUID `gorm:"primaryKey"`
	BoardID   uuid.UUID `gorm:"index;not null"`
	Name      string    `gorm:"size:100;not null"`
	Key       string    `gorm:"size:50;not null"`
	SortOrder int       `gorm:"default:0"`
	Color     string    `gorm:"size:7"`
}

// TableName 指定表名
func (KanbanSwimlane) TableName() string { return "kanban_swimlanes" }

// KanbanCardMoveHistory 卡片移动历史
type KanbanCardMoveHistory struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	CardID   uuid.UUID `gorm:"index;not null"`
	FromColID uuid.UUID `gorm:"not null"`
	ToColID   uuid.UUID `gorm:"not null"`
	MovedBy  uuid.UUID `gorm:"not null"`
	MoveTime time.Time `gorm:"not null"`
}

// TableName 指定表名
func (KanbanCardMoveHistory) TableName() string { return "kanban_card_move_histories" }
