package dto

// CreateProductLineRequest 创建产品线请求
type CreateProductLineRequest struct {
	Name        string `json:"name" binding:"required,max=200"` // 产品线名称
	Description string `json:"description"`                     // 描述
	SortOrder   int    `json:"sort_order"`                      // 排序
}

// UpdateProductLineRequest 更新产品线请求
type UpdateProductLineRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=200"`
	Description *string `json:"description"`
	SortOrder   *int    `json:"sort_order"`
	Status      *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// ProductLineInfo 产品线信息
type ProductLineInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
