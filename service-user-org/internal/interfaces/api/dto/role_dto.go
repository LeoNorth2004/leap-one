package dto

import (
	"github.com/google/uuid"
)

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=50"`         // 角色名称
	Code        string `json:"code" binding:"required,min=2,max=50,alphanum"` // 角色编码
	Type        *int8  `json:"type" binding:"omitempty,oneof=1 2"`            // 类型：1-系统角色 2-自定义角色
	Description string `json:"description" binding:"max=500"`                // 描述
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=50"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	Status      *int8   `json:"status" binding:"omitempty,oneof=0 1"`
}

// AssignPermissionsRequest 为角色分配权限请求
type AssignPermissionsRequest struct {
	PermissionIDs []uuid.UUID `json:"permission_ids" binding:"required,min=1"` // 权限ID列表
}

// RoleListResponse 角色列表响应
type RoleListResponse struct {
	List  []RoleDetailInfo `json:"list"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}

// RoleDetailResponse 角色详情响应
type RoleDetailResponse struct {
	RoleDetailInfo
	Permissions []PermissionInfo `json:"permissions"` // 该角色拥有的权限列表
}

// RoleDetailInfo 角色详细信息
type RoleDetailInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Type        int8      `json:"type"`
	Description string    `json:"description"`
	Status      int8      `json:"status"`
	CreatedAt   string    `json:"created_at"`
}

// PermissionInfo 权限简要信息
type PermissionInfo struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Code     string    `json:"code"`
	Resource string    `json:"resource"`
	Action   string    `json:"action"`
}
