package dto

import "github.com/google/uuid"

// NotificationInfo 通知消息信息
type NotificationInfo struct {
	ID         string `json:"id"`
	ReceiverID string `json:"receiver_id"`
	SenderID   string `json:"sender_id,omitempty"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Type       string `json:"type"`
	Channel    string `json:"channel"`
	ActionURL  string `json:"action_url,omitempty"`
	IsRead     bool   `json:"is_read"`
	ReadAt     string `json:"read_at,omitempty"`
	SentAt     string `json:"sent_at"`
	CreatedAt  string `json:"created_at"`
}

// NotificationListResponse 通知列表响应
type NotificationListResponse struct {
	List        []NotificationInfo `json:"list"`
	Total       int64              `json:"total"`
	Page        int                `json:"page"`
	Size        int                `json:"size"`
	UnreadCount int64              `json:"unread_count"`
}

// UnreadCountResponse 未读数量响应
type UnreadCountResponse struct {
	UnreadCount int64 `json:"unread_count"`
}

// TemplateInfo 模板信息
type TemplateInfo struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	Channels  string `json:"channels"`
	EventType string `json:"event_type"`
	IsSystem  bool   `json:"is_system"`
	CreatedAt string `json:"created_at"`
}

// TemplateListResponse 模板列表响应
type TemplateListResponse struct {
	List  []TemplateInfo `json:"list"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

// EmailLogInfo 邮件日志信息
type EmailLogInfo struct {
	ID         string `json:"id"`
	ToAddress  string `json:"to_address"`
	Subject    string `json:"subject"`
	Status     string `json:"status"`
	ErrorMsg   string `json:"error_msg,omitempty"`
	SentAt     string `json:"sent_at,omitempty"`
	RetryCount int    `json:"retry_count"`
	CreatedAt  string `json:"created_at"`
}

// EmailLogListResponse 邮件日志列表响应
type EmailLogListResponse struct {
	List  []EmailLogInfo `json:"list"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

// WebhookInfo Webhook配置信息
type WebhookInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Events    string `json:"events"`
	IsActive  bool   `json:"is_active"`
	CreatorID string `json:"creator_id"`
	CreatedAt string `json:"created_at"`
}

// WebhookLogInfo Webhook调用日志信息
type WebhookLogInfo struct {
	ID         string `json:"id"`
	EventType  string `json:"event_type"`
	RequestURL string `json:"request_url"`
	StatusCode int    `json:"status_code"`
	DurationMs int64  `json:"duration_ms"`
	IsSuccess  bool   `json:"is_success"`
	CreatedAt  string `json:"created_at"`
}

// SubscriptionInfo 订阅设置信息
type SubscriptionInfo struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	EventType string `json:"event_type"`
	Channel   string `json:"channel"`
	Enabled   bool   `json:"enabled"`
}

// UpdateSettingsRequest 更新订阅设置请求
type UpdateSettingsRequest struct {
	Subscriptions []SubscriptionInfo `json:"subscriptions" binding:"required"`
}

// CreateTemplateRequest 创建模板请求
type CreateTemplateRequest struct {
	Code      string `json:"code" binding:"required,max=100"`
	Name      string `json:"name" binding:"required,max=200"`
	Subject   string `json:"subject" binding:"required,max=500"`
	Body      string `json:"body"`
	Channels  string `json:"channels" binding:"omitempty,max=100"`
	EventType string `json:"event_type" binding:"omitempty,max=100"`
}

// UpdateTemplateRequest 更新模板请求
type UpdateTemplateRequest struct {
	Name      *string `json:"name" binding:"omitempty,max=200"`
	Subject   *string `json:"subject" binding:"omitempty,max=500"`
	Body      *string `json:"body"`
	Channels  *string `json:"channels" binding:"omitempty,max=100"`
	EventType *string `json:"event_type" binding:"omitempty,max=100"`
}

// CreateWebhookRequest 创建Webhook请求
type CreateWebhookRequest struct {
	Name      string    `json:"name" binding:"required,max=200"`
	URL       string    `json:"url" binding:"required,max=500"`
	Secret    string    `json:"secret" binding:"omitempty,max=200"`
	Events    string    `json:"events"` // JSON事件列表
	CreatorID uuid.UUID `json:"creator_id" binding:"required"`
}

// UpdateWebhookRequest 更新Webhook请求
type UpdateWebhookRequest struct {
	Name     *string `json:"name" binding:"omitempty,max=200"`
	URL      *string `json:"url" binding:"omitempty,max=500"`
	Secret   *string `json:"secret" binding:"omitempty,max=200"`
	Events   *string `json:"events"`
	IsActive *bool   `json:"is_active"`
}
