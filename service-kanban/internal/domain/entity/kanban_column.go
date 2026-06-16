package entity

import (
	"github.com/google/uuid"
)

// KanbanColumn 看板�?type KanbanColumn struct {
	BaseModel
	ID        uuid.UUID `gorm:"primaryKey"`
	BoardID   uuid.UUID `gorm:"index;not null"`
	Name      string    `gorm:"size:100;not null"` // �?"待办"�?进行�?�?已完�?
	Key       string    `gorm:"size:50;not null"` // �?"todo", "doing", "done"
	WIPLimit  *int      // WIP限制（在制品数量上限�?	Color     string    `gorm:"size:7"` // 列颜�?	SortOrder int       `gorm:"default:0"`
	Type      string    `gorm:"size:20;default:'normal'"` // normal/backlog/done(完成�?
}

// TableName 指定表名
func (KanbanColumn) TableName() string { return "kanban_columns" }
