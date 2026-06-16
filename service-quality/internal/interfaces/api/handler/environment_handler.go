package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-quality/internal/domain/entity"
	"leap-one/service-quality/internal/domain/repository"
	"leap-one/service-quality/internal/interfaces/api/dto"
)

// EnvironmentHandler 测试环境管理Handler
type EnvironmentHandler struct {
	envRepo repository.EnvironmentRepository
	logger  *zap.Logger
}

// NewEnvironmentHandler 创建测试环境管理Handler实例
func NewEnvironmentHandler(envRepo repository.EnvironmentRepository, logger *zap.Logger) *EnvironmentHandler {
	return &EnvironmentHandler{
		envRepo: envRepo,
		logger:  logger,
	}
}

// CreateEnvironment 创建环境（POST /api/v1/environments�?func (h *EnvironmentHandler) CreateEnvironment(c *gin.Context) {
	var req dto.CreateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	env := &entity.TestEnvironment{
		Name:        req.Name,
		URL:         req.URL,
		Type:        req.Type,
		OS:          req.OS,
		Browser:     req.Browser,
		Description: req.Description,
	}
	if env.Type == "" {
		env.Type = "dev"
	}
	if req.IsActive != nil {
		env.IsActive = *req.IsActive
	} else {
		env.IsActive = true
	}

	if err := h.envRepo.Create(ctx, env); err != nil {
		h.logger.Error("创建测试环境失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建测试环境失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "测试环境创建成功",
		"env_id":  env.ID.String(),
	})
}

// ListEnvironments 环境列表（GET /api/v1/environments�?func (h *EnvironmentHandler) ListEnvironments(c *gin.Context) {
	includeInactive := c.DefaultQuery("include_inactive", "false") == "true"

	ctx := c.Request.Context()
	envs, err := h.envRepo.List(ctx, includeInactive)
	if err != nil {
		h.logger.Error("查询测试环境列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询测试环境列表失败"})
		return
	}

	list := make([]dto.EnvironmentInfo, len(envs))
	for i, e := range envs {
		list[i] = dto.EnvironmentInfo{
			ID:          e.ID.String(),
			Name:        e.Name,
			URL:         e.URL,
			Type:        e.Type,
			OS:          e.OS,
			Browser:     e.Browser,
			Description: e.Description,
			IsActive:    e.IsActive,
			CreatedAt:   e.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   e.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, dto.EnvironmentListResponse{List: list})
}

// GetEnvironment 环境详情（GET /api/v1/environments/:id�?func (h *EnvironmentHandler) GetEnvironment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的环境ID格式"})
		return
	}

	ctx := c.Request.Context()
	env, err := h.envRepo.GetByID(ctx, id)
	if err != nil || env == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试环境不存�?})
		return
	}

	c.JSON(http.StatusOK, dto.EnvironmentInfo{
		ID:          env.ID.String(),
		Name:        env.Name,
		URL:         env.URL,
		Type:        env.Type,
		OS:          env.OS,
		Browser:     env.Browser,
		Description: env.Description,
		IsActive:    env.IsActive,
		CreatedAt:   env.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   env.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

// UpdateEnvironment 更新环境（PUT /api/v1/environments/:id�?func (h *EnvironmentHandler) UpdateEnvironment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的环境ID格式"})
		return
	}

	var req dto.UpdateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	env, err := h.envRepo.GetByID(ctx, id)
	if err != nil || env == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试环境不存�?})
		return
	}

	if req.Name != nil {
		env.Name = *req.Name
	}
	if req.URL != nil {
		env.URL = *req.URL
	}
	if req.Type != nil {
		env.Type = *req.Type
	}
	if req.OS != nil {
		env.OS = *req.OS
	}
	if req.Browser != nil {
		env.Browser = *req.Browser
	}
	if req.Description != nil {
		env.Description = *req.Description
	}
	if req.IsActive != nil {
		env.IsActive = *req.IsActive
	}

	if err := h.envRepo.Update(ctx, env); err != nil {
		h.logger.Error("更新测试环境失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新测试环境失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "测试环境更新成功"})
}

// DeleteEnvironment 删除环境（DELETE /api/v1/environments/:id�?func (h *EnvironmentHandler) DeleteEnvironment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的环境ID格式"})
		return
	}

	ctx := c.Request.Context()
	env, getErr := h.envRepo.GetByID(ctx, id)
	if getErr != nil || env == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试环境不存�?})
		return
	}

	if err := h.envRepo.Delete(ctx, id); err != nil {
		h.logger.Error("删除测试环境失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除测试环境失败"})
		return
	}

	h.logger.Info("删除测试环境成功", zap.String("env_id", id.String()))
	c.JSON(http.StatusOK, gin.H{"message": "测试环境删除成功"})
}
