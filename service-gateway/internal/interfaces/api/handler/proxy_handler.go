package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"leap-one/service-gateway/internal/config"
	"leap-one/service-gateway/internal/interfaces/api/dto"
	"leap-one/service-gateway/pkg/errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ProxyHandler 反向代理处理器
// 使用httputil.ReverseProxy实现请求转发到下游微服务
type ProxyHandler struct {
	cfg    *config.Config
	logger *zap.Logger
}

// NewProxyHandler 创建反向代理处理器实例
func NewProxyHandler(cfg *config.Config, logger *zap.Logger) *ProxyHandler {
	return &ProxyHandler{
		cfg:    cfg,
		logger: logger,
	}
}

// Proxy 创建指定目标服务的代理处理器
func (h *ProxyHandler) Proxy(serviceName string) gin.HandlerFunc {
	targetURL := h.cfg.GetServiceURL(serviceName)
	if targetURL == "" {
		return func(c *gin.Context) {
			c.JSON(http.StatusServiceUnavailable, dto.Fail(
				errors.ErrServiceUnavailable.Code,
				fmt.Sprintf("服务 %s 未配置", serviceName),
			))
		}
	}

	parsed, err := url.Parse(targetURL)
	if err != nil {
		h.logger.Error("解析服务地址失败",
			zap.String("service", serviceName),
			zap.String("url", targetURL),
			zap.Error(err),
		)
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, dto.Fail(errors.ErrInternalServer.Code, "服务配置错误"))
		}
	}

	return func(c *gin.Context) {
		h.logger.Debug("代理请求",
			zap.String("service", serviceName),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)

		reverseProxy := httputil.NewSingleHostReverseProxy(parsed)

		// 自定义错误处理：上游不可用时返回友好错误信息
		originalDirector := reverseProxy.Director
		reverseProxy.Director = func(req *http.Request) {
			originalDirector(req)

			// 转发时附加用户信息Header（从JWT中间件注入的上下文中获取）
			if userID, exists := c.Get("userId"); exists {
				req.Header.Set("X-User-ID", fmt.Sprintf("%v", userID))
			}
			if username, exists := c.Get("username"); exists {
				req.Header.Set("X-User-Username", fmt.Sprintf("%v", username))
			}
			if roles, exists := c.Get("roles"); exists {
				if roleSlice, ok := roles.([]string); ok {
					req.Header.Set("X-User-Roles", strings.Join(roleSlice, ","))
				}
			}

			// 标记网关转发的请求
			req.Header.Set("X-Forwarded-By", "leap-one-gateway")
			req.Header.Set("X-Real-IP", c.ClientIP())

			// 保留原始请求头中的关键信息
			if host := c.GetHeader("X-Forwarded-Host"); host != "" {
				req.Header.Set("X-Forwarded-Host", host)
			} else {
				req.Header.Set("X-Forwarded-Host", req.Host)
			}
			if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
				req.Header.Set("X-Forwarded-Proto", proto)
			}
		}

		reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			h.logger.Error("反向代理错误",
				zap.String("service", serviceName),
				zap.String("url", r.URL.String()),
				zap.Error(err),
			)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			errorResp := dto.Fail(errors.ErrBadGateway.Code,
				fmt.Sprintf("上游服务 [%s] 暂时不可用，请稍后重试", serviceName))
			jsonBytes, _ := json.Marshal(errorResp)
			w.Write(jsonBytes)
		}

		reverseProxy.Transport = &http.Transport{
			ResponseHeaderTimeout: 30 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
		}

		reverseProxy.ServeHTTP(c.Writer, c.Request)
	}
}
