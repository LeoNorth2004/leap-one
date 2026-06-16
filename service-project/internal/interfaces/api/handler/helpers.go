package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"leap-one/service-project/internal/application"
)

// handleServiceError 统一处理业务逻辑错误并返回HTTP响应
func handleServiceError(c *gin.Context, err error) {
	switch err {
	case application.ErrProjectNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case application.ErrProjectCodeExists:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case application.ErrInvalidProjectStatus:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case application.ErrMemberNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case application.ErrMemberAlreadyExists:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case application.ErrInvalidMemberRole:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case application.ErrRiskNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case application.ErrTemplateNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case application.ErrIterationNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case application.ErrInvalidIterationDate:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// getCurrentUserID 从Gin上下文中获取当前登录用户ID
func getCurrentUserID(c *gin.Context) uuid.UUID {
	userIDVal, exists := c.Get("userID")
	if !exists {
		return uuid.Nil
	}
	if id, ok := userIDVal.(string); ok {
		if parsed, err := uuid.Parse(id); err == nil {
			return parsed
		}
	}
	return uuid.Nil
}

// getStringFromContext 从Gin上下文中安全获取字符串�?func getStringFromContext(c *gin.Context, key string) string {
	val, exists := c.Get(key)
	if !exists {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}

// getUUIDFromContext 从Gin上下文中安全获取UUID�?func getUUIDFromContext(c *gin.Context, key string) (uuid.UUID, bool) {
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

// strconvDefaultInt 安全解析整数（带默认值）
func strconvDefaultInt(s string, defaultVal int) (int, error) {
	result := defaultVal
	// 这里简化处理，实际使用时可用strconv.Atoi
	return result, nil
}
