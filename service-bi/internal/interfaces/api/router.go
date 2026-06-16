package api

import (
	"github.com/gin-gonic/gin"
	"leap-one/service-bi/internal/interfaces/api/handler"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(
	r *gin.Engine,
	dashboardHandler *handler.DashboardHandler,
	reportHandler *handler.ReportHandler,
	statsHandler *handler.StatsHandler,
) {
	r.Use(gin.Recovery())

	// 健康检�?
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "leap-one-bi",
		})
	})

	// API路由�?
	v1 := r.Group("/api/v1")
	{
		// ---- BI大屏 ----
		dashboards := v1.Group("/dashboards")
		{
			dashboards.GET("/company-overview", dashboardHandler.GetCompanyOverview)
			dashboards.GET("/annual-overview", dashboardHandler.GetAnnualOverview)
			dashboards.GET("/ranking", dashboardHandler.GetRanking)
			dashboards.GET("/sprint-burndown", dashboardHandler.GetSprintBurndown)
			dashboards.GET("/annual-summary", dashboardHandler.GetAnnualSummary)
			dashboards.GET("/:id", dashboardHandler.GetDashboardByID)
		}

		// ---- 报表管理 ----
		reports := v1.Group("/reports")
		{
			reports.POST("", reportHandler.CreateReport)
			reports.GET("", reportHandler.ListReports)
			reports.GET("/:id", reportHandler.GetReport)
			reports.PUT("/:id", reportHandler.UpdateReport)
			reports.DELETE("/:id", reportHandler.DeleteReport)
			reports.GET("/:id/data", reportHandler.GetReportData)
			reports.GET("/:id/export", reportHandler.ExportReport)
		}

		// ---- 统计API ----
		stats := v1.Group("/stats")
		{
			stats.GET("/project-progress", statsHandler.ProjectProgress)
			stats.GET("/workload", statsHandler.Workload)
			stats.GET("/quality", statsHandler.Quality)
			stats.GET("/requirement-completion", statsHandler.RequirementCompletion)
			stats.GET("/bug-trends", statsHandler.BugTrends)
		}
	}
}
