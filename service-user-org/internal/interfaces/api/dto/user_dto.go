package dto

import "github.com/google/uuid"

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username     string     `json:"username" binding:"required,min=3,max=50,alphanum"` // 用户名
	Password     string     `json:"password" binding:"required,min=8,max=50"`          // 密码
	Email        string     `json:"email" binding:"omitempty,email"`                   // 邮箱
	Phone        string     `json:"phone" binding:"omitempty,len=11"`                  // 手机号
	RealName     string     `json:"real_name" binding:"max=50"`                        // 真实姓名
	Avatar       string     `json:"avatar" binding:"omitempty,url"`                    // 头像URL
	DepartmentID *uuid.UUID `json:"department_id"`                                     // 部门ID
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email        *string    `json:"email" binding:"omitempty,email"`
	Phone        *string    `json:"phone" binding:"omitempty,len=11"`
	RealName     *string    `json:"real_name" binding:"omitempty,max=50"`
	Avatar       *string    `json:"avatar" binding:"omitempty,url"`
	DepartmentID *uuid.UUID `json:"department_id"`
	Status       *int8      `json:"status" binding:"omitempty,oneof=0 1"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`              // 旧密码
	NewPassword string `json:"new_password" binding:"required,min=8,max=50"` // 新密码
}

// UserListResponse 用户列表响应（分页）
type UserListResponse struct {
	List  []UserInfo `json:"list"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

// UserDetailResponse 用户详情响应
type UserDetailResponse struct {
	UserInfo
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
	Roles     []RoleInfo  `json:"roles"`
	Groups    []GroupInfo `json:"groups"`
}

// RoleInfo 角色简要信息
type RoleInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// GroupInfo 用户组简要信息
type GroupInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
