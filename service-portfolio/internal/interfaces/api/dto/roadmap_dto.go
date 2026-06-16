package dto

// CreateRoadmapItemRequest 创建路线图项请求
type CreateRoadmapItemRequest struct {
	Title       string `json:"title" binding:"required,max=300"`              // 标题
	Description string `json:"description"`                                   // 描述
	Quarter     string `json:"quarter" binding:"omitempty,oneof=Q1 Q2 Q3 Q4"` // 季度
	Year        int    `json:"year"`                                          // 年份
	Status      string `json:"status" binding:"omitempty,oneof=planning in_progress done cancelled"`
	Priority    int    `json:"priority" binding:"min=1,max=5"` // 优先级 1-5
	SortOrder   int    `json:"sort_order"`                     // 排序
}

// UpdateRoadmapItemRequest 更新路线图项请求
type UpdateRoadmapItemRequest struct {
	Title       *string `json:"title" binding:"omitempty,max=300"`
	Description *string `json:"description"`
	Quarter     *string `json:"quarter" binding:"omitempty,oneof=Q1 Q2 Q3 Q4"`
	Year        *int    `json:"year"`
	Status      *string `json:"status" binding:"omitempty,oneof=planning in_progress done cancelled"`
	Priority    *int    `json:"priority" binding:"omitempty,min=1,max=5"`
	SortOrder   *int    `json:"sort_order"`
}

// RoadmapItemInfo 路线图项信息
type RoadmapItemInfo struct {
	ID          string `json:"id"`
	ProductID   string `json:"product_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Quarter     string `json:"quarter"`
	Year        int    `json:"year"`
	Status      string `json:"status"`
	Priority    int    `json:"priority"`
	SortOrder   int    `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ReorderRoadmapRequest 路线图排序请求
type ReorderRoadmapRequest struct {
	ItemIDs []string `json:"item_ids" binding:"required,min=1"` // 按新顺序排列的路线图项ID列表
}
