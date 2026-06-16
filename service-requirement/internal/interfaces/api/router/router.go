package router

import (
	"github.com/gin-gonic/gin"

	"leap-one/service-requirement/internal/interfaces/api/handler"
)

// SetupRouter 配置API路由
func SetupRouter(h *handler.RequirementHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 健康检�?
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "leap-one-requirement",
			"db":      "connected",
		})
	})

	// API v1 路由�?
	v1 := r.Group("/api/v1")
	{
		requirements := v1.Group("/requirements")
		{
			requirements.POST("", h.CreateRequirement)         // 创建需�?
			requirements.GET("", h.ListRequirements)           // 需求列表（分页+高级筛选）
			requirements.GET("/tree", h.GetRequirementTree)    // 需求树（产品维度）
			requirements.GET("/matrix", h.GetMatrix)           // 需求跟踪矩�?
			requirements.GET("/export", func(c *gin.Context) { // 导出需�?
				c.JSON(200, gin.H{"message": "导出功能待实�?})
			})

			// 单个需求操�?
			requirements.GET("/:id", h.GetRequirement)       // 需求详�?
			requirements.PUT("/:id", h.UpdateRequirement)    // 更新需�?
			requirements.DELETE("/:id", h.DeleteRequirement) // 删除需�?
			requirements.PUT("/:id/status", h.UpdateStatus)  // 更改状�?

			// 评审相关
			requirements.POST("/:id/review", h.SubmitReview) // 提交评审
			requirements.GET("/:id/reviews", h.GetReviews)   // 评审记录

			// 变更日志
			requirements.POST("/:id/change", h.CreateChangeLog) // 发起变更
			requirements.GET("/:id/changes", h.GetChangeLogs)   // 变更日志

			// 关联关系
			requirements.POST("/:id/relations", h.AddRelation) // 添加关联
			requirements.GET("/:id/relations", h.GetRelations) // 关联列表
		}
	}

	return r
}
