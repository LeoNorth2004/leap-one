package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-bi/internal/domain/entity"
	"leap-one/service-bi/internal/domain/repository"
	"leap-one/service-bi/internal/interfaces/api/dto"
)

// ReportHandler 报表管理Handler
type ReportHandler struct {
	reportRepo repository.ReportTemplateRepository
	snapshotRepo repository.DataSnapshotRepository
	logger     *zap.Logger
}

// NewReportHandler 创建报表管理Handler实例
func NewReportHandler(reportRepo repository.ReportTemplateRepository, snapshotRepo repository.DataSnapshotRepository, logger *zap.Logger) *ReportHandler {
	return &ReportHandler{
		reportRepo: reportRepo,
		snapshotRepo: snapshotRepo,
		logger:     logger,
	}
}

// CreateReport 创建自定义报�?(POST /api/v1/reports)
func (h *ReportHandler) CreateReport(c *gin.Context) {
	var req dto.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	tpl := &entity.ReportTemplate{
		Name:      req.Name,
		Type:      req.Type,
		Config:    req.Config,
		ChartType: req.ChartType,
		CreatorID: req.CreatorID,
	}
	if tpl.ChartType == "" {
		tpl.ChartType = "table"
	}

	ctx := c.Request.Context()
	if err := h.reportRepo.Create(ctx, tpl); err != nil {
		h.logger.Error("创建报表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建报表失败"})
		return
	}

	h.logger.Info("创建报表成功",
		zap.String("report_id", tpl.ID.String()),
		zap.String("report_name", req.Name),
	)

	c.JSON(http.StatusCreated, gin.H{
		"message":   "报表创建成功",
		"report_id": tpl.ID.String(),
	})
}

// ListReports 报表列表 (GET /api/v1/reports)
func (h *ReportHandler) ListReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	creatorIDStr := c.Query("creator_id")
	reportType := c.Query("type")

	if page < 1 { page = 1 }
	if size < 1 || size > 100 { size = 20 }

	var creatorID uuid.UUID
	if creatorIDStr != "" {
		if pid, err := uuid.Parse(creatorIDStr); err == nil {
			creatorID = pid
		}
	}

	ctx := c.Request.Context()
	reports, total, err := h.reportRepo.List(ctx, page, size, creatorID, reportType)
	if err != nil {
		h.logger.Error("查询报表列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询报表列表失败"})
		return
	}

	list := make([]dto.ReportInfo, len(reports))
	for i, r := range reports {
		list[i] = dto.ReportInfo{
			ID:        r.ID.String(),
			Name:      r.Name,
			Type:      r.Type,
			ChartType: r.ChartType,
			CreatorID: r.CreatorID.String(),
			CreatedAt: r.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, dto.ReportListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// GetReport 获取报表数据 (GET /api/v1/reports/:id)
func (h *ReportHandler) GetReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的报表ID格式"})
		return
	}

	ctx := c.Request.Context()
	tpl, err := h.reportRepo.GetByID(ctx, id)
	if err != nil || tpl == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "报表不存�?})
		return
	}

	resp := dto.ReportDetailResponse{
		ReportInfo: dto.ReportInfo{
			ID:        tpl.ID.String(),
			Name:      tpl.Name,
			Type:      tpl.Type,
			ChartType: tpl.ChartType,
			CreatorID: tpl.CreatorID.String(),
			CreatedAt: tpl.CreatedAt.Format("2006-01-02 15:04:05"),
		},
		Config:    tpl.Config,
		UpdatedAt: tpl.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateReport 更新报表 (PUT /api/v1/reports/:id)
func (h *ReportHandler) UpdateReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的报表ID格式"})
		return
	}

	var req dto.UpdateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	tpl, err := h.reportRepo.GetByID(ctx, id)
	if err != nil || tpl == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "报表不存�?})
		return
	}

	if req.Name != nil { tpl.Name = *req.Name }
	if req.Type != nil { tpl.Type = *req.Type }
	if req.Config != nil { tpl.Config = *req.Config }
	if req.ChartType != nil { tpl.ChartType = *req.ChartType }

	if err := h.reportRepo.Update(ctx, tpl); err != nil {
		h.logger.Error("更新报表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新报表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "报表更新成功"})
}

// DeleteReport 删除报表 (DELETE /api/v1/reports/:id)
func (h *ReportHandler) DeleteReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的报表ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.reportRepo.Delete(ctx, id); err != nil {
		h.logger.Error("删除报表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除报表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "报表删除成功"})
}

// GetReportData 获取报表数据 (GET /api/v1/reports/:id/data)
func (h *ReportHandler) GetReportData(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的报表ID格式"})
		return
	}

	ctx := c.Request.Context()
	tpl, err := h.reportRepo.GetByID(ctx, id)
	if err != nil || tpl == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "报表不存�?})
		return
	}

	// 根据报表类型返回模拟统计数据
	data := h.generateMockStatsData(tpl.Type)

	c.JSON(http.StatusOK, gin.H{
		"report_id": id.String(),
		"report_name": tpl.Name,
		"type": tpl.Type,
		"chart_type": tpl.ChartType,
		"data": data,
	})
}

// ExportReport 导出报表 (GET /api/v1/reports/:id/export)
func (h *ReportHandler) ExportReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的报表ID格式"})
		return
	}

	format := c.DefaultQuery("format", "excel")
	if format != "excel" && format != "csv" && format != "pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的导出格式，支�? excel, csv, pdf"})
		return
	}

	ctx := c.Request.Context()
	tpl, err := h.reportRepo.GetByID(ctx, id)
	if err != nil || tpl == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "报表不存�?})
		return
	}

	h.logger.Info("导出报表",
		zap.String("report_id", id.String()),
		zap.String("format", format),
	)

	c.JSON(http.StatusOK, gin.H{
		"message":    "报表导出任务已提�?,
		"report_id":  id.String(),
		"report_name": tpl.Name,
		"format":     format,
		"download_url": "/api/v1/reports/" + id.String() + "/download?format=" + format,
	})
}

// generateMockStatsData 根据指标类型生成模拟统计数据
func (h *ReportHandler) generateMockStatsData(metricType string) interface{} {
	switch metricType {
	case "project_progress":
		return gin.H{
			"total_projects": 25,
			"completed": 18,
			"in_progress": 5,
			"not_started": 2,
			"completion_rate": 72.0,
		}
	case "workload":
		return gin.H{
			"total_hours": 12500,
			"completed_hours": 9800,
			"avg_per_person": 156.3,
		}
	case "quality":
		return gin.H{
			"bug_count": 45,
			"resolved": 38,
			"critical": 2,
			"resolution_rate": 84.4,
		}
	case "requirement_completion":
		return gin.H{
			"total_requirements": 120,
			"completed": 95,
			"in_progress": 15,
			"completion_rate": 79.2,
		}
	case "bug_trend":
		return []map[string]interface{}{
			{"date": "2026-01", "opened": 12, "resolved": 10},
			{"date": "2026-02", "opened": 8,  "resolved": 11},
			{"date": "2026-03", "opened": 15, "resolved": 13},
			{"date": "2026-04", "opened": 6,  "resolved": 9},
			{"date": "2026-05", "opened": 10, "resolved": 8},
			{"date": "2026-06", "opened": 7,  "resolved": 12},
		}
	default:
		return gin.H{"metric": metricType, "value": 0}
	}
}
