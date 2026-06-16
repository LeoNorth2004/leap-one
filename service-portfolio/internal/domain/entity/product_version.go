package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductVersion 产品版本/发布记录
type ProductVersion struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProductID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"product_id"`    // 所属产品ID
	Name        string         `gorm:"size:200;not null" json:"name"`                // 版本号如 v1.0.0
	ReleaseDate *time.Time     `json:"release_date"`                                 // 发布日期
	Status      string         `gorm:"size:20;default:'planning'" json:"status"`     // planning/developing/testing/released/archived
	Description string         `gorm:"type:text" json:"description"`                  // 版本说明
	Plan        string         `gorm:"type:text" json:"plan"`                         // 发布计划
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (ProductVersion) TableName() string {
	return "product_versions"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (pv *ProductVersion) BeforeCreate(tx *gorm.DB) error {
	if pv.ID == uuid.Nil {
		pv.ID = uuid.New()
	}
	return nil
}

// ProductRoadmapItem 产品路线图项
type ProductRoadmapItem struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProductID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"product_id"`    // 所属产品ID
	Title       string         `gorm:"size:300;not null" json:"title"`                // 路线图标题
	Description string         `gorm:"type:text" json:"description"`                   // 描述
	Quarter     string         `gorm:"size:20" json:"quarter"`                         // 季度 Q1/Q2/Q3/Q4
	Year        int            `json:"year"`                                           // 年份
	Status      string         `gorm:"size:20;default:'planning'" json:"status"`       // planning/in_progress/done/cancelled
	Priority    int            `gorm:"default:3" json:"priority"`                      // 优先级 1-5
	SortOrder   int            `gorm:"default:0" json:"sort_order"`                    // 排序
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (ProductRoadmapItem) TableName() string {
	return "product_roadmap_items"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (ri *ProductRoadmapItem) BeforeCreate(tx *gorm.DB) error {
	if ri.ID == uuid.Nil {
		ri.ID = uuid.New()
	}
	return nil
}
