package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-task/internal/application"
	"leap-one/service-task/internal/domain/entity"
	"leap-one/service-task/internal/interfaces/api/dto"
)

// WorkflowHandler 工作流管理Handler
type WorkflowHandler struct {
	workflowSvc *application.WorkflowService
	slaSvc      *application.SLAConfigService
	logger      *zap.Logger
}

// NewWorkflowHandler 创建工作流管理Handler实例
func NewWorkflowHandler(workflowSvc *application.WorkflowService, slaSvc *application.SLAConfigService, logger *zap.Logger) *WorkflowHandler {
	return &WorkflowHandler{
		workflowSvc: workflowSvc,
		slaSvc:      slaSvc,
		logger:      logger,
	}
}

// ==================== 工作流管理 ====================

// CreateWorkflow 创建工作流（POST /api/v1/workflows）
func (h *WorkflowHandler) CreateWorkflow(c *gin.Context) {
	var req dto.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	wf, err := h.workflowSvc.CreateWorkflow(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "工作流创建成功",
		"workflow_id": wf.ID.String(),
	})
}

// GetWorkflow 获取工作流详情（GET /api/v1/workflows/:id）
func (h *WorkflowHandler) GetWorkflow(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工作流ID格式"})
		return
	}

	ctx := c.Request.Context()
	wf, svcErr := h.workflowSvc.GetWorkflow(ctx, id)
	if svcErr != nil || wf == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "工作流不存在"})
		return
	}

	resp := buildWorkflowInfo(wf)
	c.JSON(http.StatusOK, resp)
}

// UpdateWorkflow 更新工作流（PUT /api/v1/workflows/:id）
func (h *WorkflowHandler) UpdateWorkflow(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工作流ID格式"})
		return
	}

	var req dto.UpdateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	svcErr := h.workflowSvc.UpdateWorkflow(ctx, id, &req)
	if svcErr != nil {
		if svcErr == application.ErrWorkflowNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "工作流不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "工作流更新成功"})
}

// DeleteWorkflow 删除工作流（DELETE /api/v1/workflows/:id）
func (h *WorkflowHandler) DeleteWorkflow(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工作流ID格式"})
		return
	}

	ctx := c.Request.Context()
	svcErr := h.workflowSvc.DeleteWorkflow(ctx, id)
	if svcErr != nil {
		if svcErr == application.ErrWorkflowNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "工作流不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "工作流删除成功"})
}

// ListWorkflows 工作流列表（GET /api/v1/workflows）
func (h *WorkflowHandler) ListWorkflows(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	wfType := c.Query("type")

	ctx := c.Request.Context()
	workflows, total, err := h.workflowSvc.ListWorkflows(ctx, page, size, wfType)
	if err != nil {
		h.logger.Error("查询工作流列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询工作流列表失败"})
		return
	}

	list := make([]dto.WorkflowInfo, len(workflows))
	for i, w := range workflows {
		list[i] = buildWorkflowInfo(w)
	}

	c.JSON(http.StatusOK, dto.WorkflowListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// AddTransition 添加状态转换（POST /api/v1/workflows/:id/transitions）
func (h *WorkflowHandler) AddTransition(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工作流ID格式"})
		return
	}

	var req dto.CreateTransitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	transition, svcErr := h.workflowSvc.AddTransition(ctx, id, &req)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "状态转换规则添加成功",
		"transition_id": transition.ID.String(),
	})
}

// ==================== SLA配置管理 ====================

// CreateSLAConfig 创建SLA配置（POST /api/v1/sla-configs）
func (h *WorkflowHandler) CreateSLAConfig(c *gin.Context) {
	var req dto.CreateSLAConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	cfg, err := h.slaSvc.CreateSLAConfig(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "SLA配置创建成功",
		"config_id": cfg.ID.String(),
	})
}

// GetSLAConfig 获取SLA配置详情（GET /api/v1/sla-configs/:id）
func (h *WorkflowHandler) GetSLAConfig(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的配置ID格式"})
		return
	}

	ctx := c.Request.Context()
	cfg, svcErr := h.slaSvc.GetSLAConfig(ctx, id)
	if svcErr != nil || cfg == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SLA配置不存在"})
		return
	}

	c.JSON(http.StatusOK, dto.SLAConfigInfo{
		ID:                cfg.ID.String(),
		Type:              cfg.Type,
		Priority:          cfg.Priority,
		ResponseSLA:       cfg.ResponseSLA,
		ResolveSLA:        cfg.ResolveSLA,
		BusinessHoursOnly: cfg.BusinessHoursOnly,
		CreatedAt:         cfg.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         cfg.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

// UpdateSLAConfig 更新SLA配置（PUT /api/v1/sla-configs/:id）
func (h *WorkflowHandler) UpdateSLAConfig(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的配置ID格式"})
		return
	}

	var req dto.UpdateSLAConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	svcErr := h.slaSvc.UpdateSLAConfig(ctx, id, &req)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SLA配置更新成功"})
}

// DeleteSLAConfig 删除SLA配置（DELETE /api/v1/sla-configs/:id）
func (h *WorkflowHandler) DeleteSLAConfig(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的配置ID格式"})
		return
	}

	ctx := c.Request.Context()
	svcErr := h.slaSvc.DeleteSLAConfig(ctx, id)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SLA配置删除成功"})
}

// ListSLAConfigs SLA配置列表（GET /api/v1/sla-configs）
func (h *WorkflowHandler) ListSLAConfigs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	ctx := c.Request.Context()
	configs, total, err := h.slaSvc.ListSLAConfigs(ctx, page, size)
	if err != nil {
		h.logger.Error("查询SLA配置列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询SLA配置列表失败"})
		return
	}

	list := make([]dto.SLAConfigInfo, len(configs))
	for i, cfg := range configs {
		list[i] = dto.SLAConfigInfo{
			ID:                cfg.ID.String(),
			Type:              cfg.Type,
			Priority:          cfg.Priority,
			ResponseSLA:       cfg.ResponseSLA,
			ResolveSLA:        cfg.ResolveSLA,
			BusinessHoursOnly: cfg.BusinessHoursOnly,
			CreatedAt:         cfg.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:         cfg.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, dto.SLAConfigListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// ==================== 辅助方法 ====================

func buildWorkflowInfo(wf *entity.IssueWorkflow) dto.WorkflowInfo {
	info := dto.WorkflowInfo{
		ID:            wf.ID.String(),
		Name:          wf.Name,
		Type:          wf.Type,
		InitialStatus: wf.InitialStatus,
		Description:   wf.Description,
		CreatedAt:     wf.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     wf.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if wf.Transitions != nil {
		info.Transitions = make([]dto.TransitionInfo, len(wf.Transitions))
		for i, t := range wf.Transitions {
			info.Transitions[i] = dto.TransitionInfo{
				ID:         t.ID.String(),
				FromStatus: t.FromStatus,
				ToStatus:   t.ToStatus,
				Condition:  t.Condition,
				Name:       t.Name,
				SortOrder:  t.SortOrder,
			}
		}
	}

	return info
}
