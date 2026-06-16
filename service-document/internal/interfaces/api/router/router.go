package router

import (
	"github.com/gin-gonic/gin"

	"leap-one/service-document/internal/interfaces/api/handler"
)

// SetupRouter 配置API路由
func SetupRouter(h *handler.DocumentHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "leap-one-document", "db": "connected"})
	})

	v1 := r.Group("/api/v1")
	{
		// 文档管理
		docs := v1.Group("/documents")
		{
			docs.POST("", h.Create)
			docs.GET("", h.List)
			docs.GET("/tree", h.GetTree)
			docs.GET("/search", h.Search)
			docs.GET("/:id", h.GetByID)
			docs.PUT("/:id", h.Update)
			docs.DELETE("/:id", h.Delete)
			docs.POST("/:id/publish", h.Publish)
			docs.GET("/:id/versions", h.ListVersions)
			docs.GET("/:id/versions/:vid", h.GetVersion)
			docs.POST("/:id/restore", h.Restore)
			docs.POST("/:id/comments", h.AddComment)
			docs.GET("/:id/comments", h.ListComments)
			docs.POST("/:id/favorite", h.Favorite)
			docs.DELETE("/:id/favorite", h.Unfavorite)
			docs.POST("/:id/attachments", h.UploadAttachment)
			docs.GET("/export/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "导出功能待实�?}) })
		}

		// 文档分类
		cats := v1.Group("/categories")
		{
			cats.POST("", h.CreateCategory)
			cats.GET("", h.ListCategories)
			cats.PUT("/:id", h.UpdateCategory)
			cats.DELETE("/:id", h.DeleteCategory)
		}

		// 知识�?		kbs := v1.Group("/knowledge-bases")
		{
			kbs.POST("", h.CreateKB)
			kbs.GET("", h.ListKBs)
			kbs.GET("/:id", h.GetKB)
			kbs.PUT("/:id", h.UpdateKB)
			kbs.DELETE("/:id", h.DeleteKB)
		}

		// 模板�?		v1.GET("/templates", h.ListTemplates)
	}

	return r
}
