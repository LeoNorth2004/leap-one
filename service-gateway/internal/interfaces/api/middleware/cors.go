package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"leap-one/service-gateway/internal/config"

	"github.com/gin-gonic/gin"
)

// CORS 跨域资源共享中间件
// 支持配置化的允许来源、方法、头部列表
func CORS(cfg *config.CORSConfig) gin.HandlerFunc {
	allowOriginMap := make(map[string]bool)
	for _, origin := range cfg.AllowOrigins {
		allowOriginMap[origin] = true
	}

	allowMethods := strings.Join(cfg.AllowMethods, ", ")
	allowHeaders := strings.Join(cfg.AllowHeaders, ", ")
	exposeHeaders := strings.Join(cfg.ExposeHeaders, ", ")

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// 判断是否允许该来源
		if allowOriginMap["*"] || allowOriginMap[origin] {
			if allowOriginMap["*"] {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", allowMethods)
		c.Writer.Header().Set("Access-Control-Allow-Headers", allowHeaders)
		c.Writer.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
		c.Writer.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))

		if cfg.AllowCredentials && !allowOriginMap["*"] {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// 预检请求直接返回204
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
