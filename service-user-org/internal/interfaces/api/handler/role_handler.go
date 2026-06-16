package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-user-org/internal/domain/entity"
	"leap-one/service-user-org/internal/domain/repository"
	"leap-one/service-user-org/internal/interfaces/api/dto"
)

// RoleHandler 角色权限管理Handler（RBAC核心）
type RoleHandler struct {
	roleRepo repository.RoleRepository
	permRepo repository.PermissionRepository
	logger   *zap.Logger
}

// NewRoleHandler 创建角色管理Handler实例
func NewRoleHandler(
	roleRepo repository.RoleRepository,
	permRepo repository.PermissionRepository,
	logger *zap.Logger,
) *RoleHandler {
	return &RoleHandler{
		roleRepo: roleRepo,
		permRepo: permRepo,
		logger:   logger,
	}
}

// CreateRole 创建角色（POST /api/v1/roles）
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 校验角色编码唯一性
	if existing, _ := h.roleRepo.GetByCode(ctx, req.Code); existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "角色编码已存在"})
		return
	}

	roleType := int8(1) // 默认系统角色
	if req.Type != nil {
		roleType = *req.Type
	}

	role := &entity.Role{
		Name:        req.Name,
		Code:        req.Code,
		Type:        roleType,
		Description: req.Description,
		Status:      1,
	}

	if err := h.roleRepo.Create(ctx, role); err != nil {
		h.logger.Error("创建角色失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建角色失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "角色创建成功",
		"role_id": role.ID.String(),
	})
}

// GetRole 获取角色详情（GET /api/v1/roles/:id）- 含权限列表
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID格式"})
		return
	}

	ctx := c.Request.Context()
	role, err := h.roleRepo.GetByID(ctx, id)
	if err != nil || role == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
		return
	}

	// 获取角色的权限列表
	permissions, _ := h.roleRepo.GetRolePermissions(ctx, id)
	permList := make([]dto.PermissionInfo, len(permissions))
	for i, p := range permissions {
		permList[i] = dto.PermissionInfo{
			ID:       p.ID,
			Name:     p.Name,
			Code:     p.Code,
			Resource: p.Resource,
			Action:   p.Action,
		}
	}

	resp := dto.RoleDetailResponse{
		RoleDetailInfo: dto.RoleDetailInfo{
			ID:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Type:        role.Type,
			Description: role.Description,
			Status:      role.Status,
			CreatedAt:   role.CreatedAt.Format("2006-01-02 15:04:05"),
		},
		Permissions: permList,
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateRole 更新角色（PUT /api/v1/roles/:id）
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID格式"})
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	role, err := h.roleRepo.GetByID(ctx, id)
	if err != nil || role == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
		return
	}

	// 系统预置角色不允许修改名称和编码
	if role.Type == 1 {
		if req.Name != nil && *req.Name != role.Name {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不允许修改系统预置角色的名称"})
			return
		}
	}

	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = *req.Description
	}
	if req.Status != nil {
		role.Status = *req.Status
	}

	if err := h.roleRepo.Update(ctx, role); err != nil {
		h.logger.Error("更新角色失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新角色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "角色更新成功"})
}

// DeleteRole 删除角色（DELETE /api/v1/roles/:id）- 检查是否有关联用户
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID格式"})
		return
	}

	ctx := c.Request.Context()
	role, err := h.roleRepo.GetByID(ctx, id)
	if err != nil || role == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
		return
	}

	// 系统预置角色不允许删除
	if role.Type == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不允许删除系统预置角色"})
		return
	}

	// TODO: 可在此处检查是否有关联用户

	if err := h.roleRepo.Delete(ctx, id); err != nil {
		h.logger.Error("删除角色失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除角色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "角色删除成功"})
}

// ListRoles 分页查询角色列表（GET /api/v1/roles）
func (h *RoleHandler) ListRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	ctx := c.Request.Context()
	roles, total, err := h.roleRepo.List(ctx, page, size)
	if err != nil {
		h.logger.Error("查询角色列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询角色列表失败"})
		return
	}

	list := make([]dto.RoleDetailInfo, len(roles))
	for i, r := range roles {
		list[i] = dto.RoleDetailInfo{
			ID:          r.ID,
			Name:        r.Name,
			Code:        r.Code,
			Type:        r.Type,
			Description: r.Description,
			Status:      r.Status,
			CreatedAt:   r.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, dto.RoleListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// AssignPermissions 为角色分配权限（POST /api/v1/roles/:id/permissions）
func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID格式"})
		return
	}

	var req dto.AssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 验证角色存在
	role, rErr := h.roleRepo.GetByID(ctx, id)
	if rErr != nil || role == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
		return
	}

	// 先清除旧权限再批量添加新权限
	if err := h.roleRepo.AssignPermissions(ctx, id, req.PermissionIDs); err != nil {
		h.logger.Error("分配权限失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "分配权限失败"})
		return
	}

	h.logger.Info("为角色分配权限",
		zap.String("operator", getStringFromContext(c, "username")),
		zap.String("role_code", role.Code),
		zap.Int("permission_count", len(req.PermissionIDs)),
	)

	c.JSON(http.StatusOK, gin.H{
		"role_id":          id.String(),
		"permission_count": len(req.PermissionIDs),
		"message":          "权限分配成功",
	})
}

// GetPermissions 获取所有权限列表（GET /api/v1/permissions）
func (h *RoleHandler) GetPermissions(c *gin.Context) {
	resource := c.Query("resource") // 可按资源类型筛选

	ctx := c.Request.Context()
	permissions, _, err := h.permRepo.List(ctx, 1, 1000, resource)
	if err != nil {
		h.logger.Error("获取权限列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取权限列表失败"})
		return
	}

	list := make([]dto.PermissionInfo, len(permissions))
	for i, p := range permissions {
		list[i] = dto.PermissionInfo{
			ID:       p.ID,
			Name:     p.Name,
			Code:     p.Code,
			Resource: p.Resource,
			Action:   p.Action,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"list":     list,
		"resource": resource,
		"total":    len(list),
	})
}

// GetRoleUsers 获取角色下的用户列表（GET /api/v1/roles/:id/users）
func (h *RoleHandler) GetRoleUsers(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID格式"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	ctx := c.Request.Context()
	users, total, err := h.permRepo.GetRoleUsers(ctx, id, page, size)
	if err != nil {
		h.logger.Error("获取角色用户列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	list := make([]dto.UserInfo, len(users))
	for i, u := range users {
		list[i] = dto.UserInfo{
			ID:       u.ID.String(),
			Username: u.Username,
			RealName: u.RealName,
			Email:    u.Email,
			Avatar:   u.Avatar,
			Status:   u.Status,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
		"page":  page,
		"size":  size,
	})
}
