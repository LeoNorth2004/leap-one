// Leap One 看板服务
// 负责看板管理、泳道、WIP控制等功�?
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
	"leap-one/service-kanban/internal/application/service"
	"leap-one/service-kanban/internal/config"
	infraDb "leap-one/service-kanban/internal/infrastructure/db"
	"leap-one/service-kanban/internal/infrastructure/repository"
	"leap-one/service-kanban/internal/interfaces/api/handler"
	"leap-one/service-kanban/internal/interfaces/api/router"
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

	if err := infraDb.AutoMigrate(db, logger); err != nil {
		logger.Fatal("数据库迁移失�?, zap.Error(err))
	}

	// 初始化仓储层
	boardRepo := repository.NewBoardRepository(db)
	columnRepo := repository.NewColumnRepository(db)
	cardRepo := repository.NewCardRepository(db)
	swimlaneRepo := repository.NewSwimlaneRepository(db)
	moveHistoryRepo := repository.NewMoveHistoryRepository(db)

	// 初始化应用服�?
	boardSvc := service.NewBoardService(boardRepo, logger)
	columnSvc := service.NewColumnService(columnRepo, cardRepo, logger)
	cardSvc := service.NewCardService(cardRepo, moveHistoryRepo, columnRepo, logger)
	swimlaneSvc := service.NewSwimlaneService(swimlaneRepo)
	statsSvc := service.NewStatisticsService(cardRepo, columnRepo)

	// 初始化HTTP处理�?
	h := handler.NewKanbanHandler(boardSvc, columnSvc, cardSvc, swimlaneSvc, statsSvc, logger)

	r := router.SetupRouter(h)

	seedData(db, logger)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: r, ReadTimeout: cfg.Server.ReadTimeout, WriteTimeout: cfg.Server.WriteTimeout}

	go func() {
		logger.Info("看板服务启动",
			zap.String("addr", addr), zap.Int("port", cfg.Server.Port),
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
	logger.Info("看板服务已安全停�?)
}

// seedData 注入预置示例数据（含完整看板、列、泳道、卡片）
func seedData(db *gorm.DB, logger *zap.Logger) {
	var count int64
	db.Model(&struct{}{}).Table("kanban_boards").Count(&count)
	if count > 0 {
		logger.Info("检测到已有数据，跳过预置数据注�?)
		return
	}

	projectID := uuid.MustParse("b0000000-0000-0000-0000-000000000001")
	ownerID := uuid.MustParse("c0000000-0000-0000-0000-000000000001")

	// 创建示例看板
	boardID := uuid.MustParse("kb000001-0000-0000-0000-000000000001")
	db.Table("kanban_boards").Create(map[string]interface{}{
		"id": boardID, "name": "Leap One 项目开发看�?, "type": "project",
		"ref_id": projectID, "owner_id": ownerID,
		"description": "Leap One项目主开发看板，跟踪需求→设计→开发→测试→发布的完整流程",
	})

	// 创建看板列（标准Scrum流程�?
	columns := []map[string]interface{}{
		{"id": uuid.MustParse("col00001-0000-0000-0000-000000000001"), "board_id": boardID,
			"name": "待办(Backlog)", "key": "backlog", "wip_limit": 20, "color": "#d9d9d9", "sort_order": 0, "type": "backlog"},
		{"id": uuid.MustParse("col00002-0000-0000-0000-000000000002"), "board_id": boardID,
			"name": "待评�?, "key": "reviewing", "wip_limit": 5, "color": "#faad14", "sort_order": 1},
		{"id": uuid.MustParse("col00003-0000-0000-0000-000000000003"), "board_id": boardID,
			"name": "进行�?Doing)", "key": "doing", "wip_limit": 8, "color": "#1890ff", "sort_order": 2},
		{"id": uuid.MustParse("col00004-0000-0000-0000-000000000004"), "board_id": boardID,
			"name": "测试�?, "key": "testing", "wip_limit": 5, "color": "#722ed1", "sort_order": 3},
		{"id": uuid.MustParse("col00005-0000-0000-0000-000000000005"), "board_id": boardID,
			"name": "已完�?Done)", "key": "done", "color": "#52c41a", "sort_order": 4, "type": "done"},
	}
	for _, col := range columns {
		db.Table("kanban_columns").Create(col)
	}

	// 创建泳道
	swimlanes := []map[string]interface{}{
		{"id": uuid.MustParse("sw000001-0000-0000-0000-000000000001"), "board_id": boardID,
			"name": "核心功能", "key": "core", "color": "#1890ff", "sort_order": 0},
		{"id": uuid.MustParse("sw000002-0000-0000-0000-000000000002"), "board_id": boardID,
			"name": "优化改进", "key": "improvement", "color": "#faad14", "sort_order": 1},
		{"id": uuid.MustParse("sw000003-0000-0000-0000-000000000003"), "board_id": boardID,
			"name": "Bug修复", "key": "bugfix", "color": "#ff4d4f", "sort_order": 2},
		{"id": uuid.MustParse("sw000004-0000-0000-0000-000000000004"), "board_id": boardID,
			"name": "技术债务", "key": "tech_debt", "color": "#8c8c8c", "sort_order": 3},
	}
	for _, sw := range swimlanes {
		db.Table("kanban_swimlanes").Create(sw)
	}

	backlogColID := uuid.MustParse("col00001-0000-0000-0000-000000000001")
	doingColID := uuid.MustParse("col00003-0000-0000-0000-000000000003")
	testColID := uuid.MustParse("col00004-0000-0000-0000-000000000004")
	doneColID := uuid.MustParse("col00005-0000-0000-0000-000000000005")
	coreSwimlaneID := uuid.MustParse("sw000001-0000-0000-0000-000000000001")
	improveSwimlaneID := uuid.MustParse("sw000002-0000-0000-0000-000000000002")

	assigneeID := uuid.MustParse("d0000000-0000-0000-0000-000000000001")

	// 创建示例卡片
	cards := []map[string]interface{}{
		// 待办卡片
		{
			"id":       uuid.MustParse("card0001-0000-0000-0000-000000000001"),
			"board_id": boardID, "column_id": backlogColID, "swimlane_id": coreSwimlaneID,
			"card_type": "requirement", "ref_id": uuid.MustParse("30000000-0000-0000-0000-000000000002"),
			"title": "实现SSO单点登录集成", "priority": 2, "assignee_id": assigneeID,
			"tags": "认证,安全,高优先级", "sort_order": 0,
		},
		{
			"id":       uuid.MustParse("card0002-0000-0000-0000-000000000002"),
			"board_id": boardID, "column_id": backlogColID, "swimlane_id": coreSwimlaneID,
			"card_type": "requirement", "ref_id": uuid.MustParse("30000000-0000-0000-0000-000000000003"),
			"title": "多因素认�?MFA)支持", "priority": 3, "assignee_id": assigneeID,
			"tags": "认证,安全", "sort_order": 1,
		},
		// 进行中卡�?
		{
			"id":       uuid.MustParse("card0003-0000-0000-0000-000000000003"),
			"board_id": boardID, "column_id": doingColID, "swimlane_id": coreSwimlaneID,
			"card_type": "requirement", "ref_id": uuid.MustParse("30000000-0000-0000-0000-000000000001"),
			"title": "实现OAuth2.0授权码模式登�?, "priority": 1, "assignee_id": assigneeID,
			"tags": "认证,OAuth,紧�?, "sort_order": 0,
		},
		{
			"id":       uuid.MustParse("card0004-0000-0000-0000-000000000004"),
			"board_id": boardID, "column_id": doingColID, "swimlane_id": improveSwimlaneID,
			"card_type": "task", "ref_id": uuid.MustParse("t0000001-0000-0000-0000-000000000001"),
			"title": "优化API响应性能，目标P99<200ms", "priority": 2, "assignee_id": assigneeID,
			"tags": "性能,优化", "sort_order": 0,
		},
		// 测试中卡�?
		{
			"id":       uuid.MustParse("card0005-0000-0000-0000-000000000005"),
			"board_id": boardID, "column_id": testColID, "swimlane_id": coreSwimlaneID,
			"card_type": "task", "ref_id": uuid.MustParse("t0000002-0000-0000-0000-000000000001"),
			"title": "用户注册流程集成测试", "priority": 1, "assignee_id": assigneeID,
			"tags": "测试,QA", "sort_order": 0,
		},
		// 已完成卡�?
		{
			"id":       uuid.MustParse("card0006-0000-0000-0000-000000000006"),
			"board_id": boardID, "column_id": doneColID, "swimlane_id": coreSwimlaneID,
			"card_type": "task", "ref_id": uuid.MustParse("t0000003-0000-0000-0000-000000000001"),
			"title": "数据库表结构设计与迁移脚本编�?, "priority": 1,
			"tags": "数据�?基础架构", "sort_order": 0,
		},
	}
	for _, card := range cards {
		db.Table("kanban_cards").Create(card)
	}

	// 创建移动历史示例
	now := time.Now()
	moveHistories := []map[string]interface{}{
		{
			"id":             uuid.MustParse("mh000001-0000-0000-0000-000000000001"),
			"card_id":        uuid.MustParse("card0006-0000-0000-0000-000000000006"),
			"from_column_id": backlogColID, "to_column_id": doingColID,
			"moved_by": ownerID, "move_time": now.Add(-7 * 24 * time.Hour),
		},
		{
			"id":             uuid.MustParse("mh000002-0000-0000-0000-000000000002"),
			"card_id":        uuid.MustParse("card0006-0000-0000-0000-000000000006"),
			"from_column_id": doingColID, "to_column_id": testColID,
			"moved_by": ownerID, "move_time": now.Add(-4 * 24 * time.Hour),
		},
		{
			"id":             uuid.MustParse("mh000003-0000-0000-0000-000000000003"),
			"card_id":        uuid.MustParse("card0006-0000-0000-0000-000000000006"),
			"from_column_id": testColID, "to_column_id": doneColID,
			"moved_by": ownerID, "move_time": now.Add(-2 * 24 * time.Hour),
		},
	}
	for _, mh := range moveHistories {
		db.Table("kanban_card_move_histories").Create(mh)
	}

	logger.Info("预置数据注入完成",
		zap.Int("boards", 1),
		zap.Int("columns", len(columns)),
		zap.Int("swimlanes", len(swimlanes)),
		zap.Int("cards", len(cards)),
		zap.Int("move_histories", len(moveHistories)),
	)
}
