package dto

import "github.com/google/uuid"

// CreateConversationRequest 创建对话请求
type CreateConversationRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Title  string    `json:"title" binding:"max=300"`
	Model  string    `json:"model" binding:"omitempty,max=50"`
}

// SendMessageRequest 发送消息请�?
type SendMessageRequest struct {
	Content string `json:"content" binding:"required,max=10000"`
}

// ConversationInfo 对话简要信�?
type ConversationInfo struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	Title        string `json:"title"`
	Model        string `json:"model"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	MessageCount int    `json:"message_count"`
}

// ConversationDetailResponse 对话详情（含消息历史�?
type ConversationDetailResponse struct {
	ConversationInfo
	Messages []MessageInfo `json:"messages"`
}

// MessageInfo 消息信息
type MessageInfo struct {
	ID         string `json:"id"`
	Role       string `json:"role"`
	Content    string `json:"content"`
	TokenCount int    `json:"token_count"`
	CreatedAt  string `json:"created_at"`
}

// ConversationListResponse 对话列表响应
type ConversationListResponse struct {
	List  []ConversationInfo `json:"list"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
}

// AssistRequirementRequest AI辅助编写需求请�?
type AssistRequirementRequest struct {
	Title       string `json:"title" binding:"required,max=500"`
	Description string `json:"description" binding:"required,max=10000"`
	Priority    string `json:"priority" binding:"omitempty,oneof=high medium low"`
}

// AssistTestCaseRequest AI辅助编写测试用例请求
type AssistTestCaseRequest struct {
	RequirementDesc string `json:"requirement_desc" binding:"required,max=10000"`
	Type            string `json:"type" binding:"omitempty,oneof=functional performance security ui"`
}

// TaskAssignSuggestRequest 智能任务分配建议请求
type TaskAssignSuggestRequest struct {
	TaskDescription string `json:"task_description" binding:"required,max=5000"`
	SkillsRequired  string `json:"skills_required,omitempty"` // JSON数组格式
}

// PredictRequirementsRequest 需求预测请�?
type PredictRequirementsRequest struct {
	ProjectID      uuid.UUID `json:"project_id" binding:"required"`
	HistoricalData string    `json:"historical_data,omitempty"`
}

// IdentifyRisksRequest 风险识别请求
type IdentifyRisksRequest struct {
	ProjectID   uuid.UUID `json:"project_id" binding:"required"`
	Description string    `json:"description" binding:"max=5000"`
}

// PredictionInfo 预测记录信息
type PredictionInfo struct {
	ID         string  `json:"id"`
	Type       string  `json:"type"`
	TargetID   string  `json:"target_id"`
	Result     string  `json:"result"`
	Confidence float64 `json:"confidence"`
	Model      string  `json:"model"`
	CreatedAt  string  `json:"created_at"`
}

// PredictionListResponse 预测记录列表响应
type PredictionListResponse struct {
	List  []PredictionInfo `json:"list"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}

// AIConfigResponse AI配置响应（不含APIKey�?
type AIConfigResponse struct {
	ID          string  `json:"id"`
	Provider    string  `json:"provider"`
	APIEndpoint string  `json:"api_endpoint"`
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	IsActive    bool    `json:"is_active"`
}

// UpdateAIConfigRequest 更新AI配置请求
type UpdateAIConfigRequest struct {
	Provider    *string  `json:"provider" binding:"omitempty,oneof=openai azure anthropic local"`
	APIKey      *string  `json:"api_key"`
	APIEndpoint *string  `json:"api_endpoint"`
	Model       *string  `json:"model"`
	MaxTokens   *int     `json:"max_tokens"`
	Temperature *float64 `json:"temperature"`
	IsActive    *bool    `json:"is_active"`
}

// TestConnectionRequest 测试连接请求
type TestConnectionRequest struct {
	Provider    string `json:"provider" binding:"required"`
	APIKey      string `json:"api_key" binding:"required"`
	APIEndpoint string `json:"api_endpoint"`
	Model       string `json:"model"`
}

// AIAssistantResponse AI辅助功能通用响应
type AIAssistantResponse struct {
	Content     string                 `json:"content"`
	Suggestions []string               `json:"suggestions,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}
