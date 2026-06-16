package api

import (
	"github.com/gin-gonic/gin"
	"leap-one/service-notification/internal/interfaces/api/handler"
)

func RegisterRoutes(
	r *gin.Engine,
	notiHandler *handler.NotificationHandler,
	tplHandler *handler.TemplateHandler,
	emailHandler *handler.EmailLogHandler,
	webhookHandler *handler.WebhookHandler,
) {
	r.Use(gin.Recovery())
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "service": "leap-one-notification"}) })

	v1 := r.Group("/api/v1")
	{
		// 通知消息
		notis := v1.Group("/notifications")
		{
			notis.GET("", notiHandler.ListNotifications)
			notis.GET("/unread-count", notiHandler.GetUnreadCount)
			notis.PUT("/:id/read", notiHandler.MarkAsRead)
			notis.PUT("/read-all", notiHandler.MarkAllAsRead)
			notis.DELETE("/:id", notiHandler.DeleteNotification)
			notis.GET("/settings", notiHandler.GetSettings)
			notis.PUT("/settings", notiHandler.UpdateSettings)
		}
		// 消息模板
		tpls := v1.Group("/notification-templates")
		{
			tpls.POST("", tplHandler.CreateTemplate)
			tpls.GET("", tplHandler.ListTemplates)
			tpls.GET("/:id", tplHandler.GetTemplate)
			tpls.PUT("/:id", tplHandler.UpdateTemplate)
			tpls.DELETE("/:id", tplHandler.DeleteTemplate)
		}
		// 邮件日志
		emailLogs := v1.Group("/email-logs")
		{
			emailLogs.GET("", emailHandler.ListEmailLogs)
			emailLogs.GET("/:id", emailHandler.GetEmailLog)
			emailLogs.POST("/:id/resend", emailHandler.ResendEmail)
		}
		// Webhook
		webhooks := v1.Group("/webhooks")
		{
			webhooks.POST("", webhookHandler.CreateWebhook)
			webhooks.GET("", webhookHandler.ListWebhooks)
			webhooks.GET("/:id", webhookHandler.GetWebhook)
			webhooks.PUT("/:id", webhookHandler.UpdateWebhook)
			webhooks.DELETE("/:id", webhookHandler.DeleteWebhook)
			webhooks.POST("/:id/test", webhookHandler.TestWebhook)
			webhooks.GET("/:id/logs", webhookHandler.ListWebhookLogs)
		}
	}
}
