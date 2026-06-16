package repository

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-ai/internal/domain/entity"
)

// ConversationRepository 对话仓库接口
type ConversationRepository interface {
	Create(ctx context.Context, conv *entity.AIConversation) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.AIConversation, error)
	ListByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*entity.AIConversation, int64, error)
	Update(ctx context.Context, conv *entity.AIConversation) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetWithMessages(ctx context.Context, id uuid.UUID) (*entity.AIConversation, error)
}

// MessageRepository 消息仓库接口
type MessageRepository interface {
	Create(ctx context.Context, msg *entity.AIMessage) error
	BatchCreate(ctx context.Context, msgs []*entity.AIMessage) error
	ListByConversation(ctx context.Context, conversationID uuid.UUID, page, pageSize int) ([]*entity.AIMessage, int64, error)
	GetRecentByConversation(ctx context.Context, conversationID uuid.UUID, limit int) ([]*entity.AIMessage, error)
}

// PredictionRepository 预测记录仓库接口
type PredictionRepository interface {
	Create(ctx context.Context, pred *entity.AIPrediction) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.AIPrediction, error)
	List(ctx context.Context, page, pageSize int, predType string, targetID uuid.UUID) ([]*entity.AIPrediction, int64, error)
	ListByTarget(ctx context.Context, targetID uuid.UUID) ([]*entity.AIPrediction, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// AIConfigRepository AI配置仓库接口
type AIConfigRepository interface {
	GetActive(ctx context.Context) (*entity.AIConfig, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.AIConfig, error)
	Update(ctx context.Context, cfg *entity.AIConfig) error
	Create(ctx context.Context, cfg *entity.AIConfig) error
	ListAll(ctx context.Context) ([]*entity.AIConfig, error)
}
