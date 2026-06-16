package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"leap-one/service-gateway/internal/config"
	"leap-one/service-gateway/pkg/errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter 基于IP的限流器（使用内存存储）
// 生产环境可替换为基于Redis的分布式实现
type IPRateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*rate.Limiter
	rate     rate.Limit
	burst    int
}

// NewIPRateLimiter 创建基于IP的令牌桶限流器
func NewIPRateLimiter(r rate.Limit, burst int) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    burst,
	}
}

// GetLimiter 获取或创建指定IP的限流器实例
func (l *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, exists := l.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(l.rate, l.burst)
		l.limiters[ip] = limiter
	}

	return limiter
}

// CleanupLimiter 清理指定IP的限流器（用于连接断开后释放内存）
func (l *IPRateLimiter) CleanupLimiter(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.limiters, ip)
}

// RateLimit 限流中间件
// 基于令牌桶算法实现，支持按IP维度限流
// 默认配置：100请求/分钟/IP，超出返回429状态码和Retry-After头
func RateLimit(cfg *config.RateLimitConfig) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) { c.Next() }
	}

	// 计算每秒速率：requestsPerMinute / 60
	rps := float64(cfg.RequestsPerMinute) / 60.0
	if cfg.RPS > 0 {
		rps = cfg.RPS
	}
	burst := cfg.Burst
	if burst <= 0 {
		burst = cfg.RequestsPerMinute
	}

	limiter := NewIPRateLimiter(rate.Limit(rps), burst)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// 尝试从令牌桶获取令牌
		if !limiter.GetLimiter(clientIP).Allow() {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    errors.ErrTooManyRequests.Code,
				"message": errors.ErrTooManyRequests.Message,
				"hint":    fmt.Sprintf("每分钟最多允许 %d 次请求，请稍后重试", cfg.RequestsPerMinute),
			})
			return
		}

		c.Next()
	}
}

// UserRateLimit 基于用户ID的限流中间件
// 需要在JWT中间件之后使用，从上下文中获取用户ID进行限流
func UserRateLimit(cfg *config.RateLimitConfig) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) { c.Next() }
	}

	rps := float64(cfg.RequestsPerMinute) / 60.0
	if cfg.RPS > 0 {
		rps = cfg.RPS * 2 // 用户维度的限制可以适当放宽
	}
	burst := cfg.Burst * 2

	limiter := NewIPRateLimiter(rate.Limit(rps), burst)

	return func(c *gin.Context) {
		userID, exists := c.Get("userId")
		if !exists {
			// 未认证用户回退到IP限流
			clientIP := c.ClientIP()
			if !limiter.GetLimiter(clientIP).Allow() {
				c.Header("Retry-After", "60")
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"code":    errors.ErrTooManyRequests.Code,
					"message": errors.ErrTooManyRequests.Message,
				})
				return
			}
			c.Next()
			return
		}

		key := fmt.Sprintf("user:%v", userID)
		if !limiter.GetLimiter(key).Allow() {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    errors.ErrTooManyRequests.Code,
				"message": errors.ErrTooManyRequests.Message,
			})
			return
		}

		c.Next()
	}
}

// 定期清理过期限流器的后台任务（可选）
func (l *IPRateLimiter) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			l.mu.Lock()
			// 简单策略：清空所有（生产环境应按最后访问时间清理）
			// 此处保留基本实现，避免内存泄漏
			if len(l.limiters) > 10000 {
				l.limiters = make(map[string]*rate.Limiter)
			}
			l.mu.Unlock()
		}
	}()
}
