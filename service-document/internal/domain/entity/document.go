package entity

import (
	"github.com/google/uuid"
)

// Document 文档主表
type Document struct {
	BaseModel
	ID              uuid.UUID        `gorm:"primaryKey"`
	Title           string           `gorm:"size:500;not null"`
	Content         string           `gorm:"type:text"` // Markdown/HTML内容
	Type            string           `gorm:"size:30;default:'markdown'"` // markdown/html/wiki
	CategoryID      *uuid.UUID       // 分类
	ParentID        *uuid.UUID       // 父文档（支持文档树）
	ProductID       *uuid.UUID
	ProjectID       *uuid.UUID
	OwnerID         uuid.UUID        `gorm:"not null"`
	Status          string           `gorm:"size:20;default:'draft'"` // draft/published/archived
	Visibility      string           `gorm:"size:20;default:'public'"` // public/private/team/custom
	PermissionUsers string          `gorm:"type:text"` // JSON自定义权限用户列�?	Version         int              `gorm:"default:1"`
	TemplateID      *uuid.UUID
	Tags            string           `gorm:"type:text"`
	IsTemplate      bool             `gorm:"default:false"` // 是否为模�?	OrderIndex      int              `gorm:"default:0"`
	Children        []Document       `gorm:"foreignKey:ParentID"`
	Versions        []DocumentVersion `gorm:"foreignKey:DocumentID"`
	Attachments     []DocumentAttachment `gorm:"foreignKey:DocumentID"`
}

// TableName 指定表名
func (Document) TableName() string { return "documents" }
