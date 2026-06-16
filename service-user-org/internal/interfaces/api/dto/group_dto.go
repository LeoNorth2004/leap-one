package dto

import "github.com/google/uuid"

// CreateGroupRequest 创建用户组请求
type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`         // 用户组名称
	Code        string `json:"code" binding:"required,min=2,max=50,alphanum"` // 用户组编码
	Description string `json:"description" binding:"max=500"`                 // 描述
}

// UpdateGroupRequest 更新用户组请求
type UpdateGroupRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	Status      *int8   `json:"status" binding:"omitempty,oneof=0 1"`
}

// AddMemberRequest 添加成员请求
type AddMemberRequest struct {
	UserIDs []uuid.UUID `json:"user_ids" binding:"required,min=1"` // 要添加的用户ID列表
}

// GroupListResponse 用户组列表响应
type GroupListResponse struct {
	List  []GroupDetailInfo `json:"list"`
	Total int64             `json:"total"`
	Page  int               `json="page"`
	Size  int               `json="size"`
}

// GroupDetailResponse 用户组详情响应
type GroupDetailResponse struct {
	GroupDetailInfo
	Members []UserInfo `json:"members"` // 成员列表
}

// GroupDetailInfo 用户组详细信息
type GroupDetailInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	MemberCount int       `json:"member_count"`
	Status      int8      `json:"status"`
	CreatedAt   string    `json:"created_at"`
}
