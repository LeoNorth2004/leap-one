package api

import (
	"leap-one/service-user-org/internal/interfaces/api/handler"
	"leap-one/service-user-org/internal/interfaces/api/middleware"

	"github.com/gin-gonic/gin"
)

// RouterConfig 路由配置参数
type RouterConfig struct {
	JWTSecret string
}

// RegisterRoutes 注册所有API路由
func RegisterRoutes(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	roleHandler *handler.RoleHandler,
	deptHandler *handler.DepartmentHandler,
	groupHandler *handler.UserGroupHandler,
	config RouterConfig,
) {
	// 全局中间件
	r.Use(gin.Recovery())

	// ==================== 健康检查（公开）====================
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "leap-one-user-org",
		})
	})

	// ==================== 公开路由（无需认证）====================
	public := r.Group("/api/v1")
	{
		// 认证相关
		authGroup := public.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}
	}

	// ==================== 需要认证的路由（JWT中间件）====================
	// 使用AuthMiddleware保护以下所有路由
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(config.JWTSecret))
	{
		// ---- 认证模块 ----
		authProtected := protected.Group("/auth")
		{
			authProtected.POST("/logout", authHandler.Logout)
			authProtected.GET("/profile", authHandler.GetProfile)
			authProtected.PUT("/profile", authHandler.UpdateProfile)
			authProtected.PUT("/password", authHandler.ChangePassword)
			authProtected.POST("/refresh", authHandler.RefreshToken)
		}

		// ---- 用户管理模块 ----
		users := protected.Group("/users")
		{
			users.GET("/me", userHandler.GetCurrentUser)
			users.PUT("/me/password", userHandler.ChangePassword)
			users.GET("/search", userHandler.SearchUsers)
			users.GET("", userHandler.ListUsers)
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.PUT("/:id/status", userHandler.ToggleStatus)
			users.PUT("/:id/reset-password", userHandler.ResetPassword)
		}

		// ---- 角色权限管理模块 ----
		roles := protected.Group("/roles")
		{
			roles.GET("", roleHandler.ListRoles)
			roles.POST("", roleHandler.CreateRole)
			roles.GET("/:id", roleHandler.GetRole)
			roles.PUT("/:id", roleHandler.UpdateRole)
			roles.DELETE("/:id", roleHandler.DeleteRole)
			roles.POST("/:id/permissions", roleHandler.AssignPermissions)
			roles.GET("/:id/users", roleHandler.GetRoleUsers)
		}
		// 权限列表
		protected.GET("/permissions", roleHandler.GetPermissions)

		// ---- 部门管理模块 ----
		depts := protected.Group("/departments")
		{
			depts.GET("", deptHandler.ListDepartments)
			depts.GET("/tree", deptHandler.GetDepartmentTree)
			depts.POST("", deptHandler.CreateDepartment)
			depts.GET("/:id", deptHandler.GetDepartment)
			depts.PUT("/:id", deptHandler.UpdateDepartment)
			depts.DELETE("/:id", deptHandler.DeleteDepartment)
			depts.PUT("/:id/move", deptHandler.MoveDepartment)
			depts.GET("/:id/members", deptHandler.GetDepartmentMembers)
		}

		// ---- 用户组管理模块 ----
		groups := protected.Group("/groups")
		{
			groups.GET("", groupHandler.ListGroups)
			groups.POST("", groupHandler.CreateGroup)
			groups.GET("/:id", groupHandler.GetGroup)
			groups.PUT("/:id", groupHandler.UpdateGroup)
			groups.DELETE("/:id", groupHandler.DeleteGroup)
			groups.POST("/:id/members", groupHandler.AddMembers)
			groups.DELETE("/:id/members", groupHandler.RemoveMembers)
			groups.GET("/:id/members", groupHandler.ListGroupMembers)
		}
	}
}
