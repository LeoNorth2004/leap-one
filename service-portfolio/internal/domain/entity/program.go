package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Program 项目集实体 - 支持多级嵌套的项目集合管理
type Program struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"size:200;not null" json:"name"`                          // 项目集名称
	Code        string         `gorm:"size:50;uniqueIndex;not null" json:"code"`               // 编号（唯一）
	Description string         `gorm:"type:text" json:"description"`                           // 描述
	ParentID    *uuid.UUID     `gorm:"type:uuid" json:"parent_id"`                             // 父项目集ID，支持多级嵌套
	OwnerID     uuid.UUID      `gorm:"type:uuid;not null" json:"owner_id"`                     // 负责人ID
	Status      string         `gorm:"size:20;default:'active'" json:"status"`                 // active/paused/completed/cancelled
	Budget      *float64       `json:"budget"`                                                 // 预算
	StartDate   *time.Time     `json:"start_date"`                                             // 开始日期
	EndDate     *time.Time     `json:"end_date"`                                               // 结束日期
	Priority    int            `gorm:"default:3" json:"priority"`                              // 优先级 1-5，1最高
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Children    []Program      `gorm:"foreignKey:ParentID" json:"children,omitempty"` // 子项目集
}

// TableName 指定数据库表名
func (Program) TableName() string {
	return "programs"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (p *Program) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// Milestone 项目集里程碑
type Milestone struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProgramID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"program_id"`       // 所属项目集
	Name        string         `gorm:"size:200;not null" json:"name"`                     // 里程碑名称
	Description string         `gorm:"type:text" json:"description"`                      // 描述
	DueDate     *time.Time     `json:"due_date"`                                          // 预计完成日期
	Status      string         `gorm:"size:20;default:'pending'" json:"status"`           // pending/completed/overdue
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (Milestone) TableName() string {
	return "milestones"
}

// Risk 项目集风险项
type Risk struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProgramID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"program_id"`        // 所属项目集
	Title       string         `gorm:"size:300;not null" json:"title"`                    // 风险标题
	Description string         `gorm:"type:text" json:"description"`                       // 风险描述
	Probability string         `gorm:"size:20" json:"probability"`                         // 概率：low/medium/high
	Impact      string         `gorm:"size:20" json:"impact"`                              // 影响程度：low/medium/high
	Status      string         `gorm:"size:20;default:'open'" json:"status"`              // open/mitigating/closed
	OwnerID     *uuid.UUID     `gorm:"type:uuid" json:"owner_id"`                         // 风险负责人
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (Risk) TableName() string {
	return "risks"
}
