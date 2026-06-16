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
	"leap-one/service-notification/internal/application"
	"leap-one/service-notification/internal/config"
	"leap-one/service-notification/internal/infrastructure/db"
	"leap-one/service-notification/internal/infrastructure/repository_impl"
	"leap-one/service-notification/internal/interfaces/api"
	"leap-one/service-notification/internal/interfaces/api/handler"
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
	sqlDB, _ := database.DB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Fatal("数据库健康检查失�?, zap.Error(err))
	}
	notiRepo := repository_impl.NewNotificationRepository(database)
	tplRepo := repository_impl.NewTemplateRepository(database)
	emailLogRepo := repository_impl.NewEmailLogRepository(database)
	webhookCfgRepo := repository_impl.NewWebhookConfigRepository(database)
	webhookLogRepo := repository_impl.NewWebhookLogRepository(database)
	subRepo := repository_impl.NewSubscriptionRepository(database)
	_ = application.NewNotificationService(notiRepo, tplRepo, subRepo, logger)
	notiH := handler.NewNotificationHandler(notiRepo, subRepo, logger)
	tplH := handler.NewTemplateHandler(tplRepo, logger)
	emailH := handler.NewEmailLogHandler(emailLogRepo, logger)
	webhookH := handler.NewWebhookHandler(webhookCfgRepo, webhookLogRepo, logger)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	api.RegisterRoutes(r, notiH, tplH, emailH, webhookH)
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: r, ReadTimeout: cfg.Server.ReadTimeout, WriteTimeout: cfg.Server.WriteTimeout}
	go func() {
		logger.Info("消息通知服务启动", zap.String("addr", addr), zap.Int("port", cfg.Server.Port))
		if e := srv.ListenAndServe(); e != nil && e != http.ErrServerClosed {
			logger.Fatal("服务器启动失�?, zap.Error(e))
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sc, scCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer scCancel()
	srv.Shutdown(sc)
	sqlDB.Close()
	logger.Info("消息通知服务已停�?)
}
