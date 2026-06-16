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

// TemplateHandler 消息模板管理Handler
type TemplateHandler struct{ tplRepo repository.TemplateRepository; logger *zap.Logger }

func NewTemplateHandler(tplRepo repository.TemplateRepository, logger *zap.Logger) *TemplateHandler {
	return &TemplateHandler{tplRepo: tplRepo, logger: logger}
}

func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	var req dto.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()}); return }
	tpl := &entity.NotificationTemplate{Code: req.Code, Name: req.Name, Subject: req.Subject, Body: req.Body, Channels: req.Channels, EventType: req.EventType}
	ctx := c.Request.Context()
	if err := h.tplRepo.Create(ctx, tpl); err != nil { h.logger.Error("创建模板失败", zap.Error(err)); c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"}); return }
	c.JSON(http.StatusCreated, gin.H{"message": "模板创建成功", "template_id": tpl.ID.String()})
}
func (h *TemplateHandler) ListTemplates(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1")); size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 { page = 1 }; if size < 1 || size > 100 { size = 20 }
	ctx := c.Request.Context()
	list, total, err := h.tplRepo.List(ctx, page, size)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"}); return }
	items := make([]dto.TemplateInfo, len(list))
	for i, t := range list { items[i] = dto.TemplateInfo{ID: t.ID.String(), Code: t.Code, Name: t.Name, Subject: t.Subject, Body: t.Body, Channels: t.Channels, EventType: t.EventType, IsSystem: t.IsSystem, CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05")} }
	c.JSON(http.StatusOK, dto.TemplateListResponse{List: items, Total: total, Page: page, Size: size})
}
func (h *TemplateHandler) GetTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID"}); return }
	ctx := c.Request.Context()
	t, err := h.tplRepo.GetByID(ctx, id)
	if err != nil || t == nil { c.JSON(http.StatusNotFound, gin.H{"error": "模板不存�?}); return }
	c.JSON(http.StatusOK, dto.TemplateInfo{ID: t.ID.String(), Code: t.Code, Name: t.Name, Subject: t.Subject, Body: t.Body, Channels: t.Channels, EventType: t.EventType, IsSystem: t.IsSystem, CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05")})
}
func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID"}); return }
	var req dto.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"}); return }
	ctx := c.Request.Context()
	t, err := h.tplRepo.GetByID(ctx, id)
	if err != nil || t == nil { c.JSON(http.StatusNotFound, gin.H{"error": "模板不存�?}); return }
	if req.Name != nil { t.Name = *req.Name }; if req.Subject != nil { t.Subject = *req.Subject }; if req.Body != nil { t.Body = *req.Body }; if req.Channels != nil { t.Channels = *req.Channels }; if req.EventType != nil { t.EventType = *req.EventType }
	if err := h.tplRepo.Update(ctx, t); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"}); return }
	c.JSON(http.StatusOK, gin.H{"message": "模板更新成功"})
}
func (h *TemplateHandler) DeleteTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID"}); return }
	ctx := c.Request.Context()
	if err := h.tplRepo.Delete(ctx, id); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"}); return }
	c.JSON(http.StatusOK, gin.H{"message": "模板删除成功"})
}

// EmailLogHandler 邮件日志Handler
type EmailLogHandler struct{ logRepo repository.EmailLogRepository; logger *zap.Logger }

func NewEmailLogHandler(logRepo repository.EmailLogRepository, logger *zap.Logger) *EmailLogHandler {
	return &EmailLogHandler{logRepo: logRepo, logger: logger}
}

func (h *EmailLogHandler) ListEmailLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1")); size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	status := c.Query("status"); if page < 1 { page = 1 }; if size < 1 || size > 100 { size = 20 }
	ctx := c.Request.Context()
	list, total, err := h.logRepo.List(ctx, page, size, status)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"}); return }
	items := make([]dto.EmailLogInfo, len(list))
	for i, l := range list {
		items[i] = dto.EmailLogInfo{ID: l.ID.String(), ToAddress: l.ToAddress, Subject: l.Subject, Status: l.Status, ErrorMsg: l.ErrorMsg, RetryCount: l.RetryCount, CreatedAt: l.CreatedAt.Format("2006-01-02 15:04:05")}
		if l.SentAt != nil { items[i].SentAt = l.SentAt.Format("2006-01-02 15:04:05") }
	}
	c.JSON(http.StatusOK, dto.EmailLogListResponse{List: items, Total: total, Page: page, Size: size})
}
func (h *EmailLogHandler) GetEmailLog(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "无效的日志ID"}); return }
	ctx := c.Request.Context()
	log, err := h.logRepo.GetByID(ctx, id)
	if err != nil || log == nil { c.JSON(http.StatusNotFound, gin.H{"error": "日志不存�?}); return }
	resp := dto.EmailLogInfo{ID: log.ID.String(), ToAddress: log.ToAddress, Subject: log.Subject, Status: log.Status, Content: log.Content, ErrorMsg: log.ErrorMsg, RetryCount: log.RetryCount, CreatedAt: log.CreatedAt.Format("2006-01-02 15:04:05")}
	if log.SentAt != nil { resp.SentAt = log.SentAt.Format("2006-01-02 15:04:05") }
	c.JSON(http.StatusOK, resp)
}
func (h *EmailLogHandler) ResendEmail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "无效的日志ID"}); return }
	ctx := c.Request.Context()
	h.logRepo.IncrementRetry(ctx, id)
	h.logRepo.UpdateStatus(ctx, id, "pending", "")
	h.logger.Info("重发邮件请求", zap.String("log_id", id.String()))
	c.JSON(http.StatusOK, gin.H{"message": "邮件重发任务已加入队�?})
}
