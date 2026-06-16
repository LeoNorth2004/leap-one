package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"leap-one/service-gateway/internal/config"
	"leap-one/service-gateway/internal/interfaces/api/dto"
	"leap-one/service-gateway/pkg/errors"
	"leap-one/service-gateway/pkg/jwtutil"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	cfg    *config.Config
	jwtGen *jwtutil.Generator
	logger *zap.Logger
	client *http.Client // 用于转发请求到user-org-service的HTTP客户端
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(cfg *config.Config, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		cfg:    cfg,
		jwtGen: jwtutil.NewGenerator(&cfg.JWT),
		logger: logger,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// LoginHandler 登录处理
// 转发登录请求到user-org-service进行验证，成功后由网关签发JWT Token对
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(errors.ErrBadRequest.Code, "参数校验失败: "+err.Error()))
		return
	}

	// 构造转发到user-org服务的请求
	targetURL := h.cfg.GetServiceURL("user-org") + "/api/v1/auth/login"
	bodyBytes, _ := json.Marshal(req)

	h.logger.Debug("转发登录请求", zap.String("target", targetURL))

	resp, err := h.client.Post(targetURL, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		h.logger.Error("用户服务不可达", zap.Error(err), zap.String("target", targetURL))
		c.JSON(http.StatusBadGateway, dto.Fail(
			errors.ErrBadGateway.Code,
			"认证服务暂时不可用，请稍后重试",
		))
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// 上游返回非2xx，直接透传错误
	if resp.StatusCode != http.StatusOK {
		c.Data(resp.StatusCode, "application/json", respBody)
		return
	}

	// 解析上游响应获取用户信息（期望包含userId, username, roles）
	var userResp struct {
		Code int `json:"code"`
		Data struct {
			UserID   string   `json:"userId"`
			Username string   `json:"username"`
			Roles    []string `json:"roles"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &userResp); err != nil || userResp.Code != 200 {
		// 如果无法解析用户信息但上游返回成功，使用默认角色签发Token
		h.logger.Warn("解析用户信息失败，使用默认配置签发Token", zap.Error(err))
		tokenPair, err := h.jwtGen.GenerateToken(req.Username, req.Username, []string{"user"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.Fail(errors.ErrInternalServer.Code, "Token签发失败"))
			return
		}
		c.JSON(http.StatusOK, dto.Success(dto.LoginResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresIn:    tokenPair.ExpiresIn,
		}))
		return
	}

	// 网关签发JWT Token对
	tokenPair, err := h.jwtGen.GenerateToken(
		userResp.Data.UserID,
		userResp.Data.Username,
		userResp.Data.Roles,
	)
	if err != nil {
		h.logger.Error("Token签发失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.Fail(errors.ErrInternalServer.Code, "Token签发失败"))
		return
	}

	h.logger.Info("用户登录成功",
		zap.String("username", userResp.Data.Username),
		zap.String("userId", userResp.Data.UserID),
	)

	c.JSON(http.StatusOK, dto.Success(dto.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}))
}

// RegisterHandler 注册处理
// 转发注册请求到user-org-service
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(errors.ErrBadRequest.Code, "参数校验失败: "+err.Error()))
		return
	}

	targetURL := h.cfg.GetServiceURL("user-org") + "/api/v1/auth/register"
	bodyBytes, _ := json.Marshal(req)

	resp, err := h.client.Post(targetURL, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		h.logger.Error("注册服务不可达", zap.Error(err), zap.String("target", targetURL))
		c.JSON(http.StatusBadGateway, dto.Fail(errors.ErrBadGateway.Code, "注册服务暂时不可用"))
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", respBody)
}

// RefreshTokenHandler 刷新Token处理
// 使用有效的RefreshToken换取新的Access Token和Refresh Token对
func (h *AuthHandler) RefreshTokenHandler(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(errors.ErrBadRequest.Code, "参数校验失败: "+err.Error()))
		return
	}

	// 验证并刷新Token
	tokenPair, err := h.jwtGen.RefreshToken(req.RefreshToken)
	if err != nil {
		h.logger.Warn("RefreshToken刷新失败", zap.Error(err))
		c.JSON(http.StatusUnauthorized, dto.Fail(
			errors.ErrTokenExpired.Code,
			"RefreshToken无效或已过期，请重新登录",
		))
		return
	}

	c.JSON(http.StatusOK, dto.Success(dto.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}))
}

// forwardRequest 通用请求转发方法
func (h *AuthHandler) forwardRequest(targetURL, method string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, targetURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("构造请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gateway-Forwarded-By", "leap-one-gateway")

	return h.client.Do(req)
}
