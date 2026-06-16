package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserGroup 用户组实体 - 用于批量管理用户权限
type UserGroup struct {
	ID          uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string        `gorm:"type:varchar(100);not null" json:"name"`                                    // 用户组名称
	Code        string        `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`                          // 用户组编码
	Description string        `gorm:"type:text" json:"description"`                                              // 描述
	MemberCount int           `gorm:"type:int;default:0" json:"member_count"`                                     // 成员数量
	Status      int8          `gorm:"type:smallint;default:1;comment:'1-正常 0-禁用'" json:"status"`              // 状态
	CreatedBy   uuid.UUID     `gorm:"type:uuid" json:"created_by"`                                               // 创建人
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (UserGroup) TableName() string {
	return "user_groups"
}

// UserGroupMember 用户组成员关联表
type UserGroupMember struct {
	UserGroupID uuid.UUID `gorm:"type:uuid;primary_key;column:user_group_id" json:"user_group_id"`
	UserID      uuid.UUID `gorm:"type:uuid;primary_key;column:user_id" json:"user_id"`
	JoinedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"joined_at"`
}

// TableName 指定成员关联表名
func (UserGroupMember) TableName() string {
	return "user_group_members"
}
