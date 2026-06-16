package entity

import (
	"github.com/google/uuid"
)

// DocumentAttachment 文档附件
type DocumentAttachment struct {
	BaseModel
	ID         uuid.UUID `gorm:"primaryKey"`
	DocumentID uuid.UUID `gorm:"index;not null"`
	FileName   string    `gorm:"size:255;not null"`
	FileSize   int64
	FileType   string    `gorm:"size:100"`
	FileURL    string    `gorm:"size:500;not null"`
	UploadedBy uuid.UUID
}

// TableName 指定表名
func (DocumentAttachment) TableName() string { return "document_attachments" }
