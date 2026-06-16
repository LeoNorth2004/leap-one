package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-user-org/internal/domain/entity"
	"leap-one/service-user-org/internal/domain/repository"
	"leap-one/service-user-org/internal/infrastructure/cache"
	"leap-one/service-user-org/internal/infrastructure/repository_impl"
	"leap-one/service-user-org/internal/interfaces/api/dto"
)

// AuthHandler 认证Handler（登录、登出、注册等）
type AuthHandler struct {
	userRepo   repository.UserRepository
	roleRepo   repository.RoleRepository
	permRepo   repository.PermissionRepository
	cache      *cache.RedisClient
	logger     *zap.Logger
	jwtSecret  string
	jwtIssuer  string
	expireTime time.Duration
}

// NewAuthHandler 创建认证Handler实例
func NewAuthHandler(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permRepo repository.PermissionRepository,
	redisCache *cache.RedisClient,
	logger *zap.Logger,
	jwtSecret string,
	jwtIssuer string,
	expireTime time.Duration,
) *AuthHandler {
	return &AuthHandler{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		permRepo:   permRepo,
		cache:      redisCache,
		logger:     logger,
		jwtSecret:  jwtSecret,
		jwtIssuer:  jwtIssuer,
		expireTime: expireTime,
	}
}

// Register 用户注册（POST /api/v1/auth/register）
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.LoginRequest // 复用LoginRequest结构体，包含username和password
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 校验用户名唯一性
	existingUser, _ := h.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 密码bcrypt加密
	hashedPassword, err := repository_impl.HashPassword(req.Password)
	if err != nil {
		h.logger.Error("密码加密失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败，请稍后重试"})
		return
	}

	user := &entity.User{
		Username: req.Username,
		Password: hashedPassword,
		Status:   1, // 默认启用
	}

	if err := h.userRepo.Create(ctx, user); err != nil {
		h.logger.Error("创建用户失败", zap.Error(err), zap.String("username", req.Username))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败，请稍后重试"})
		return
	}

	h.logger.Info("用户注册成功",
		zap.String("user_id", user.ID.String()),
		zap.String("username", req.Username),
	)

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"user_id": user.ID.String(),
	})
}

// Login 用户登录（POST /api/v1/auth/login）
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 根据用户名查询用户
	user, err := h.userRepo.GetByUsername(ctx, req.Username)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if err := repository_impl.ValidatePassword(req.Password, user.Password); err != nil {
		h.logger.Warn("登录密码验证失败",
			zap.String("username", req.Username),
			zap.String("user_id", user.ID.String()),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 检查用户状态是否正常
	if user.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "账号已被禁用，请联系管理员"})
		return
	}

	// 获取用户的角色编码列表
	roleCodes := h.getUserRoleCodes(ctx, user.ID)

	// 生成JWT令牌
	token, expiresAt, err := h.generateToken(user.ID, user.Username, roleCodes)
	if err != nil {
		h.logger.Error("生成JWT令牌失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败，请稍后重试"})
		return
	}

	// 更新最后登录时间和IP
	if loginErr := h.userRepo.UpdateLastLogin(ctx, user.ID, c.ClientIP()); loginErr != nil {
		h.logger.Warn("更新最后登录时间失败", zap.Error(loginErr))
	}

	h.logger.Info("用户登录成功",
		zap.String("user_id", user.ID.String()),
		zap.String("username", req.Username),
		zap.String("client_ip", c.ClientIP()),
	)

	c.JSON(http.StatusOK, dto.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		UserInfo:  h.buildUserInfo(user, roleCodes),
	})
}

// Logout 用户登出（POST /api/v1/auth/logout）
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := splitBearerToken(authHeader)
		if len(parts) == 2 && parts[0] == "Bearer" {
			// 将Token加入黑名单（TTL设为JWT剩余过期时间）
			if h.cache != nil {
				if err := h.cache.AddToBlacklist(c.Request.Context(), parts[1], h.expireTime); err != nil {
					h.logger.Warn("将Token加入黑名单失败", zap.Error(err))
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// GetProfile 获取当前用户信息（GET /api/v1/auth/profile）
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	ctx := c.Request.Context()
	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	roleCodes := h.getUserRoleCodes(ctx, user.ID)

	c.JSON(http.StatusOK, gin.H{
		"user": h.buildUserInfo(user, roleCodes),
	})
}

// UpdateProfile 修改个人资料（PUT /api/v1/auth/profile）
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 更新允许修改的字段
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.RealName != nil {
		user.RealName = *req.RealName
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}

	if err := h.userRepo.Update(ctx, user); err != nil {
		h.logger.Error("更新用户资料失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败，请稍后重试"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "资料更新成功"})
}

// ChangePassword 修改密码（PUT /api/v1/auth/password）
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	ctx := c.Request.Context()
	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 验证旧密码
	if err := repository_impl.ValidatePassword(req.OldPassword, user.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "旧密码不正确"})
		return
	}

	// 加密新密码
	hashedNewPwd, err := repository_impl.HashPassword(req.NewPassword)
	if err != nil {
		h.logger.Error("新密码加密失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码修改失败"})
		return
	}

	if err := h.userRepo.UpdatePassword(ctx, userID, hashedNewPwd); err != nil {
		h.logger.Error("更新密码失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码修改失败"})
		return
	}

	h.logger.Info("用户修改密码成功", zap.String("user_id", userID.String()))

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// RefreshToken 刷新令牌（POST /api/v1/auth/refresh）
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "令牌刷新功能由网关统一处理",
	})
}

// generateToken 生成JWT令牌
func (h *AuthHandler) generateToken(userID uuid.UUID, username string, roleCodes []string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(h.expireTime)

	claims := &dto.JWTClaims{
		UserID:   userID.String(),
		Username: username,
		Roles:    joinStrings(roleCodes, ","),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    h.jwtIssuer,
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, nil
}

// getUserRoleCodes 获取用户的角色编码列表
func (h *AuthHandler) getUserRoleCodes(ctx context.Context, userID uuid.UUID) []string {
	roleCodes, _ := h.permRepo.GetUserPermissions(ctx, userID)
	if len(roleCodes) > 0 {
		// 从权限编码中提取角色信息（这里简化处理，实际可从user_roles表直接查）
		return roleCodes[:min(len(roleCodes), 10)] // 限制返回数量
	}
	return []string{}
}

// buildUserInfo 构建用户基本信息DTO（不含密码）
func (h *AuthHandler) buildUserInfo(user *entity.User, roleCodes []string) dto.UserInfo {
	info := dto.UserInfo{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		RealName: user.RealName,
		Avatar:   user.Avatar,
		Status:   user.Status,
		Roles:    roleCodes,
	}
	if user.DepartmentID != nil {
		info.DepartmentID = user.DepartmentID.String()
	}
	return info
}
