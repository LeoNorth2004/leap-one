package dto

import "github.com/google/uuid"

// CreateDepartmentRequest 创建部门请求
type CreateDepartmentRequest struct {
	Name        string     `json:"name" binding:"required,min=2,max=100"`    // 部门名称
	Code        string     `json:"code" binding:"required,min=2,max=50"`     // 部门编码
	ParentID    *uuid.UUID `json:"parent_id"`                                // 上级部门ID
	SortOrder   *int       `json:"sort_order"`                               // 排序
	Leader      *string    `json:"leader" binding:"omitempty,max=50"`        // 负责人
	Phone       *string    `json:"phone" binding:"omitempty,len=20"`         // 联系电话
	Email       *string    `json:"email" binding:"omitempty,email"`          // 部门邮箱
	Description *string    `json:"description" binding:"omitempty,max=1000"` // 描述
}

// UpdateDepartmentRequest 更新部门请求
type UpdateDepartmentRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=100"`
	SortOrder   *int    `json:"sort_order"`
	Leader      *string `json:"leader" binding:"omitempty,max=50"`
	Phone       *string `json:"phone" binding:"omitempty,len=20"`
	Email       *string `json:"email" binding:"omitempty,email"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
	Status      *int8   `json:"status" binding:"omitempty,oneof=0 1"`
}

// DepartmentTreeResponse 部门树形结构响应
type DepartmentTreeResponse struct {
	ID        uuid.UUID                 `json:"id"`
	Name      string                    `json:"name"`
	Code      string                    `json:"code"`
	ParentID  *uuid.UUID                `json:"parent_id"`
	Level     int                       `json:"level"`
	SortOrder int                       `json:"sort_order"`
	Leader    string                    `json:"leader"`
	Status    int8                      `json:"status"`
	Children  []*DepartmentTreeResponse `json:"children"` // 子部门列表
}

// DepartmentListResponse 部门列表响应
type DepartmentListResponse struct {
	List  []*DepartmentInfo `json:"list"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Size  int               `json:"size"`
}

// DepartmentInfo 部门简要信息
type DepartmentInfo struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	ParentID    *uuid.UUID `json:"parent_id"`
	Level       int        `json:"level"`
	Leader      string     `json:"leader"`
	Status      int8       `json:"status"`
	MemberCount int        `json:"member_count"`
}
