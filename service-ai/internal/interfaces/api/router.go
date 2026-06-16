package api

import (
	"github.com/gin-gonic/gin"
	"leap-one/service-ai/internal/interfaces/api/handler"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(
	r *gin.Engine,
	convHandler *handler.ConversationHandler,
	assistHandler *handler.AIAssistHandler,
	configHandler *handler.AIConfigHandler,
) {
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "leap-one-ai"})
	})

	v1 := r.Group("/api/v1")
	{
		// ---- AI对话 ----
		convs := v1.Group("/ai/conversations")
		{
			convs.POST("", convHandler.CreateConversation)
			convs.GET("", convHandler.ListConversations)
			convs.GET("/:id", convHandler.GetConversation)
			convs.DELETE("/:id", convHandler.DeleteConversation)
			convs.POST("/:id/messages", convHandler.SendMessage)
			convs.GET("/:id/stream", convHandler.StreamConnection)
		}

		// ---- AI辅助功能 ----
		aiAssist := v1.Group("/ai")
		{
			aiAssist.POST("/assist/requirement", assistHandler.AssistRequirement)
			aiAssist.POST("/assist/test-case", assistHandler.AssistTestCase)
			aiAssist.POST("/suggest/task-assign", assistHandler.SuggestTaskAssign)
			aiAssist.POST("/predict/requirements", assistHandler.PredictRequirements)
			aiAssist.POST("/identify/risks", assistHandler.IdentifyRisks)
			aiAssist.GET("/predictions", assistHandler.ListPredictions)
		}

		// ---- AI配置 ----
		aiConfig := v1.Group("/ai/config")
		{
			aiConfig.GET("", configHandler.GetConfig)
			aiConfig.PUT("", configHandler.UpdateConfig)
			aiConfig.POST("/test", configHandler.TestConnection)
		}
	}
}
