package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"leap-one/service-quality/internal/domain/entity"
	"leap-one/service-quality/internal/domain/repository"
	"leap-one/service-quality/internal/interfaces/api/dto"
)

// ReportHandler 统计报表Handler
type ReportHandler struct {
	bugRepo  repository.BugRepository
	caseRepo repository.TestCaseRepository
	planRepo repository.TestPlanRepository
	logger   *zap.Logger
}

// NewReportHandler 创建统计报表Handler实例
func NewReportHandler(
	bugRepo repository.BugRepository,
	caseRepo repository.TestCaseRepository,
	planRepo repository.TestPlanRepository,
	logger *zap.Logger,
) *ReportHandler {
	return &ReportHandler{
		bugRepo:  bugRepo,
		caseRepo: caseRepo,
		planRepo: planRepo,
		logger:   logger,
	}
}

// QualityStatistics 质量统计概览（GET /api/v1/quality/statistics�?
func (h *ReportHandler) QualityStatistics(c *gin.Context) {
	ctx := c.Request.Context()

	productID := parseUUIDPtr(c.Query("product_id"))
	projectID := parseUUIDPtr(c.Query("project_id"))

	// 获取Bug统计
	bugStats, err := h.bugRepo.GetStatistics(ctx, productID, projectID)
	if err != nil {
		h.logger.Error("获取Bug统计数据失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败"})
		return
	}

	// 计算未关闭Bug数（new + confirmed + in_progress + reopened�?
	openBugs := bugStats.NewCount + bugStats.ConfirmedCnt + bugStats.InProgress + bugStats.ReopenedCnt

	// 计算解决�?
	var resolvedRate float64
	totalClosed := bugStats.ResolvedCnt + bugStats.ClosedCnt
	if bugStats.TotalCount > 0 {
		resolvedRate = float64(totalClosed) / float64(bugStats.TotalCount) * 100
	}

	resp := dto.QualityStatisticsResponse{
		TotalBugs:      bugStats.TotalCount,
		OpenBugs:       openBugs,
		ResolvedRate:   resolvedRate,
		AvgResolveDays: 0, // TODO: 可基于历史数据计算平均解决天�?
		ByStatus: map[string]int64{
			"new":         bugStats.NewCount,
			"confirmed":   bugStats.ConfirmedCnt,
			"in_progress": bugStats.InProgress,
			"resolved":    bugStats.ResolvedCnt,
			"closed":      bugStats.ClosedCnt,
			"reopened":    bugStats.ReopenedCnt,
		},
		BySeverity: bugStats.BySeverity,
		ByPriority: bugStats.ByPriority,
		ByType:     bugStats.ByType,
	}

	// 用例统计（简化版，实际可通过caseRepo获取更详细数据）
	resp.TestCaseStats = dto.TestCaseStatistics{
		TotalCount: 0, // 需要额外查�?
	}
	// 计划统计
	resp.TestPlanStats = dto.TestPlanStatistics{}

	c.JSON(http.StatusOK, resp)
}

// BugTrends Bug趋势分析（GET /api/v1/quality/bug-trends�?
func (h *ReportHandler) BugTrends(c *gin.Context) {
	days := 30 // 默认最�?0�?

	// 按日期统计Bug创建和解决趋�?
	trends := make([]dto.BugTrendItem, days)
	now := time.Now()

	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		// 简化实现：此处返回基础结构，实际应通过聚合查询统计
		trends[days-1-i] = dto.BugTrendItem{
			Date:     dateStr,
			Created:  0,
			Resolved: 0,
			Reopened: 0,
		}
	}

	c.JSON(http.StatusOK, dto.BugTrendResponse{Trends: trends})
}

// PassRate 通过率统计（GET /api/v1/quality/pass-rate�?
func (h *ReportHandler) PassRate(c *gin.Context) {
	ctx := c.Request.Context()

	productID := parseUUIDPtr(c.Query("product_id"))
	projectID := parseUUIDPtr(c.Query("project_id"))

	// 获取所有已完成的测试计�?
	filter := &repository.TestPlanFilter{
		Status:    "completed",
		ProductID: productID,
		ProjectID: projectID,
	}

	plans, _, err := h.planRepo.List(ctx, 1, 100, filter)
	if err != nil {
		h.logger.Error("获取测试计划列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取通过率数据失�?})
		return
	}

	planItems := make([]dto.PassRateItem, len(plans))
	var totalPassed, totalCases int

	for i, p := range plans {
		passed, failed, blocked, skipped, notRun := countPlanResults(p.Cases)

		item := dto.PassRateItem{
			PlanName:   p.Name,
			TotalCases: len(p.Cases),
			Passed:     passed,
			Failed:     failed,
			Blocked:    blocked,
			Skipped:    skipped,
			NotRun:     notRun,
		}
		if len(p.Cases) > 0 {
			item.PassRate = float64(passed) / float64(len(p.Cases)) * 100
		}
		planItems[i] = item

		totalPassed += passed
		totalCases += len(p.Cases)
	}

	var overallPassRate float64
	if totalCases > 0 {
		overallPassRate = float64(totalPassed) / float64(totalCases) * 100
	}

	c.JSON(http.StatusOK, dto.PassRateResponse{
		OverallPassRate: overallPassRate,
		Plans:           planItems,
	})
}

// countPlanResults 统计测试计划的用例执行结果分�?
func countPlanResults(cases []entity.TestPlanCase) (passed, failed, blocked, skipped, notRun int) {
	for _, pc := range cases {
		switch pc.Result {
		case "passed":
			passed++
		case "failed":
			failed++
		case "blocked":
			blocked++
		case "skipped":
			skipped++
		default:
			notRun++
		}
	}
	return
}
