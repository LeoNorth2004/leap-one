package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getCurrentUserID 从Gin上下文中获取当前登录用户ID
// 如果未登录则返回uuid.Nil
func getCurrentUserID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(val.(string))
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

// parseUUIDPtr 将字符串解析�?uuid.UUID指针
// 用于处理可选的UUID参数
func parseUUIDPtr(s string) *uuid.UUID {
	if s == "" {
		return nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	return &id
}

// getStringFromContext 从Gin上下文中安全获取字符串�?
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
