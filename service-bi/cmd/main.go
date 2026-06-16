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
	"leap-one/service-bi/internal/application"
	"leap-one/service-bi/internal/config"
	"leap-one/service-bi/internal/infrastructure/db"
	"leap-one/service-bi/internal/infrastructure/repository_impl"
	"leap-one/service-bi/internal/interfaces/api"
	"leap-one/service-bi/internal/interfaces/api/handler"
	"go.uber.org/zap"
)

func main() {
	// 初始化日�?
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// 加载配置
	cfg, err := config.Load("")
	if err != nil {
		logger.Fatal("加载配置失败", zap.Error(err))
	}

	// 初始化PostgreSQL数据库连�?
	database, err := db.InitPostgreSQL(cfg, logger)
	if err != nil {
		logger.Fatal("数据库初始化失败", zap.Error(err))
	}

	// 自动迁移所有实体表
	if err := db.AutoMigrate(database); err != nil {
		logger.Fatal("数据库迁移失�?, zap.Error(err))
	}
	logger.Info("数据库自动迁移完�?)

	// 启动时检查数据库健康状�?
	sqlDB, _ := database.DB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Fatal("数据库健康检查失�?, zap.Error(err))
	}
	logger.Info("数据库连接正�?)

	// ==================== 依赖注入：创建Repository实现 ====================
	dashboardRepo := repository_impl.NewDashboardConfigRepository(database)
	reportRepo := repository_impl.NewReportTemplateRepository(database)
	snapshotRepo := repository_impl.NewDataSnapshotRepository(database)

	// ==================== 创建应用服务 ====================
	biSvc := application.NewBIStatService(dashboardRepo, reportRepo, snapshotRepo, logger)
	_ = biSvc

	// ==================== 创建Handler实例 ====================
	dashboardHandler := handler.NewDashboardHandler(dashboardRepo, logger)
	reportHandler := handler.NewReportHandler(reportRepo, snapshotRepo, logger)
	statsHandler := handler.NewStatsHandler(logger)

	// ==================== 设置Gin引擎和注册路�?====================
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api.RegisterRoutes(r, dashboardHandler, reportHandler, statsHandler)

	// ==================== 创建HTTP服务器并启动 ====================
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("BI统计服务启动",
			zap.String("addr", addr),
			zap.Int("port", cfg.Server.Port),
			zap.String("database", cfg.Database.DBName),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失�?, zap.Error(err))
		}
	}()

	// 优雅退�?
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

	logger.Info("BI统计服务已安全停�?)
}
