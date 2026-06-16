// Leap One 项目管理服务
// 负责项目管理、迭代/Sprint管理、里程碑、风险、自定义字段等功能
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"leap-one/service-project/internal/application"
	"leap-one/service-project/internal/config"
	"leap-one/service-project/internal/domain/entity"
	"leap-one/service-project/internal/infrastructure/db"
	"leap-one/service-project/internal/infrastructure/repository_impl"
	"leap-one/service-project/internal/interfaces/api"
	"leap-one/service-project/internal/interfaces/api/handler"

	"github.com/gin-gonic/gin"
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

	database, err := db.InitPostgreSQL(cfg, logger)
	if err != nil {
		logger.Fatal("数据库初始化失败", zap.Error(err))
	}

	if err := db.AutoMigrate(database); err != nil {
		logger.Fatal("数据库迁移失败", zap.Error(err))
	}
	logger.Info("数据库自动迁移完成")

	if err := db.CheckHealth(database); err != nil {
		logger.Fatal("数据库健康检查失败", zap.Error(err))
	}
	logger.Info("数据库连接正常")

	initPreseededData(database, logger)

	projectRepo := repository_impl.NewProjectRepository(database)
	memberRepo := repository_impl.NewProjectMemberRepository(database)
	templateRepo := repository_impl.NewProjectTemplateRepository(database)
	milestoneRepo := repository_impl.NewProjectMilestoneRepository(database)
	riskRepo := repository_impl.NewProjectRiskRepository(database)
	customFieldRepo := repository_impl.NewCustomFieldRepository(database)
	iterRepo := repository_impl.NewIterationRepository(database)

	projectSvc := application.NewProjectService(projectRepo, memberRepo, logger)
	memberSvc := application.NewProjectMemberService(memberRepo, projectRepo, logger)
	milestoneSvc := application.NewMilestoneService(milestoneRepo, projectRepo, logger)
	riskSvc := application.NewRiskService(riskRepo, projectRepo, logger)
	customFieldSvc := application.NewCustomFieldService(customFieldRepo, projectRepo, logger)
	iterSvc := application.NewIterationService(iterRepo, projectRepo, logger)
	templateSvc := application.NewTemplateService(templateRepo, logger)
	statsSvc := application.NewProjectStatisticsService(
		projectRepo, memberRepo, milestoneRepo, riskRepo, iterRepo, logger,
	)

	projectHandler := handler.NewProjectHandler(projectSvc, logger)
	memberHandler := handler.NewMemberHandler(memberSvc, logger)
	milestoneHandler := handler.NewMilestoneHandler(milestoneSvc, logger)
	riskHandler := handler.NewRiskHandler(riskSvc, logger)
	customFieldHandler := handler.NewCustomFieldHandler(customFieldSvc, logger)
	iterationHandler := handler.NewIterationHandler(iterSvc, logger)
	templateHandler := handler.NewTemplateHandler(templateSvc, logger)
	statsHandler := handler.NewStatisticsHandler(statsSvc, logger)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api.RegisterRoutes(r,
		projectHandler,
		memberHandler,
		milestoneHandler,
		riskHandler,
		customFieldHandler,
		iterationHandler,
		templateHandler,
		statsHandler,
		api.RouterConfig{JWTSecret: cfg.JWT.Secret},
	)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("项目管理服务启动",
			zap.String("addr", addr),
			zap.Int("port", cfg.Server.Port),
			zap.String("database", cfg.Database.DBName),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失败", zap.Error(err))
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

	sqlDB, _ := database.DB()
	if sqlDB != nil {
		if sqlErr := sqlDB.Close(); sqlErr != nil {
			logger.Warn("关闭数据库连接失败", zap.Error(sqlErr))
		}
	}

	logger.Info("项目管理服务已安全停止")
}

func initPreseededData(dbConn *gorm.DB, logger *zap.Logger) {
	ctx := context.Background()

	predefinedTemplates := []struct {
		Name        string
		Description string
		Type        string
		Config      string
	}{
		{
			Name:        "敏捷开发模板",
			Description: "适用于Scrum/敏捷开发团队，包含Sprint迭代、每日站会、回顾会议等敏捷实践",
			Type:        "agile",
			Config: `{
				"phases": [
					{"key": "sprint_planning", "name": "Sprint计划", "order": 1},
					{"key": "development", "name": "开发阶段", "order": 2},
					{"key": "testing", "name": "测试验证", "order": 3},
					{"key": "sprint_review", "name": "Sprint评审", "order": 4},
					{"key": "retrospective", "name": "回顾总结", "order": 5}
				],
				"defaultFields": [
					{"key": "story_points", "name": "故事点", "type": "number"},
					{"key": "priority", "name": "优先级", "type": "select", "options": ["P0","P1","P2","P3"]},
					{"key": "epic_link", "name": "关联Epic", "type": "user"}
				],
				"sprintDuration": 14,
				"ceremonies": ["daily_standup", "sprint_planning", "sprint_review", "retrospective"]
			}`,
		},
		{
			Name:        "瀑布模型模板",
			Description: "适用于传统瀑布开发模式，强调需求、设计、开发、测试的阶段性交付",
			Type:        "waterfall",
			Config: `{
				"phases": [
					{"key": "requirements", "name": "需求分析", "order": 1},
					{"key": "design", "name": "系统设计", "order": 2},
					{"key": "implementation", "name": "编码实现", "order": 3},
					{"key": "verification", "name": "验证测试", "order": 4},
					{"key": "deployment", "name": "部署上线", "order": 5},
					{"key": "maintenance", "name": "运维维护", "order": 6}
				],
				"defaultFields": [
					{"key": "phase", "name": "所属阶段", "type": "select"},
					{"key": "completion_rate", "name": "完成率", "type": "number"}
				],
				"gates": [
					{"name": "需求评审", "afterPhase": "requirements"},
					{"name": "设计评审", "afterPhase": "design"},
					{"name": "UAT验收", "afterPhase": "verification"}
				]
			}`,
		},
		{
			Name:        "轻量管理模板",
			Description: "适用于小型团队或轻量级项目管理，简化流程，快速交付",
			Type:        "lightweight",
			Config: `{
				"phases": [
					{"key": "todo", "name": "待办", "order": 1},
					{"key": "doing", "name": "进行中", "order": 2},
					{"key": "done", "name": "已完成", "order": 3}
				],
				"defaultFields": [
					{"key": "due_date", "name": "截止日期", "type": "date"},
					{"key": "assignee", "name": "负责人", "type": "user"}
				],
				"simpleMode": true
			}`,
		},
	}

	for _, pt := range predefinedTemplates {
		var count int64
		dbConn.WithContext(ctx).Model(&entity.ProjectTemplate{}).
			Where("name = ? AND is_system = ?", pt.Name, true).
			Count(&count)
		if count == 0 {
			template := &entity.ProjectTemplate{
				Name:        pt.Name,
				Description: pt.Description,
				Type:        pt.Type,
				Config:      pt.Config,
				IsSystem:    true,
			}
			if err := dbConn.Create(template).Error; err != nil {
				logger.Warn("创建预置模板失败", zap.String("name", pt.Name), zap.Error(err))
			} else {
				logger.Info("预置模板已创建",
					zap.String("name", pt.Name),
					zap.String("id", template.ID.String()),
				)
			}
		}
	}

	logger.Info("预置数据初始化完成")
}
