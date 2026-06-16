package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Iteration 迭代/Sprint实体 - 敏捷开发迭代管理核心模型
type Iteration struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProjectID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"project_id"`   // 所属项目ID
	Name        string         `gorm:"type:varchar(200);not null" json:"name"`       // 迭代名称（如"Sprint 1"）
	Description string         `gorm:"type:text" json:"description"`                 // 迭代描述
	Status      string         `gorm:"type:varchar(20);default:'planning'" json:"status"` // 状态：planning/active/completed/cancelled
	StartDate   time.Time      `json:"start_date"`                                    // 开始日期
	EndDate     time.Time      `json:"end_date"`                                      // 结束日期
	Capacity    *float64       `json:"capacity"`                                       // 迭代容量（故事点或工时）
	Goal        string         `gorm:"type:varchar(500)" json:"goal"`                 // 迭代目标
	SortOrder   int            `gorm:"default:0" json:"sort_order"`                   // 排序序号
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (Iteration) TableName() string {
	return "iterations"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (i *Iteration) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// IsActive 判断迭代是否处于活跃状态
func (i *Iteration) IsActive() bool {
	return i.Status == "active"
}

// IsCompleted 判断迭代是否已完成
func (i *Iteration) IsCompleted() bool {
	return i.Status == "completed"
}
