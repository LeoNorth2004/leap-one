package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-project/internal/application"
	"leap-one/service-project/internal/domain/entity"
	"leap-one/service-project/internal/interfaces/api/dto"
)

// ProjectHandler 项目管理Handler
type ProjectHandler struct {
	projectSvc *application.ProjectService
	logger     *zap.Logger
}

// NewProjectHandler 创建项目管理Handler实例
func NewProjectHandler(projectSvc *application.ProjectService, logger *zap.Logger) *ProjectHandler {
	return &ProjectHandler{
		projectSvc: projectSvc,
		logger:     logger,
	}
}

// CreateProject 创建项目（POST /api/v1/projects�?
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	userID := getCurrentUserID(c)

	input := &application.CreateProjectInput{
		Name:        req.Name,
		Code:        req.Code,
		ProgramID:   req.ProgramID,
		Description: req.Description,
		PMID:        req.PMID,
		Type:        req.Type,
		Priority:    req.Priority,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Budget:      req.Budget,
		TemplateID:  req.TemplateID,
		CreatedByID: userID,
	}

	project, err := h.projectSvc.CreateProject(ctx, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "项目创建成功",
		"project_id": project.ID.String(),
	})
}

// GetProject 获取项目详情（GET /api/v1/projects/:id�?
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}

	ctx := c.Request.Context()
	project, err := h.projectSvc.GetProjectDetail(ctx, id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	resp := buildProjectDetailResponse(project)
	c.JSON(http.StatusOK, resp)
}

// UpdateProject 更新项目（PUT /api/v1/projects/:id�?
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}

	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	userID := getCurrentUserID(c)

	input := &application.UpdateProjectInput{
		Name:        req.Name,
		Description: req.Description,
		PMID:        req.PMID,
		Type:        req.Type,
		Priority:    req.Priority,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Budget:      req.Budget,
		UpdatedByID: userID,
	}

	project, err := h.projectSvc.UpdateProject(ctx, id, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "项目更新成功",
		"data":    buildProjectInfo(project),
	})
}

// DeleteProject 删除项目（DELETE /api/v1/projects/:id�?
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.projectSvc.DeleteProject(ctx, id); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "项目删除成功"})
}

// ListProjects 分页查询项目列表（GET /api/v1/projects�?
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")
	programID := c.Query("program_id")
	pmID := c.Query("pm_id")
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	ctx := c.Request.Context()
	projects, total, err := h.projectSvc.ListProjects(
		ctx, page, size, keyword, status, programID, pmID, sortBy, sortOrder,
	)
	if err != nil {
		h.logger.Error("查询项目列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询项目列表失败"})
		return
	}

	list := make([]dto.ProjectInfo, len(projects))
	for i, p := range projects {
		list[i] = buildProjectInfo(p)
	}

	c.JSON(http.StatusOK, dto.ProjectListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// ArchiveProject 归档项目（POST /api/v1/projects/:id/archive�?
func (h *ProjectHandler) ArchiveProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.projectSvc.ChangeProjectStatus(ctx, id, "archived"); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "项目已归�?})
}

// CancelProject 取消项目（POST /api/v1/projects/:id/cancel�?
func (h *ProjectHandler) CancelProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.projectSvc.ChangeProjectStatus(ctx, id, "cancelled"); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "项目已取�?})
}

// ==================== 辅助方法 ====================

// buildProjectInfo 构建项目简要信�?
func buildProjectInfo(p *entity.Project) dto.ProjectInfo {
	info := dto.ProjectInfo{
		ID:          p.ID.String(),
		Name:        p.Name,
		Code:        p.Code,
		Description: p.Description,
		PMID:        p.PMID.String(),
		Status:      p.Status,
		Type:        p.Type,
		Priority:    p.Priority,
		CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if p.ProgramID != nil {
		info.ProgramID = p.ProgramID.String()
	}
	if p.StartDate != nil {
		info.StartDate = p.StartDate.Format("2006-01-02")
	}
	if p.EndDate != nil {
		info.EndDate = p.EndDate.Format("2006-01-02")
	}
	if p.Budget != nil {
		info.Budget = *p.Budget
	}
	return info
}

// buildProjectDetailResponse 构建项目详情响应
func buildProjectDetailResponse(p *entity.Project) dto.ProjectDetailResponse { // Need to fix this - should use entity.Project
	info := dto.ProjectInfo{
		ID:          p.ID.String(),
		Name:        p.Name,
		Code:        p.Code,
		Description: p.Description,
		PMID:        p.PMID.String(),
		Status:      p.Status,
		Type:        p.Type,
		Priority:    p.Priority,
		CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if p.ProgramID != nil {
		info.ProgramID = p.ProgramID.String()
	}
	if p.StartDate != nil {
		info.StartDate = p.StartDate.Format("2006-01-02")
	}
	if p.EndDate != nil {
		info.EndDate = p.EndDate.Format("2006-01-02")
	}
	if p.Budget != nil {
		info.Budget = *p.Budget
	}

	return dto.ProjectDetailResponse{
		ProjectInfo:  info,
		Version:      p.Version,
		CreatedByID:  p.CreatedByID.String(),
		UpdatedAt:    p.UpdatedAt.Format("2006-01-02 15:04:05"),
		Members:      []dto.MemberInfo{},      // TODO: 从成员服务获�?
		Milestones:   []dto.MilestoneInfo{},   // TODO: 从里程碑服务获取
		Risks:        []dto.RiskInfo{},        // TODO: 从风险服务获�?
		CustomFields: []dto.CustomFieldInfo{}, // TODO: 从自定义字段服务获取
	}
}
