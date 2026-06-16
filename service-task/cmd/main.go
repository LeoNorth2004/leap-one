// Leap One 任务与工单服务
// 负责任务管理、工单/事项管理、工作流、SLA配置等功能
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"leap-one/service-task/internal/application"
	"leap-one/service-task/internal/config"
	"leap-one/service-task/internal/domain/entity"
	"leap-one/service-task/internal/infrastructure/cache"
	"leap-one/service-task/internal/infrastructure/db"
	"leap-one/service-task/internal/infrastructure/repository_impl"
	"leap-one/service-task/internal/interfaces/api"
	"leap-one/service-task/internal/interfaces/api/handler"

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

	sqlDB, _ := database.DB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Fatal("数据库健康检查失败", zap.Error(err))
	}
	logger.Info("数据库连接正常")

	var redisCache *cache.RedisClient
	if redisClient, redisErr := cache.InitRedis(&cfg.Redis, logger); redisErr != nil {
		logger.Warn("Redis初始化失败，将使用无缓存模式运行", zap.Error(redisErr))
	} else {
		redisCache = redisClient
		defer redisCache.Close()
	}

	taskRepo := repository_impl.NewTaskRepository(database)
	assignmentRepo := repository_impl.NewTaskAssignmentRepository(database)
	commentRepo := repository_impl.NewTaskCommentRepository(database)
	taskAttachmentRepo := repository_impl.NewTaskAttachmentRepository(database)
	linkRepo := repository_impl.NewTaskLinkRepository(database)
	worklogRepo := repository_impl.NewTaskWorklogRepository(database)

	issueRepo := repository_impl.NewIssueRepository(database)
	issueCommentRepo := repository_impl.NewIssueCommentRepository(database)
	issueAttachmentRepo := repository_impl.NewIssueAttachmentRepository(database)
	templateRepo := repository_impl.NewIssueTemplateRepository(database)
	workflowRepo := repository_impl.NewIssueWorkflowRepository(database)
	slaConfigRepo := repository_impl.NewIssueSLAConfigRepository(database)

	initPreseededData(database, logger)

	taskSvc := application.NewTaskService(taskRepo, assignmentRepo, commentRepo,
		taskAttachmentRepo, linkRepo, worklogRepo, logger)
	issueSvc := application.NewIssueService(issueRepo, issueCommentRepo, issueAttachmentRepo,
		templateRepo, workflowRepo, slaConfigRepo, logger)
	templateSvc := application.NewTemplateService(templateRepo)
	workflowSvc := application.NewWorkflowService(workflowRepo)
	slaSvc := application.NewSLAConfigService(slaConfigRepo)

	_ = redisCache

	taskHandler := handler.NewTaskHandler(taskSvc, logger)
	issueHandler := handler.NewIssueHandler(issueSvc, logger)
	templateHandler := handler.NewTemplateHandler(templateSvc, logger)
	workflowHandler := handler.NewWorkflowHandler(workflowSvc, slaSvc, logger)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api.RegisterRoutes(r, taskHandler, issueHandler, templateHandler, workflowHandler)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("任务与工单服务启动",
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

	if sqlDBErr := sqlDB.Close(); sqlDBErr != nil {
		logger.Warn("关闭数据库连接失败", zap.Error(sqlDBErr))
	}

	logger.Info("任务与工单服务已安全停止")
}

func initPreseededData(dbConn *gorm.DB, logger *zap.Logger) {
	ctx := context.Background()

	predefinedWorkflows := []struct {
		Name          string
		Type          string
		InitialStatus string
		Description   string
		Transitions   []struct {
			FromStatus string
			ToStatus   string
			Name       string
			SortOrder  int
		}
	}{
		{
			Name: "Bug处理工作流", Type: "bug", InitialStatus: "new",
			Description: "标准Bug修复流程：新建→处理中→已解决→已关闭",
			Transitions: []struct {
				FromStatus string
				ToStatus   string
				Name       string
				SortOrder  int
			}{
				{"new", "in_progress", "开始处理", 1},
				{"in_progress", "resolved", "标记为已解决", 2},
				{"in_progress", "waiting", "挂起等待", 3},
				{"waiting", "in_progress", "恢复处理", 4},
				{"resolved", "closed", "关闭工单", 5},
				{"resolved", "in_progress", "重新打开", 6},
				{"new", "cancelled", "取消", 7},
				{"in_progress", "cancelled", "取消", 8},
				{"waiting", "cancelled", "取消", 9},
			},
		},
		{
			Name: "需求处理工作流", Type: "feature", InitialStatus: "new",
			Description: "功能需求流程：新建→处理中→等待→已解决→已关闭",
			Transitions: []struct {
				FromStatus string
				ToStatus   string
				Name       string
				SortOrder  int
			}{
				{"new", "in_progress", "开始开发", 1},
				{"in_progress", "waiting", "等待确认", 2},
				{"waiting", "in_progress", "继续开发", 3},
				{"in_progress", "resolved", "已完成", 4},
				{"resolved", "closed", "发布关闭", 5},
				{"resolved", "in_progress", "返回修改", 6},
				{"new", "cancelled", "取消需求", 7},
			},
		},
	}

	for _, pw := range predefinedWorkflows {
		var count int64
		dbConn.WithContext(ctx).Model(&entity.IssueWorkflow{}).Where("name = ?", pw.Name).Count(&count)
		if count == 0 {
			wf := &entity.IssueWorkflow{
				Name:          pw.Name,
				Type:          pw.Type,
				InitialStatus: pw.InitialStatus,
				Description:   pw.Description,
			}
			dbConn.WithContext(ctx).Create(wf)

			for _, t := range pw.Transitions {
				trans := &entity.IssueWorkflowTransition{
					WorkflowID: wf.ID,
					FromStatus: t.FromStatus,
					ToStatus:   t.ToStatus,
					Name:       t.Name,
					SortOrder:  t.SortOrder,
				}
				dbConn.WithContext(ctx).Create(trans)
			}

			logger.Debug("预置工作流已创建",
				zap.String("name", pw.Name),
				zap.String("id", wf.ID.String()),
			)
		}
	}

	predefinedTemplates := []struct {
		Name   string
		Type   string
		Fields string
	}{
		{Name: "Bug报告模板", Type: "bug", Fields: `{"fields":[{"key":"title","label":"标题","type":"text","required":true},{"key":"description","label":"复现步骤","type":"textarea","required":true},{"key":"environment","label":"环境信息","type":"textarea"},{"key":"expected","label":"期望结果","type":"textarea"},{"key":"actual","label":"实际结果","type":"textarea"},{"key":"severity","label":"严重程度","type":"select","options":["1-致命","2-严重","3-一般","4-轻微"]},{"key":"priority","label":"优先级","type":"select","options":["1-紧急","2-高","3-中","4-低","5-最低"]}]}`},
		{Name: "功能需求模板", Type: "feature", Fields: `{"fields":[{"key":"title","label":"需求标题","type":"text","required":true},{"key":"description","label":"需求描述","type":"textarea","required":true},{"key":"acceptance_criteria","label":"验收标准","type":"textarea"},{"key":"user_story","label":"用户故事","type":"textarea"},{"key":"priority","label":"优先级","type":"select"}]}`},
		{Name: "服务请求模板", Type: "request", Fields: `{"fields":[{"key":"title","label":"请求标题","type":"text","required":true},{"key":"description","label":"详细描述","type":"textarea"},{"key":"request_type","label":"请求类型","type":"select","options":["咨询","权限申请","数据查询","系统配置"]},{"key":"urgency","label":"紧急程度","type":"select"}]}`},
		{Name: "事件报告模板", Type: "incident", Fields: `{"fields":[{"key":"title","label":"事件标题","type":"text","required":true},{"key":"description","label":"事件描述","type":"textarea","required":true},{"key":"impact_scope","label":"影响范围","type":"select","options":["个人","团队","部门","全公司"]},{"key":"occurrence_time","label":"发生时间","type":"datetime"}]}`},
	}

	for _, pt := range predefinedTemplates {
		var count int64
		dbConn.WithContext(ctx).Model(&entity.IssueTemplate{}).Where("name = ?", pt.Name).Count(&count)
		if count == 0 {
			tmpl := &entity.IssueTemplate{
				Name:     pt.Name,
				Type:     pt.Type,
				Fields:   pt.Fields,
				IsSystem: true,
			}
			dbConn.WithContext(ctx).Create(tmpl)
			logger.Debug("预置模板已创建",
				zap.String("name", pt.Name),
				zap.String("type", pt.Type),
			)
		}
	}

	slaMatrix := map[string]map[int][2]int{
		"bug":      {1: {15, 240}, 2: {30, 480}, 3: {60, 1440}, 4: {120, 2880}, 5: {240, 4320}},
		"feature":  {1: {30, 2880}, 2: {60, 5760}, 3: {120, 10080}, 4: {240, 20160}, 5: {480, 30240}},
		"request":  {1: {10, 120}, 2: {20, 360}, 3: {30, 720}, 4: {60, 1440}, 5: {120, 2880}},
		"incident": {1: {5, 60}, 2: {15, 180}, 3: {30, 480}, 4: {60, 1440}, 5: {120, 2880}},
	}

	issueTypes := []string{"bug", "feature", "request", "incident"}
	priorities := []int{1, 2, 3, 4, 5}

	for _, itype := range issueTypes {
		for _, prio := range priorities {
			if sla, ok := slaMatrix[itype]; ok {
				if slaVal, exists := sla[prio]; exists {
					var count int64
					dbConn.WithContext(ctx).Model(&entity.IssueSLAConfig{}).
						Where("type = ? AND priority = ?", itype, prio).Count(&count)
					if count == 0 {
						slaCfg := &entity.IssueSLAConfig{
							Type:              itype,
							Priority:          prio,
							ResponseSLA:       slaVal[0],
							ResolveSLA:        slaVal[1],
							BusinessHoursOnly: (itype == "request"),
						}
						dbConn.WithContext(ctx).Create(slaCfg)
					}
				}
			}
		}
	}

	logger.Info("预置数据初始化完成")
}
