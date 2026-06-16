package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"leap-one/service-project/internal/interfaces/api/handler"
	"leap-one/service-project/internal/interfaces/api/middleware"
)

// RouterConfig 路由配置参数
type RouterConfig struct {
	JWTSecret string
}

// RegisterRoutes 注册所有API路由
func RegisterRoutes(
	r *gin.Engine,
	projectHandler *handler.ProjectHandler,
	memberHandler *handler.MemberHandler,
	milestoneHandler *handler.MilestoneHandler,
	riskHandler *handler.RiskHandler,
	customFieldHandler *handler.CustomFieldHandler,
	iterationHandler *handler.IterationHandler,
	templateHandler *handler.TemplateHandler,
	statsHandler *handler.StatisticsHandler,
	config RouterConfig,
) {
	// 全局中间�?
	r.Use(gin.Recovery())

	// ==================== 健康检查（公开�?===================
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "leap-one-project",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// ==================== 需要认证的路由（JWT中间件）====================
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(config.JWTSecret))
	{

		// ==================== 项目管理模块 ====================
		projects := protected.Group("/projects")
		{
			projects.GET("", projectHandler.ListProjects)                // 项目列表（分�?筛�?搜索�?
			projects.POST("", projectHandler.CreateProject)              // 创建项目
			projects.GET("/:id", projectHandler.GetProject)              // 项目详情
			projects.PUT("/:id", projectHandler.UpdateProject)           // 更新项目
			projects.DELETE("/:id", projectHandler.DeleteProject)        // 删除项目(软删�?
			projects.POST("/:id/archive", projectHandler.ArchiveProject) // 归档项目
			projects.POST("/:id/cancel", projectHandler.CancelProject)   // 取消项目

			// 项目成员管理
			projects.GET("/:id/members", memberHandler.ListMembers)           // 成员列表
			projects.POST("/:id/members", memberHandler.AddMember)            // 添加成员
			projects.DELETE("/:id/members/:uid", memberHandler.RemoveMember)  // 移除成员
			projects.PUT("/:id/members/:uid", memberHandler.UpdateMemberRole) // 更新成员角色

			// 项目里程碑管�?
			projects.GET("/:id/milestones", milestoneHandler.ListMilestones)                  // 里程碑列�?
			projects.POST("/:id/milestones", milestoneHandler.CreateMilestone)                // 创建里程�?
			projects.PUT("/:id/milestones/:mid", milestoneHandler.UpdateMilestone)            // 更新里程�?
			projects.DELETE("/:id/milestones/:mid", milestoneHandler.DeleteMilestone)         // 删除里程�?
			projects.PUT("/:id/milestones/:mid/complete", milestoneHandler.CompleteMilestone) // 完成里程�?

			// 项目风险管理
			projects.GET("/:id/risks", riskHandler.ListRisks)          // 风险列表
			projects.POST("/:id/risks", riskHandler.CreateRisk)        // 创建风险
			projects.PUT("/:id/risks/:rid", riskHandler.UpdateRisk)    // 更新风险
			projects.DELETE("/:id/risks/:rid", riskHandler.DeleteRisk) // 删除风险

			// 自定义字段管�?
			projects.GET("/:id/custom-fields", customFieldHandler.ListCustomFields)          // 自定义字段列�?
			projects.POST("/:id/custom-fields", customFieldHandler.AddCustomField)           // 添加自定义字�?
			projects.PUT("/:id/custom-fields/:fid", customFieldHandler.UpdateCustomField)    // 更新自定义字�?
			projects.DELETE("/:id/custom-fields/:fid", customFieldHandler.DeleteCustomField) // 删除自定义字�?

			// 项目统计
			projects.GET("/:id/statistics", statsHandler.GetProjectStatistics) // 项目统计数据
			projects.GET("/:id/burndown", statsHandler.GetProjectBurndown)     // 燃尽图数�?
			projects.GET("/:id/gantt", statsHandler.GetProjectGantt)           // 甘特图数�?
		}

		// ==================== 迭代/Sprint管理模块 ====================
		iterations := protected.Group("/iterations")
		{
			iterations.POST("", iterationHandler.CreateIteration)                  // 创建迭代
			iterations.GET("", iterationHandler.ListIterations)                    // 迭代列表
			iterations.GET("/:id", iterationHandler.GetIteration)                  // 迭代详情
			iterations.PUT("/:id", iterationHandler.UpdateIteration)               // 更新迭代
			iterations.DELETE("/:id", iterationHandler.DeleteIteration)            // 删除迭代
			iterations.POST("/:id/start", iterationHandler.StartIteration)         // 开始迭�?
			iterations.POST("/:id/complete", iterationHandler.CompleteIteration)   // 完成迭代
			iterations.GET("/:id/board", iterationHandler.GetIterationBoard)       // 迭代看板数据
			iterations.GET("/:id/burndown", iterationHandler.GetIterationBurndown) // 迭代燃尽图数�?
		}

		// ==================== 项目模板管理模块 ====================
		templates := protected.Group("/templates")
		{
			templates.GET("", templateHandler.ListTemplates)         // 模板列表
			templates.POST("", templateHandler.CreateTemplate)       // 创建模板
			templates.GET("/:id", templateHandler.GetTemplate)       // 模板详情
			templates.PUT("/:id", templateHandler.UpdateTemplate)    // 更新模板
			templates.DELETE("/:id", templateHandler.DeleteTemplate) // 删除模板
		}
	}
}
