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

// ==================== 风险Handler ====================

// RiskHandler 风险管理Handler
type RiskHandler struct {
	riskSvc *application.RiskService
	logger  *zap.Logger
}

// NewRiskHandler 创建风险管理Handler实例
func NewRiskHandler(riskSvc *application.RiskService, logger *zap.Logger) *RiskHandler {
	return &RiskHandler{
		riskSvc: riskSvc,
		logger:  logger,
	}
}

// ListRisks 获取风险列表（GET /api/v1/projects/:id/risks�?
func (h *RiskHandler) ListRisks(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}

	ctx := c.Request.Context()
	risks, err := h.riskSvc.ListRisks(ctx, projectID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	list := make([]dto.RiskInfo, len(risks))
	for i, r := range risks {
		list[i] = buildRiskInfo(r)
	}

	c.JSON(http.StatusOK, dto.RiskListResponse{List: list, Total: int64(len(list))})
}

// CreateRisk 创建风险（POST /api/v1/projects/:id/risks�?
func (h *RiskHandler) CreateRisk(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}

	var req dto.CreateRiskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	input := &application.CreateRiskInput{
		Title:       req.Title,
		Description: req.Description,
		Probability: req.Probability,
		Impact:      req.Impact,
		OwnerID:     req.OwnerID,
		Mitigation:  req.Mitigation,
	}

	risk, err := h.riskSvc.CreateRisk(ctx, projectID, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "风险创建成功",
		"risk_id": risk.ID.String(),
	})
}

// UpdateRisk 更新风险（PUT /api/v1/projects/:id/risks/:rid�?
func (h *RiskHandler) UpdateRisk(c *gin.Context) {
	id, err := uuid.Parse(c.Param("rid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的风险ID格式"})
		return
	}

	var req dto.UpdateRiskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	input := &application.UpdateRiskInput{
		Title:       req.Title,
		Description: req.Description,
		Probability: req.Probability,
		Impact:      req.Impact,
		OwnerID:     req.OwnerID,
		Mitigation:  req.Mitigation,
		Status:      req.Status,
	}

	risk, err := h.riskSvc.UpdateRisk(ctx, id, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "风险更新成功",
		"data":    buildRiskInfo(risk),
	})
}

// DeleteRisk 删除风险（DELETE /api/v1/projects/:id/risks/:rid�?
func (h *RiskHandler) DeleteRisk(c *gin.Context) {
	id, err := uuid.Parse(c.Param("rid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的风险ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.riskSvc.DeleteRisk(ctx, id); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "风险删除成功"})
}

// buildRiskInfo 构建风险信息
func buildRiskInfo(r *entity.ProjectRisk) dto.RiskInfo {
	return dto.RiskInfo{
		ID:          r.ID.String(),
		ProjectID:   r.ProjectID.String(),
		Title:       r.Title,
		Description: r.Description,
		Probability: r.Probability,
		Impact:      r.Impact,
		Severity:    r.Severity,
		Status:      r.Status,
		OwnerID:     r.OwnerID.String(),
		Mitigation:  r.Mitigation,
		CreatedAt:   r.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   r.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ==================== 自定义字段Handler ====================

// CustomFieldHandler 自定义字段Handler
type CustomFieldHandler struct {
	fieldSvc *application.CustomFieldService
	logger   *zap.Logger
}

// NewCustomFieldHandler 创建自定义字段Handler实例
func NewCustomFieldHandler(fieldSvc *application.CustomFieldService, logger *zap.Logger) *CustomFieldHandler {
	return &CustomFieldHandler{
		fieldSvc: fieldSvc,
		logger:   logger,
	}
}

// ListCustomFields 获取自定义字段列表（GET /api/v1/projects/:id/custom-fields�?
func (h *CustomFieldHandler) ListCustomFields(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}

	ctx := c.Request.Context()
	fields, err := h.fieldSvc.ListCustomFields(ctx, projectID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	list := make([]dto.CustomFieldInfo, len(fields))
	for i, f := range fields {
		list[i] = buildCustomFieldInfo(f)
	}

	c.JSON(http.StatusOK, dto.CustomFieldListResponse{List: list})
}

// AddCustomField 添加自定义字段（POST /api/v1/projects/:id/custom-fields�?
func (h *CustomFieldHandler) AddCustomField(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID格式"})
		return
	}

	var req dto.CreateCustomFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	input := &application.AddCustomFieldInput{
		Name:      req.Name,
		FieldKey:  req.FieldKey,
		FieldType: req.FieldType,
		Options:   req.Options,
		Required:  req.Required,
		SortOrder: req.SortOrder,
	}

	field, err := h.fieldSvc.AddCustomField(ctx, projectID, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "自定义字段添加成�?,
		"field_id": field.ID.String(),
	})
}

// UpdateCustomField 更新自定义字段（PUT /api/v1/projects/:id/custom-fields/:fid�?
func (h *CustomFieldHandler) UpdateCustomField(c *gin.Context) {
	id, err := uuid.Parse(c.Param("fid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的字段ID格式"})
		return
	}

	var req dto.UpdateCustomFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	input := &application.UpdateCustomFieldInput{
		Name:      req.Name,
		FieldKey:  req.FieldKey,
		FieldType: req.FieldType,
		Options:   req.Options,
		Required:  req.Required,
		SortOrder: req.SortOrder,
	}

	field, err := h.fieldSvc.UpdateCustomField(ctx, id, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "自定义字段更新成�?,
		"data":    buildCustomFieldInfo(field),
	})
}

// DeleteCustomField 删除自定义字段（DELETE /api/v1/projects/:id/custom-fields/:fid�?
func (h *CustomFieldHandler) DeleteCustomField(c *gin.Context) {
	id, err := uuid.Parse(c.Param("fid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的字段ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.fieldSvc.DeleteCustomField(ctx, id); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "自定义字段删除成�?})
}

// buildCustomFieldInfo 构建自定义字段信�?
func buildCustomFieldInfo(f *entity.CustomField) dto.CustomFieldInfo {
	return dto.CustomFieldInfo{
		ID:        f.ID.String(),
		ProjectID: f.ProjectID.String(),
		Name:      f.Name,
		FieldKey:  f.FieldKey,
		FieldType: f.FieldType,
		Options:   f.Options,
		Required:  f.Required,
		SortOrder: f.SortOrder,
		CreatedAt: f.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
