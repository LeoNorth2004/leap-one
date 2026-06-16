package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TestEnvironment 测试环境实体 - 环境配置管理
type TestEnvironment struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name        string         `gorm:"type:varchar(200);not null" json:"name"` // 环境名称
	URL         string         `gorm:"type:varchar(500)" json:"url"`           // 环境访问地址
	Type        string         `gorm:"size:30;default:'dev'" json:"type"`      // dev/test/staging/prod
	OS          string         `gorm:"type:varchar(100)" json:"os"`            // 操作系统
	Browser     string         `gorm:"type:varchar(100)" json:"browser"`       // 默认浏览�?
	Description string         `gorm:"type:text" json:"description"`           // 环境描述
	IsActive    bool           `gorm:"default:true" json:"is_active"`          // 是否启用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表�?
func (TestEnvironment) TableName() string {
	return "test_environments"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (e *TestEnvironment) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
