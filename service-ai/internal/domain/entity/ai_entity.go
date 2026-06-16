package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AIConversation AI对话会话实体
type AIConversation struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID      `gorm:"index;not null" json:"user_id"`      // 所属用�?
	Title     string         `gorm:"size:300" json:"title"`              // 会话标题
	Model     string         `gorm:"size:50;default:gpt-4" json:"model"` // 使用的AI模型
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Messages  []AIMessage    `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
}

func (AIConversation) TableName() string { return "ai_conversations" }

func (a *AIConversation) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// AIMessage 对话消息实体
type AIMessage struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ConversationID uuid.UUID `gorm:"index;not null" json:"conversation_id"` // 所属会�?
	Role           string    `gorm:"size:20;not null" json:"role"`          // user/assistant/system
	Content        string    `gorm:"type:text;not null" json:"content"`     // 消息内容
	TokenCount     int       `json:"token_count"`                           // Token消�?
	Model          string    `gorm:"size:50" json:"model"`                  // 生成模型
	CreatedAt      time.Time `json:"created_at"`
}

func (AIMessage) TableName() string { return "ai_messages" }

func (a *AIMessage) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// AIPrediction AI预测记录实体
type AIPrediction struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Type       string         `gorm:"size:50;not null" json:"type"`    // 预测类型：requirement_prediction/task_assignment/risk_identification
	TargetID   uuid.UUID      `gorm:"index;not null" json:"target_id"` // 目标对象ID
	Result     string         `gorm:"type:text" json:"result"`         // JSON预测结果
	Confidence float64        `json:"confidence"`                      // 置信�?0-1)
	Model      string         `gorm:"size:50" json:"model"`            // 预测模型
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AIPrediction) TableName() string { return "ai_predictions" }

func (a *AIPrediction) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// AIConfig AI配置实体
type AIConfig struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Provider    string    `gorm:"size:50;not null" json:"provider"`    // openai/azure/anthropic/local
	APIKey      string    `gorm:"size:500" json:"-"`                   // 加密存储（不序列化）
	APIEndpoint string    `gorm:"size:500" json:"api_endpoint"`        // API端点
	Model       string    `gorm:"size:100;default:gpt-4" json:"model"` // 默认模型
	MaxTokens   int       `gorm:"default:2048" json:"max_tokens"`      // 最大Token�?
	Temperature float64   `gorm:"default:0.7" json:"temperature"`      // 温度参数
	IsActive    bool      `gorm:"default:true" json:"is_active"`       // 是否激�?
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (AIConfig) TableName() string { return "ai_configs" }

func (a *AIConfig) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
