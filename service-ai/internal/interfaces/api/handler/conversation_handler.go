package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-ai/internal/domain/entity"
	"leap-one/service-ai/internal/domain/repository"
	"leap-one/service-ai/internal/interfaces/api/dto"
)

// ConversationHandler AI对话管理Handler
type ConversationHandler struct {
	convRepo repository.ConversationRepository
	msgRepo  repository.MessageRepository
	logger   *zap.Logger
}

// NewConversationHandler 创建对话管理Handler实例
func NewConversationHandler(convRepo repository.ConversationRepository, msgRepo repository.MessageRepository, logger *zap.Logger) *ConversationHandler {
	return &ConversationHandler{convRepo: convRepo, msgRepo: msgRepo, logger: logger}
}

// CreateConversation 创建对话 (POST /api/v1/ai/conversations)
func (h *ConversationHandler) CreateConversation(c *gin.Context) {
	var req dto.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	conv := &entity.AIConversation{
		UserID: req.UserID,
		Title:  req.Title,
		Model:  req.Model,
	}
	if conv.Title == "" {
		conv.Title = "新对�?
	}
	if conv.Model == "" {
		conv.Model = "gpt-4"
	}

	ctx := c.Request.Context()
	if err := h.convRepo.Create(ctx, conv); err != nil {
		h.logger.Error("创建对话失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建对话失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "对话创建成功", "conversation_id": conv.ID.String()})
}

// ListConversations 对话列表 (GET /api/v1/ai/conversations)
func (h *ConversationHandler) ListConversations(c *gin.Context) {
	userIDStr := c.Query("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil || userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的user_id"})
		return
	}

	ctx := c.Request.Context()
	convs, total, err := h.convRepo.ListByUserID(ctx, userID, page, size)
	if err != nil {
		h.logger.Error("查询对话列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询对话列表失败"})
		return
	}

	list := make([]dto.ConversationInfo, len(convs))
	for i, conv := range convs {
		list[i] = dto.ConversationInfo{
			ID:           conv.ID.String(),
			UserID:       conv.UserID.String(),
			Title:        conv.Title,
			Model:        conv.Model,
			CreatedAt:    conv.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    conv.UpdatedAt.Format("2006-01-02 15:04:05"),
			MessageCount: len(conv.Messages),
		}
	}

	c.JSON(http.StatusOK, dto.ConversationListResponse{List: list, Total: total, Page: page, Size: size})
}

// GetConversation 对话详情（含消息历史�?GET /api/v1/ai/conversations/:id)
func (h *ConversationHandler) GetConversation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的对话ID格式"})
		return
	}

	ctx := c.Request.Context()
	conv, err := h.convRepo.GetWithMessages(ctx, id)
	if err != nil || conv == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "对话不存�?})
		return
	}

	messages := make([]dto.MessageInfo, len(conv.Messages))
	for i, msg := range conv.Messages {
		messages[i] = dto.MessageInfo{
			ID:         msg.ID.String(),
			Role:       msg.Role,
			Content:    msg.Content,
			TokenCount: msg.TokenCount,
			CreatedAt:  msg.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, dto.ConversationDetailResponse{
		ConversationInfo: dto.ConversationInfo{
			ID:           conv.ID.String(),
			UserID:       conv.UserID.String(),
			Title:        conv.Title,
			Model:        conv.Model,
			CreatedAt:    conv.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    conv.UpdatedAt.Format("2006-01-02 15:04:05"),
			MessageCount: len(conv.Messages),
		},
		Messages: messages,
	})
}

// DeleteConversation 删除对话 (DELETE /api/v1/ai/conversations/:id)
func (h *ConversationHandler) DeleteConversation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的对话ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.convRepo.Delete(ctx, id); err != nil {
		h.logger.Error("删除对话失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除对话失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "对话删除成功"})
}

// SendMessage 发送消息（流式SSE响应�?POST /api/v1/ai/conversations/:id/messages)
func (h *ConversationHandler) SendMessage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的对话ID格式"})
		return
	}

	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 保存用户消息
	userMsg := &entity.AIMessage{
		ConversationID: id,
		Role:           "user",
		Content:        req.Content,
	}
	if err := h.msgRepo.Create(ctx, userMsg); err != nil {
		h.logger.Error("保存用户消息失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送消息失�?})
		return
	}

	// 模拟AI回复（实际应调用AI API�?
	assistantMsg := &entity.AIMessage{
		ConversationID: id,
		Role:           "assistant",
		Content:        "感谢您的提问。我是Leap One AI助手，正在为您处理请求。这是一个模拟回复，实际部署后将会连接真实的AI模型服务�?,
		TokenCount:     42,
		Model:          "gpt-4",
	}
	if err := h.msgRepo.Create(ctx, assistantMsg); err != nil {
		h.logger.Error("保存AI回复失败", zap.Error(err))
	}

	c.JSON(http.StatusOK, gin.H{
		"message_id":  assistantMsg.ID.String(),
		"role":        "assistant",
		"content":     assistantMsg.Content,
		"token_count": assistantMsg.TokenCount,
	})
}

// StreamConnection SSE流式连接 (GET /api/v1/ai/conversations/:id/stream)
func (h *ConversationHandler) StreamConnection(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的对话ID格式"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// SSE握手事件
	c.SSEvent("connected", gin.H{"conversation_id": id.String()})

	// 发送模拟数据流（实际应从AI API获取�?
	c.SSEvent("message", gin.H{
		"role":    "assistant",
		"content": "这是来自AI的流式回复数据�?,
	})

	c.SSEvent("done", gin.H{"status": "completed"})
}
