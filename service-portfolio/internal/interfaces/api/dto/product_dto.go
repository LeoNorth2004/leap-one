package dto

import "github.com/google/uuid"

// CreateProductRequest 创建产品请求
type CreateProductRequest struct {
	Name          string  `json:"name" binding:"required,max=200"`                       // 产品名称
	Code          string  `json:"code" binding:"required,max=50"`                        // 产品编码
	ProgramID     *string `json:"program_id"`                                            // 关联项目集ID
	ProductLineID *string `json:"product_line_id"`                                       // 关联产品线ID
	Description   string  `json:"description"`                                           // 描述
	OwnerID       string  `json:"owner_id" binding:"required"`                           // PO负责人ID
	Type          string  `json:"type" binding:"omitempty,oneof=normal branch platform"` // normal/branch/platform
	Platform      string  `json:"platform" binding:"max=50"`                             // 支持平台
}

// UpdateProductRequest 更新产品请求
type UpdateProductRequest struct {
	Name          *string `json:"name" binding:"omitempty,max=200"`
	Code          *string `json:"code" binding:"omitempty,max=50"`
	ProgramID     *string `json:"program_id"`
	ProductLineID *string `json:"product_line_id"`
	Description   *string `json:"description"`
	OwnerID       *string `json:"owner_id"`
	Status        *string `json:"status" binding:"omitempty,oneof=active released archived"`
	Type          *string `json:"type" binding:"omitempty,oneof=normal branch platform"`
	Platform      *string `json:"platform" binding:"omitempty,max=50"`
}

// ProductInfo 产品基本信息
type ProductInfo struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Code          string  `json:"code"`
	ProgramID     *string `json:"program_id"`
	ProductLineID *string `json:"product_line_id"`
	Description   string  `json:"description"`
	OwnerID       string  `json:"owner_id"`
	Status        string  `json:"status"`
	Type          string  `json:"type"`
	Platform      string  `json:"platform"`
}

// ProductListResponse 产品列表响应（分页）
type ProductListResponse struct {
	List  []ProductInfo `json:"list"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
}

// ProductDetailResponse 产品详情响应
type ProductDetailResponse struct {
	ProductInfo
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
	Versions  []VersionInfo     `json:"versions,omitempty"`
	Roadmap   []RoadmapItemInfo `json:"roadmap,omitempty"`
	Plans     []PlanInfo        `json:"plans,omitempty"`
}

// parseOptionalUUID 解析可选的UUID字符串指针
func parseOptionalUUID(s *string) (*uuid.UUID, bool) {
	if s == nil || *s == "" {
		return nil, false
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return nil, false
	}
	return &id, true
}
