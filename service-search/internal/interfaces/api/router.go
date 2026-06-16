package api

import (
	"github.com/gin-gonic/gin"
	"leap-one/service-search/internal/interfaces/api/handler"
)

func RegisterRoutes(r *gin.Engine, searchH *handler.SearchHandler) {
	r.Use(gin.Recovery())
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "service": "leap-one-search"}) })
	v1 := r.Group("/api/v1")
	{
		search := v1.Group("/search")
		{
			search.GET("", searchH.GlobalSearch)
			search.GET("/advanced", searchH.AdvancedSearch)
			search.POST("/save", searchH.SaveSearch)
			search.GET("/saved", searchH.ListSavedSearches)
			search.DELETE("/saved/:id", searchH.DeleteSavedSearch)
			search.GET("/history", searchH.SearchHistory)
			search.DELETE("/history", searchH.ClearHistory)
			search.GET("/suggestions", searchH.Suggestions)
			search.POST("/index", searchH.TriggerIndex)
			search.GET("/index/status", searchH.IndexStatus)
		}
	}
}
