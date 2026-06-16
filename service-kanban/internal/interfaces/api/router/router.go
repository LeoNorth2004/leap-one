package router

import (
	"github.com/gin-gonic/gin"

	"leap-one/service-kanban/internal/interfaces/api/handler"
)

// SetupRouter 配置API路由
func SetupRouter(h *handler.KanbanHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "leap-one-kanban", "db": "connected"})
	})

	v1 := r.Group("/api/v1")
	{
		// 看板管理
		boards := v1.Group("/boards")
		{
			boards.POST("", h.CreateBoard)
			boards.GET("", h.ListBoards)
			boards.GET("/:id", h.GetBoard)
			boards.PUT("/:id", h.UpdateBoard)
			boards.DELETE("/:id", h.DeleteBoard)

			// 看板�?
			boards.POST("/:id/columns", h.CreateColumn)
			boards.PUT("/:id/columns/:cid", h.UpdateColumn)
			boards.DELETE("/:id/columns/:cid", h.DeleteColumn)
			boards.PUT("/:id/columns/order", h.ReorderColumns)

			// 泳道
			boards.POST("/:id/swimlanes", h.CreateSwimlane)
			boards.PUT("/:id/swimlanes/:sid", h.UpdateSwimlane)
			boards.DELETE("/:id/swimlanes/:sid", h.DeleteSwimlane)

			// 卡片操作
			boards.POST("/:id/cards", h.CreateCard)
			boards.PUT("/:id/cards/:cardId", h.UpdateCard)
			boards.DELETE("/:id/cards/:cardId", h.DeleteCard)
			boards.PUT("/:id/cards/:cardId/move", h.MoveCard)
			boards.GET("/:id/cards/:cardId/history", h.GetMoveHistory)

			// 统计
			boards.GET("/:id/statistics", h.GetStatistics)
			boards.GET("/:id/cfd", h.GetCFD)
			boards.GET("/:id/export", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "导出功能待实�?})
			})
		}
	}

	return r
}
