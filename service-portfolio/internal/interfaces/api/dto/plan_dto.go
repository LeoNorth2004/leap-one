package dto

// CreatePlanRequest 创建计划请求
type CreatePlanRequest struct {
	ProductID string `json:"product_id" binding:"required"`   // 产品ID
	Name      string `json:"name" binding:"required,max=200"` // 计划名称
	Content   string `json:"content"`                         // 计划内容
	Status    string `json:"status" binding:"omitempty,oneof=active completed cancelled"`
	StartDate string `json:"start_date"` // 开始日期 YYYY-MM-DD
	EndDate   string `json:"end_date"`   // 结束日期 YYYY-MM-DD
}

// UpdatePlanRequest 更新计划请求
type UpdatePlanRequest struct {
	Name      *string `json:"name" binding:"omitempty,max=200"`
	Content   *string `json:"content"`
	Status    *string `json:"status" binding:"omitempty,oneof=active completed cancelled"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
}

// PlanInfo 计划信息
type PlanInfo struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Content   string  `json:"content"`
	Status    string  `json:"status"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// PlanListResponse 计划列表响应（分页）
type PlanListResponse struct {
	List  []PlanInfo `json:"list"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}
