package handler

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"leap-one/service-gateway/internal/config"
	"leap-one/service-gateway/internal/interfaces/api/dto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HealthHandler 健康检查处理器
// 检查网关自身状态及各下游服务的连通性
type HealthHandler struct {
	cfg        *config.Config
	logger     *zap.Logger
	httpClient *http.Client
}

// NewHealthHandler 创建健康检查处理器实例
func NewHealthHandler(cfg *config.Config, logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		cfg:    cfg,
		logger: logger,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// ServiceHealth 单个服务的健康状态
type ServiceHealth struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Status    string `json:"status"` // "healthy", "unhealthy", "unknown"
	LatencyMs int64  `json:"latencyMs,omitempty"`
	Error     string `json:"error,omitempty"`
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string          `json:"status"`    // "healthy", "degraded", "unhealthy"
	Timestamp string          `json:"timestamp"` // 检查时间
	Gateway   ServiceHealth   `json:"gateway"`   // 网关自身状态
	Services  []ServiceHealth `json:"services"`  // 下游服务状态列表
}

// Check 健康检查端点处理
func (h *HealthHandler) Check(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Gateway: ServiceHealth{
			Name:   "leap-one-gateway",
			Status: "healthy",
		},
	}

	var wg sync.WaitGroup
	serviceNames := []string{
		"user-org", "portfolio", "project", "task",
		"requirement", "quality", "devops", "document",
		"kanban", "bi", "ai", "notification", "search", "config",
	}

	services := make([]ServiceHealth, len(serviceNames))

	for i, name := range serviceNames {
		wg.Add(1)
		go func(idx int, svcName string) {
			defer wg.Done()
			services[idx] = h.checkService(svcName)
		}(i, name)
	}
	wg.Wait()

	response.Services = services

	// 综合判断整体健康状态
	unhealthyCount := 0
	for _, svc := range response.Services {
		if svc.Status == "unhealthy" {
			unhealthyCount++
		}
	}
	if unhealthyCount > 0 {
		response.Status = "degraded"
	}
	if unhealthyCount > len(serviceNames)/2 {
		response.Status = "unhealthy"
	}

	c.JSON(http.StatusOK, dto.Success(response))
}

// checkService 检查单个下游服务的连通性
func (h *HealthHandler) checkService(serviceName string) ServiceHealth {
	svcURL := h.cfg.GetServiceURL(serviceName)
	health := ServiceHealth{
		Name:   serviceName,
		URL:    svcURL,
		Status: "unknown",
	}

	if svcURL == "" {
		health.Status = "unhealthy"
		health.Error = "服务未配置"
		return health
	}

	start := time.Now()
	targetURL := fmt.Sprintf("%s/healthz", svcURL)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, targetURL, nil)
	if err != nil {
		health.Status = "unhealthy"
		health.Error = err.Error()
		return health
	}

	resp, err := h.httpClient.Do(req)
	latency := time.Since(start).Milliseconds()
	health.LatencyMs = latency

	if err != nil {
		// 区分网络不可达和连接拒绝
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			health.Status = "unhealthy"
			health.Error = "连接超时"
		} else {
			health.Status = "unhealthy"
			health.Error = fmt.Sprintf("连接失败: %v", err)
		}
		return health
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		health.Status = "healthy"
	} else {
		health.Status = "unhealthy"
		health.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}

	return health
}

// Liveness 存活探针（K8s livenessProbe）
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Readiness 就绪探针（K8s readinessProbe）
func (h *HealthHandler) Readiness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
