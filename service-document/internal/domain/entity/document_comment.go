package entity

import (
	"github.com/google/uuid"
)

// DocumentComment 文档评论
type DocumentComment struct {
	BaseModel
	ID        uuid.UUID  `gorm:"primaryKey"`
	DocumentID uuid.UUID `gorm:"index;not null"`
	UserID    uuid.UUID  `gorm:"not null"`
	Content   string     `gorm:"type:text;not null"`
	Position  string     `gorm:"size:50"` // 定位位置（段落引用等�?	ParentID  *uuid.UUID
}

// TableName 指定表名
func (DocumentComment) TableName() string { return "document_comments" }
