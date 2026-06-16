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

// ProductLineHandler 产品线管理Handler
type ProductLineHandler struct {
	productSvc *application.ProductService
	logger     *zap.Logger
}

// NewProductLineHandler 创建产品线管理Handler实例
func NewProductLineHandler(productSvc *application.ProductService, logger *zap.Logger) *ProductLineHandler {
	return &ProductLineHandler{
		productSvc: productSvc,
		logger:     logger,
	}
}

// CreateProductLine 创建产品线（POST /api/v1/product-lines）
func (h *ProductLineHandler) CreateProductLine(c *gin.Context) {
	var req dto.CreateProductLineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	line := &entity.ProductLine{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      "active",
	}

	result, svcErr := h.productSvc.CreateProductLine(ctx, line)
	if svcErr != nil {
		h.logger.Error("创建产品线失败", zap.Error(svcErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建产品线失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":         "产品线创建成功",
		"product_line_id": result.ID.String(),
	})
}

// ListProductLines 获取全部产品线列表（GET /api/v1/product-lines）
func (h *ProductLineHandler) ListProductLines(c *gin.Context) {
	ctx := c.Request.Context()
	lines, svcErr := h.productSvc.ListProductLines(ctx)
	if svcErr != nil {
		h.logger.Error("查询产品线列表失败", zap.Error(svcErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询产品线列表失败"})
		return
	}

	list := make([]dto.ProductLineInfo, len(lines))
	for i, l := range lines {
		list[i] = dto.ProductLineInfo{
			ID:          l.ID.String(),
			Name:        l.Name,
			Description: l.Description,
			SortOrder:   l.SortOrder,
			Status:      l.Status,
			CreatedAt:   l.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   l.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"list": list, "count": len(list)})
}

// GetProductLine 获取产品线详情（GET /api/v1/product-lines/:id）
func (h *ProductLineHandler) GetProductLine(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品线ID格式"})
		return
	}

	ctx := c.Request.Context()
	line, svcErr := h.productSvc.GetProductLineDetail(ctx, id)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ProductLineInfo{
		ID:          line.ID.String(),
		Name:        line.Name,
		Description: line.Description,
		SortOrder:   line.SortOrder,
		Status:      line.Status,
		CreatedAt:   line.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   line.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

// UpdateProductLine 更新产品线（PUT /api/v1/product-lines/:id）
func (h *ProductLineHandler) UpdateProductLine(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品线ID格式"})
		return
	}

	var req dto.UpdateProductLineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	line, svcErr := h.productSvc.GetProductLineDetail(ctx, id)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	if req.Name != nil {
		line.Name = *req.Name
	}
	if req.Description != nil {
		line.Description = *req.Description
	}
	if req.SortOrder != nil {
		line.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		line.Status = *req.Status
	}

	if svcErr := h.productSvc.UpdateProductLine(ctx, line); svcErr != nil {
		h.logger.Error("更新产品线失败", zap.Error(svcErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新产品线失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "产品线更新成功"})
}

// DeleteProductLine 删除产品线（DELETE /api/v1/product-lines/:id）
func (h *ProductLineHandler) DeleteProductLine(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品线ID格式"})
		return
	}

	ctx := c.Request.Context()
	if svcErr := h.productSvc.DeleteProductLine(ctx, id); svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "产品线删除成功"})
}
