package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-project/internal/application"
	"leap-one/service-project/internal/domain/entity"
	"leap-one/service-project/internal/interfaces/api/dto"
)

// MemberHandler 项目成员Handler
type MemberHandler struct {
	memberSvc *application.ProjectMemberService
	logger    *zap.Logger
}

// NewMemberHandler 创建成员管理Handler实例
func NewMemberHandler(memberSvc *application.ProjectMemberService, logger *zap.Logger) *MemberHandler {
	return &MemberHandler{
		memberSvc: memberSvc,
		logger:    logger,
	}
}

// ListMembers 获取项目成员列表（GET /api/v1/projects/:id/members�?
func (h *MemberHandler) ListMembers(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}

	ctx := c.Request.Context()
	members, total, err := h.memberSvc.ListMembers(ctx, projectID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	list := make([]dto.MemberInfo, len(members))
	for i, m := range members {
		list[i] = dto.MemberInfo{
			ID:       m.ID.String(),
			UserID:   m.UserID.String(),
			Role:     m.Role,
			JoinTime: m.JoinTime.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, dto.MemberListResponse{
		List:  list,
		Total: total,
	})
}

// AddMember 添加项目成员（POST /api/v1/projects/:id/members�?
func (h *MemberHandler) AddMember(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}

	var req dto.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	member, err := h.memberSvc.AddMember(ctx, projectID, req.UserID, req.Role)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "成员添加成功",
		"member_id": member.ID.String(),
	})
}

// RemoveMember 移除项目成员（DELETE /api/v1/projects/:id/members/:uid�?
func (h *MemberHandler) RemoveMember(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}
	userID, err := uuid.Parse(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.memberSvc.RemoveMember(ctx, projectID, userID); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "成员已移�?})
}

// UpdateMemberRole 更新成员角色（PUT /api/v1/projects/:id/members/:uid�?
func (h *MemberHandler) UpdateMemberRole(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}
	userID, err := uuid.Parse(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID格式"})
		return
	}

	var req dto.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	if err := h.memberSvc.UpdateMemberRole(ctx, projectID, userID, req.Role); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "角色更新成功"})
}

// ==================== 里程碑Handler ====================

// MilestoneHandler 里程碑管理Handler
type MilestoneHandler struct {
	milestoneSvc *application.MilestoneService
	logger       *zap.Logger
}

// NewMilestoneHandler 创建里程碑管理Handler实例
func NewMilestoneHandler(milestoneSvc *application.MilestoneService, logger *zap.Logger) *MilestoneHandler {
	return &MilestoneHandler{
		milestoneSvc: milestoneSvc,
		logger:       logger,
	}
}

// ListMilestones 获取里程碑列表（GET /api/v1/projects/:id/milestones�?
func (h *MilestoneHandler) ListMilestones(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}

	ctx := c.Request.Context()
	milestones, err := h.milestoneSvc.ListMilestones(ctx, projectID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	list := make([]dto.MilestoneInfo, len(milestones))
	for i, ms := range milestones {
		list[i] = buildMilestoneInfo(ms)
	}

	c.JSON(http.StatusOK, gin.H{"list": list})
}

// CreateMilestone 创建里程碑（POST /api/v1/projects/:id/milestones�?
func (h *MilestoneHandler) CreateMilestone(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}

	var req dto.CreateMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	input := &application.CreateMilestoneInput{
		Name:        req.Name,
		Description: req.Description,
		DueDate:     req.DueDate,
		SortOrder:   req.SortOrder,
	}

	ms, err := h.milestoneSvc.CreateMilestone(ctx, projectID, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "里程碑创建成�?,
		"milestone_id": ms.ID.String(),
	})
}

// UpdateMilestone 更新里程碑（PUT /api/v1/projects/:id/milestones/:mid�?
func (h *MilestoneHandler) UpdateMilestone(c *gin.Context) {
	id, err := uuid.Parse(c.Param("mid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的里程碑ID格式"})
		return
	}

	var req dto.UpdateMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	input := &application.UpdateMilestoneInput{
		Name:        req.Name,
		Description: req.Description,
		DueDate:     req.DueDate,
		SortOrder:   req.SortOrder,
	}

	ms, err := h.milestoneSvc.UpdateMilestone(ctx, id, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "里程碑更新成�?,
		"data":    buildMilestoneInfo(ms),
	})
}

// DeleteMilestone 删除里程碑（DELETE /api/v1/projects/:id/milestones/:mid�?
func (h *MilestoneHandler) DeleteMilestone(c *gin.Context) {
	id, err := uuid.Parse(c.Param("mid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的里程碑ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.milestoneSvc.DeleteMilestone(ctx, id); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "里程碑删除成�?})
}

// CompleteMilestone 完成里程碑（PUT /api/v1/projects/:id/milestones/:mid/complete�?
func (h *MilestoneHandler) CompleteMilestone(c *gin.Context) {
	id, err := uuid.Parse(c.Param("mid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的里程碑ID格式"})
		return
	}

	ctx := c.Request.Context()
	completedBy := getCurrentUserID(c)

	if err := h.milestoneSvc.CompleteMilestone(ctx, id, completedBy); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "里程碑已完成"})
}

// buildMilestoneInfo 构建里程碑信�?
func buildMilestoneInfo(ms *entity.ProjectMilestone) dto.MilestoneInfo {
	info := dto.MilestoneInfo{
		ID:          ms.ID.String(),
		ProjectID:   ms.ProjectID.String(),
		Name:        ms.Name,
		Description: ms.Description,
		DueDate:     ms.DueDate.Format("2006-01-02"),
		Status:      ms.Status,
		SortOrder:   ms.SortOrder,
		CreatedAt:   ms.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if ms.CompletedAt != nil {
		info.CompletedAt = ms.CompletedAt.Format("2006-01-02 15:04:05")
	}
	if ms.CompletedBy != nil {
		info.CompletedBy = ms.CompletedBy.String()
	}
	return info
}
