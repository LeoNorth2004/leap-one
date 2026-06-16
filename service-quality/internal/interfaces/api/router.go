package api

import (
	"github.com/gin-gonic/gin"
	"leap-one/service-quality/internal/interfaces/api/handler"
)

// RegisterRoutes 注册所有API路由
// 按照DDD分层架构组织，包含测试用例、套件、计划、Bug、环境和报表六大模块
func RegisterRoutes(
	r *gin.Engine,
	caseHandler *handler.CaseHandler,
	suiteHandler *handler.SuiteHandler,
	planHandler *handler.PlanHandler,
	bugHandler *handler.BugHandler,
	envHandler *handler.EnvironmentHandler,
	reportHandler *handler.ReportHandler,
) {
	// 全局中间�?
	r.Use(gin.Recovery())

	// ==================== 健康检查（公开�?===================
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "leap-one-quality",
		})
	})

	// ==================== API v1 路由�?====================
	v1 := r.Group("/api/v1")
	{
		// ---- 测试用例模块 ----
		cases := v1.Group("/test-cases")
		{
			cases.POST("", caseHandler.CreateCase)         // 创建用例
			cases.GET("", caseHandler.ListCases)           // 用例列表（分�?筛选）
			cases.GET("/:id", caseHandler.GetCase)         // 用例详情
			cases.PUT("/:id", caseHandler.UpdateCase)      // 更新用例
			cases.DELETE("/:id", caseHandler.DeleteCase)   // 删除用例
			cases.POST("/import", caseHandler.ImportCases) // 导入用例
			cases.GET("/export", func(c *gin.Context) {    // 导出用例（TODO: 实现导出逻辑�?
				c.JSON(200, gin.H{"message": "导出功能待实�?})
			})
			cases.POST("/:id/review", caseHandler.ReviewCase) // 评审用例
		}

		// ---- 测试套件模块 ----
		suites := v1.Group("/test-suites")
		{
			suites.POST("", suiteHandler.CreateSuite)                          // 创建套件
			suites.GET("", suiteHandler.ListSuites)                            // 套件列表
			suites.GET("/:id", suiteHandler.GetSuite)                          // 套件详情+用例列表
			suites.PUT("/:id", suiteHandler.UpdateSuite)                       // 更新套件
			suites.DELETE("/:id", suiteHandler.DeleteSuite)                    // 删除套件
			suites.POST("/:id/cases", suiteHandler.AddCasesToSuite)            // 添加用例到套�?
			suites.DELETE("/:id/cases/:cid", suiteHandler.RemoveCaseFromSuite) // 移除用例
		}

		// ---- 测试计划模块 ----
		plans := v1.Group("/test-plans")
		{
			plans.POST("", planHandler.CreatePlan)                          // 创建计划
			plans.GET("", planHandler.ListPlans)                            // 计划列表
			plans.GET("/:id", planHandler.GetPlan)                          // 计划详情
			plans.PUT("/:id", planHandler.UpdatePlan)                       // 更新计划
			plans.DELETE("/:id", planHandler.DeletePlan)                    // 删除计划
			plans.POST("/:id/start", planHandler.StartPlan)                 // 开始执�?
			plans.POST("/:id/complete", planHandler.CompletePlan)           // 完成计划
			plans.POST("/:id/cases", planHandler.AddCasesToPlan)            // 添加用例到计�?
			plans.POST("/:id/cases/:pcid/execute", planHandler.ExecuteCase) // 执行用例
		}

		// ---- Bug管理模块 ----
		bugs := v1.Group("/bugs")
		{
			bugs.POST("", bugHandler.CreateBug)                        // 创建Bug
			bugs.GET("", bugHandler.ListBugs)                          // Bug列表（高级筛选）
			bugs.GET("/:id", bugHandler.GetBug)                        // Bug详情（含历史�?
			bugs.PUT("/:id", bugHandler.UpdateBug)                     // 更新Bug
			bugs.DELETE("/:id", bugHandler.DeleteBug)                  // 删除Bug
			bugs.POST("/:id/confirm", bugHandler.ConfirmBug)           // 确认Bug
			bugs.POST("/:id/resolve", bugHandler.ResolveBug)           // 解决Bug
			bugs.POST("/:id/close", bugHandler.CloseBug)               // 关闭Bug
			bugs.POST("/:id/reopen", bugHandler.ReopenBug)             // 激活Bug
			bugs.POST("/:id/comments", bugHandler.AddComment)          // 添加评论
			bugs.GET("/:id/comments", bugHandler.ListComments)         // 评论历史
			bugs.POST("/:id/attachments", bugHandler.UploadAttachment) // 上传附件
			bugs.GET("/:id/history", bugHandler.ListHistory)           // 变更历史
			bugs.GET("/my", bugHandler.MyBugs)                         // 我的Bug
			bugs.GET("/export", func(c *gin.Context) {                 // 导出Bug（TODO�?
				c.JSON(200, gin.H{"message": "导出功能待实�?})
			})
		}

		// ---- 测试环境模块 ----
		envs := v1.Group("/environments")
		{
			envs.POST("", envHandler.CreateEnvironment)       // 创建环境
			envs.GET("", envHandler.ListEnvironments)         // 环境列表
			envs.GET("/:id", envHandler.GetEnvironment)       // 环境详情
			envs.PUT("/:id", envHandler.UpdateEnvironment)    // 更新环境
			envs.DELETE("/:id", envHandler.DeleteEnvironment) // 删除环境
		}

		// ---- 报表统计模块 ----
		quality := v1.Group("/quality")
		{
			quality.GET("/statistics", reportHandler.QualityStatistics) // 质量统计概览
			quality.GET("/bug-trends", reportHandler.BugTrends)         // Bug趋势分析
			quality.GET("/pass-rate", reportHandler.PassRate)           // 通过率统�?
		}
	}
}
