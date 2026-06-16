package api

import (
	"leap-one/service-task/internal/interfaces/api/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(
	r *gin.Engine,
	taskHandler *handler.TaskHandler,
	issueHandler *handler.IssueHandler,
	templateHandler *handler.TemplateHandler,
	workflowHandler *handler.WorkflowHandler,
) {
	// 全局中间件
	r.Use(gin.Recovery())

	// ==================== 健康检查（公开）===================
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "leap-one-task",
		})
	})

	// ==================== API路由组 ====================
	v1 := r.Group("/api/v1")
	{
		// ---- 任务管理模块 ----
		tasks := v1.Group("/tasks")
		{
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("", taskHandler.ListTasks)
			tasks.GET("/my", taskHandler.MyTasks)
			tasks.GET("/export", func(c *gin.Context) { c.JSON(200, gin.H{"message": "导出功能开发中"}) })
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
			tasks.PUT("/:id/status", taskHandler.UpdateStatus)
			tasks.POST("/:id/start", taskHandler.StartTask)
			tasks.POST("/:id/complete", taskHandler.CompleteTask)
			tasks.POST("/:id/assign", taskHandler.AssignTask)
			tasks.DELETE("/:id/assignments/:uid", taskHandler.RemoveAssignment)
			tasks.POST("/:id/comments", taskHandler.AddComment)
			tasks.GET("/:id/comments", taskHandler.ListComments)
			tasks.DELETE("/:id/comments/:cid", taskHandler.DeleteComment)
			tasks.POST("/:id/attachments", taskHandler.AddAttachment)
			tasks.GET("/:id/attachments", taskHandler.ListAttachments)
			tasks.DELETE("/:id/attachments/:aid", taskHandler.DeleteAttachment)
			tasks.POST("/:id/worklogs", taskHandler.AddWorklog)
			tasks.GET("/:id/worklogs", taskHandler.ListWorklogs)
			tasks.POST("/:id/subtasks", taskHandler.CreateSubTask)
			tasks.GET("/:id/subtasks", taskHandler.ListSubTasks)
			tasks.POST("/:id/links", taskHandler.AddTaskLink)
			tasks.GET("/:id/links", taskHandler.ListTaskLinks)
			tasks.DELETE("/:id/links/:lid", taskHandler.DeleteTaskLink)
		}

		// ---- 工单管理模块 ----
		issues := v1.Group("/issues")
		{
			issues.POST("", issueHandler.CreateIssue)
			issues.GET("", issueHandler.ListIssues)
			issues.GET("/my", issueHandler.MyIssues)
			issues.GET("/export", func(c *gin.Context) { c.JSON(200, gin.H{"message": "导出功能开发中"}) })
			issues.GET("/:id", issueHandler.GetIssue)
			issues.PUT("/:id", issueHandler.UpdateIssue)
			issues.DELETE("/:id", issueHandler.DeleteIssue)
			issues.POST("/:id/transition", issueHandler.TransitionIssue)
			issues.POST("/:id/comments", issueHandler.AddComment)
			issues.GET("/:id/comments", issueHandler.ListComments)
			issues.POST("/:id/attachments", issueHandler.AddAttachment)
			issues.GET("/:id/attachments", issueHandler.ListAttachments)
			issues.GET("/:id/sla", issueHandler.GetSLAInfo)
			issues.POST("/:id/satisfaction", issueHandler.RateSatisfaction)
		}

		// ---- 工单模板管理模块 ----
		templates := v1.Group("/issue-templates")
		{
			templates.POST("", templateHandler.CreateTemplate)
			templates.GET("", templateHandler.ListTemplates)
			templates.GET("/:id", templateHandler.GetTemplate)
			templates.PUT("/:id", templateHandler.UpdateTemplate)
			templates.DELETE("/:id", templateHandler.DeleteTemplate)
		}

		// ---- 工作流管理模块 ----
		workflows := v1.Group("/workflows")
		{
			workflows.POST("", workflowHandler.CreateWorkflow)
			workflows.GET("", workflowHandler.ListWorkflows)
			workflows.GET("/:id", workflowHandler.GetWorkflow)
			workflows.PUT("/:id", workflowHandler.UpdateWorkflow)
			workflows.DELETE("/:id", workflowHandler.DeleteWorkflow)
			workflows.POST("/:id/transitions", workflowHandler.AddTransition)
		}

		// ---- SLA配置管理模块 ----
		slaConfigs := v1.Group("/sla-configs")
		{
			slaConfigs.POST("", workflowHandler.CreateSLAConfig)
			slaConfigs.GET("", workflowHandler.ListSLAConfigs)
			slaConfigs.GET("/:id", workflowHandler.GetSLAConfig)
			slaConfigs.PUT("/:id", workflowHandler.UpdateSLAConfig)
			slaConfigs.DELETE("/:id", workflowHandler.DeleteSLAConfig)
		}
	}
}
