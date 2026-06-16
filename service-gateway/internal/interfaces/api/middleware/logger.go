package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 请求日志中间件
// 记录请求方法、路径、状态码、耗时、客户端IP等关键信息
// 使用结构化日志格式输出
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		clientIP := c.ClientIP()

		// 读取并恢复请求体（避免影响后续处理）
		var body []byte
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			body, _ = io.ReadAll(io.LimitReader(c.Request.Body, 4096))
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		// 处理请求
		c.Next()

		// 计算耗时和状态码
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		errorMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 构建结构化日志
		attrs := []slog.Attr{
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", statusCode),
			slog.Duration("latency", latency),
			slog.String("ip", clientIP),
			slog.String("user_agent", c.GetHeader("User-Agent")),
		}

		if query != "" {
			attrs = append(attrs, slog.String("query", query))
		}
		if errorMsg != "" {
			attrs = append(attrs, slog.String("error", errorMsg))
		}
		if statusCode >= 400 {
			// 仅在错误时记录请求体（脱敏处理）
			if len(body) > 0 {
				bodyStr := string(body)
				if len(bodyStr) > 512 {
					bodyStr = bodyStr[:512] + "...(truncated)"
				}
				attrs = append(attrs, slog.String("body", bodyStr))
			}
		}

		msg := fmt.Sprintf("[API] %s %s -> %d (%v)", method, path, statusCode, latency)
		slog.LogAttrs(c.Request.Context(), logLevel(statusCode), msg, attrs...)
	}
}

// logLevel 根据HTTP状态码返回对应的日志级别
func logLevel(code int) slog.Level {
	switch {
	case code >= 500:
		return slog.LevelError
	case code >= 400:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
