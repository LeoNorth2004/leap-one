package api

import("github.com/gin-gonic/gin";"leap-one/service-config/internal/interfaces/api/handler")

func RegisterRoutes(r*gin.Engine,cfgH*handler.ConfigHandler,flagH*handler.FlagHandler,auditH*handler.AuditLogHandler){
r.Use(gin.Recovery()); r.GET("/healthz",func(c*gin.Context){c.JSON(200,gin.H{"status":"ok","service":"leap-one-config"})})
v1:=r.Group("/api/v1")
{
	configs:=v1.Group("/configs")
	{configs.GET("",cfgH.ListConfigs);configs.GET("/:category/:key",cfgH.GetConfig);configs.PUT("/:category/:key",cfgH.UpdateConfig);configs.POST("",cfgH.BatchUpdateConfigs);configs.GET("/groups",cfgH.GetConfigGroups)}
	flags:=v1.Group("/feature-flags")
	{flags.GET("",flagH.ListFlags);flags.GET("/:key",flagH.GetFlag);flags.PUT("/:key",flagH.UpdateFlag);flags.POST("",flagH.CreateFlag);flags.DELETE("/:key",flagH.DeleteFlag)}
	audits:=v1.Group("/audit-logs")
	{audits.GET("",auditH.ListLogs);audits.GET("/:id",auditH.GetLog)}
}
}
