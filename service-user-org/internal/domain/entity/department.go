package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Department 部门实体 - 组织架构树形结构
type Department struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`                    // 部门名称
	Code        string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`         // 部门编码
	ParentID    *uuid.UUID     `gorm:"type:uuid" json:"parent_id"`                                // 上级部门ID
	Level       int            `gorm:"type:int;default:1;comment:'层级深度'" json:"level"`            // 层级
	SortOrder   int            `gorm:"type:int;default:0" json:"sort_order"`                      // 排序
	Leader      string         `gorm:"type:varchar(50)" json:"leader"`                            // 部门负责人
	Phone       string         `gorm:"type:varchar(20)" json:"phone"`                             // 联系电话
	Email       string         `gorm:"type:varchar(100)" json:"email"`                            // 部门邮箱
	Status      int8           `gorm:"type:smallint;default:1;comment:'1-正常 0-禁用'" json:"status"` // 状态
	Description string         `gorm:"type:text" json:"description"`                              // 描述
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Children    []Department   `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName 指定数据库表名
func (Department) TableName() string {
	return "departments"
}
