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

// ProductHandler 产品管理Handler
type ProductHandler struct {
	productSvc *application.ProductService
	logger     *zap.Logger
}

// NewProductHandler 创建产品管理Handler实例
func NewProductHandler(productSvc *application.ProductService, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{
		productSvc: productSvc,
		logger:     logger,
	}
}

// CreateProduct 创建产品（POST /api/v1/products）
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	ownerID, err := uuid.Parse(req.OwnerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的负责人ID格式"})
		return
	}

	product := &entity.Product{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		OwnerID:     ownerID,
		Type:        req.Type,
		Platform:    req.Platform,
		Status:      "active",
	}

	if req.Type == "" {
		product.Type = "normal"
	}

	// 解析可选关联ID
	if pid, ok := parseOptionalUUID(req.ProgramID); ok {
		product.ProgramID = pid
	}
	if plid, ok := parseOptionalUUID(req.ProductLineID); ok {
		product.ProductLineID = plid
	}

	result, svcErr := h.productSvc.CreateProduct(ctx, product)
	if svcErr != nil {
		switch svcErr {
		case application.ErrProductCodeExists:
			c.JSON(http.StatusConflict, gin.H{"error": "产品编码已存在"})
		case application.ErrProductLineNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的产品线不存在"})
		case application.ErrProgramNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"error": "指定的项目集不存在"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "产品创建成功",
		"product_id": result.ID.String(),
	})
}

// GetProduct 获取产品详情（GET /api/v1/products/:id）
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID格式"})
		return
	}

	ctx := c.Request.Context()
	product, svcErr := h.productSvc.GetProductDetail(ctx, id)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	resp := h.buildProductDetailResponse(ctx, product)
	c.JSON(http.StatusOK, resp)
}

// UpdateProduct 更新产品（PUT /api/v1/products/:id）
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID格式"})
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	product, svcErr := h.productSvc.GetProductDetail(ctx, id)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	// 更新字段
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Code != nil {
		product.Code = *req.Code
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.OwnerID != nil {
		if oid, oErr := uuid.Parse(*req.OwnerID); oErr == nil {
			product.OwnerID = oid
		}
	}
	if req.Status != nil {
		product.Status = *req.Status
	}
	if req.Type != nil {
		product.Type = *req.Type
	}
	if req.Platform != nil {
		product.Platform = *req.Platform
	}
	if pid, ok := parseOptionalUUID(req.ProgramID); ok || req.ProgramID != nil {
		product.ProgramID = pid
	}
	if plid, ok := parseOptionalUUID(req.ProductLineID); ok || req.ProductLineID != nil {
		product.ProductLineID = plid
	}

	if svcErr := h.productSvc.UpdateProduct(ctx, product); svcErr != nil {
		if svcErr == application.ErrProductCodeExists {
			c.JSON(http.StatusConflict, gin.H{"error": "产品编码已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "产品更新成功"})
}

// DeleteProduct 删除产品（DELETE /api/v1/products/:id）
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID格式"})
		return
	}

	ctx := c.Request.Context()
	if svcErr := h.productSvc.DeleteProduct(ctx, id); svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "产品删除成功"})
}

// ListProducts 分页查询产品列表（GET /api/v1/products）
func (h *ProductHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")
	productLineIDStr := c.Query("product_line_id")

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	var productLineID *uuid.UUID
	if productLineIDStr != "" {
		if plid, err := uuid.Parse(productLineIDStr); err == nil {
			productLineID = &plid
		}
	}

	ctx := c.Request.Context()
	products, total, svcErr := h.productSvc.ListProducts(ctx, page, size, keyword, status, productLineID)
	if svcErr != nil {
		h.logger.Error("查询产品列表失败", zap.Error(svcErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询产品列表失败"})
		return
	}

	list := make([]dto.ProductInfo, len(products))
	for i, p := range products {
		list[i] = buildProductInfo(p)
	}

	c.JSON(http.StatusOK, dto.ProductListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// ==================== 辅助方法 ====================

// buildProductDetailResponse 构建产品详情响应
func (h *ProductHandler) buildProductDetailResponse(_ interface{}, product *entity.Product) dto.ProductDetailResponse {
	info := buildProductInfo(product)
	return dto.ProductDetailResponse{
		ProductInfo: info,
		CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// buildProductInfo 构建产品基本信息
func buildProductInfo(p *entity.Product) dto.ProductInfo {
	info := dto.ProductInfo{
		ID:          p.ID.String(),
		Name:        p.Name,
		Code:        p.Code,
		Description: p.Description,
		OwnerID:     p.OwnerID.String(),
		Status:      p.Status,
		Type:        p.Type,
		Platform:    p.Platform,
	}
	if p.ProgramID != nil {
		ps := p.ProgramID.String()
		info.ProgramID = &ps
	}
	if p.ProductLineID != nil {
		pls := p.ProductLineID.String()
		info.ProductLineID = &pls
	}
	return info
}
