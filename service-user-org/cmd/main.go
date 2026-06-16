// Leap One 用户与组织服务
// 负责用户管理、部门管理、角色权限(RBAC)、用户组等功能
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"leap-one/service-user-org/internal/application"
	"leap-one/service-user-org/internal/config"
	"leap-one/service-user-org/internal/domain/entity"
	"leap-one/service-user-org/internal/infrastructure/cache"
	"leap-one/service-user-org/internal/infrastructure/db"
	"leap-one/service-user-org/internal/infrastructure/repository_impl"
	"leap-one/service-user-org/internal/interfaces/api"
	"leap-one/service-user-org/internal/interfaces/api/handler"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	userRepo := repository_impl.NewUserRepository(database)
	roleRepo := repository_impl.NewRoleRepository(database)
	permRepo := repository_impl.NewPermissionRepository(database)
	deptRepo := repository_impl.NewDepartmentRepository(database)
	groupRepo := repository_impl.NewUserGroupRepository(database)

	// ==================== 预置数据初始化（首次启动自动创建）====================
	initPreseededData(database, logger)

	// ==================== 创建应用服务 ====================
	userSvc := application.NewUserService(userRepo, roleRepo, logger)
	authSvc := application.NewAuthService(userRepo, permRepo, logger)
	_ = userSvc // 应用服务可被Handler使用
	_ = authSvc

	// ==================== 创建Handler实例 ====================
	authHandler := handler.NewAuthHandler(userRepo, roleRepo, permRepo, redisCache, logger,
		cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.ExpireTime)
	userHandler := handler.NewUserHandler(userRepo, roleRepo, permRepo, groupRepo, logger)
	roleHandler := handler.NewRoleHandler(roleRepo, permRepo, logger)
	deptHandler := handler.NewDepartmentHandler(deptRepo, userRepo, logger)
	groupHandler := handler.NewUserGroupHandler(groupRepo, logger)

	// ==================== 设置Gin引擎和注册路由 =====================
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api.RegisterRoutes(r, authHandler, userHandler, roleHandler, deptHandler, groupHandler,
		api.RouterConfig{JWTSecret: cfg.JWT.Secret})

	// ==================== 创建HTTP服务器并启动 ====================
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("用户与组织服务启动",
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

	logger.Info("用户与组织服务已安全停止")
}

// initPreseededData 初始化预置数据（角色、权限、超级管理员账号）
// 仅在数据不存在时插入，支持幂等执行
func initPreseededData(dbConn *gorm.DB, logger *zap.Logger) {
	ctx := context.Background()

	// ---- 1. 创建预置角色 ----
	predefinedRoles := []struct {
		Name        string
		Code        string
		Type        int8
		Description string
	}{
		{"超级管理员", "super_admin", 1, "拥有系统所有权限"},
		{"管理员", "admin", 1, "拥有大部分管理权限"},
		{"项目经理", "pm", 1, "项目管理相关权限"},
		{"产品经理", "po", 1, "产品需求相关权限"},
		{"开发工程师", "dev", 1, "代码开发相关权限"},
		{"测试工程师", "qa", 1, "测试相关权限"},
		{"查看者", "viewer", 1, "只读查看权限"},
	}

	for _, pr := range predefinedRoles {
		existing, _ := findRoleByCode(ctx, dbConn, pr.Code)
		if existing == nil {
			role := &entity.Role{
				Name:        pr.Name,
				Code:        pr.Code,
				Type:        pr.Type,
				Description: pr.Description,
				Status:      1,
			}
			if err := dbConn.Create(role).Error; err != nil {
				logger.Warn("创建预置角色失败", zap.String("code", pr.Code), zap.Error(err))
			} else {
				logger.Debug("预置角色已创建", zap.String("code", pr.Code), zap.String("id", role.ID.String()))
			}
		}
	}

	// ---- 2. 创建预置权限 ----
	predefinedPermissions := []struct {
		Name     string
		Code     string
		Resource string
		Action   string
		Module   string
	}{
		// 用户管理模块
		{"创建用户", "user:create", "user", "create", "用户管理"},
		{"查看用户", "user:read", "user", "read", "用户管理"},
		{"编辑用户", "user:update", "user", "update", "用户管理"},
		{"删除用户", "user:delete", "user", "delete", "用户管理"},
		{"重置密码", "user:reset_password", "user", "reset_password", "用户管理"},
		// 角色权限模块
		{"创建角色", "role:create", "role", "create", "角色管理"},
		{"查看角色", "role:read", "role", "read", "角色管理"},
		{"编辑角色", "role:update", "role", "update", "角色管理"},
		{"删除角色", "role:delete", "role", "delete", "角色管理"},
		{"分配权限", "role:assign_permission", "role", "assign_permission", "角色管理"},
		// 部门管理模块
		{"创建部门", "department:create", "department", "create", "部门管理"},
		{"查看部门", "department:read", "department", "read", "部门管理"},
		{"编辑部门", "department:update", "department", "update", "部门管理"},
		{"删除部门", "department:delete", "department", "delete", "部门管理"},
		// 项目管理模块
		{"创建项目", "project:create", "project", "create", "项目管理"},
		{"查看项目", "project:read", "project", "read", "项目管理"},
		{"编辑项目", "project:update", "project", "update", "项目管理"},
		{"删除项目", "project:delete", "project", "delete", "项目管理"},
		// 需求管理模块
		{"创建需求", "requirement:create", "requirement", "create", "需求管理"},
		{"查看需求", "requirement:read", "requirement", "read", "需求管理"},
		{"编辑需求", "requirement:update", "requirement", "update", "需求管理"},
		{"删除需求", "requirement:delete", "requirement", "delete", "需求管理"},
		// 用户组模块
		{"创建用户组", "group:create", "group", "create", "用户组管理"},
		{"查看用户组", "group:read", "group", "read", "用户组管理"},
		{"编辑用户组", "group:update", "group", "update", "用户组管理"},
		{"删除用户组", "group:delete", "group", "delete", "用户组管理"},
	}

	for _, pp := range predefinedPermissions {
		var count int64
		dbConn.WithContext(ctx).Model(&entity.Permission{}).Where("code = ?", pp.Code).Count(&count)
		if count == 0 {
			perm := &entity.Permission{
				Name:     pp.Name,
				Code:     pp.Code,
				Resource: pp.Resource,
				Action:   pp.Action,
				Module:   pp.Module,
			}
			if err := dbConn.Create(perm).Error; err != nil {
				logger.Warn("创建预置权限失败", zap.String("code", pp.Code), zap.Error(err))
			}
		}
	}

	// ---- 3. 为super_admin角色分配全部权限 ----
	adminRole, _ := findRoleByCode(ctx, dbConn, "super_admin")
	if adminRole != nil {
		var allPermIDs []string
		dbConn.WithContext(ctx).Model(&entity.Permission{}).Pluck("id", &allPermIDs)
		if len(allPermIDs) > 0 {
			// 清除旧关联
			dbConn.Where("role_id = ?", adminRole.ID).Delete(&entity.RolePermission{})
			// 批量插入新关联
			rps := make([]entity.RolePermission, len(allPermIDs))
			for i, pidStr := range allPermIDs {
				pid, parseErr := uuid.Parse(pidStr)
				if parseErr == nil {
					rps[i] = entity.RolePermission{RoleID: adminRole.ID, PermissionID: pid}
				}
			}
			dbConn.Create(&rps)
			logger.Debug("为super_admin角色分配了全部权限", zap.Int("count", len(allPermIDs)))
		}
	}

	// ---- 4. 创建预置超级管理员账号 ----
	existingAdmin, _ := findUserByUsername(ctx, dbConn, "admin")
	if existingAdmin == nil {
		hashedPwd, hashErr := repository_impl.HashPassword("Admin@123456")
		if hashErr != nil {
			logger.Fatal("加密管理员密码失败", zap.Error(hashErr))
		}

		adminUser := &entity.User{
			Username: "admin",
			Password: hashedPwd,
			RealName: "系统管理员",
			Status:   1,
		}

		if err := dbConn.Create(adminUser).Error; err != nil {
			logger.Warn("创建超级管理员账号失败", zap.Error(err))
		} else {
			logger.Info("超级管理员账号已创建",
				zap.String("username", "admin"),
				zap.String("password", "Admin@123456"),
				zap.String("user_id", adminUser.ID.String()),
			)
		}

		// 为管理员分配super_admin角色
		if adminRole != nil {
			userRole := entity.UserRole{UserID: adminUser.ID, RoleID: adminRole.ID}
			dbConn.Create(&userRole)
		}
	}

	logger.Info("预置数据初始化完成")
}

// findRoleByCode 根据编码查找角色（辅助函数）
func findRoleByCode(ctx context.Context, dbConn *gorm.DB, code string) (*entity.Role, error) {
	var role entity.Role
	err := dbConn.WithContext(ctx).Where("code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// findUserByUsername 根据用户名查找用户（辅助函数）
func findUserByUsername(ctx context.Context, dbConn *gorm.DB, username string) (*entity.User, error) {
	var user entity.User
	err := dbConn.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
