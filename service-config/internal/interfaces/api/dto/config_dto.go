package dto

// ConfigItem 配置项信�?
type ConfigItem struct {
	ID          string
	Category    string
	Key         string
	Value       string
	ValueType   string
	IsEncrypted bool
	IsPublic    bool
	Description string
	SortOrder   int
	CreatedAt   string
	UpdatedAt   string
}

// ConfigGroupResponse 配置分组响应
type ConfigGroupResponse struct {
	Groups map[string][]ConfigItem `json:"groups"`
}

// UpdateConfigRequest 更新配置请求
type UpdateConfigRequest struct {
	Value string `json:"value"`
}

// BatchUpdateConfigRequest 批量更新配置请求
type BatchUpdateConfigRequest struct {
	Configs []ConfigItemUpdate `json:"configs" binding:"required"`
}
type ConfigItemUpdate struct {
	Category string `json:"category" binding:"required"`
	Key      string `json:"key" binding:"required"`
	Value    string `json:"value"`
}

// FeatureFlagInfo 功能开关信�?
type FeatureFlagInfo struct {
	ID          string
	Key         string
	Name        string
	Description string
	Enabled     bool
	Rules       string
	CreatedAt   string
	UpdatedAt   string
}

// CreateFeatureFlagRequest 创建开关请�?
type CreateFeatureFlagRequest struct {
	Key         string `json:"key" binding:"required,max=200"`
	Name        string `json:"name" binding:"required,max=200"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Rules       string `json:"rules"`
}

// UpdateFeatureFlagRequest 更新开关请�?
type UpdateFeatureFlagRequest struct {
	Name        *string
	Description *string
	Enabled     *bool
	Rules       *string
}

// AuditLogInfo 审计日志信息
type AuditLogInfo struct {
	ID         string
	UserID     string
	Action     string
	Resource   string
	ResourceID string
	Detail     string
	IPAddress  string
	UserAgent  string
	CreatedAt  string
}

// AuditLogListResponse 审计日志列表响应
type AuditLogListResponse struct {
	List  []AuditLogInfo
	Total int64
	Page  int
	Size  int
}
