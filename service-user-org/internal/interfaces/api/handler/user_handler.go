package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-user-org/internal/domain/entity"
	"leap-one/service-user-org/internal/domain/repository"
	"leap-one/service-user-org/internal/infrastructure/repository_impl"
	"leap-one/service-user-org/internal/interfaces/api/dto"
)

// UserHandler 用户管理Handler
type UserHandler struct {
	userRepo  repository.UserRepository
	roleRepo  repository.RoleRepository
	permRepo  repository.PermissionRepository
	groupRepo repository.UserGroupRepository
	logger    *zap.Logger
}

// NewUserHandler 创建用户管理Handler实例
func NewUserHandler(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permRepo repository.PermissionRepository,
	groupRepo repository.UserGroupRepository,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		permRepo:  permRepo,
		groupRepo: groupRepo,
		logger:    logger,
	}
}

// CreateUser 管理员创建用户（POST /api/v1/users）
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 校验用户名唯一性
	if existing, _ := h.userRepo.GetByUsername(ctx, req.Username); existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 密码加密
	hashedPwd, err := repository_impl.HashPassword(req.Password)
	if err != nil {
		h.logger.Error("密码加密失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	user := &entity.User{
		Username:     req.Username,
		Password:     hashedPwd,
		Email:        req.Email,
		Phone:        req.Phone,
		RealName:     req.RealName,
		Avatar:       req.Avatar,
		DepartmentID: req.DepartmentID,
		Status:       1,
	}

	if err := h.userRepo.Create(ctx, user); err != nil {
		h.logger.Error("创建用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	h.logger.Info("管理员创建用户成功",
		zap.String("operator", getStringFromContext(c, "username")),
		zap.String("new_user", req.Username),
	)

	c.JSON(http.StatusCreated, gin.H{
		"message": "用户创建成功",
		"user_id": user.ID.String(),
	})
}

// GetUser 获取用户详情（GET /api/v1/users/:id）
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID格式"})
		return
	}

	ctx := c.Request.Context()
	user, err := h.userRepo.GetByID(ctx, id)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	resp := h.buildUserDetailResponse(ctx, user)
	c.JSON(http.StatusOK, resp)
}

// UpdateUser 更新用户信息（PUT /api/v1/users/:id）
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID格式"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	user, err := h.userRepo.GetByID(ctx, id)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 更新字段
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
	if req.DepartmentID != nil {
		user.DepartmentID = req.DepartmentID
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := h.userRepo.Update(ctx, user); err != nil {
		h.logger.Error("更新用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户更新成功"})
}

// DeleteUser 删除用户（DELETE /api/v1/users/:id）
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID格式"})
		return
	}

	ctx := c.Request.Context()

	// 检查用户是否存在
	user, uErr := h.userRepo.GetByID(ctx, id)
	if uErr != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	if err := h.userRepo.Delete(ctx, id); err != nil {
		h.logger.Error("删除用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	h.logger.Info("删除用户成功",
		zap.String("operator", getStringFromContext(c, "username")),
		zap.String("deleted_user", user.Username),
	)

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// ListUsers 分页查询用户列表（GET /api/v1/users）
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	ctx := c.Request.Context()
	users, total, err := h.userRepo.List(ctx, page, size, keyword)
	if err != nil {
		h.logger.Error("查询用户列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户列表失败"})
		return
	}

	list := make([]dto.UserInfo, len(users))
	for i, user := range users {
		list[i] = dto.UserInfo{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
			Phone:    user.Phone,
			RealName: user.RealName,
			Avatar:   user.Avatar,
			Status:   user.Status,
		}
		if user.DepartmentID != nil {
			list[i].DepartmentID = user.DepartmentID.String()
		}
	}

	c.JSON(http.StatusOK, dto.UserListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// ChangePassword 修改当前用户密码（PUT /api/v1/users/me/password）
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	userID, err := uuid.Parse(userIDVal.(string))
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码修改失败"})
		return
	}

	if err := h.userRepo.UpdatePassword(ctx, userID, hashedNewPwd); err != nil {
		h.logger.Error("更新密码失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码修改失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// GetCurrentUser 获取当前登录用户信息（GET /api/v1/users/me）
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	userID, err := uuid.Parse(userIDVal.(string))
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

	resp := h.buildUserDetailResponse(ctx, user)
	c.JSON(http.StatusOK, resp)
}

// ResetPassword 重置用户密码（PUT /api/v1/users/:id/reset-password）- 管理员操作
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID格式"})
		return
	}

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=8,max=50"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	user, err := h.userRepo.GetByID(ctx, id)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	hashedPwd, hashErr := repository_impl.HashPassword(req.NewPassword)
	if hashErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码重置失败"})
		return
	}

	if err := h.userRepo.UpdatePassword(ctx, id, hashedPwd); err != nil {
		h.logger.Error("重置密码失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码重置失败"})
		return
	}

	h.logger.Info("管理员重置用户密码",
		zap.String("operator", getStringFromContext(c, "username")),
		zap.String("target_user", user.Username),
	)

	c.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
}

// ToggleStatus 启用/禁用用户（PUT /api/v1/users/:id/status）
func (h *UserHandler) ToggleStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID格式"})
		return
	}

	var req struct {
		Status int8 `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	user, err := h.userRepo.GetByID(ctx, id)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	user.Status = req.Status
	if err := h.userRepo.Update(ctx, user); err != nil {
		h.logger.Error("更新用户状态失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}

	action := "禁用"
	if req.Status == 1 {
		action = "启用"
	}
	h.logger.Info(action+"用户",
		zap.String("operator", getStringFromContext(c, "username")),
		zap.String("target_user", user.Username),
	)

	c.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}

// SearchUsers 全局搜索用户（GET /api/v1/users/search）- 供其他服务调用
func (h *UserHandler) SearchUsers(c *gin.Context) {
	keyword := c.Query("keyword")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	ctx := c.Request.Context()
	users, err := h.permRepo.SearchUsers(ctx, keyword, limit)
	if err != nil {
		h.logger.Error("搜索用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索用户失败"})
		return
	}

	list := make([]dto.UserInfo, len(users))
	for i, user := range users {
		list[i] = dto.UserInfo{
			ID:       user.ID.String(),
			Username: user.Username,
			RealName: user.RealName,
			Email:    user.Email,
			Avatar:   user.Avatar,
		}
	}

	c.JSON(http.StatusOK, gin.H{"list": list, "count": len(list)})
}

// buildUserDetailResponse 构建用户详情响应（含角色、组信息）
func (h *UserHandler) buildUserDetailResponse(ctx context.Context, user *entity.User) dto.UserDetailResponse {
	info := dto.UserInfo{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		RealName: user.RealName,
		Avatar:   user.Avatar,
		Status:   user.Status,
	}
	if user.DepartmentID != nil {
		info.DepartmentID = user.DepartmentID.String()
	}

	// 获取用户角色
	roles := make([]dto.RoleInfo, 0)
	// 可通过关联查询获取角色信息

	// 获取用户组
	groups := make([]dto.GroupInfo, 0)
	// 可通过关联查询获取组信息

	return dto.UserDetailResponse{
		UserInfo:  info,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		Roles:     roles,
		Groups:    groups,
	}
}
