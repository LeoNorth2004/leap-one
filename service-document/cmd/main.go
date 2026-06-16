// Leap One 文档与知识库服务
// 负责文档管理、附件存储、知识库等功�?
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
	"leap-one/service-document/internal/application/service"
	"leap-one/service-document/internal/config"
	infraDb "leap-one/service-document/internal/infrastructure/db"
	"leap-one/service-document/internal/infrastructure/repository"
	"leap-one/service-document/internal/interfaces/api/handler"
	"leap-one/service-document/internal/interfaces/api/router"
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
	docRepo := repository.NewDocumentRepository(db)
	versionRepo := repository.NewDocumentVersionRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	kbRepo := repository.NewKnowledgeBaseRepository(db)
	favRepo := repository.NewFavoriteRepository(db)
	tagRepo := repository.NewTagRepository(db)

	// 初始化应用服�?
	docSvc := service.NewDocumentService(docRepo, versionRepo, favRepo, logger)
	versionSvc := service.NewVersionService(versionRepo)
	commentSvc := service.NewCommentService(commentRepo)
	categorySvc := service.NewCategoryService(categoryRepo)
	kbSvc := service.NewKnowledgeBaseService(kbRepo)
	attachSvc := service.NewAttachmentService(attachmentRepo)
	tagSvc := service.NewTagService(tagRepo)

	// 初始化HTTP处理�?
	h := handler.NewDocumentHandler(docSvc, versionSvc, commentSvc, categorySvc, kbSvc, attachSvc, tagSvc, logger)

	r := router.SetupRouter(h)

	seedData(db, logger)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: r, ReadTimeout: cfg.Server.ReadTimeout, WriteTimeout: cfg.Server.WriteTimeout}

	go func() {
		logger.Info("文档与知识库服务启动",
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
	logger.Info("文档与知识库服务已安全停�?)
}

// seedData 注入预置示例数据
func seedData(db *gorm.DB, logger *zap.Logger) {
	var count int64
	db.Model(&struct{}{}).Table("documents").Count(&count)
	if count > 0 {
		logger.Info("检测到已有数据，跳过预置数据注�?)
		return
	}

	projectID := uuid.MustParse("b0000000-0000-0000-0000-000000000001")
	ownerID := uuid.MustParse("c0000000-0000-0000-0000-000000000001")

	docs := []map[string]interface{}{
		{
			"id":    uuid.MustParse("d0000000-0000-0000-0000-000000000001"),
			"title": "Leap One 项目开发规范指�?, "type": "markdown", "status": "published",
			"project_id": projectID, "owner_id": ownerID, "visibility": "public",
			"version": 3, "tags": "规范,开�?入门",
			"content": "# Leap One 项目开发规范\n\n## 代码规范\n- 遵循Go官方编码规范\n- 使用DDD分层架构\n\n## Git工作流\n- 采用GitFlow分支模型\n- 提交信息遵循Conventional Commits",
		},
		{
			"id":    uuid.MustParse("d0000000-0000-0000-0000-000000000002"),
			"title": "系统架构设计文档", "type": "markdown", "status": "published",
			"project_id": projectID, "owner_id": ownerID, "visibility": "team",
			"version": 2, "tags": "架构,设计,技�?,
			"content": "# 系统架构设计\n\n## 整体架构\n采用微服务架构，包含以下核心服务：\n\n### 服务清单\n1. **用户服务** - 用户认证与权限管理\n2. **项目管理服务** - 项目全生命周期管理\n3. **需求管理服�?* - 需求Epic/Feature/Story管理\n4. **任务管理服务** - 任务分配与跟踪\n5. **看板服务** - 可视化任务看板\n6. **文档服务** - 文档与知识库管理\n\n## 技术栈\n- 后端: Go 1.23+ / Gin / GORM v2\n- 数据�? PostgreSQL\n- 缓存: Redis\n- 消息队列: RabbitMQ",
		},
		{
			"id":    uuid.MustParse("d0000000-0000-0000-0000-000000000003"),
			"title": "API接口设计规范", "type": "markdown", "status": "draft",
			"project_id": projectID, "owner_id": ownerID, "visibility": "private",
			"version": 1, "tags": "API,接口,RESTful",
			"parent_id": uuid.MustParse("d0000000-0000-0000-0000-000000000001"),
			"content":   "# API接口设计规范\n\n## RESTful规范\n- 使用名词复数形式\n- 使用HTTP方法表达操作语义\n- 统一返回格式\n\n## 版本控制\n所有API通过URL路径版本化：`/api/v1/...`",
		},
		{
			"id":    uuid.MustParse("d0000000-0000-0000-0000-000000000004"),
			"title": "需求分析报告模�?, "type": "markdown", "status": "published",
			"project_id": projectID, "owner_id": ownerID, "is_template": true,
			"version": 1, "tags": "模板,需�?报告",
			"content": "# [项目名称] 需求分析报告\n\n## 1. 项目背景\n（描述项目背景和目标）\n\n## 2. 用户角色\n| 角色 | 描述 | 权限 |\n|------|------|------|\n|      |      |      |\n\n## 3. 功能需求\n### 3.1 核心功能\n...\n\n## 4. 非功能需求\n- 性能要求\n- 安全要求\n- 可用性要�?,
		},
	}

	for _, d := range docs {
		db.Model(&map[string]interface{}{}).Table("documents").Create(d)
	}

	// 预置分类
	categories := []map[string]interface{}{
		{"id": uuid.MustParse("cat00001-0000-0000-0000-000000000001"), "name": "技术文�?, "sort_order": 1},
		{"id": uuid.MustParse("cat00001-0000-0000-0000-000000000002"), "name": "产品文档", "sort_order": 2},
		{"id": uuid.MustParse("cat00001-0000-0000-0000-000000000003"), "name": "运维文档", "sort_order": 3},
	}
	for _, cat := range categories {
		db.Model(&map[string]interface{}{}).Table("document_categories").Create(cat)
	}

	// 预置标签
	tags := []map[string]interface{}{
		{"id": uuid.MustParse("tag00001-0000-0000-0000-000000000001"), "name": "技�?, "color": "#1890ff"},
		{"id": uuid.MustParse("tag00001-0000-0000-0000-000000000002"), "name": "产品", "color": "#52c41a"},
		{"id": uuid.MustParse("tag00001-0000-0000-0000-000000000003"), "name": "规范", "color": "#faad14"},
	}
	for _, t := range tags {
		db.Model(&map[string]interface{}{}).Table("document_tags").Create(t)
	}

	// 预置知识�?
	db.Model(&map[string]interface{}{}).Table("knowledge_bases").Create(map[string]interface{}{
		"id":          uuid.MustParse("kb000001-0000-0000-0000-000000000001"),
		"name":        "Leap One 技术知识库",
		"description": "收录Leap One平台相关的技术文档、最佳实践和解决方案",
		"owner_id":    ownerID,
		"is_public":   true,
	})

	logger.Info("预置数据注入完成",
		zap.Int("documents", len(docs)),
		zap.Int("categories", len(categories)),
		zap.Int("tags", len(tags)),
	)
}
