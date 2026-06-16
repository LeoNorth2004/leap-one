package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role 角色实体 - RBAC角色模型
type Role struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(50);not null" json:"name"`                        // 角色名称
	Code        string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`            // 角色编码
	Type        int8           `gorm:"type:smallint;default:1;comment:'1-系统角色 2-自定义角色'" json:"type"` // 角色类型
	Status      int8           `gorm:"type:smallint;default:1;comment:'1-正常 0-禁用'" json:"status"`    // 状态
	Description string         `gorm:"type:text" json:"description"`                                 // 描述
	CreatedBy   uuid.UUID      `gorm:"type:uuid" json:"created_by"`                                  // 创建人
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (Role) TableName() string {
	return "roles"
}
