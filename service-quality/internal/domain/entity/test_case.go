package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TestCase 测试用例实体 - 质量管理核心模型
type TestCase struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Title          string         `gorm:"type:varchar(500);not null" json:"title"` // 用例标题
	Module         string         `gorm:"type:varchar(200)" json:"module"`         // 所属模�?目录
	Precondition   string         `gorm:"type:text" json:"precondition"`           // 前置条件
	Steps          string         `gorm:"type:text" json:"steps"`                  // JSON数组测试步骤
	ExpectedResult string         `gorm:"type:text" json:"expected_result"`        // 预期结果
	Priority       int            `gorm:"default:3" json:"priority"`               // 优先�?1-5
	Type           string         `gorm:"size:30;default:'manual'" json:"type"`    // manual/automated/smoke/regression
	Status         string         `gorm:"size:20;default:'draft'" json:"status"`   // draft/active/archived
	Automation     bool           `gorm:"default:false" json:"automation"`         // 是否自动�?
	ProductID      *uuid.UUID     `gorm:"type:uuid" json:"product_id"`             // 关联产品ID
	ProjectID      *uuid.UUID     `gorm:"type:uuid" json:"project_id"`             // 关联项目ID
	RequirementID  *uuid.UUID     `gorm:"type:uuid" json:"requirement_id"`         // 关联需求ID
	Version        int            `gorm:"default:1" json:"version"`                // 用例版本�?
	CreatorID      uuid.UUID      `gorm:"type:uuid;not null" json:"creator_id"`    // 创建人ID
	ReviewerID     *uuid.UUID     `gorm:"type:uuid" json:"reviewer_id"`            // 评审人ID
	ReviewedAt     *time.Time     `json:"reviewed_at"`                             // 评审时间
	Tags           string         `gorm:"type:text" json:"tags"`                   // 标签（逗号分隔或JSON�?
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表�?
func (TestCase) TableName() string {
	return "test_cases"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (t *TestCase) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
