package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-project/internal/application"
	"leap-one/service-project/internal/domain/entity"
	"leap-one/service-project/internal/interfaces/api/dto"
)

// ==================== 迭代Handler ====================

// IterationHandler 迭代/Sprint管理Handler
type IterationHandler struct {
	iterSvc *application.IterationService
	logger  *zap.Logger
}

// NewIterationHandler 创建迭代管理Handler实例
func NewIterationHandler(iterSvc *application.IterationService, logger *zap.Logger) *IterationHandler {
	return &IterationHandler{
		iterSvc: iterSvc,
		logger:  logger,
	}
}

// CreateIteration 创建迭代（POST /api/v1/iterations�?func (h *IterationHandler) CreateIteration(c *gin.Context) {
	var req dto.CreateIterationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	input := &application.CreateIterationInput{
		ProjectID:   req.ProjectID,
		Name:        req.Name,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Capacity:    req.Capacity,
		Goal:        req.Goal,
		SortOrder:   req.SortOrder,
	}

	iteration, err := h.iterSvc.CreateIteration(ctx, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "迭代创建成功",
		"iteration_id": iteration.ID.String(),
	})
}

// GetIteration 获取迭代详情（GET /api/v1/iterations/:id�?func (h *IterationHandler) GetIteration(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的迭代ID格式"})
		return
	}

	ctx := c.Request.Context()
	iteration, err := h.iterSvc.GetIterationDetail(ctx, id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	resp := buildIterationDetail(iteration)
	c.JSON(http.StatusOK, resp)
}

// ListIterations 迭代列表（GET /api/v1/iterations�?func (h *IterationHandler) ListIterations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	status := c.Query("status")
	projectIDStr := c.Query("project_id")

	var projectID uuid.UUID
	if projectIDStr != "" {
		if pid, parseErr := uuid.Parse(projectIDStr); parseErr == nil {
			projectID = pid
		}
	}

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	ctx := c.Request.Context()
	iterations, total, err := h.iterSvc.ListIterations(ctx, page, size, projectID, status)
	if err != nil {
		h.logger.Error("查询迭代列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询迭代列表失败"})
		return
	}

	list := make([]dto.IterationInfo, len(iterations))
	for i, iter := range iterations {
		list[i] = buildIterationInfo(iter)
	}

	c.JSON(http.StatusOK, dto.IterationListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// UpdateIteration 更新迭代（PUT /api/v1/iterations/:id�?func (h *IterationHandler) UpdateIteration(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的迭代ID格式"})
		return
	}

	var req dto.UpdateIterationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	input := &application.UpdateIterationInput{
		Name:        req.Name,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Capacity:    req.Capacity,
		Goal:        req.Goal,
		SortOrder:   req.SortOrder,
	}

	iteration, err := h.iterSvc.UpdateIteration(ctx, id, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "迭代更新成功",
		"data":    buildIterationInfo(iteration),
	})
}

// DeleteIteration 删除迭代（DELETE /api/v1/iterations/:id�?func (h *IterationHandler) DeleteIteration(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的迭代ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.iterSvc.DeleteIteration(ctx, id); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "迭代删除成功"})
}

// StartIteration 开始迭代（POST /api/v1/iterations/:id/start�?func (h *IterationHandler) StartIteration(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的迭代ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.iterSvc.StartIteration(ctx, id); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "迭代已启�?})
}

// CompleteIteration 完成迭代（POST /api/v1/iterations/:id/complete�?func (h *IterationHandler) CompleteIteration(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的迭代ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.iterSvc.CompleteIteration(ctx, id); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "迭代已完�?})
}

// GetIterationBoard 迭代看板数据（GET /api/v1/iterations/:id/board�?func (h *IterationHandler) GetIterationBoard(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的迭代ID格式"})
		return
	}

	ctx := c.Request.Context()
	iteration, err := h.iterSvc.GetIterationDetail(ctx, id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// TODO: 集成任务服务后，从任务服务获取看板数�?	boardData := map[string]interface{}{
		"iteration_id": id.String(),
		"name":         iteration.Name,
		"status":       iteration.Status,
		"columns": []map[string]interface{}{
			{"key": "todo", "name": "待办", "tasks": []interface{}{}},
			{"key": "in_progress", "name": "进行�?, "tasks": []interface{}{}},
			{"key": "done", "name": "已完�?, "tasks": []interface{}{}},
		},
	}

	c.JSON(http.StatusOK, boardData)
}

// GetIterationBurndown 迭代燃尽图数据（GET /api/v1/iterations/:id/burndown�?func (h *IterationHandler) GetIterationBurndown(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的迭代ID格式"})
		return
	}

	ctx := c.Request.Context()
	iteration, err := h.iterSvc.GetIterationDetail(ctx, id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// 生成理想燃尽线数�?	burndown := generateBurndownData(iteration)

	c.JSON(http.StatusOK, burndown)
}

// ==================== 辅助方法 ====================

// buildIterationInfo 构建迭代简要信�?func buildIterationInfo(i *entity.Iteration) dto.IterationInfo {
	info := dto.IterationInfo{
		ID:          i.ID.String(),
		ProjectID:   i.ProjectID.String(),
		Name:        i.Name,
		Description: i.Description,
		Status:      i.Status,
		StartDate:   i.StartDate.Format("2006-01-02"),
		EndDate:     i.EndDate.Format("2006-01-02"),
		Goal:        i.Goal,
		SortOrder:   i.SortOrder,
		CreatedAt:   i.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   i.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if i.Capacity != nil {
		info.Capacity = *i.Capacity
	}
	return info
}

// buildIterationDetail 构建迭代详情响应
func buildIterationDetail(i *entity.Iteration) dto.IterationDetailResponse {
	return dto.IterationDetailResponse{
		IterationInfo: buildIterationInfo(i),
		BoardData:     nil, // TODO: 从任务服务获�?		Burndown:      dto.BurndownData{}, // TODO: 从统计服务获�?	}
}

// generateBurndownData 生成燃尽图数�?func generateBurndownData(iter *entity.Iteration) dto.BurndownData {
	totalDays := int(iter.EndDate.Sub(iter.StartDate).Hours()/24) + 1
	if totalDays <= 0 {
		totalDays = 1
	}

	points := make([]dto.BurndownPoint, totalDays)
	for d := 0; d < totalDays; d++ {
		date := iter.StartDate.AddDate(0, 0, d)
		ratio := float64(totalDays-1-d) / float64(totalDays-1)
		if ratio < 0 {
			ratio = 0
		}
		points[d] = dto.BurndownPoint{
			Date:  date.Format("2006-01-02"),
			Value: ratio,
		}
	}

	return dto.BurndownData{
		SprintName:  iter.Name,
		TotalPoints: 0, // TODO: 从任务服务获取实际总故事点
		IdealLine:   points,
		ActualLine:  []dto.BurndownPoint{}, // TODO: 从任务服务获取实际数�?		Remaining:   0,
	}
}
