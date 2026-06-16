package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-portfolio/internal/application"
	"leap-one/service-portfolio/internal/domain/entity"
	"leap-one/service-portfolio/internal/interfaces/api/dto"
)

// PlanHandler 产品计划管理Handler
type PlanHandler struct {
	productSvc *application.ProductService
	logger     *zap.Logger
}

// NewPlanHandler 创建计划管理Handler实例
func NewPlanHandler(productSvc *application.ProductService, logger *zap.Logger) *PlanHandler {
	return &PlanHandler{
		productSvc: productSvc,
		logger:     logger,
	}
}

// CreatePlan 创建计划（POST /api/v1/plans）
func (h *PlanHandler) CreatePlan(c *gin.Context) {
	var req dto.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID格式"})
		return
	}

	ctx := c.Request.Context()

	startDate, _ := parseDate(req.StartDate)
	endDate, _ := parseDate(req.EndDate)

	plan := &entity.ProductPlan{
		ProductID: productID,
		Name:      req.Name,
		Content:   req.Content,
		Status:    "active",
		StartDate: startDate,
		EndDate:   endDate,
	}
	if req.Status != "" {
		plan.Status = req.Status
	}

	if svcErr := h.productSvc.CreatePlan(ctx, plan); svcErr != nil {
		switch svcErr {
		case application.ErrProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "产品不存在"})
		default:
			h.logger.Error("创建计划失败", zap.Error(svcErr))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建计划失败"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "计划创建成功",
		"plan_id": plan.ID.String(),
	})
}

// ListPlans 分页查询计划列表（GET /api/v1/plans）
func (h *PlanHandler) ListPlans(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	ctx := c.Request.Context()
	plans, total, svcErr := h.productSvc.ListAllPlans(ctx, page, size)
	if svcErr != nil {
		h.logger.Error("查询计划列表失败", zap.Error(svcErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询计划列表失败"})
		return
	}

	list := make([]dto.PlanInfo, len(plans))
	for i, p := range plans {
		list[i] = buildPlanInfo(p)
	}

	c.JSON(http.StatusOK, dto.PlanListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// ListProductPlans 获取某产品的计划列表（GET /api/v1/products/:id/plans）
func (h *PlanHandler) ListProductPlans(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID格式"})
		return
	}

	ctx := c.Request.Context()
	plans, svcErr := h.productSvc.ListPlans(ctx, productID)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	list := make([]dto.PlanInfo, len(plans))
	for i, p := range plans {
		list[i] = buildPlanInfo(p)
	}

	c.JSON(http.StatusOK, gin.H{"list": list, "count": len(list)})
}

// GetPlan 获取计划详情（GET /api/v1/plans/:id）
func (h *PlanHandler) GetPlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的计划ID格式"})
		return
	}

	ctx := c.Request.Context()
	plan, svcErr := h.productSvc.GetPlanDetail(ctx, id)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, buildPlanInfo(plan))
}

// UpdatePlan 更新计划（PUT /api/v1/plans/:id）
func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的计划ID格式"})
		return
	}

	var req dto.UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	plan, svcErr := h.productSvc.GetPlanDetail(ctx, id)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	if req.Name != nil {
		plan.Name = *req.Name
	}
	if req.Content != nil {
		plan.Content = *req.Content
	}
	if req.Status != nil {
		plan.Status = *req.Status
	}
	if req.StartDate != nil {
		plan.StartDate, _ = parseDate(*req.StartDate)
	}
	if req.EndDate != nil {
		plan.EndDate, _ = parseDate(*req.EndDate)
	}

	if svcErr := h.productSvc.UpdatePlan(ctx, plan); svcErr != nil {
		h.logger.Error("更新计划失败", zap.Error(svcErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新计划失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "计划更新成功"})
}

// DeletePlan 删除计划（DELETE /api/v1/plans/:id）
func (h *PlanHandler) DeletePlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的计划ID格式"})
		return
	}

	ctx := c.Request.Context()
	if svcErr := h.productSvc.DeletePlan(ctx, id); svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "计划删除成功"})
}

// ==================== 辅助方法 ====================

// buildPlanInfo 构建计划信息
func buildPlanInfo(p *entity.ProductPlan) dto.PlanInfo {
	info := dto.PlanInfo{
		ID:        p.ID.String(),
		ProductID: p.ProductID.String(),
		Name:      p.Name,
		Content:   p.Content,
		Status:    p.Status,
		CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if p.StartDate != nil {
		ss := p.StartDate.Format("2006-01-02")
		info.StartDate = &ss
	}
	if p.EndDate != nil {
		es := p.EndDate.Format("2006-01-02")
		info.EndDate = &es
	}
	return info
}
