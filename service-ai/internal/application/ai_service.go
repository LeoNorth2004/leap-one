package application

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"leap-one/service-ai/internal/domain/entity"
	"leap-one/service-ai/internal/domain/repository"
	"go.uber.org/zap"
)

// AIService AI应用服务 - 协调AI相关的业务逻辑
type AIService struct {
	convRepo  repository.ConversationRepository
	msgRepo   repository.MessageRepository
	predRepo  repository.PredictionRepository
	cfgRepo   repository.AIConfigRepository
	logger    *zap.Logger
}

// NewAIService 创建AI应用服务实例
func NewAIService(
	convRepo repository.ConversationRepository,
	msgRepo repository.MessageRepository,
	predRepo repository.PredictionRepository,
	cfgRepo repository.AIConfigRepository,
	logger *zap.Logger,
) *AIService {
	return &AIService{
		convRepo: convRepo,
		msgRepo:  msgRepo,
		predRepo: predRepo,
		cfgRepo:  cfgRepo,
		logger:   logger,
	}
}

// CreateConversationUseCase 创建对话用例
func (s *AIService) CreateConversationUseCase(ctx context.Context, userID uuid.UUID, title, model string) (*entity.AIConversation, error) {
	conv := &entity.AIConversation{
		UserID: userID,
		Title:  title,
		Model:  model,
	}
	if conv.Title == "" { conv.Title = "新对�? }
	if conv.Model == "" { conv.Model = "gpt-4" }

	if err := s.convRepo.Create(ctx, conv); err != nil {
		s.logger.Error("创建对话失败", zap.Error(err))
		return nil, err
	}
	return conv, nil
}

// SendMessageUseCase 发送消息用例（含AI回复生成�?func (s *AIService) SendMessageUseCase(ctx context.Context, conversationID uuid.UUID, content string) (*entity.AIMessage, error) {
	// 保存用户消息
	userMsg := &entity.AIMessage{ConversationID: conversationID, Role: "user", Content: content}
	if err := s.msgRepo.Create(ctx, userMsg); err != nil {
		return nil, err
	}

	// 获取活跃AI配置
	aiCfg, _ := s.cfgRepo.GetActive(ctx)

	// TODO: 实际调用AI API生成回复，此处为模拟
	replyContent := "感谢您的提问。我是Leap One AI助手�?
	model := "gpt-4"
	if aiCfg != nil { model = aiCfg.Model }

	assistantMsg := &entity.AIMessage{
		ConversationID: conversationID,
		Role:           "assistant",
		Content:        replyContent,
		Model:          model,
	}

	if err := s.msgRepo.Create(ctx, assistantMsg); err != nil {
		return nil, err
	}
	return assistantMsg, nil
}

// SavePredictionUseCase 保存预测记录用例
func (s *AIService) SavePredictionUseCase(ctx context.Context, predType string, targetID uuid.UUID, result interface{}, confidence float64) (*entity.AIPrediction, error) {
	resultJSON, _ := json.Marshal(result)

	prediction := &entity.AIPrediction{
		Type:       predType,
		TargetID:   targetID,
		Result:     string(resultJSON),
		Confidence: confidence,
		Model:      "gpt-4",
	}

	if err := s.predRepo.Create(ctx, prediction); err != nil {
		return nil, err
	}
	return prediction, nil
}
