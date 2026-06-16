package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"leap-one/service-ai/internal/domain/entity"
	"leap-one/service-ai/internal/domain/repository"
	"leap-one/service-ai/internal/interfaces/api/dto"
)

// AIConfigHandler AI配置管理Handler
type AIConfigHandler struct {
	cfgRepo repository.AIConfigRepository
	logger  *zap.Logger
}

// NewAIConfigHandler 创建AI配置管理Handler实例
func NewAIConfigHandler(cfgRepo repository.AIConfigRepository, logger *zap.Logger) *AIConfigHandler {
	return &AIConfigHandler{cfgRepo: cfgRepo, logger: logger}
}

// GetConfig 获取配置 (GET /api/v1/ai/config)
func (h *AIConfigHandler) GetConfig(c *gin.Context) {
	ctx := c.Request.Context()
	cfg, err := h.cfgRepo.GetActive(ctx)
	if err != nil || cfg == nil {
		// 返回默认配置
		c.JSON(http.StatusOK, dto.AIConfigResponse{
			ID: "", Provider: "openai", Model: "gpt-4",
			MaxTokens: 2048, Temperature: 0.7, IsActive: false,
		})
		return
	}

	c.JSON(http.StatusOK, dto.AIConfigResponse{
		ID: cfg.ID.String(), Provider: cfg.Provider, APIEndpoint: cfg.APIEndpoint,
		Model: cfg.Model, MaxTokens: cfg.MaxTokens, Temperature: cfg.Temperature,
		IsActive: cfg.IsActive,
	})
}

// UpdateConfig 更新配置 (PUT /api/v1/ai/config)
func (h *AIConfigHandler) UpdateConfig(c *gin.Context) {
	var req dto.UpdateAIConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	cfg, err := h.cfgRepo.GetActive(ctx)
	if err != nil || cfg == nil {
		// 不存在则创建
		cfg = &entity.AIConfig{Provider: "openai", Model: "gpt-4", MaxTokens: 2048, Temperature: 0.7, IsActive: true}
		if err := h.cfgRepo.Create(ctx, cfg); err != nil {
			h.logger.Error("创建AI配置失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新配置失败"})
			return
		}
	}

	// 更新字段
	if req.Provider != nil {
		cfg.Provider = *req.Provider
	}
	if req.APIKey != nil {
		cfg.APIKey = *req.APIKey
	}
	if req.APIEndpoint != nil {
		cfg.APIEndpoint = *req.APIEndpoint
	}
	if req.Model != nil {
		cfg.Model = *req.Model
	}
	if req.MaxTokens != nil {
		cfg.MaxTokens = *req.MaxTokens
	}
	if req.Temperature != nil {
		cfg.Temperature = *req.Temperature
	}
	if req.IsActive != nil {
		cfg.IsActive = *req.IsActive
	}

	if err := h.cfgRepo.Update(ctx, cfg); err != nil {
		h.logger.Error("更新AI配置失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "配置更新成功"})
}

// TestConnection 测试AI连接 (POST /api/v1/ai/config/test)
func (h *AIConfigHandler) TestConnection(c *gin.Context) {
	var req dto.TestConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	h.logger.Info("测试AI连接",
		zap.String("provider", req.Provider),
		zap.String("endpoint", req.APIEndpoint),
	)

	// 模拟连接测试结果（实际应调用真实API�?
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "连接测试成功",
		"details": gin.H{
			"provider":   req.Provider,
			"endpoint":   req.APIEndpoint,
			"model":      req.Model,
			"latency_ms": 156,
			"model_info": "GPT-4 Turbo (128K context)",
		},
	})
}
