package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification 通知消息实体
type Notification struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	ReceiverID uuid.UUID     `gorm:"index;not null" json:"receiver_id"`       // 接收�?	SenderID  *uuid.UUID     `json:"sender_id"`                               // 发送人（系统通知为空�?	Title     string         `gorm:"size:300;not null" json:"title"`          // 标题
	Content   string         `gorm:"type:text" json:"content"`                // 内容
	Type      string         `gorm:"size:30;default:'system'" json:"type"`    // system/task/requirement/bug/issue/document/mention
	Channel   string         `gorm:"size:20;default:'site'" json:"channel"`   // site/email/webhook
	ActionURL string         `gorm:"size:500" json:"action_url"`              // 点击跳转地址
	IsRead    bool           `gorm:"default:false" json:"is_read"`             // 是否已读
	ReadAt    *time.Time     `json:"read_at"`                                 // 已读时间
	SentAt    time.Time      `gorm:"default:NOW()" json:"sent_at"`            // 发送时�?	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
func (Notification) TableName() string { return "notifications" }
func (n *Notification) BeforeCreate(tx *gorm.DB) error { if n.ID == uuid.Nil { n.ID = uuid.New() }; return nil }

// NotificationTemplate 消息模板实体
type NotificationTemplate struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Code      string         `gorm:"size:100;uniqueIndex;not null" json:"code"` // 模板编码
	Name      string         `gorm:"size:200;not null" json:"name"`             // 模板名称
	Subject   string         `gorm:"size:500;not null" json:"subject"`          // 标题模板（支持变量{{var}}�?	Body      string         `gorm:"type:text" json:"body"`                     // 内容模板
	Channels  string         `gorm:"size:100" json:"channels"`                  // 支持通道 site,email,webhook
	EventType string         `gorm:"size:100" json:"event_type"`               // 触发事件类型
	IsSystem  bool           `gorm:"default:false" json:"is_system"`            // 是否系统内置
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
func (NotificationTemplate) TableName() string { return "notification_templates" }
func (n *NotificationTemplate) BeforeCreate(tx *gorm.DB) error { if n.ID == uuid.Nil { n.ID = uuid.New() }; return nil }

// EmailLog 邮件发送日志实�?type EmailLog struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	ToAddress  string         `gorm:"size:200;not null" json:"to_address"`     // 收件人地址
	Subject    string         `gorm:"size:500;not null" json:"subject"`         // 邮件主题
	Content    string         `gorm:"type:text" json:"content"`                 // 邮件内容
	Status     string         `gorm:"size:20;default:'pending'" json:"status"` // pending/sent/failed
	ErrorMsg   string         `gorm:"type:text" json:"error_msg"`              // 错误信息
	SentAt     *time.Time     `json:"sent_at"`                                  // 发送时�?	RetryCount int            `gorm:"default:0" json:"retry_count"`             // 重试次数
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
func (EmailLog) TableName() string { return "email_logs" }
func (e *EmailLog) BeforeCreate(tx *gorm.DB) error { if e.ID == uuid.Nil { e.ID = uuid.New() }; return nil }

// WebhookConfig Webhook配置实体
type WebhookConfig struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name      string         `gorm:"size:200;not null" json:"name"`             // 配置名称
	URL       string         `gorm:"size:500;not null" json:"url"`              // 回调URL
	Secret    string         `gorm:"size:200" json:"secret"`                    // 签名密钥
	Events    string         `gorm:"type:text" json:"events"`                   // JSON监听事件列表
	IsActive  bool           `gorm:"default:true" json:"is_active"`             // 是否激�?	CreatorID uuid.UUID      `gorm:"not null" json:"creator_id"`                // 创建�?	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
func (WebhookConfig) TableName() string { return "webhook_configs" }
func (w *WebhookConfig) BeforeCreate(tx *gorm.DB) error { if w.ID == uuid.Nil { w.ID = uuid.New() }; return nil }

// NotificationSubscription 通知订阅设置实体
type NotificationSubscription struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID      `gorm:"index;not null" json:"user_id"`             // 用户ID
	EventType string         `gorm:"size:100;not null" json:"event_type"`       // 事件类型
	Channel   string         `gorm:"size:20;not null" json:"channel"`           // 通道：site/email/webhook
	Enabled   bool           `gorm:"default:true" json:"enabled"`              // 是否启用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
func (NotificationSubscription) TableName() string { return "notification_subscriptions" }
func (n *NotificationSubscription) BeforeCreate(tx *gorm.DB) error { if n.ID == uuid.Nil { n.ID = uuid.New() }; return nil }

// WebhookLog Webhook调用日志实体
type WebhookLog struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	WebhookID   uuid.UUID `gorm:"index;not null" json:"webhook_id"`  // 关联Webhook配置
	EventType   string    `gorm:"size:100" json:"event_type"`        // 触发事件
	RequestURL  string    `gorm:"size:500" json:"request_url"`       // 请求URL
	StatusCode  int       `json:"status_code"`                      // HTTP状态码
	Response    string    `gorm:"type:text" json:"response"`        // 响应内容
	DurationMs  int64     `json:"duration_ms"`                      // 耗时(毫秒)
	IsSuccess   bool      `json:"is_success"`                       // 是否成功
	CreatedAt   time.Time `json:"created_at"`
}
func (WebhookLog) TableName() string { return "webhook_logs" }
func (w *WebhookLog) BeforeCreate(tx *gorm.DB) error { if w.ID == uuid.Nil { w.ID = uuid.New() }; return nil }
