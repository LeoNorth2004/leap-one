package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-portfolio/internal/application"
	"leap-one/service-portfolio/internal/domain/entity"
	"leap-one/service-portfolio/internal/interfaces/api/dto"
)

// ProgramHandler 项目集管理Handler
type ProgramHandler struct {
	programSvc *application.ProgramService
	logger     *zap.Logger
}

// NewProgramHandler 创建项目集管理Handler实例
func NewProgramHandler(programSvc *application.ProgramService, logger *zap.Logger) *ProgramHandler {
	return &ProgramHandler{
		programSvc: programSvc,
		logger:     logger,
	}
}

// CreateProgram 创建项目集（POST /api/v1/programs）
func (h *ProgramHandler) CreateProgram(c *gin.Context) {
	var req dto.CreateProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 解析日期字段
	startDate, _ := parseDatePtr(req.StartDate)
	endDate, _ := parseDatePtr(req.EndDate)
	parentID, _ := parseOptionalUUID(req.ParentID)
	ownerID, err := uuid.Parse(req.OwnerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的负责人ID格式"})
		return
	}

	program := &entity.Program{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		ParentID:    parentID,
		OwnerID:     ownerID,
		Budget:      req.Budget,
		StartDate:   startDate,
		EndDate:     endDate,
		Priority:    req.Priority,
		Status:      "active",
	}

	result, svcErr := h.programSvc.CreateProgram(ctx, program)
	if svcErr != nil {
		if svcErr == application.ErrProgramCodeExists {
			c.JSON(http.StatusConflict, gin.H{"error": "项目集编号已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "项目集创建成功",
		"program_id": result.ID.String(),
	})
}

// GetProgram 获取项目集详情（GET /api/v1/programs/:id）
func (h *ProgramHandler) GetProgram(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目集ID格式"})
		return
	}

	ctx := c.Request.Context()
	program, svcErr := h.programSvc.GetProgramDetail(ctx, id)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	resp := h.buildProgramDetailResponse(program)
	c.JSON(http.StatusOK, resp)
}

// UpdateProgram 更新项目集（PUT /api/v1/programs/:id）
func (h *ProgramHandler) UpdateProgram(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目集ID格式"})
		return
	}

	var req dto.UpdateProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	program, svcErr := h.programSvc.GetProgramDetail(ctx, id)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	// 更新字段（仅更新非空字段）
	if req.Name != nil {
		program.Name = *req.Name
	}
	if req.Code != nil {
		program.Code = *req.Code
	}
	if req.Description != nil {
		program.Description = *req.Description
	}
	if req.ParentID != nil {
		pid, pErr := parseOptionalUUID(req.ParentID)
		if pErr {
			program.ParentID = pid
		} else {
			program.ParentID = nil
		}
	}
	if req.OwnerID != nil {
		oid, oErr := uuid.Parse(*req.OwnerID)
		if oErr == nil {
			program.OwnerID = oid
		}
	}
	if req.Status != nil {
		program.Status = *req.Status
	}
	if req.Budget != nil {
		program.Budget = req.Budget
	}
	if req.StartDate != nil {
		program.StartDate, _ = parseDate(*req.StartDate)
	}
	if req.EndDate != nil {
		program.EndDate, _ = parseDate(*req.EndDate)
	}
	if req.Priority != nil {
		program.Priority = *req.Priority
	}

	if svcErr := h.programSvc.UpdateProgram(ctx, program); svcErr != nil {
		if svcErr == application.ErrProgramCodeExists {
			c.JSON(http.StatusConflict, gin.H{"error": "项目集编号已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "项目集更新成功"})
}

// DeleteProgram 删除项目集（DELETE /api/v1/programs/:id）
func (h *ProgramHandler) DeleteProgram(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目集ID格式"})
		return
	}

	ctx := c.Request.Context()
	if svcErr := h.programSvc.DeleteProgram(ctx, id); svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "项目集删除成功"})
}

// ListPrograms 分页查询项目集列表（GET /api/v1/programs）
func (h *ProgramHandler) ListPrograms(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	ctx := c.Request.Context()
	programs, total, svcErr := h.programSvc.ListPrograms(ctx, page, size, keyword, status)
	if svcErr != nil {
		h.logger.Error("查询项目集列表失败", zap.Error(svcErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询项目集列表失败"})
		return
	}

	list := make([]dto.ProgramInfo, len(programs))
	for i, p := range programs {
		list[i] = buildProgramInfo(p)
	}

	c.JSON(http.StatusOK, dto.ProgramListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// GetProgramTree 获取项目集树形结构（GET /api/v1/programs/tree）
func (h *ProgramHandler) GetProgramTree(c *gin.Context) {
	ctx := c.Request.Context()
	tree, svcErr := h.programSvc.GetProgramTree(ctx)
	if svcErr != nil {
		h.logger.Error("获取项目集树失败", zap.Error(svcErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取项目集树失败"})
		return
	}

	list := make([]dto.ProgramInfo, len(tree))
	for i, p := range tree {
		list[i] = buildProgramInfoWithChildren(p)
	}

	c.JSON(http.StatusOK, gin.H{"tree": list})
}

// CreateMilestone 添加里程碑（POST /api/v1/programs/:id/milestones）
func (h *ProgramHandler) CreateMilestone(c *gin.Context) {
	programID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目集ID格式"})
		return
	}

	var req dto.CreateMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	dueDate, _ := parseDate(req.DueDate)

	milestone := &entity.Milestone{
		ProgramID:   programID,
		Name:        req.Name,
		Description: req.Description,
		DueDate:     dueDate,
		Status:      "pending",
	}
	if req.Status != "" {
		milestone.Status = req.Status
	}

	if svcErr := h.programSvc.CreateMilestone(ctx, milestone); svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "里程碑创建成功",
		"milestone_id": milestone.ID.String(),
	})
}

// ListMilestones 获取里程碑列表（GET /api/v1/programs/:id/milestones）
func (h *ProgramHandler) ListMilestones(c *gin.Context) {
	programID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目集ID格式"})
		return
	}

	ctx := c.Request.Context()
	milestones, svcErr := h.programSvc.ListMilestones(ctx, programID)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	list := make([]dto.MilestoneInfo, len(milestones))
	for i, m := range milestones {
		list[i] = dto.MilestoneInfo{
			ID:          m.ID.String(),
			Name:        m.Name,
			Description: m.Description,
			Status:      m.Status,
			CreatedAt:   m.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if m.DueDate != nil {
			ds := m.DueDate.Format("2006-01-02")
			list[i].DueDate = &ds
		}
	}

	c.JSON(http.StatusOK, gin.H{"list": list, "count": len(list)})
}

// CreateRisk 添加风险项（POST /api/v1/programs/:id/risks）
func (h *ProgramHandler) CreateRisk(c *gin.Context) {
	programID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目集ID格式"})
		return
	}

	var req dto.CreateRiskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	ownerID, _ := parseOptionalUUID(req.OwnerID)

	risk := &entity.Risk{
		ProgramID:   programID,
		Title:       req.Title,
		Description: req.Description,
		Probability: req.Probability,
		Impact:      req.Impact,
		OwnerID:     ownerID,
		Status:      "open",
	}
	if req.Status != "" {
		risk.Status = req.Status
	}

	if svcErr := h.programSvc.CreateRisk(ctx, risk); svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "风险项创建成功",
		"risk_id": risk.ID.String(),
	})
}

// ListRisks 获取风险列表（GET /api/v1/programs/:id/risks）
func (h *ProgramHandler) ListRisks(c *gin.Context) {
	programID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目集ID格式"})
		return
	}

	ctx := c.Request.Context()
	risks, svcErr := h.programSvc.ListRisks(ctx, programID)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	list := make([]dto.RiskInfo, len(risks))
	for i, r := range risks {
		list[i] = dto.RiskInfo{
			ID:          r.ID.String(),
			Title:       r.Title,
			Description: r.Description,
			Probability: r.Probability,
			Impact:      r.Impact,
			Status:      r.Status,
			CreatedAt:   r.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if r.OwnerID != nil {
			os := r.OwnerID.String()
			list[i].OwnerID = &os
		}
	}

	c.JSON(http.StatusOK, gin.H{"list": list, "count": len(list)})
}

// GetStatistics 获取项目集统计信息（GET /api/v1/programs/:id/statistics）
func (h *ProgramHandler) GetStatistics(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目集ID格式"})
		return
	}

	ctx := c.Request.Context()
	stats, svcErr := h.programSvc.GetProgramStatistics(ctx, id)
	if svcErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ==================== 辅助方法 ====================

// buildProgramDetailResponse 构建项目集详情响应
func (h *ProgramHandler) buildProgramDetailResponse(p *entity.Program) dto.ProgramDetailResponse {
	info := buildProgramInfo(p)
	return dto.ProgramDetailResponse{
		ProgramInfo: info,
		CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// buildProgramInfo 构建项目集基本信息（不含子节点）
func buildProgramInfo(p *entity.Program) dto.ProgramInfo {
	info := dto.ProgramInfo{
		ID:          p.ID.String(),
		Name:        p.Name,
		Code:        p.Code,
		Description: p.Description,
		OwnerID:     p.OwnerID.String(),
		Status:      p.Status,
		Budget:      p.Budget,
		Priority:    p.Priority,
	}
	if p.ParentID != nil {
		ps := p.ParentID.String()
		info.ParentID = &ps
	}
	if p.StartDate != nil {
		ss := p.StartDate.Format("2006-01-02")
		info.StartDate = &ss
	}
	if p.EndDate != nil {
		es := p.EndDate.Format("2006-01-02")
		info.EndDate = &es
	}
	return info
}

// buildProgramInfoWithChildren 构建含子节点的项目集信息（用于树形结构）
func buildProgramInfoWithChildren(p *entity.Program) dto.ProgramInfo {
	info := buildProgramInfo(p)
	if len(p.Children) > 0 {
		children := make([]dto.ProgramInfo, len(p.Children))
		for i := range p.Children {
			children[i] = buildProgramInfoWithChildren(&p.Children[i])
		}
		info.Children = children
	}
	return info
}

// parseDate 解析日期字符串为time.Time指针
func parseDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// parseDatePtr 解析可选的日期字符串指针
func parseDatePtr(dateStr *string) (*time.Time, error) {
	if dateStr == nil || *dateStr == "" {
		return nil, nil
	}
	return parseDate(*dateStr)
}

// parseOptionalUUID 解析可选的UUID字符串指针
func parseOptionalUUID(s *string) (*uuid.UUID, bool) {
	if s == nil || *s == "" {
		return nil, false
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return nil, false
	}
	return &id, true
}
