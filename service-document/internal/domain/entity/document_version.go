package entity

import (
	"time"

	"github.com/google/uuid"
)

// DocumentVersion 文档版本
type DocumentVersion struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	DocumentID uuid.UUID `gorm:"index;not null"`
	VersionNo  int       `gorm:"not null"`
	Title      string    `gorm:"size:500"`
	Content    string    `gorm:"type:text"`
	ChangeNote string    `gorm:"size:500"` // 版本说明
	CreatedBy  uuid.UUID `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

// TableName 指定表名
func (DocumentVersion) TableName() string { return "document_versions" }
