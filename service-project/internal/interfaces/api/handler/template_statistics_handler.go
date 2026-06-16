package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-project/internal/application"
	"leap-one/service-project/internal/domain/entity"
	"leap-one/service-project/internal/interfaces/api/dto"
)

// ==================== 模板Handler ====================

// TemplateHandler 项目模板Handler
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

// ListTemplates 模板列表（GET /api/v1/templates�?func (h *TemplateHandler) ListTemplates(c *gin.Context) {
	page, _ := strconvDefaultInt(c.DefaultQuery("page", "1"), 1)
	size, _ := strconvDefaultInt(c.DefaultQuery("size", "20"), 20)
	templateType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	ctx := c.Request.Context()
	templates, total, err := h.templateSvc.ListTemplates(ctx, page, size, templateType)
	if err != nil {
		h.logger.Error("查询模板列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询模板列表失败"})
		return
	}

	list := make([]dto.TemplateInfo, len(templates))
	for i, t := range templates {
		list[i] = buildTemplateInfo(t)
	}

	c.JSON(http.StatusOK, dto.TemplateListResponse{
		List:  list,
		Total: total,
	})
}

// CreateTemplate 创建模板（POST /api/v1/templates�?func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	var req dto.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	input := &application.CreateTemplateInput{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Config:      req.Config,
	}

	template, err := h.templateSvc.CreateTemplate(ctx, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "模板创建成功",
		"template_id": template.ID.String(),
	})
}

// GetTemplate 获取模板详情（GET /api/v1/templates/:id�?func (h *TemplateHandler) GetTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID格式"})
		return
	}

	ctx := c.Request.Context()
	template, err := h.templateSvc.GetTemplateByID(ctx, id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, buildTemplateInfo(template))
}

// UpdateTemplate 更新模板（PUT /api/v1/templates/:id�?func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
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
	input := &application.UpdateTemplateInput{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Config:      req.Config,
	}

	template, err := h.templateSvc.UpdateTemplate(ctx, id, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "模板更新成功",
		"data":    buildTemplateInfo(template),
	})
}

// DeleteTemplate 删除模板（DELETE /api/v1/templates/:id�?func (h *TemplateHandler) DeleteTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模板ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.templateSvc.DeleteTemplate(ctx, id); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "模板删除成功"})
}

// buildTemplateInfo 构建模板信息
func buildTemplateInfo(t *entity.ProjectTemplate) dto.TemplateInfo {
	return dto.TemplateInfo{
		ID:          t.ID.String(),
		Name:        t.Name,
		Description: t.Description,
		Type:        t.Type,
		Config:      t.Config,
		IsSystem:    t.IsSystem,
		CreatedAt:   t.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   t.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ==================== 统计Handler ====================

// StatisticsHandler 统计分析Handler
type StatisticsHandler struct {
	statsSvc *application.ProjectStatisticsService
	logger   *zap.Logger
}

// NewStatisticsHandler 创建统计分析Handler实例
func NewStatisticsHandler(statsSvc *application.ProjectStatisticsService, logger *zap.Logger) *StatisticsHandler {
	return &StatisticsHandler{
		statsSvc: statsSvc,
		logger:   logger,
	}
}

// GetProjectStatistics 项目统计数据（GET /api/v1/projects/:id/statistics�?func (h *StatisticsHandler) GetProjectStatistics(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}

	ctx := c.Request.Context()
	stats, err := h.statsSvc.GetProjectStatistics(ctx, projectID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetProjectBurndown 燃尽图数据（GET /api/v1/projects/:id/burndown�?func (h *StatisticsHandler) GetProjectBurndown(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}

	ctx := c.Request.Context()
	data, err := h.statsSvc.GetBurndownData(ctx, projectID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetProjectGantt 甘特图数据（GET /api/v1/projects/:id/gantt�?func (h *StatisticsHandler) GetProjectGantt(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}

	ctx := c.Request.Context()
	data, err := h.statsSvc.GetGanttData(ctx, projectID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, data)
}
