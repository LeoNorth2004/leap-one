package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductPlan 产品计划 - 用于管理产品的阶段性计划
type ProductPlan struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProductID uuid.UUID      `gorm:"type:uuid;index;not null" json:"product_id"` // 所属产品ID
	Name      string         `gorm:"size:200;not null" json:"name"`              // 计划名称
	Content   string         `gorm:"type:text" json:"content"`                   // 计划内容
	Status    string         `gorm:"size:20;default:'active'" json:"status"`     // active/completed/cancelled
	StartDate *time.Time     `json:"start_date"`                                 // 开始日期
	EndDate   *time.Time     `json:"end_date"`                                   // 结束日期
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (ProductPlan) TableName() string {
	return "product_plans"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (pp *ProductPlan) BeforeCreate(tx *gorm.DB) error {
	if pp.ID == uuid.Nil {
		pp.ID = uuid.New()
	}
	return nil
}
