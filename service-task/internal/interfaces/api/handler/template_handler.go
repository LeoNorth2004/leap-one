package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-task/internal/application"
	"leap-one/service-task/internal/interfaces/api/dto"
)

// TemplateHandler 工单模板管理Handler
type TemplateHandler struct {
	templateSvc *application.TemplateService
	logger      *zap.Logger
}

// NewTemplateHandler 创建模板管理Handler实例
func NewTemplateHandler(templateSvc *application.TemplateService, logger *zap.Logger) *TemplateHandler {
	return &TemplateHandler{
		templateSvc: templateSvc,
		logger:      logger,
	}
}

// CreateTemplate 创建模板（POST /api/v1/issue-templates）
func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	var req dto.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	tmpl, err := h.templateSvc.CreateTemplate(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "模板创建成功",
		"template_id": tmpl.ID.String(),
	})
}

// GetTemplate 获取模板详情（GET /api/v1/issue-templates/:id）
func (h *TemplateHandler) GetTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID格式"})
		return
	}

	ctx := c.Request.Context()
	tmpl, svcErr := h.templateSvc.GetTemplate(ctx, id)
	if svcErr != nil || tmpl == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
		return
	}

	resp := dto.TemplateInfo{
		ID:        tmpl.ID.String(),
		Name:      tmpl.Name,
		Type:      tmpl.Type,
		Fields:    tmpl.Fields,
		IsSystem:  tmpl.IsSystem,
		CreatedAt: tmpl.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: tmpl.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if tmpl.WorkflowID != nil {
		resp.WorkflowID = tmpl.WorkflowID.String()
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateTemplate 更新模板（PUT /api/v1/issue-templates/:id）
func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID格式"})
		return
	}

	var req dto.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	svcErr := h.templateSvc.UpdateTemplate(ctx, id, &req)
	if svcErr != nil {
		if svcErr == application.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "模板更新成功"})
}

// DeleteTemplate 删除模板（DELETE /api/v1/issue-templates/:id）
func (h *TemplateHandler) DeleteTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID格式"})
		return
	}

	ctx := c.Request.Context()
	svcErr := h.templateSvc.DeleteTemplate(ctx, id)
	if svcErr != nil {
		if svcErr == application.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "模板删除成功"})
}

// ListTemplates 模板列表（GET /api/v1/issue-templates）
func (h *TemplateHandler) ListTemplates(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	tmplType := c.Query("type")

	ctx := c.Request.Context()
	templates, total, err := h.templateSvc.ListTemplates(ctx, page, size, tmplType)
	if err != nil {
		h.logger.Error("查询模板列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询模板列表失败"})
		return
	}

	list := make([]dto.TemplateInfo, len(templates))
	for i, t := range templates {
		list[i] = dto.TemplateInfo{
			ID:        t.ID.String(),
			Name:      t.Name,
			Type:      t.Type,
			Fields:    t.Fields,
			IsSystem:  t.IsSystem,
			CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: t.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		if t.WorkflowID != nil {
			list[i].WorkflowID = t.WorkflowID.String()
		}
	}

	c.JSON(http.StatusOK, dto.TemplateListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}
