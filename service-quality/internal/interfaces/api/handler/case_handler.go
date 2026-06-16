package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-quality/internal/domain/entity"
	"leap-one/service-quality/internal/domain/repository"
	"leap-one/service-quality/internal/interfaces/api/dto"
)

// CaseHandler 测试用例管理Handler
type CaseHandler struct {
	caseRepo repository.TestCaseRepository
	logger   *zap.Logger
}

// NewCaseHandler 创建测试用例管理Handler实例
func NewCaseHandler(caseRepo repository.TestCaseRepository, logger *zap.Logger) *CaseHandler {
	return &CaseHandler{
		caseRepo: caseRepo,
		logger:   logger,
	}
}

// CreateCase 创建用例（POST /api/v1/test-cases�?
func (h *CaseHandler) CreateCase(c *gin.Context) {
	var req dto.CreateTestCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	currentUserID, _ := getCurrentUserID(c)

	testCase := &entity.TestCase{
		Title:          req.Title,
		Module:         req.Module,
		Precondition:   req.Precondition,
		Steps:          req.Steps,
		ExpectedResult: req.ExpectedResult,
		Priority:       req.Priority,
		Type:           req.Type,
		Automation:     req.Automation,
		ProductID:      req.ProductID,
		ProjectID:      req.ProjectID,
		RequirementID:  req.RequirementID,
		Tags:           req.Tags,
		CreatorID:      currentUserID,
	}
	if testCase.Type == "" {
		testCase.Type = "manual"
	}
	if testCase.Priority == 0 {
		testCase.Priority = 3
	}

	if err := h.caseRepo.Create(ctx, testCase); err != nil {
		h.logger.Error("创建测试用例失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建测试用例失败"})
		return
	}

	h.logger.Info("创建测试用例成功",
		zap.String("case_id", testCase.ID.String()),
		zap.String("title", testCase.Title),
	)

	c.JSON(http.StatusCreated, gin.H{
		"message": "测试用例创建成功",
		"case_id": testCase.ID.String(),
	})
}

// ListCases 用例列表（分�?筛选）（GET /api/v1/test-cases�?
func (h *CaseHandler) ListCases(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	filter := &repository.TestCaseFilter{
		Keyword:       c.Query("keyword"),
		Type:          c.Query("type"),
		Status:        c.Query("status"),
		ProductID:     parseUUIDPtr(c.Query("product_id")),
		ProjectID:     parseUUIDPtr(c.Query("project_id")),
		CreatorID:     parseUUIDPtr(c.Query("creator_id")),
		RequirementID: parseUUIDPtr(c.Query("requirement_id")),
	}

	if prioStr := c.Query("priority"); prioStr != "" {
		prio, err := strconv.Atoi(prioStr)
		if err == nil {
			filter.Priority = &prio
		}
	}
	if autoStr := c.Query("automation"); autoStr != "" {
		auto := autoStr == "true" || autoStr == "1"
		filter.Automation = &auto
	}

	ctx := c.Request.Context()
	cases, total, err := h.caseRepo.List(ctx, page, size, filter)
	if err != nil {
		h.logger.Error("查询测试用例列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询测试用例列表失败"})
		return
	}

	list := make([]dto.TestCaseInfo, len(cases))
	for i, tc := range cases {
		list[i] = buildTestCaseInfo(tc)
	}

	c.JSON(http.StatusOK, dto.TestCaseListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// GetCase 用例详情（GET /api/v1/test-cases/:id�?
func (h *CaseHandler) GetCase(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用例ID格式"})
		return
	}

	ctx := c.Request.Context()
	tc, err := h.caseRepo.GetByID(ctx, id)
	if err != nil || tc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试用例不存�?})
		return
	}

	resp := buildTestCaseDetail(tc)
	c.JSON(http.StatusOK, resp)
}

// UpdateCase 更新用例（PUT /api/v1/test-cases/:id�?
func (h *CaseHandler) UpdateCase(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用例ID格式"})
		return
	}

	var req dto.UpdateTestCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	tc, err := h.caseRepo.GetByID(ctx, id)
	if err != nil || tc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试用例不存�?})
		return
	}

	// 更新字段
	applyTestCaseUpdate(tc, &req)

	if err := h.caseRepo.Update(ctx, tc); err != nil {
		h.logger.Error("更新测试用例失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新测试用例失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "测试用例更新成功"})
}

// DeleteCase 删除用例（DELETE /api/v1/test-cases/:id�?
func (h *CaseHandler) DeleteCase(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用例ID格式"})
		return
	}

	ctx := c.Request.Context()

	tc, getErr := h.caseRepo.GetByID(ctx, id)
	if getErr != nil || tc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试用例不存�?})
		return
	}

	if err := h.caseRepo.Delete(ctx, id); err != nil {
		h.logger.Error("删除测试用例失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除测试用例失败"})
		return
	}

	h.logger.Info("删除测试用例成功", zap.String("case_id", id.String()))
	c.JSON(http.StatusOK, gin.H{"message": "测试用例删除成功"})
}

// ImportCases 导入用例（POST /api/v1/test-cases/import�?
func (h *CaseHandler) ImportCases(c *gin.Context) {
	var req dto.ImportTestCasesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	currentUserID, _ := getCurrentUserID(c)

	cases := make([]*entity.TestCase, len(req.Cases))
	for i, item := range req.Cases {
		tcType := item.Type
		if tcType == "" {
			tcType = "manual"
		}
		priority := item.Priority
		if priority == 0 {
			priority = 3
		}
		cases[i] = &entity.TestCase{
			Title:          item.Title,
			Module:         item.Module,
			Precondition:   item.Precondition,
			Steps:          item.Steps,
			ExpectedResult: item.ExpectedResult,
			Priority:       priority,
			Type:           tcType,
			Automation:     item.Automation,
			ProductID:      item.ProductID,
			ProjectID:      item.ProjectID,
			RequirementID:  item.RequirementID,
			Tags:           item.Tags,
			CreatorID:      currentUserID,
		}
	}

	if err := h.caseRepo.BatchCreate(ctx, cases); err != nil {
		h.logger.Error("导入测试用例失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "导入测试用例失败"})
		return
	}

	h.logger.Info("导入测试用例成功", zap.Int("count", len(cases)))
	c.JSON(http.StatusOK, gin.H{
		"message": "导入成功",
		"count":   len(cases),
	})
}

// ReviewCase 评审用例（POST /api/v1/test-cases/:id/review�?
func (h *CaseHandler) ReviewCase(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用例ID格式"})
		return
	}

	var req dto.ReviewTestCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	tc, err := h.caseRepo.GetByID(ctx, id)
	if err != nil || tc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试用例不存�?})
		return
	}

	if err := h.caseRepo.Review(ctx, id, req.ReviewerID); err != nil {
		h.logger.Error("评审测试用例失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评审测试用例失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用例评审完成"})
}

// buildTestCaseInfo 构建用例简要信�?
func buildTestCaseInfo(tc *entity.TestCase) dto.TestCaseInfo {
	info := dto.TestCaseInfo{
		ID:         tc.ID.String(),
		Title:      tc.Title,
		Module:     tc.Module,
		Priority:   tc.Priority,
		Type:       tc.Type,
		Status:     tc.Status,
		Automation: tc.Automation,
		CreatorID:  tc.CreatorID.String(),
		Version:    tc.Version,
		CreatedAt:  tc.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  tc.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if tc.ProductID != nil {
		s := tc.ProductID.String()
		info.ProductID = &s
	}
	if tc.ProjectID != nil {
		s := tc.ProjectID.String()
		info.ProjectID = &s
	}
	return info
}

// buildTestCaseDetail 构建用例详情响应
func buildTestCaseDetail(tc *entity.TestCase) dto.TestCaseDetailResponse {
	info := buildTestCaseInfo(tc)
	detail := dto.TestCaseDetailResponse{
		TestCaseInfo:   info,
		Precondition:   tc.Precondition,
		Steps:          tc.Steps,
		ExpectedResult: tc.ExpectedResult,
		Tags:           tc.Tags,
	}
	if tc.RequirementID != nil {
		s := tc.RequirementID.String()
		detail.RequirementID = &s
	}
	if tc.ReviewerID != nil {
		s := tc.ReviewerID.String()
		detail.ReviewerID = &s
	}
	if tc.ReviewedAt != nil {
		s := tc.ReviewedAt.Format("2006-01-02 15:04:05")
		detail.ReviewedAt = &s
	}
	return detail
}

// applyTestCaseUpdate 应用更新字段到用例实�?
func applyTestCaseUpdate(tc *entity.TestCase, req *dto.UpdateTestCaseRequest) {
	if req.Title != nil {
		tc.Title = *req.Title
	}
	if req.Module != nil {
		tc.Module = *req.Module
	}
	if req.Precondition != nil {
		tc.Precondition = *req.Precondition
	}
	if req.Steps != nil {
		tc.Steps = *req.Steps
	}
	if req.ExpectedResult != nil {
		tc.ExpectedResult = *req.ExpectedResult
	}
	if req.Priority != nil {
		tc.Priority = *req.Priority
	}
	if req.Type != nil {
		tc.Type = *req.Type
	}
	if req.Status != nil {
		tc.Status = *req.Status
	}
	if req.Automation != nil {
		tc.Automation = *req.Automation
	}
	if req.ProductID != nil {
		tc.ProductID = req.ProductID
	}
	if req.ProjectID != nil {
		tc.ProjectID = req.ProjectID
	}
	if req.RequirementID != nil {
		tc.RequirementID = req.RequirementID
	}
	if req.Tags != nil {
		tc.Tags = *req.Tags
	}
	// 更新时版本号递增
	tc.Version++
}
