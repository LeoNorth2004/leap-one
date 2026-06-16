package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-portfolio/internal/application"
	"leap-one/service-portfolio/internal/domain/entity"
	"leap-one/service-portfolio/internal/interfaces/api/dto"
)

// VersionHandler 产品版本管理Handler
type VersionHandler struct {
	productSvc *application.ProductService
	logger     *zap.Logger
}

// NewVersionHandler 创建版本管理Handler实例
func NewVersionHandler(productSvc *application.ProductService, logger *zap.Logger) *VersionHandler {
	return &VersionHandler{
		productSvc: productSvc,
		logger:     logger,
	}
}

// CreateVersion 创建版本（POST /api/v1/products/:id/versions）
func (h *VersionHandler) CreateVersion(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID格式"})
		return
	}

	var req dto.CreateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	releaseDate, _ := parseDate(req.ReleaseDate)

	version := &entity.ProductVersion{
		ProductID:   productID,
		Name:        req.Name,
		ReleaseDate: releaseDate,
		Description: req.Description,
		Plan:        req.Plan,
		Status:      "planning",
	}
	if req.Status != "" {
		version.Status = req.Status
	}

	if svcErr := h.productSvc.CreateVersion(ctx, version); svcErr != nil {
		switch svcErr {
		case application.ErrProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "产品不存在"})
		default:
			h.logger.Error("创建版本失败", zap.Error(svcErr))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建版本失败"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "版本创建成功",
		"version_id": version.ID.String(),
	})
}

// ListVersions 获取产品的版本列表（GET /api/v1/products/:id/versions）
func (h *VersionHandler) ListVersions(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID格式"})
		return
	}

	ctx := c.Request.Context()
	versions, svcErr := h.productSvc.ListVersions(ctx, productID)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	list := make([]dto.VersionInfo, len(versions))
	for i, v := range versions {
		list[i] = buildVersionInfo(v)
	}

	c.JSON(http.StatusOK, gin.H{"list": list, "count": len(list)})
}

// UpdateVersion 更新版本信息（PUT /api/v1/products/:pid/versions/:vid）
func (h *VersionHandler) UpdateVersion(c *gin.Context) {
	versionID, err := uuid.Parse(c.Param("vid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的版本ID格式"})
		return
	}

	var req dto.UpdateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	version, svcErr := h.productSvc.GetVersionDetail(ctx, versionID)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	if req.Name != nil {
		version.Name = *req.Name
	}
	if req.ReleaseDate != nil {
		version.ReleaseDate, _ = parseDate(*req.ReleaseDate)
	}
	if req.Status != nil {
		version.Status = *req.Status
	}
	if req.Description != nil {
		version.Description = *req.Description
	}
	if req.Plan != nil {
		version.Plan = *req.Plan
	}

	if svcErr := h.productSvc.UpdateVersion(ctx, version); svcErr != nil {
		h.logger.Error("更新版本失败", zap.Error(svcErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新版本失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "版本更新成功"})
}

// ReleaseVersion 发布版本（POST /api/v1/products/:pid/versions/:vid/release）
func (h *VersionHandler) ReleaseVersion(c *gin.Context) {
	versionID, err := uuid.Parse(c.Param("vid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的版本ID格式"})
		return
	}

	var req dto.ReleaseVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	releaseDate, dateErr := parseDate(req.ReleaseDate)
	if dateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的发布日期格式，请使用 YYYY-MM-DD"})
		return
	}

	now := time.Now()
	if releaseDate == nil {
		releaseDate = &now
	}

	if svcErr := h.productSvc.ReleaseVersion(ctx, versionID, releaseDate); svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "版本已发布"})
}

// ==================== 辅助方法 ====================

// buildVersionInfo 构建版本信息
func buildVersionInfo(v *entity.ProductVersion) dto.VersionInfo {
	info := dto.VersionInfo{
		ID:          v.ID.String(),
		ProductID:   v.ProductID.String(),
		Name:        v.Name,
		Status:      v.Status,
		Description: v.Description,
		Plan:        v.Plan,
		CreatedAt:   v.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   v.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if v.ReleaseDate != nil {
		rs := v.ReleaseDate.Format("2006-01-02")
		info.ReleaseDate = &rs
	}
	return info
}
