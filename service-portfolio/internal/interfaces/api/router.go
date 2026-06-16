package api

import (
	"leap-one/service-portfolio/internal/interfaces/api/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(
	r *gin.Engine,
	programHandler *handler.ProgramHandler,
	productHandler *handler.ProductHandler,
	productLineHandler *handler.ProductLineHandler,
	versionHandler *handler.VersionHandler,
	roadmapHandler *handler.RoadmapHandler,
	planHandler *handler.PlanHandler,
) {
	// 全局中间件
	r.Use(gin.Recovery())

	// ==================== 健康检查（公开）====================
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "leap-one-portfolio",
		})
	})

	// ==================== API路由组 ====================
	v1 := r.Group("/api/v1")
	{
		// ---- 项目集管理模块 ----
		programs := v1.Group("/programs")
		{
			programs.POST("", programHandler.CreateProgram)
			programs.GET("", programHandler.ListPrograms)
			programs.GET("/tree", programHandler.GetProgramTree)
			programs.GET("/:id", programHandler.GetProgram)
			programs.PUT("/:id", programHandler.UpdateProgram)
			programs.DELETE("/:id", programHandler.DeleteProgram)
			// 里程碑管理
			programs.POST("/:id/milestones", programHandler.CreateMilestone)
			programs.GET("/:id/milestones", programHandler.ListMilestones)
			// 风险管理
			programs.POST("/:id/risks", programHandler.CreateRisk)
			programs.GET("/:id/risks", programHandler.ListRisks)
			// 统计信息
			programs.GET("/:id/statistics", programHandler.GetStatistics)
		}

		// ---- 产品管理模块 ----
		products := v1.Group("/products")
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("", productHandler.ListProducts)
			products.GET("/:id", productHandler.GetProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)
			// 路线图管理
			roadmap := products.Group("/:id/roadmap")
			{
				roadmap.GET("", roadmapHandler.ListRoadmapItems)
				roadmap.POST("", roadmapHandler.CreateRoadmapItem)
				roadmap.PUT("/reorder", roadmapHandler.ReorderRoadmapItems)
				roadmap.PUT("/:rid", roadmapHandler.UpdateRoadmapItem)
				roadmap.DELETE("/:rid", roadmapHandler.DeleteRoadmapItem)
			}
			// 版本管理
			versions := products.Group("/:id/versions")
			{
				versions.GET("", versionHandler.ListVersions)
				versions.POST("", versionHandler.CreateVersion)
				versions.PUT("/:vid", versionHandler.UpdateVersion)
				versions.POST("/:vid/release", versionHandler.ReleaseVersion)
			}
			// 计划管理（按产品）
			products.GET("/:id/plans", planHandler.ListProductPlans)
		}

		// ---- 产品线管理模块 ----
		productLines := v1.Group("/product-lines")
		{
			productLines.POST("", productLineHandler.CreateProductLine)
			productLines.GET("", productLineHandler.ListProductLines)
			productLines.GET("/:id", productLineHandler.GetProductLine)
			productLines.PUT("/:id", productLineHandler.UpdateProductLine)
			productLines.DELETE("/:id", productLineHandler.DeleteProductLine)
		}

		// ---- 产品计划管理模块 ----
		plans := v1.Group("/plans")
		{
			plans.POST("", planHandler.CreatePlan)
			plans.GET("", planHandler.ListPlans)
			plans.GET("/:id", planHandler.GetPlan)
			plans.PUT("/:id", planHandler.UpdatePlan)
			plans.DELETE("/:id", planHandler.DeletePlan)
		}
	}
}
