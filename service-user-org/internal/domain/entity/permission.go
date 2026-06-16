package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permission 权限实体 - RBAC权限模型
type Permission struct {
	ID          uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string        `gorm:"type:varchar(100);not null" json:"name"`                                   // 权限名称
	Code        string        `gorm:"type:varchar(100);uniqueIndex;not null" json:"code"`                        // 权限编码（如 user:create）
	Resource    string        `gorm:"type:varchar(50);not null" json:"resource"`                                // 资源类型（如 user, project）
	Action      string        `gorm:"type:varchar(50);not null" json:"action"`                                  // 操作类型（如 create, read, update, delete）
	Module      string        `gorm:"type:varchar(50)" json:"module"`                                           // 所属模块
	Description string        `gorm:"type:text" json:"description"`                                             // 描述
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (Permission) TableName() string {
	return "permissions"
}

// RolePermission 角色权限关联表（多对多中间表）
type RolePermission struct {
	RoleID       uuid.UUID `gorm:"type:uuid;primary_key;column:role_id" json:"role_id"`
	PermissionID uuid.UUID `gorm:"type:uuid;primary_key;column:permission_id" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 指定关联表名
func (RolePermission) TableName() string {
	return "role_permissions"
}

// UserRole 用户角色关联表（多对多中间表）
type UserRole struct {
	UserID    uuid.UUID `gorm:"type:uuid;primary_key;column:user_id" json:"user_id"`
	RoleID    uuid.UUID `gorm:"type:uuid;primary_key;column:role_id" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定关联表名
func (UserRole) TableName() string {
	return "user_roles"
}
