package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims JWT令牌声明结构
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Roles    string `json:"roles"` // 逗号分隔的角色编码列表
	jwt.RegisteredClaims
}

// AuthMiddleware JWT认证中间件
// 从请求头Authorization中提取并验证JWT令牌
func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "未提供认证令牌",
				"code":  "UNAUTHORIZED",
			})
			return
		}

		// 解析Bearer Token格式：Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "认证格式错误，请使用 Bearer <token>",
				"code":  "INVALID_TOKEN_FORMAT",
			})
			return
		}

		tokenString := parts[1]

		// 解析和验证JWT令牌
		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "令牌无效或已过期",
				"code":  "INVALID_TOKEN",
			})
			return
		}

		// 将用户信息注入上下文，供后续Handler使用
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件
// 如果有令牌则解析用户信息，没有则放行（用于公开接口）
func OptionalAuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]
			claims := &JWTClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})

			if err == nil && token.Valid {
				c.Set("userID", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("roles", claims.Roles)
			}
		}

		c.Next()
	}
}

// RoleBasedAccessControl 基于角色的访问控制中间件
// requiredRoles: 允许访问此路由的角色编码列表（满足任一即可）
func RoleBasedAccessControl(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesValue, exists := c.Get("roles")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "无权限访问",
				"code":  "FORBIDDEN",
			})
			return
		}

		userRoles := rolesValue.(string)
		if userRoles == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "无权限访问",
				"code":  "FORBIDDEN",
			})
			return
		}

		// 检查用户是否拥有所需角色之一
		userRoleList := strings.Split(userRoles, ",")
		for _, userRole := range userRoleList {
			for _, reqRole := range requiredRoles {
				if strings.TrimSpace(userRole) == reqRole {
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "权限不足，需要以下角色之一: " + strings.Join(requiredRoles, ", "),
			"code":  "INSUFFICIENT_PERMISSIONS",
		})
	}
}
