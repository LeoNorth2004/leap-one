// Leap One API网关服务
// 负责路由转发、JWT验证、限流、CORS等网关功能
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"leap-one/service-gateway/internal/config"
	"leap-one/service-gateway/internal/interfaces/api/handler"
	"leap-one/service-gateway/internal/interfaces/api/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// 加载配置
	cfg, err := config.Load("")
	if err != nil {
		logger.Fatal("加载配置失败", zap.Error(err))
	}

	logger.Info("配置加载完成",
		zap.String("host", cfg.Server.Host),
		zap.Int("port", cfg.Server.Port),
		zap.Int("services", len(cfg.Services)),
	)

	// 设置Gin为Release模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎
	r := gin.New()

	// ===== 注册中间件链（顺序重要）=====
	// 1. CORS跨域
	r.Use(middleware.CORS(&cfg.CORS))
	// 2. Recovery异常恢复
	r.Use(gin.Recovery())
	// 3. 结构化请求日志
	r.Use(middleware.Logger())
	// 4. IP维度限流
	r.Use(middleware.RateLimit(&cfg.RateLimit))
	// 5. JWT认证（带白名单跳过）
	r.Use(middleware.JWTAuth(&cfg.JWT))

	// 初始化处理器
	authHandler := handler.NewAuthHandler(cfg, logger)
	proxyHandler := handler.NewProxyHandler(cfg, logger)
	healthHandler := handler.NewHealthHandler(cfg, logger)

	// ===== 注册路由 =====

	// 健康检查端点（无需认证，已在JWT白名单中）
	r.GET("/healthz", healthHandler.Check)
	r.GET("/readyz", healthHandler.Readiness)
	r.GET("/livez", healthHandler.Liveness)

	// API v1 路由组
	api := r.Group("/api/v1")
	{
		// --- 认证路由（白名单路径，无需Token） ---
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.LoginHandler)
			auth.POST("/register", authHandler.RegisterHandler)
			auth.POST("/refresh-token", authHandler.RefreshTokenHandler)
		}

		// --- 微服务代理路由 ---
		registerServiceRoutes(api, proxyHandler)
	}

	// 创建HTTP服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 启动服务器（协程）
	go func() {
		logger.Info("API网关服务启动",
			zap.String("addr", addr),
			zap.Int("port", cfg.Server.Port),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 优雅退出：监听系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("收到退出信号，正在关闭服务...", zap.String("signal", sig.String()))

	// 设置关闭超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("强制关闭服务", zap.Error(err))
	}

	logger.Info("API网关服务已安全停止")
}

// registerServiceRoutes 注册所有微服务的反向代理路由
func registerServiceRoutes(api *gin.RouterGroup, ph *handler.ProxyHandler) {
	// 服务名到URL路径前缀的映射
	serviceRoutes := map[string]string{
		"user-org":     "/user",
		"portfolio":    "/portfolio",
		"project":      "/project",
		"task":         "/task",
		"requirement":  "/requirement",
		"quality":      "/quality",
		"devops":       "/devops",
		"document":     "/document",
		"kanban":       "/kanban",
		"bi":           "/bi",
		"ai":           "/ai",
		"notification": "/notification",
		"search":       "/search",
		"config":       "/config",
	}

	for serviceName, pathPrefix := range serviceRoutes {
		svcGroup := api.Group(pathPrefix)
		svcGroup.Any("/*path", ph.Proxy(serviceName))
	}
}
