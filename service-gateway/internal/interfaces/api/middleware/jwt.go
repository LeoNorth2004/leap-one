package middleware

import (
	"net/http"
	"strings"

	"leap-one/service-gateway/internal/config"
	"leap-one/service-gateway/pkg/errors"
	"leap-one/service-gateway/pkg/jwtutil"

	"github.com/gin-gonic/gin"
)

// 上下文键常量
const (
	UserIDKey   = "userId"
	UsernameKey = "username"
	RolesKey    = "roles"
)

// JWTAuth JWT认证中间件
// 从Authorization头提取Bearer Token，验证有效性后注入用户信息到上下文
// 支持SkipPaths白名单配置，白名单路径跳过认证
func JWTAuth(cfg *config.JWTConfig) gin.HandlerFunc {
	generator := jwtutil.NewGenerator(cfg)

	// 构建路径白名单的快速查找map
	skipPaths := make(map[string]bool, len(cfg.SkipPaths))
	for _, path := range cfg.SkipPaths {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 白名单路径直接放行
		if skipPaths[path] {
			c.Next()
			return
		}

		// 提取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    errors.ErrUnauthorized.Code,
				"message": "缺少认证令牌，请在请求头中提供 Authorization: Bearer <token>",
			})
			return
		}

		// 解析Bearer Token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    errors.ErrUnauthorized.Code,
				"message": "认证格式错误，请使用 Bearer <token> 格式",
			})
			return
		}

		tokenString := parts[1]

		// 验证Token
		claims, err := generator.ParseToken(tokenString)
		if err != nil {
			// 区分过期和其他错误
			if strings.Contains(err.Error(), "token_expired") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":    errors.ErrTokenExpired.Code,
					"message": errors.ErrTokenExpired.Message,
					"hint":    "请使用RefreshToken刷新或重新登录",
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    errors.ErrTokenInvalid.Code,
				"message": errors.ErrTokenInvalid.Message,
			})
			return
		}

		// 将用户信息注入上下文
		c.Set(UserIDKey, claims.UserID)
		c.Set(UsernameKey, claims.Username)
		c.Set(RolesKey, claims.Roles)

		c.Next()
	}
}
