// Leap One 需求管理服�?
// 负责需求全生命周期管理(Epic→Feature→Story)
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"leap-one/service-requirement/internal/application/service"
	"leap-one/service-requirement/internal/config"
	infraDb "leap-one/service-requirement/internal/infrastructure/db"
	"leap-one/service-requirement/internal/infrastructure/repository"
	"leap-one/service-requirement/internal/interfaces/api/handler"
	"leap-one/service-requirement/internal/interfaces/api/router"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.Load("")
	if err != nil {
		logger.Fatal("加载配置失败", zap.Error(err))
	}

	db, err := infraDb.NewDatabase(&cfg.Database, logger)
	if err != nil {
		logger.Fatal("数据库初始化失败", zap.Error(err))
	}

	// 自动迁移表结�?
	if err := infraDb.AutoMigrate(db, logger); err != nil {
		logger.Fatal("数据库迁移失�?, zap.Error(err))
	}

	// 初始化仓储层
	reqRepo := repository.NewRequirementRepository(db)
	relationRepo := repository.NewRequirementRelationRepository(db)
	reviewRepo := repository.NewRequirementReviewRepository(db)
	changeLogRepo := repository.NewRequirementChangeLogRepository(db)

	// 初始化应用服务层
	reqService := service.NewRequirementService(reqRepo, relationRepo, changeLogRepo, logger)
	reviewService := service.NewReviewService(reviewRepo, logger)
	changeService := service.NewChangeLogService(changeLogRepo, logger)
	relationSvc := service.NewRelationService(relationRepo, logger)

	// 初始化HTTP处理�?
	h := handler.NewRequirementHandler(reqService, reviewService, changeService, relationSvc, logger)

	// 配置路由
	r := router.SetupRouter(h)

	// 注入预置数据
	seedData(db, logger)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("需求管理服务启�?,
			zap.String("addr", addr),
			zap.Int("port", cfg.Server.Port),
			zap.String("database", cfg.Database.DBName),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失�?, zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("收到退出信号，正在关闭服务...", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("强制关闭服务", zap.Error(err))
	}

	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}

	logger.Info("需求管理服务已安全停止")
}

// seedData 注入预置示例数据（仅用于演示和开发环境）
func seedData(db *gorm.DB, logger *zap.Logger) {
	var count int64
	db.Model(&struct{}{}).Table("requirements").Count(&count)
	if count > 0 {
		logger.Info("检测到已有数据，跳过预置数据注�?)
		return
	}

	productID := uuid.MustParse("a0000000-0000-0000-0000-000000000001")
	projectID := uuid.MustParse("b0000000-0000-0000-0000-000000000001")
	ownerID := uuid.MustParse("c0000000-0000-0000-0000-000000000001")

	// 预置Epic需�?
	epic1ID := uuid.MustParse("10000000-0000-0000-0000-000000000001")
	epic1 := map[string]interface{}{
		"id":              epic1ID,
		"title":           "用户管理系统升级",
		"code":            "REQ-001",
		"description":     "对现有用户管理系统进行全面升级，提升用户体验和系统安全�?,
		"type":            "epic",
		"level":           1,
		"product_id":      productID,
		"project_id":      projectID,
		"status":          "developing",
		"priority":        1,
		"category":        "业务需�?,
		"owner_id":        ownerID,
		"stage":           "dev",
		"release_version": "v2.0",
	}
	db.Model(&map[string]interface{}{}).Table("requirements").Create(epic1)

	// 预置Feature需�?
	feature1ID := uuid.MustParse("20000000-0000-0000-0000-000000000001")
	feature1 := map[string]interface{}{
		"id":              feature1ID,
		"title":           "用户认证模块重构",
		"code":            "REQ-002",
		"description":     "重构用户认证模块，支持OAuth2.0、SSO单点登录等多种认证方�?,
		"type":            "feature",
		"level":           2,
		"parent_id":       epic1ID,
		"product_id":      productID,
		"project_id":      projectID,
		"status":          "developing",
		"priority":        2,
		"category":        "研发需�?,
		"owner_id":        ownerID,
		"stage":           "dev",
		"release_version": "v2.0",
	}
	db.Model(&map[string]interface{}{}).Table("requirements").Create(feature1)

	// 预置Story需�?
	stories := []map[string]interface{}{
		{
			"id":              uuid.MustParse("30000000-0000-0000-0000-000000000001"),
			"title":           "实现OAuth2.0授权码模式登�?,
			"code":            "REQ-003",
			"description":     "实现基于OAuth2.0授权码模式的第三方登录功能，支持微信、支付宝等平�?,
			"type":            "story",
			"level":           3,
			"parent_id":       feature1ID,
			"product_id":      productID,
			"project_id":      projectID,
			"status":          "testing",
			"priority":        1,
			"story_points":    float64Ptr(5),
			"estimated_hours": float64Ptr(16),
			"category":        "用户需�?,
			"owner_id":        ownerID,
			"stage":           "test",
			"release_version": "v2.0",
		},
		{
			"id":              uuid.MustParse("30000000-0000-0000-0000-000000000002"),
			"title":           "实现SSO单点登录集成",
			"code":            "REQ-004",
			"description":     "集成企业级SSO单点登录系统，支持CAS、SAML协议",
			"type":            "story",
			"level":           3,
			"parent_id":       feature1ID,
			"product_id":      productID,
			"project_id":      projectID,
			"status":          "planning",
			"priority":        2,
			"story_points":    float64Ptr(8),
			"estimated_hours": float64Ptr(24),
			"category":        "用户需�?,
			"owner_id":        ownerID,
			"stage":           "requirement",
			"release_version": "v2.0",
		},
		{
			"id":              uuid.MustParse("30000000-0000-0000-0000-000000000003"),
			"title":           "多因素认�?MFA)支持",
			"code":            "REQ-005",
			"description":     "增加TOTP验证码、短信验证码等多因素认证方式",
			"type":            "story",
			"level":           3,
			"parent_id":       feature1ID,
			"product_id":      productID,
			"project_id":      projectID,
			"status":          "draft",
			"priority":        3,
			"story_points":    float64Ptr(3),
			"estimated_hours": float64Ptr(12),
			"category":        "安全需�?,
			"owner_id":        ownerID,
			"stage":           "requirement",
			"release_version": "v2.1",
		},
	}
	for _, s := range stories {
		db.Model(&map[string]interface{}{}).Table("requirements").Create(s)
	}

	// 第二个Epic
	epic2ID := uuid.MustParse("10000000-0000-0000-0000-000000000002")
	epic2 := map[string]interface{}{
		"id":              epic2ID,
		"title":           "数据分析看板建设",
		"code":            "REQ-006",
		"description":     "构建企业级数据分析看板，提供多维度数据可视化能力",
		"type":            "epic",
		"level":           1,
		"product_id":      productID,
		"project_id":      projectID,
		"status":          "planning",
		"priority":        2,
		"category":        "业务需�?,
		"owner_id":        ownerID,
		"stage":           "requirement",
		"release_version": "v2.5",
	}
	db.Model(&map[string]interface{}{}).Table("requirements").Create(epic2)

	logger.Info("预置数据注入完成",
		zap.Int("epics", 2),
		zap.Int("features", 1),
		zap.Int("stories", 3),
	)
}

func float64Ptr(v float64) *float64 { return &v }
