package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TestSuite 测试套件实体 - 用例集合管理
type TestSuite struct {
	ID          uuid.UUID               `gorm:"type:uuid;primary_key" json:"id"`
	Name        string                  `gorm:"type:varchar(200);not null" json:"name"`    // 套件名称
	Description string                  `gorm:"type:text" json:"description"`              // 套件描述
	ProductID   *uuid.UUID              `gorm:"type:uuid" json:"product_id"`               // 关联产品ID
	ProjectID   *uuid.UUID              `gorm:"type:uuid" json:"project_id"`               // 关联项目ID
	CreatorID   uuid.UUID               `gorm:"type:uuid;not null" json:"creator_id"`      // 创建人ID
	Cases       []TestCaseSuiteRelation `gorm:"foreignKey:SuiteID" json:"cases,omitempty"` // 套件中的用例列表
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
	DeletedAt   gorm.DeletedAt          `gorm:"index" json:"-"`
}

// TableName 指定数据库表�?
func (TestSuite) TableName() string {
	return "test_suites"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (s *TestSuite) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// TestCaseSuiteRelation 用例-套件多对多关联中间表实体
type TestCaseSuiteRelation struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	SuiteID   uuid.UUID `gorm:"type:uuid;index;not null" json:"suite_id"` // 套件ID
	CaseID    uuid.UUID `gorm:"type:uuid;index;not null" json:"case_id"`  // 用例ID
	SortOrder int       `gorm:"default:0" json:"sort_order"`              // 排序顺序
}

// TableName 指定关联表名
func (TestCaseSuiteRelation) TableName() string {
	return "test_case_suite_relations"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (r *TestCaseSuiteRelation) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
