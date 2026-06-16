package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User 用户实体 - 系统用户核心模型
type User struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username       string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`     // 用户名（唯一）
	Password       string         `gorm:"type:varchar(255);not null" json:"-"`                       // 密码（不序列化输出）
	Email          string         `gorm:"type:varchar(100);uniqueIndex" json:"email"`                // 邮箱
	Phone          string         `gorm:"type:varchar(20)" json:"phone"`                             // 手机号
	RealName       string         `gorm:"type:varchar(50)" json:"real_name"`                         // 真实姓名
	Avatar         string         `gorm:"type:text" json:"avatar"`                                   // 头像URL
	Status         int8           `gorm:"type:smallint;default:1;comment:'1-正常 0-禁用'" json:"status"` // 用户状态
	DepartmentID   *uuid.UUID     `gorm:"type:uuid" json:"department_id"`                            // 所属部门ID
	LastLoginAt    *time.Time     `json:"last_login_at"`                                             // 最后登录时间
	LastLoginIP    string         `gorm:"type:varchar(50)" json:"last_login_ip"`                     // 最后登录IP
	PasswordExpire *time.Time     `json:"password_expire"`                                           // 密码过期时间
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定数据库表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子：自动生成UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
