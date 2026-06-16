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

// SuiteHandler 测试套件管理Handler
type SuiteHandler struct {
	suiteRepo repository.TestSuiteRepository
	caseRepo  repository.TestCaseRepository
	logger    *zap.Logger
}

// NewSuiteHandler 创建测试套件管理Handler实例
func NewSuiteHandler(suiteRepo repository.TestSuiteRepository, caseRepo repository.TestCaseRepository, logger *zap.Logger) *SuiteHandler {
	return &SuiteHandler{
		suiteRepo: suiteRepo,
		caseRepo:  caseRepo,
		logger:    logger,
	}
}

// CreateSuite 创建套件（POST /api/v1/test-suites�?
func (h *SuiteHandler) CreateSuite(c *gin.Context) {
	var req dto.CreateTestSuiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	currentUserID, _ := getCurrentUserID(c)

	suite := &entity.TestSuite{
		Name:        req.Name,
		Description: req.Description,
		ProductID:   req.ProductID,
		ProjectID:   req.ProjectID,
		CreatorID:   currentUserID,
	}

	if err := h.suiteRepo.Create(ctx, suite); err != nil {
		h.logger.Error("创建测试套件失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建测试套件失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "测试套件创建成功",
		"suite_id": suite.ID.String(),
	})
}

// ListSuites 套件列表（GET /api/v1/test-suites�?
func (h *SuiteHandler) ListSuites(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	productID := parseUUIDPtr(c.Query("product_id"))
	projectID := parseUUIDPtr(c.Query("project_id"))

	ctx := c.Request.Context()
	suites, total, err := h.suiteRepo.List(ctx, page, size, productID, projectID)
	if err != nil {
		h.logger.Error("查询测试套件列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询测试套件列表失败"})
		return
	}

	list := make([]dto.TestSuiteInfo, len(suites))
	for i, s := range suites {
		list[i] = dto.TestSuiteInfo{
			ID:          s.ID.String(),
			Name:        s.Name,
			Description: s.Description,
			CreatorID:   s.CreatorID.String(),
			CaseCount:   len(s.Cases),
			CreatedAt:   s.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   s.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, dto.TestSuiteListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// GetSuite 套件详情+用例列表（GET /api/v1/test-suites/:id�?
func (h *SuiteHandler) GetSuite(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的套件ID格式"})
		return
	}

	ctx := c.Request.Context()
	suite, err := h.suiteRepo.GetByID(ctx, id)
	if err != nil || suite == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试套件不存�?})
		return
	}

	// 获取关联的用例信�?
	caseItems := make([]dto.SuiteCaseItem, 0)
	for _, rel := range suite.Cases {
		item := dto.SuiteCaseItem{
			CaseID:    rel.CaseID.String(),
			SortOrder: rel.SortOrder,
		}
		// 查询用例标题
		tc, caseErr := h.caseRepo.GetByID(ctx, rel.CaseID)
		if caseErr == nil && tc != nil {
			item.CaseTitle = tc.Title
		}
		caseItems = append(caseItems, item)
	}

	c.JSON(http.StatusOK, dto.TestSuiteDetailResponse{
		TestSuiteInfo: dto.TestSuiteInfo{
			ID:          suite.ID.String(),
			Name:        suite.Name,
			Description: suite.Description,
			CreatorID:   suite.CreatorID.String(),
			CaseCount:   len(suite.Cases),
			CreatedAt:   suite.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   suite.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
		Cases: caseItems,
	})
}

// UpdateSuite 更新套件（PUT /api/v1/test-suites/:id�?
func (h *SuiteHandler) UpdateSuite(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的套件ID格式"})
		return
	}

	var req dto.UpdateTestSuiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	suite, err := h.suiteRepo.GetByID(ctx, id)
	if err != nil || suite == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试套件不存�?})
		return
	}

	if req.Name != nil {
		suite.Name = *req.Name
	}
	if req.Description != nil {
		suite.Description = *req.Description
	}
	if req.ProductID != nil {
		suite.ProductID = req.ProductID
	}
	if req.ProjectID != nil {
		suite.ProjectID = req.ProjectID
	}

	if err := h.suiteRepo.Update(ctx, suite); err != nil {
		h.logger.Error("更新测试套件失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新测试套件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "测试套件更新成功"})
}

// DeleteSuite 删除套件（DELETE /api/v1/test-suites/:id�?
func (h *SuiteHandler) DeleteSuite(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的套件ID格式"})
		return
	}

	ctx := c.Request.Context()
	suite, getErr := h.suiteRepo.GetByID(ctx, id)
	if getErr != nil || suite == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试套件不存�?})
		return
	}

	if err := h.suiteRepo.Delete(ctx, id); err != nil {
		h.logger.Error("删除测试套件失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除测试套件失败"})
		return
	}

	h.logger.Info("删除测试套件成功", zap.String("suite_id", id.String()))
	c.JSON(http.StatusOK, gin.H{"message": "测试套件删除成功"})
}

// AddCasesToSuite 添加用例到套件（POST /api/v1/test-suites/:id/cases�?
func (h *SuiteHandler) AddCasesToSuite(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的套件ID格式"})
		return
	}

	var req dto.AddCasesToSuiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 验证套件存在
	suite, getErr := h.suiteRepo.GetByID(ctx, id)
	if getErr != nil || suite == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "测试套件不存�?})
		return
	}

	if err := h.suiteRepo.AddCases(ctx, id, req.CaseIDs); err != nil {
		h.logger.Error("添加用例到套件失�?, zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加用例到套件失�?})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "用例添加成功",
		"added_count": len(req.CaseIDs),
	})
}

// RemoveCaseFromSuite 移除用例（DELETE /api/v1/test-suites/:id/cases/:cid�?
func (h *SuiteHandler) RemoveCaseFromSuite(c *gin.Context) {
	suiteID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的套件ID格式"})
		return
	}
	caseID, err := uuid.Parse(c.Param("cid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用例ID格式"})
		return
	}

	ctx := c.Request.Context()

	if err := h.suiteRepo.RemoveCase(ctx, suiteID, caseID); err != nil {
		h.logger.Error("从套件移除用例失�?, zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "移除用例失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用例已从套件中移�?})
}
