package entity

import (
	"github.com/google/uuid"
)

// DocumentCategory 文档分类
type DocumentCategory struct {
	BaseModel
	ID        uuid.UUID  `gorm:"primaryKey"`
	Name      string     `gorm:"size:200;not null"`
	ParentID  *uuid.UUID // 支持分类�?
	SortOrder int        `gorm:"default:0"`
}

// TableName 指定表名
func (DocumentCategory) TableName() string { return "document_categories" }
