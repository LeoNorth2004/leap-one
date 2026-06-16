package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"leap-one/service-ai/internal/application"
	"leap-one/service-ai/internal/config"
	"leap-one/service-ai/internal/infrastructure/db"
	"leap-one/service-ai/internal/infrastructure/repository_impl"
	"leap-one/service-ai/internal/interfaces/api"
	"leap-one/service-ai/internal/interfaces/api/handler"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.Load("")
	if err != nil {
		logger.Fatal("加载配置失败", zap.Error(err))
	}

	database, err := db.InitPostgreSQL(cfg, logger)
	if err != nil {
		logger.Fatal("数据库初始化失败", zap.Error(err))
	}

	if err := db.AutoMigrate(database); err != nil {
		logger.Fatal("数据库迁移失�?, zap.Error(err))
	}
	logger.Info("数据库自动迁移完�?)

	sqlDB, _ := database.DB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Fatal("数据库健康检查失�?, zap.Error(err))
	}
	logger.Info("数据库连接正�?)

	// 依赖注入
	convRepo := repository_impl.NewConversationRepository(database)
	msgRepo := repository_impl.NewMessageRepository(database)
	predRepo := repository_impl.NewPredictionRepository(database)
	cfgRepo := repository_impl.NewAIConfigRepository(database)

	aiSvc := application.NewAIService(convRepo, msgRepo, predRepo, cfgRepo, logger)
	_ = aiSvc

	convHandler := handler.NewConversationHandler(convRepo, msgRepo, logger)
	assistHandler := handler.NewAIAssistHandler(predRepo, cfgRepo, logger)
	configHandler := handler.NewAIConfigHandler(cfgRepo, logger)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api.RegisterRoutes(r, convHandler, assistHandler, configHandler)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: r, ReadTimeout: cfg.Server.ReadTimeout, WriteTimeout: cfg.Server.WriteTimeout}

	go func() {
		logger.Info("AI服务启动", zap.String("addr", addr), zap.Int("port", cfg.Server.Port), zap.String("database", cfg.Database.DBName))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失�?, zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("收到退出信号，正在关闭服务...", zap.String("signal", sig.String()))

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("强制关闭服务", zap.Error(err))
	}
	if sqlDBErr := sqlDB.Close(); sqlDBErr != nil {
		logger.Warn("关闭数据库连接失�?, zap.Error(sqlDBErr))
	}
	logger.Info("AI服务已安全停�?)
}
