// Leap One 项目集与产品服务
// 负责项目集管理、产品管理、产品路线图、版本发布等功能
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"leap-one/service-portfolio/internal/application"
	"leap-one/service-portfolio/internal/config"
	"leap-one/service-portfolio/internal/infrastructure/cache"
	"leap-one/service-portfolio/internal/infrastructure/db"
	"leap-one/service-portfolio/internal/infrastructure/repository_impl"
	"leap-one/service-portfolio/internal/interfaces/api"
	"leap-one/service-portfolio/internal/interfaces/api/handler"

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

	// 初始化PostgreSQL数据库连接
	database, err := db.InitPostgreSQL(cfg, logger)
	if err != nil {
		logger.Fatal("数据库初始化失败", zap.Error(err))
	}

	// 自动迁移所有实体表
	if err := db.AutoMigrate(database); err != nil {
		logger.Fatal("数据库迁移失败", zap.Error(err))
	}
	logger.Info("数据库自动迁移完成")

	// 启动时检查数据库健康状态
	sqlDB, _ := database.DB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Fatal("数据库健康检查失败", zap.Error(err))
	}
	logger.Info("数据库连接正常")

	// 初始化Redis缓存（可选，失败不阻塞启动）
	var redisCache *cache.RedisClient
	if redisClient, redisErr := cache.InitRedis(&cfg.Redis, logger); redisErr != nil {
		logger.Warn("Redis初始化失败，将使用无缓存模式运行", zap.Error(redisErr))
	} else {
		redisCache = redisClient
		defer redisCache.Close()
	}

	// ==================== 依赖注入：创建Repository实现 ====================
	programRepo := repository_impl.NewProgramRepository(database)
	milestoneRepo := repository_impl.NewMilestoneRepository(database)
	riskRepo := repository_impl.NewRiskRepository(database)
	productRepo := repository_impl.NewProductRepository(database)
	productLineRepo := repository_impl.NewProductLineRepository(database)
	versionRepo := repository_impl.NewProductVersionRepository(database)
	roadmapRepo := repository_impl.NewProductRoadmapRepository(database)
	planRepo := repository_impl.NewProductPlanRepository(database)

	_ = redisCache // Redis可用于后续扩展（如缓存热点数据）

	// ==================== 创建应用服务 ====================
	programSvc := application.NewProgramService(programRepo, milestoneRepo, riskRepo, productRepo, logger)
	productSvc := application.NewProductService(
		productRepo, productLineRepo, versionRepo,
		roadmapRepo, planRepo, programRepo, logger,
	)

	// ==================== 创建Handler实例 ====================
	programHandler := handler.NewProgramHandler(programSvc, logger)
	productHandler := handler.NewProductHandler(productSvc, logger)
	productLineHandler := handler.NewProductLineHandler(productSvc, logger)
	versionHandler := handler.NewVersionHandler(productSvc, logger)
	roadmapHandler := handler.NewRoadmapHandler(productSvc, logger)
	planHandler := handler.NewPlanHandler(productSvc, logger)

	// ==================== 设置Gin引擎和注册路由 =====================
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api.RegisterRoutes(r, programHandler, productHandler, productLineHandler,
		versionHandler, roadmapHandler, planHandler)

	// ==================== 创建HTTP服务器并启动 ====================
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("项目集与产品服务启动",
			zap.String("addr", addr),
			zap.Int("port", cfg.Server.Port),
			zap.String("database", cfg.Database.DBName),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("收到退出信号，正在关闭服务...", zap.String("signal", sig.String()))

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("强制关闭服务", zap.Error(err))
	}

	// 关闭数据库连接
	if sqlDBErr := sqlDB.Close(); sqlDBErr != nil {
		logger.Warn("关闭数据库连接失败", zap.Error(sqlDBErr))
	}

	logger.Info("项目集与产品服务已安全停止")
}
