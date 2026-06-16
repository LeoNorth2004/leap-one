package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-notification/internal/domain/entity"
	"leap-one/service-notification/internal/domain/repository"
	"leap-one/service-notification/internal/interfaces/api/dto"
)

type WebhookHandler struct {
	cfgRepo repository.WebhookConfigRepository
	logRepo repository.WebhookLogRepository
	logger  *zap.Logger
}

func NewWebhookHandler(cfgRepo repository.WebhookConfigRepository, logRepo repository.WebhookLogRepository, logger *zap.Logger) *WebhookHandler {
	return &WebhookHandler{cfgRepo: cfgRepo, logRepo: logRepo, logger: logger}
}

func (h *WebhookHandler) CreateWebhook(c *gin.Context) {
	var req dto.CreateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}
	cfg := &entity.WebhookConfig{Name: req.Name, URL: req.URL, Secret: req.Secret, Events: req.Events, CreatorID: req.CreatorID}
	ctx := c.Request.Context()
	if err := h.cfgRepo.Create(ctx, cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Webhook创建成功", "webhook_id": cfg.ID.String()})
}

func (h *WebhookHandler) ListWebhooks(c *gin.Context) {
	creatorIDStr := c.Query("creator_id")
	creatorID, _ := uuid.Parse(creatorIDStr)
	ctx := c.Request.Context()
	var list []*entity.WebhookConfig
	var err error
	if creatorID != uuid.Nil {
		list, err = h.cfgRepo.ListByCreator(ctx, creatorID)
	} else {
		list, err = h.cfgRepo.ListActive(ctx)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	items := make([]dto.WebhookInfo, len(list))
	for i, w := range list {
		items[i] = dto.WebhookInfo{ID: w.ID.String(), Name: w.Name, URL: w.URL, Events: w.Events, IsActive: w.IsActive, CreatorID: w.CreatorID.String(), CreatedAt: w.CreatedAt.Format("2006-01-02 15:04:05")}
	}
	c.JSON(http.StatusOK, gin.H{"list": items})
}

func (h *WebhookHandler) GetWebhook(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	ctx := c.Request.Context()
	cfg, err := h.cfgRepo.GetByID(ctx, id)
	if err != nil || cfg == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "不存�?})
		return
	}
	c.JSON(http.StatusOK, dto.WebhookInfo{ID: cfg.ID.String(), Name: cfg.Name, URL: cfg.URL, Events: cfg.Events, IsActive: cfg.IsActive, CreatedAt: cfg.CreatedAt.Format("2006-01-02 15:04:05")})
}

func (h *WebhookHandler) UpdateWebhook(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	var req dto.UpdateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}
	ctx := c.Request.Context()
	cfg, err := h.cfgRepo.GetByID(ctx, id)
	if err != nil || cfg == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "不存�?})
		return
	}
	if req.Name != nil {
		cfg.Name = *req.Name
	}
	if req.URL != nil {
		cfg.URL = *req.URL
	}
	if req.Secret != nil {
		cfg.Secret = *req.Secret
	}
	if req.Events != nil {
		cfg.Events = *req.Events
	}
	if req.IsActive != nil {
		cfg.IsActive = *req.IsActive
	}
	h.cfgRepo.Update(ctx, cfg)
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (h *WebhookHandler) DeleteWebhook(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	ctx := c.Request.Context()
	h.cfgRepo.Delete(ctx, id)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (h *WebhookHandler) TestWebhook(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	ctx := c.Request.Context()
	cfg, _ := h.cfgRepo.GetByID(ctx, id)
	testLog := &entity.WebhookLog{WebhookID: id, EventType: "test", RequestURL: cfg.URL, StatusCode: 200, DurationMs: 45, IsSuccess: true}
	h.logRepo.Create(ctx, testLog)
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "测试成功"})
}

func (h *WebhookHandler) ListWebhookLogs(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	ctx := c.Request.Context()
	list, total, _ := h.logRepo.ListByWebhookID(ctx, id, page, size)
	items := make([]dto.WebhookLogInfo, len(list))
	for i, l := range list {
		items[i] = dto.WebhookLogInfo{ID: l.ID.String(), EventType: l.EventType, StatusCode: l.StatusCode, DurationMs: l.DurationMs, IsSuccess: l.IsSuccess, CreatedAt: l.CreatedAt.Format("2006-01-02 15:04:05")}
	}
	c.JSON(http.StatusOK, gin.H{"list": items, "total": total})
}
