package dto

// CreateVersionRequest 创建版本请求
type CreateVersionRequest struct {
	Name        string `json:"name" binding:"required,max=200"` // 版本号如 v1.0.0
	ReleaseDate string `json:"release_date"`                     // 发布日期 YYYY-MM-DD
	Status      string `json:"status" binding:"omitempty,oneof=planning developing testing released archived"`
	Description string `json:"description"`                       // 版本说明
	Plan        string `json:"plan"`                              // 发布计划
}

// UpdateVersionRequest 更新版本请求
type UpdateVersionRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=200"`
	ReleaseDate *string `json:"release_date"`
	Status      *string `json:"status" binding:"omitempty,oneof=planning developing testing released archived"`
	Description *string `json:"description"`
	Plan        *string `json:"plan"`
}

// VersionInfo 版本信息
type VersionInfo struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"product_id"`
	Name        string  `json:"name"`
	ReleaseDate *string `json:"release_date"`
	Status      string  `json:"status"`
	Description string  `json:"description"`
	Plan        string  `json:"plan"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// ReleaseVersionRequest 发布版本请求
type ReleaseVersionRequest struct {
	ReleaseDate string `json:"release_date" binding:"required"` // 实际发布日期
}
