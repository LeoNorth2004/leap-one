package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// splitBearerToken 解析Authorization头中的Bearer Token
func splitBearerToken(authHeader string) []string {
	return strings.SplitN(authHeader, " ", 2)
}

// joinStrings 使用分隔符连接字符串切片
func joinStrings(strs []string, sep string) string {
	return strings.Join(strs, sep)
}

// getStringFromContext 从Gin上下文中安全获取字符串值
func getStringFromContext(c *gin.Context, key string) string {
	val, exists := c.Get(key)
	if !exists {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}

// getUUIDFromContext 从Gin上下文中安全获取UUID值
func getUUIDFromContext(c *gin.Context, key string) (uuid.UUID, bool) {
	val, exists := c.Get(key)
	if !exists {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(val.(string))
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}
