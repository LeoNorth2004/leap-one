// Leap One 质量管理服务
// 负责测试用例管理、测试套件、测试计划、Bug跟踪、测试环境配置和统计报表等功�?
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
	"leap-one/service-quality/internal/application"
	"leap-one/service-quality/internal/config"
	"leap-one/service-quality/internal/infrastructure/cache"
	"leap-one/service-quality/internal/infrastructure/db"
	"leap-one/service-quality/internal/infrastructure/repository_impl"
	"leap-one/service-quality/internal/interfaces/api"
	"leap-one/service-quality/internal/interfaces/api/handler"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

	// 初始化Redis缓存（可选，失败不阻塞启动）
	var redisCache *cache.RedisClient
	if redisClient, redisErr := cache.InitRedis(&cfg.Redis, logger); redisErr != nil {
		logger.Warn("Redis初始化失败，将使用无缓存模式运行", zap.Error(redisErr))
	} else {
		redisCache = redisClient
		defer redisCache.Close()
	}

	// ==================== 依赖注入：创建Repository实现 ====================
	caseRepo := repository_impl.NewTestCaseRepository(database)
	suiteRepo := repository_impl.NewTestSuiteRepository(database)
	planRepo := repository_impl.NewTestPlanRepository(database)
	bugRepo := repository_impl.NewBugRepository(database)
	envRepo := repository_impl.NewEnvironmentRepository(database)

	// ==================== 预置数据初始化（首次启动自动创建�?===================
	initPreseededData(database, logger)

	// ==================== 创建应用服务 ====================
	qualitySvc := application.NewQualityService(bugRepo, caseRepo, logger)
	_ = qualitySvc // 应用服务可被Handler或预置初始化使用
	_ = redisCache // Redis缓存可用于后续扩展（如会话缓存等�?

	// ==================== 创建Handler实例 ====================
	caseHandler := handler.NewCaseHandler(caseRepo, logger)
	suiteHandler := handler.NewSuiteHandler(suiteRepo, caseRepo, logger)
	planHandler := handler.NewPlanHandler(planRepo, caseRepo, logger)
	bugHandler := handler.NewBugHandler(bugRepo, logger)
	envHandler := handler.NewEnvironmentHandler(envRepo, logger)
	reportHandler := handler.NewReportHandler(bugRepo, caseRepo, planRepo, logger)

	// ==================== 设置Gin引擎和注册路�?=====================
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api.RegisterRoutes(r,
		caseHandler,
		suiteHandler,
		planHandler,
		bugHandler,
		envHandler,
		reportHandler,
	)

	// ==================== 创建HTTP服务器并启动 ====================
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("质量管理服务启动",
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

	// 关闭数据库连�?
	if sqlDBErr := sqlDB.Close(); sqlDBErr != nil {
		logger.Warn("关闭数据库连接失�?, zap.Error(sqlDBErr))
	}

	logger.Info("质量管理服务已安全停�?)
}

// initPreseededData 初始化预置数据（默认Bug工作流、测试环境）
// 仅在数据不存在时插入，支持幂等执�?
func initPreseededData(dbConn *gorm.DB, logger *zap.Logger) {
	ctx := context.Background()

	// 创建应用服务用于预置数据初始�?
	qualitySvc := application.NewQualityService(nil, nil, logger)

	// ---- 1. 创建默认Bug工作�?----
	if err := qualitySvc.InitDefaultWorkflow(ctx, dbConn); err != nil {
		logger.Warn("初始化默认Bug工作流失�?, zap.Error(err))
	}

	// ---- 2. 创建默认测试环境 ----
	if err := qualitySvc.InitDefaultEnvironments(ctx, dbConn); err != nil {
		logger.Warn("初始化默认测试环境失�?, zap.Error(err))
	}

	logger.Info("预置数据初始化完�?)
}
