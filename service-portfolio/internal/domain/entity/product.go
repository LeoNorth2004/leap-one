package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Product 产品实体 - 核心产品管理模型
type Product struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name          string         `gorm:"size:200;not null" json:"name"`                     // 产品名称
	Code          string         `gorm:"size:50;uniqueIndex;not null" json:"code"`          // 产品编码（唯一）
	ProgramID     *uuid.UUID     `gorm:"type:uuid;index" json:"program_id"`                 // 关联项目集ID
	ProductLineID *uuid.UUID     `gorm:"type:uuid;index" json:"product_line_id"`            // 关联产品线ID
	Description   string         `gorm:"type:text" json:"description"`                       // 产品描述
	OwnerID       uuid.UUID      `gorm:"type:uuid;not null" json:"owner_id"`                // PO负责人ID
	Status        string         `gorm:"size:20;default:'active'" json:"status"`             // active/released/archived
	Type          string         `gorm:"size:20;default:'normal'" json:"type"`              // normal/branch/platform
	Platform      string         `gorm:"size:50" json:"platform"`                            // 支持平台（iOS/Android/Web等）
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (Product) TableName() string {
	return "products"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// ProductLine 产品线 - 用于产品分类管理
type ProductLine struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"size:200;not null" json:"name"`        // 产品线名称
	Description string         `gorm:"type:text" json:"description"`          // 描述
	SortOrder   int            `gorm:"default:0" json:"sort_order"`           // 排序
	Status      string         `gorm:"size:20;default:'active'" json:"status"` // active/inactive
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (ProductLine) TableName() string {
	return "product_lines"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (pl *ProductLine) BeforeCreate(tx *gorm.DB) error {
	if pl.ID == uuid.Nil {
		pl.ID = uuid.New()
	}
	return nil
}

// ProgramProductRelation 项目集-产品多对多关联关系
type ProgramProductRelation struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProgramID uuid.UUID      `gorm:"type:uuid;index;not null" json:"program_id"`  // 项目集ID
	ProductID uuid.UUID      `gorm:"type:uuid;index;not null" json:"product_id"` // 产品ID
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (ProgramProductRelation) TableName() string {
	return "program_product_relations"
}
