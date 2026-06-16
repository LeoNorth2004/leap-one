package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-bi/internal/domain/entity"
	"leap-one/service-bi/internal/domain/repository"
	"leap-one/service-bi/internal/interfaces/api/dto"
)

// DashboardHandler BI大屏管理Handler
type DashboardHandler struct {
	dashboardRepo repository.DashboardConfigRepository
	logger        *zap.Logger
}

// NewDashboardHandler 创建大屏管理Handler实例
func NewDashboardHandler(dashboardRepo repository.DashboardConfigRepository, logger *zap.Logger) *DashboardHandler {
	return &DashboardHandler{
		dashboardRepo: dashboardRepo,
		logger:        logger,
	}
}

// GetCompanyOverview 获取公司数据盘点大屏 (GET /api/v1/dashboards/company-overview)
func (h *DashboardHandler) GetCompanyOverview(c *gin.Context) {
	ctx := c.Request.Context()
	dash, err := h.dashboardRepo.GetByType(ctx, "company_overview")
	if err != nil || dash == nil {
		c.JSON(http.StatusOK, gin.H{"message": "暂无公司数据盘点配置", "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, h.buildDashboardResponse(dash))
}

// GetAnnualOverview 获取年度新增数据概览 (GET /api/v1/dashboards/annual-overview)
func (h *DashboardHandler) GetAnnualOverview(c *gin.Context) {
	ctx := c.Request.Context()
	dash, err := h.dashboardRepo.GetByType(ctx, "annual_data")
	if err != nil || dash == nil {
		c.JSON(http.StatusOK, gin.H{"message": "暂无年度数据概览配置", "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, h.buildDashboardResponse(dash))
}

// GetRanking 获取年度排行�?(GET /api/v1/dashboards/ranking)
func (h *DashboardHandler) GetRanking(c *gin.Context) {
	ctx := c.Request.Context()
	dash, err := h.dashboardRepo.GetByType(ctx, "ranking")
	if err != nil || dash == nil {
		c.JSON(http.StatusOK, gin.H{"message": "暂无年度排行榜配�?, "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, h.buildDashboardResponse(dash))
}

// GetSprintBurndown 获取迭代燃尽图大�?(GET /api/v1/dashboards/sprint-burndown)
func (h *DashboardHandler) GetSprintBurndown(c *gin.Context) {
	ctx := c.Request.Context()
	dash, err := h.dashboardRepo.GetByType(ctx, "sprint_burndown")
	if err != nil || dash == nil {
		c.JSON(http.StatusOK, gin.H{"message": "暂无迭代燃尽图配�?, "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, h.buildDashboardResponse(dash))
}

// GetAnnualSummary 获取年度总结大屏 (GET /api/v1/dashboards/annual-summary)
func (h *DashboardHandler) GetAnnualSummary(c *gin.Context) {
	ctx := c.Request.Context()
	dash, err := h.dashboardRepo.GetByType(ctx, "annual_summary")
	if err != nil || dash == nil {
		c.JSON(http.StatusOK, gin.H{"message": "暂无年度总结配置", "data": gin.H{}})
		return
	}
	c.JSON(http.StatusOK, h.buildDashboardResponse(dash))
}

// GetDashboardByID 获取自定义大�?(GET /api/v1/dashboards/:id)
func (h *DashboardHandler) GetDashboardByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的大屏ID格式"})
		return
	}

	ctx := c.Request.Context()
	dash, err := h.dashboardRepo.GetByID(ctx, id)
	if err != nil || dash == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "大屏不存�?})
		return
	}
	c.JSON(http.StatusOK, h.buildDashboardResponse(dash))
}

// buildDashboardResponse 构建大屏响应
func (h *DashboardHandler) buildDashboardResponse(dash *entity.DashboardConfig) dto.DashboardDetailResponse {
	return dto.DashboardDetailResponse{
		DashboardInfo: dto.DashboardInfo{
			ID:              dash.ID.String(),
			Name:            dash.Name,
			Type:            dash.Type,
			Layout:          dash.Layout,
			RefreshInterval: dash.RefreshInterval,
			IsSystem:        dash.IsSystem,
		},
		CreatedAt: dash.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: dash.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// getStringFromContext 从Gin上下文中安全获取字符串�?
func getStringFromContext(c *gin.Context, key string) string {
	val, exists := c.Get(key)
	if !exists {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}
