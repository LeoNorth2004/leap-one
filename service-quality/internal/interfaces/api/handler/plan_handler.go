package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-quality/internal/domain/entity"
	"leap-one/service-quality/internal/domain/repository"
	"leap-one/service-quality/internal/interfaces/api/dto"
)

// PlanHandler 测试计划管理Handler
type PlanHandler struct {
	planRepo repository.TestPlanRepository
	caseRepo repository.TestCaseRepository
	logger   *zap.Logger
}

// NewPlanHandler 创建测试计划管理Handler实例
func NewPlanHandler(planRepo repository.TestPlanRepository, caseRepo repository.TestCaseRepository, logger *zap.Logger) *PlanHandler {
	return &PlanHandler{
		planRepo: planRepo,
		caseRepo: caseRepo,
		logger:   logger,
	}
}

// CreatePlan 创建计划（POST /api/v1/test-plans�?
func (h *PlanHandler) CreatePlan(c *gin.Context) {
	var req dto.CreateTestPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	currentUserID, _ := getCurrentUserID(c)

	plan := &entity.TestPlan{
		Name:         req.Name,
		Description:  req.Description,
		ProductID:    req.ProductID,
		ProjectID:    req.ProjectID,
		BuildVersion: req.BuildVersion,
		CreatorID:    currentUserID,
		ExecutorIDs:  req.ExecutorIDs,
		Status:       "planning",
	}

	// 解析日期
	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			plan.StartDate = &t
		}
	}
	if req.EndDate != "" {
		if t, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			plan.EndDate = &t
		}
	}

	if err := h.planRepo.Create(ctx, plan); err != nil {
		h.logger.Error("创建测试计划失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建测试计划失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "测试计划创建成功",
		"plan_id": plan.ID.String(),
	})
}

// ListPlans 计划列表（GET /api/v1/test-plans�?
func (h *PlanHandler) ListPlans(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	filter := &repository.TestPlanFilter{
		Status:    c.Query("status"),
		ProductID: parseUUIDPtr(c.Query("product_id")),
		ProjectID: parseUUIDPtr(c.Query("project_id")),
		CreatorID: parseUUIDPtr(c.Query("creator_id")),
	}

	ctx := c.Request.Context()
	plans, total, err := h.planRepo.List(ctx, page, size, filter)
	if err != nil {
		h.logger.Error("查询测试计划列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询测试计划列表失败"})
		return
	}

	list := make([]dto.TestPlanInfo, len(plans))
	for i, p := range plans {
		list[i] = buildTestPlanInfo(p)
	}

	c.JSON(http.StatusOK, dto.TestPlanListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// GetPlan 计划详情（GET /api/v1/test-plans/:id�?
func (h *PlanHandler) GetPlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的计划ID格式"})
		return
	}

	ctx := c.Request.Context()
	plan, err := h.planRepo.GetByID(ctx, id)
	if err != nil || plan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试计划不存�?})
		return
	}

	// 构建用例执行�?
	caseItems := make([]dto.TestPlanCaseItem, 0)
	for _, pc := range plan.Cases {
		item := dto.TestPlanCaseItem{
			PlanCaseID:   pc.ID,
			CaseID:       pc.CaseID,
			Result:       pc.Result,
			ActualResult: pc.ActualResult,
			SortOrder:    pc.SortOrder,
		}
		if pc.AssigneeID != nil {
			s := pc.AssigneeID.String()
			item.AssigneeID = &s
		}
		if pc.ExecuteTime != nil {
			s := pc.ExecuteTime.Format("2006-01-02 15:04:05")
			item.ExecuteTime = &s
		}
		// 查询用例标题
		tc, caseErr := h.caseRepo.GetByID(ctx, pc.CaseID)
		if caseErr == nil && tc != nil {
			item.CaseTitle = tc.Title
		}
		caseItems = append(caseItems, item)
	}

	info := buildTestPlanInfo(plan)
	c.JSON(http.StatusOK, dto.TestPlanDetailResponse{
		TestPlanInfo: info,
		ExecutorIDs:  plan.ExecutorIDs,
		Cases:        caseItems,
	})
}

// UpdatePlan 更新计划（PUT /api/v1/test-plans/:id�?
func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的计划ID格式"})
		return
	}

	var req dto.UpdateTestPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	plan, err := h.planRepo.GetByID(ctx, id)
	if err != nil || plan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试计划不存�?})
		return
	}

	if req.Name != nil {
		plan.Name = *req.Name
	}
	if req.Description != nil {
		plan.Description = *req.Description
	}
	if req.BuildVersion != nil {
		plan.BuildVersion = *req.BuildVersion
	}
	if req.ExecutorIDs != nil {
		plan.ExecutorIDs = *req.ExecutorIDs
	}
	if req.StartDate != nil && *req.StartDate != "" {
		if t, e := time.Parse("2006-01-02", *req.StartDate); e == nil {
			plan.StartDate = &t
		}
	}
	if req.EndDate != nil && *req.EndDate != "" {
		if t, e := time.Parse("2006-01-02", *req.EndDate); e == nil {
			plan.EndDate = &t
		}
	}

	if err := h.planRepo.Update(ctx, plan); err != nil {
		h.logger.Error("更新测试计划失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新测试计划失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "测试计划更新成功"})
}

// DeletePlan 删除计划（DELETE /api/v1/test-plans/:id�?
func (h *PlanHandler) DeletePlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的计划ID格式"})
		return
	}

	ctx := c.Request.Context()
	plan, getErr := h.planRepo.GetByID(ctx, id)
	if getErr != nil || plan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试计划不存�?})
		return
	}

	if err := h.planRepo.Delete(ctx, id); err != nil {
		h.logger.Error("删除测试计划失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除测试计划失败"})
		return
	}

	h.logger.Info("删除测试计划成功", zap.String("plan_id", id.String()))
	c.JSON(http.StatusOK, gin.H{"message": "测试计划删除成功"})
}

// StartPlan 开始执行（POST /api/v1/test-plans/:id/start�?
func (h *PlanHandler) StartPlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的计划ID格式"})
		return
	}

	ctx := c.Request.Context()
	plan, getErr := h.planRepo.GetByID(ctx, id)
	if getErr != nil || plan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试计划不存�?})
		return
	}

	if plan.Status != "planning" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只有规划中的计划才能开始执�?})
		return
	}

	if err := h.planRepo.UpdatePlanStatus(ctx, id, "executing"); err != nil {
		h.logger.Error("开始执行测试计划失�?, zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "测试计划已开始执�?})
}

// CompletePlan 完成计划（POST /api/v1/test-plans/:id/complete�?
func (h *PlanHandler) CompletePlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的计划ID格式"})
		return
	}

	ctx := c.Request.Context()
	plan, getErr := h.planRepo.GetByID(ctx, id)
	if getErr != nil || plan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试计划不存�?})
		return
	}

	if plan.Status != "executing" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只有执行中的计划才能完成"})
		return
	}

	if err := h.planRepo.UpdatePlanStatus(ctx, id, "completed"); err != nil {
		h.logger.Error("完成测试计划失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "测试计划已完�?})
}

// ExecuteCase 执行用例（POST /api/v1/test-plans/:id/cases/:pcid/execute�?
func (h *PlanHandler) ExecuteCase(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的计划ID格式"})
		return
	}
	pcaseID, err := uuid.Parse(c.Param("pcid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用例记录ID格式"})
		return
	}

	var req dto.ExecuteTestCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 验证计划存在且在执行�?
	plan, getErr := h.planRepo.GetByID(ctx, planID)
	if getErr != nil || plan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试计划不存�?})
		return
	}
	if plan.Status != "executing" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只有执行中的计划才能执行用例"})
		return
	}

	now := time.Now()
	result := &entity.TestPlanCase{
		AssigneeID:   req.AssigneeID,
		Result:       req.Result,
		ExecuteTime:  &now,
		ActualResult: req.ActualResult,
		BugIDs:       req.BugIDs,
		Comment:      req.Comment,
	}

	if err := h.planRepo.ExecuteCase(ctx, pcaseID, result); err != nil {
		h.logger.Error("执行测试用例失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "执行用例失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用例执行结果已保�?})
}

// AddCasesToPlan 添加用例到计划（内部方法，可扩展为API端点�?
func (h *PlanHandler) AddCasesToPlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的计划ID格式"})
		return
	}

	var req dto.AddCasesToPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	plan, getErr := h.planRepo.GetByID(ctx, id)
	if getErr != nil || plan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试计划不存�?})
		return
	}

	if err := h.planRepo.AddCases(ctx, id, req.CaseIDs); err != nil {
		h.logger.Error("添加用例到测试计划失�?, zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加用例到测试计划失�?})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "用例添加成功",
		"added_count": len(req.CaseIDs),
	})
}

// buildTestPlanInfo 构建测试计划简要信�?
func buildTestPlanInfo(p *entity.TestPlan) dto.TestPlanInfo {
	info := dto.TestPlanInfo{
		ID:           p.ID.String(),
		Name:         p.Name,
		Description:  p.Description,
		BuildVersion: p.BuildVersion,
		Status:       p.Status,
		CreatorID:    p.CreatorID.String(),
		CaseCount:    len(p.Cases),
		CreatedAt:    p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if p.StartDate != nil {
		s := p.StartDate.Format("2006-01-02")
		info.StartDate = &s
	}
	if p.EndDate != nil {
		s := p.EndDate.Format("2006-01-02")
		info.EndDate = &s
	}
	return info
}
