package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-portfolio/internal/application"
	"leap-one/service-portfolio/internal/domain/entity"
	"leap-one/service-portfolio/internal/interfaces/api/dto"
)

// RoadmapHandler 产品路线图管理Handler
type RoadmapHandler struct {
	productSvc *application.ProductService
	logger     *zap.Logger
}

// NewRoadmapHandler 创建路线图管理Handler实例
func NewRoadmapHandler(productSvc *application.ProductService, logger *zap.Logger) *RoadmapHandler {
	return &RoadmapHandler{
		productSvc: productSvc,
		logger:     logger,
	}
}

// CreateRoadmapItem 添加路线图项（POST /api/v1/products/:id/roadmap）
func (h *RoadmapHandler) CreateRoadmapItem(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID格式"})
		return
	}

	var req dto.CreateRoadmapItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	item := &entity.ProductRoadmapItem{
		ProductID:   productID,
		Title:       req.Title,
		Description: req.Description,
		Quarter:     req.Quarter,
		Year:        req.Year,
		Priority:    req.Priority,
		SortOrder:   req.SortOrder,
		Status:      "planning",
	}
	if req.Status != "" {
		item.Status = req.Status
	}

	if svcErr := h.productSvc.CreateRoadmapItem(ctx, item); svcErr != nil {
		switch svcErr {
		case application.ErrProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "产品不存在"})
		default:
			h.logger.Error("创建路线图项失败", zap.Error(svcErr))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建路线图项失败"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":         "路线图项创建成功",
		"roadmap_item_id": item.ID.String(),
	})
}

// ListRoadmapItems 获取产品的路线图列表（GET /api/v1/products/:id/roadmap）
func (h *RoadmapHandler) ListRoadmapItems(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID格式"})
		return
	}

	ctx := c.Request.Context()
	items, svcErr := h.productSvc.ListRoadmapItems(ctx, productID)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	list := make([]dto.RoadmapItemInfo, len(items))
	for i, ri := range items {
		list[i] = buildRoadmapItemInfo(ri)
	}

	c.JSON(http.StatusOK, gin.H{"list": list, "count": len(list)})
}

// UpdateRoadmapItem 更新路线图项（PUT /api/v1/products/:id/roadmap/:rid）
func (h *RoadmapHandler) UpdateRoadmapItem(c *gin.Context) {
	itemID, err := uuid.Parse(c.Param("rid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的路线图项ID格式"})
		return
	}

	var req dto.UpdateRoadmapItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	item, svcErr := h.productSvc.GetRoadmapItemDetail(ctx, itemID)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	if req.Title != nil {
		item.Title = *req.Title
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if req.Quarter != nil {
		item.Quarter = *req.Quarter
	}
	if req.Year != nil {
		item.Year = *req.Year
	}
	if req.Status != nil {
		item.Status = *req.Status
	}
	if req.Priority != nil {
		item.Priority = *req.Priority
	}
	if req.SortOrder != nil {
		item.SortOrder = *req.SortOrder
	}

	if svcErr := h.productSvc.UpdateRoadmapItem(ctx, item); svcErr != nil {
		h.logger.Error("更新路线图项失败", zap.Error(svcErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新路线图项失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "路线图项更新成功"})
}

// DeleteRoadmapItem 删除路线图项（DELETE /api/v1/products/:id/roadmap/:rid）
func (h *RoadmapHandler) DeleteRoadmapItem(c *gin.Context) {
	itemID, err := uuid.Parse(c.Param("rid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的路线图项ID格式"})
		return
	}

	ctx := c.Request.Context()
	if svcErr := h.productSvc.DeleteRoadmapItem(ctx, itemID); svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "路线图项删除成功"})
}

// ReorderRoadmapItems 重新排序路线图项（PUT /api/v1/products/:id/roadmap/reorder）
func (h *RoadmapHandler) ReorderRoadmapItems(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID格式"})
		return
	}

	var req dto.ReorderRoadmapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 解析UUID列表
	itemIDs := make([]uuid.UUID, len(req.ItemIDs))
	for i, idStr := range req.ItemIDs {
		id, parseErr := uuid.Parse(idStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的路线图项ID: " + idStr})
			return
		}
		itemIDs[i] = id
	}

	if svcErr := h.productSvc.ReorderRoadmapItems(ctx, productID, itemIDs); svcErr != nil {
		switch svcErr {
		case application.ErrProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "产品不存在"})
		default:
			h.logger.Error("排序路线图项失败", zap.Error(svcErr))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "排序路线图项失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "路线图排序更新成功"})
}

// ==================== 辅助方法 ====================

// buildRoadmapItemInfo 构建路线图项信息
func buildRoadmapItemInfo(ri *entity.ProductRoadmapItem) dto.RoadmapItemInfo {
	return dto.RoadmapItemInfo{
		ID:          ri.ID.String(),
		ProductID:   ri.ProductID.String(),
		Title:       ri.Title,
		Description: ri.Description,
		Quarter:     ri.Quarter,
		Year:        ri.Year,
		Status:      ri.Status,
		Priority:    ri.Priority,
		SortOrder:   ri.SortOrder,
		CreatedAt:   ri.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   ri.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
