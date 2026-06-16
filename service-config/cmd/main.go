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
	"leap-one/service-config/internal/application"
	"leap-one/service-config/internal/config"
	"leap-one/service-config/internal/infrastructure/db"
	"leap-one/service-config/internal/infrastructure/repository_impl"
	"leap-one/service-config/internal/interfaces/api"
	"leap-one/service-config/internal/interfaces/api/handler"
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
	cfgRepo := repository_impl.NewSystemConfigRepository(database)
	flagRepo := repository_impl.NewFeatureFlagRepository(database)
	auditRepo := repository_impl.NewAuditLogRepository(database)
	_ = application.NewConfigService(cfgRepo, flagRepo, auditRepo, logger)
	cfgH := handler.NewConfigHandler(cfgRepo, logger)
	flagH := handler.NewFlagHandler(flagRepo, logger)
	auditH := handler.NewAuditLogHandler(auditRepo, logger)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	api.RegisterRoutes(r, cfgH, flagH, auditH)
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: r, ReadTimeout: cfg.Server.ReadTimeout, WriteTimeout: cfg.Server.WriteTimeout}
	go func() {
		logger.Info("配置中心启动", zap.String("addr", addr), zap.Int("port", cfg.Server.Port))
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
	logger.Info("配置中心已停�?)
}
