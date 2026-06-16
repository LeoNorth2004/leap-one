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

// UserGroupHandler 用户组管理Handler
type UserGroupHandler struct {
	groupRepo repository.UserGroupRepository
	logger    *zap.Logger
}

// NewUserGroupHandler 创建用户组管理Handler实例
func NewUserGroupHandler(
	groupRepo repository.UserGroupRepository,
	logger *zap.Logger,
) *UserGroupHandler {
	return &UserGroupHandler{
		groupRepo: groupRepo,
		logger:    logger,
	}
}

// CreateGroup 创建用户组（POST /api/v1/groups）
func (h *UserGroupHandler) CreateGroup(c *gin.Context) {
	var req dto.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 校验编码唯一性
	if existing, _ := h.groupRepo.GetByCode(ctx, req.Code); existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户组编码已存在"})
		return
	}

	group := &entity.UserGroup{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      1,
		MemberCount: 0,
	}

	if err := h.groupRepo.Create(ctx, group); err != nil {
		h.logger.Error("创建用户组失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户组失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "用户组创建成功",
		"group_id": group.ID.String(),
	})
}

// GetGroup 获取用户组详情（GET /api/v1/groups/:id）
func (h *UserGroupHandler) GetGroup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户组ID格式"})
		return
	}

	ctx := c.Request.Context()
	group, err := h.groupRepo.GetByID(ctx, id)
	if err != nil || group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户组不存在"})
		return
	}

	// 获取组成员列表（第一页）
	members, _, mErr := h.groupRepo.GetMembers(ctx, id, 1, 10)
	if mErr != nil {
		h.logger.Warn("获取组成员失败", zap.Error(mErr))
	}

	memberList := make([]dto.UserInfo, len(members))
	for i, m := range members {
		memberList[i] = dto.UserInfo{
			ID:       m.ID.String(),
			Username: m.Username,
			RealName: m.RealName,
			Email:    m.Email,
			Avatar:   m.Avatar,
		}
	}

	resp := dto.GroupDetailResponse{
		GroupDetailInfo: dto.GroupDetailInfo{
			ID:          group.ID,
			Name:        group.Name,
			Code:        group.Code,
			Description: group.Description,
			MemberCount: group.MemberCount,
			Status:      group.Status,
			CreatedAt:   group.CreatedAt.Format("2006-01-02 15:04:05"),
		},
		Members: memberList,
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateGroup 更新用户组（PUT /api/v1/groups/:id）
func (h *UserGroupHandler) UpdateGroup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户组ID格式"})
		return
	}

	var req dto.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	group, err := h.groupRepo.GetByID(ctx, id)
	if err != nil || group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户组不存在"})
		return
	}

	if req.Name != nil {
		group.Name = *req.Name
	}
	if req.Description != nil {
		group.Description = *req.Description
	}
	if req.Status != nil {
		group.Status = *req.Status
	}

	if err := h.groupRepo.Update(ctx, group); err != nil {
		h.logger.Error("更新用户组失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户组失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户组更新成功"})
}

// DeleteGroup 删除用户组（DELETE /api/v1/groups/:id）
func (h *UserGroupHandler) DeleteGroup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户组ID格式"})
		return
	}

	ctx := c.Request.Context()

	// 检查是否存在
	group, gErr := h.groupRepo.GetByID(ctx, id)
	if gErr != nil || group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户组不存在"})
		return
	}

	if err := h.groupRepo.Delete(ctx, id); err != nil {
		h.logger.Error("删除用户组失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户组失败"})
		return
	}

	h.logger.Info("删除用户组成功",
		zap.String("operator", getStringFromContext(c, "username")),
		zap.String("group_name", group.Name),
	)

	c.JSON(http.StatusOK, gin.H{"message": "用户组删除成功"})
}

// ListGroups 分页查询用户组列表（GET /api/v1/groups）
func (h *UserGroupHandler) ListGroups(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	ctx := c.Request.Context()
	groups, total, err := h.groupRepo.List(ctx, page, size)
	if err != nil {
		h.logger.Error("查询用户组列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户组列表失败"})
		return
	}

	list := make([]dto.GroupDetailInfo, len(groups))
	for i, g := range groups {
		list[i] = dto.GroupDetailInfo{
			ID:          g.ID,
			Name:        g.Name,
			Code:        g.Code,
			Description: g.Description,
			MemberCount: g.MemberCount,
			Status:      g.Status,
			CreatedAt:   g.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, dto.GroupListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// AddMembers 批量添加成员到用户组（POST /api/v1/groups/:id/members）
func (h *UserGroupHandler) AddMembers(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户组ID格式"})
		return
	}

	var req dto.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 验证用户组存在
	group, gErr := h.groupRepo.GetByID(ctx, id)
	if gErr != nil || group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户组不存在"})
		return
	}

	// 批量添加成员（幂等操作：重复添加不报错）
	if err := h.groupRepo.BatchAddMembers(ctx, id, req.UserIDs); err != nil {
		h.logger.Error("批量添加成员失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加成员失败"})
		return
	}

	// 更新成员计数
	if updateErr := h.groupRepo.UpdateMemberCount(ctx, id); updateErr != nil {
		h.logger.Warn("更新成员计数失败", zap.Error(updateErr))
	}

	h.logger.Info("批量添加用户组成员",
		zap.String("operator", getStringFromContext(c, "username")),
		zap.String("group_name", group.Name),
		zap.Int("added_count", len(req.UserIDs)),
	)

	c.JSON(http.StatusOK, gin.H{
		"group_id":    id.String(),
		"added_count": len(req.UserIDs),
		"message":     "成员添加成功",
	})
}

// RemoveMembers 从用户组批量移除成员（DELETE /api/v1/groups/:id/members）
func (h *UserGroupHandler) RemoveMembers(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户组ID格式"})
		return
	}

	var req dto.AddMemberRequest // 复用结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	if err := h.groupRepo.BatchRemoveMembers(ctx, id, req.UserIDs); err != nil {
		h.logger.Error("批量移除成员失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "移除成员失败"})
		return
	}

	// 更新成员计数
	if updateErr := h.groupRepo.UpdateMemberCount(ctx, id); updateErr != nil {
		h.logger.Warn("更新成员计数失败", zap.Error(updateErr))
	}

	c.JSON(http.StatusOK, gin.H{
		"group_id":      id.String(),
		"removed_count": len(req.UserIDs),
		"message":       "成员移除成功",
	})
}

// ListGroupMembers 成员列表分页查询（GET /api/v1/groups/:id/members）
func (h *UserGroupHandler) ListGroupMembers(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户组ID格式"})
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
	members, total, err := h.groupRepo.GetMembers(ctx, id, page, size)
	if err != nil {
		h.logger.Error("查询组成员列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	list := make([]dto.UserInfo, len(members))
	for i, m := range members {
		list[i] = dto.UserInfo{
			ID:       m.ID.String(),
			Username: m.Username,
			RealName: m.RealName,
			Email:    m.Email,
			Phone:    m.Phone,
			Avatar:   m.Avatar,
			Status:   m.Status,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"list":  list,
		"total": total,
		"page":  page,
		"size":  size,
	})
}
