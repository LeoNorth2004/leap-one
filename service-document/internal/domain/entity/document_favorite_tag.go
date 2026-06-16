package entity

import (
	"time"

	"github.com/google/uuid"
)

// DocumentFavorite 文档收藏
type DocumentFavorite struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	UserID     uuid.UUID `gorm:"index;not null"`
	DocumentID uuid.UUID `gorm:"index;not null"`
	CreatedAt  time.Time
}

// TableName 指定表名
func (DocumentFavorite) TableName() string { return "document_favorites" }

// DocumentTag 文档标签
type DocumentTag struct {
	ID     uuid.UUID `gorm:"primaryKey"`
	Name   string    `gorm:"size:50;not null;uniqueIndex"`
	Color  string    `gorm:"size:7"` // 颜色值如 #ff0000
	Count  int       `gorm:"default:0"`
}

// TableName 指定表名
func (DocumentTag) TableName() string { return "document_tags" }
