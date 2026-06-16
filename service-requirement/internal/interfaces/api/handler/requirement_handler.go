package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-requirement/internal/application/service"
	"leap-one/service-requirement/internal/domain/entity"
	"leap-one/service-requirement/internal/domain/repository"
	"leap-one/service-requirement/internal/interfaces/api/dto"
)

// RequirementHandler 需求HTTP处理�?
type RequirementHandler struct {
	reqService    *service.RequirementService
	reviewService *service.ReviewService
	changeService *service.ChangeLogService
	relationSvc   *service.RelationService
	logger        *zap.Logger
}

// NewRequirementHandler 创建需求处理器实例
func NewRequirementHandler(
	reqService *service.RequirementService,
	reviewService *service.ReviewService,
	changeService *service.ChangeLogService,
	relationSvc *service.RelationService,
	logger *zap.Logger,
) *RequirementHandler {
	return &RequirementHandler{
		reqService:    reqService,
		reviewService: reviewService,
		changeService: changeService,
		relationSvc:   relationSvc,
		logger:        logger,
	}
}

// CreateRequirement 创建需�?
func (h *RequirementHandler) CreateRequirement(c *gin.Context) {
	var req dto.CreateRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest("参数校验失败: "+err.Error()))
		return
	}

	entity := &entity.Requirement{
		ID:             uuid.New(),
		Title:          req.Title,
		Description:    req.Description,
		Type:           req.Type,
		ParentID:       req.ParentID,
		ProductID:      req.ProductID,
		ProjectID:      req.ProjectID,
		Priority:       req.Priority,
		Source:         req.Source,
		Category:       req.Category,
		OwnerID:        req.OwnerID,
		ReviewerID:     req.ReviewerID,
		StoryPoints:    req.StoryPoints,
		EstimatedHours: req.EstimatedHours,
		ReleaseVersion: req.ReleaseVersion,
		Stage:          req.Stage,
		SourceURL:      req.SourceURL,
		Tags:           req.Tags,
	}

	result, err := h.reqService.CreateRequirement(entity)
	if err != nil {
		h.logger.Error("创建需求失�?, zap.Error(err))
		c.JSON(500, dto.InternalError("创建需求失�?))
		return
	}

	c.JSON(201, dto.Success(toRequirementResponse(result)))
}

// GetRequirement 获取需求详�?
func (h *RequirementHandler) GetRequirement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效的需求ID"))
		return
	}

	req, err := h.reqService.GetRequirement(id)
	if err != nil {
		c.JSON(404, dto.NotFound("需求不存在"))
		return
	}

	c.JSON(200, dto.Success(toRequirementResponse(req)))
}

// UpdateRequirement 更新需�?
func (h *RequirementHandler) UpdateRequirement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效的需求ID"))
		return
	}

	var req dto.UpdateRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest("参数校验失败: "+err.Error()))
		return
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.Stage != nil {
		updates["stage"] = *req.Stage
	}
	if req.OwnerID != nil {
		updates["owner_id"] = *req.OwnerID
	}
	if req.ReleaseVersion != nil {
		updates["release_version"] = *req.ReleaseVersion
	}

	result, err := h.reqService.UpdateRequirement(id, updates)
	if err != nil {
		h.logger.Error("更新需求失�?, zap.Error(err))
		c.JSON(500, dto.InternalError("更新需求失�?))
		return
	}

	c.JSON(200, dto.Success(toRequirementResponse(result)))
}

// DeleteRequirement 删除需�?
func (h *RequirementHandler) DeleteRequirement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效的需求ID"))
		return
	}

	if err := h.reqService.DeleteRequirement(id); err != nil {
		h.logger.Error("删除需求失�?, zap.Error(err))
		c.JSON(500, dto.InternalError("删除需求失�?))
		return
	}

	c.JSON(200, dto.Success(nil))
}

// ListRequirements 需求列表（分页+高级筛选）
func (h *RequirementHandler) ListRequirements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	params := &repository.RequirementListParams{
		Page:      page,
		PageSize:  pageSize,
		Type:      c.Query("type"),
		Status:    c.Query("status"),
		Category:  c.Query("category"),
		Stage:     c.Query("stage"),
		Keyword:   c.Query("keyword"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	}

	if productID := c.Query("product_id"); productID != "" {
		if pid, err := uuid.Parse(productID); err == nil {
			params.ProductID = &pid
		}
	}
	if projectID := c.Query("project_id"); projectID != "" {
		if pid, err := uuid.Parse(projectID); err == nil {
			params.ProjectID = &pid
		}
	}
	if ownerID := c.Query("owner_id"); ownerID != "" {
		if oid, err := uuid.Parse(ownerID); err == nil {
			params.OwnerID = &oid
		}
	}
	if priorityStr := c.Query("priority"); priorityStr != "" {
		if p, err := strconv.Atoi(priorityStr); err == nil {
			params.Priority = &p
		}
	}

	list, total, err := h.reqService.ListRequirements(params)
	if err != nil {
		h.logger.Error("查询需求列表失�?, zap.Error(err))
		c.JSON(500, dto.InternalError("查询需求列表失�?))
		return
	}

	var responses []dto.RequirementResponse
	for _, item := range list {
		responses = append(responses, toRequirementResponse(item))
	}

	c.JSON(200, dto.PageSuccess(responses, total, page, pageSize))
}

// GetRequirementTree 获取需求树（产品维度）
func (h *RequirementHandler) GetRequirementTree(c *gin.Context) {
	productID, err := uuid.Parse(c.Query("product_id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效的产品ID"))
		return
	}

	tree, err := h.reqService.GetRequirementTree(productID)
	if err != nil {
		h.logger.Error("获取需求树失败", zap.Error(err))
		c.JSON(500, dto.InternalError("获取需求树失败"))
		return
	}

	var responses []dto.RequirementResponse
	for _, item := range tree {
		responses = append(responses, toRequirementResponseWithChildren(item))
	}

	c.JSON(200, dto.Success(responses))
}

// UpdateStatus 更改状�?
func (h *RequirementHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效的需求ID"))
		return
	}

	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest("参数校验失败: "+err.Error()))
		return
	}

	if err := h.reqService.UpdateStatus(id, req.Status); err != nil {
		h.logger.Error("更新状态失�?, zap.Error(err))
		c.JSON(500, dto.InternalError("更新状态失�?))
		return
	}

	c.JSON(200, dto.Success(gin.H{"status": req.Status}))
}

// SubmitReview 提交评审
func (h *RequirementHandler) SubmitReview(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效的需求ID"))
		return
	}

	var req dto.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest("参数校验失败: "+err.Error()))
		return
	}

	var participants []entity.RequirementReviewParticipant
	for _, p := range req.Participants {
		participants = append(participants, entity.RequirementReviewParticipant{
			UserID:  p.UserID,
			Opinion: p.Opinion,
			Comment: p.Comment,
		})
	}

	review, err := h.reviewService.SubmitReview(id, req.Title, req.MeetingDate, uuid.New(), participants)
	if err != nil {
		h.logger.Error("提交评审失败", zap.Error(err))
		c.JSON(500, dto.InternalError("提交评审失败"))
		return
	}

	c.JSON(201, dto.Success(dto.ReviewResponse{
		ID:            review.ID,
		RequirementID: review.RequirementID,
		Title:         review.Title,
		MeetingDate:   review.MeetingDate,
		Status:        review.Status,
		CreatorID:     review.CreatorID,
	}))
}

// GetReviews 获取评审记录
func (h *RequirementHandler) GetReviews(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效的需求ID"))
		return
	}

	reviews, err := h.reviewService.GetReviews(id)
	if err != nil {
		h.logger.Error("获取评审记录失败", zap.Error(err))
		c.JSON(500, dto.InternalError("获取评审记录失败"))
		return
	}

	var responses []dto.ReviewResponse
	for _, r := range reviews {
		responses = append(responses, toReviewResponse(r))
	}

	c.JSON(200, dto.Success(responses))
}

// CreateChangeLog 发起变更
func (h *RequirementHandler) CreateChangeLog(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效的需求ID"))
		return
	}

	var req dto.CreateChangeLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest("参数校验失败: "+err.Error()))
		return
	}

	log := &entity.RequirementChangeLog{
		RequirementID: id,
		ChangeType:    req.ChangeType,
		FieldName:     req.FieldName,
		OldValue:      req.OldValue,
		NewValue:      req.NewValue,
		Reason:        req.Reason,
		ChangeUserID:  req.ChangeUserID,
	}

	if err := h.changeService.CreateChangeLog(log); err != nil {
		h.logger.Error("发起变更失败", zap.Error(err))
		c.JSON(500, dto.InternalError("发起变更失败"))
		return
	}

	c.JSON(201, dto.Success(log))
}

// GetChangeLogs 获取变更日志
func (h *RequirementHandler) GetChangeLogs(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效的需求ID"))
		return
	}

	logs, err := h.changeService.GetChangeLogs(id)
	if err != nil {
		h.logger.Error("获取变更日志失败", zap.Error(err))
		c.JSON(500, dto.InternalError("获取变更日志失败"))
		return
	}

	var responses []dto.ChangeLogResponse
	for _, l := range logs {
		responses = append(responses, dto.ChangeLogResponse{
			ID:            l.ID,
			RequirementID: l.RequirementID,
			ChangeType:    l.ChangeType,
			FieldName:     l.FieldName,
			OldValue:      l.OldValue,
			NewValue:      l.NewValue,
			Reason:        l.Reason,
			ChangeUserID:  l.ChangeUserID,
			ReviewStatus:  l.ReviewStatus,
			CreatedAt:     l.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(200, dto.Success(responses))
}

// AddRelation 添加关联
func (h *RequirementHandler) AddRelation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效的需求ID"))
		return
	}

	var req dto.CreateRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest("参数校验失败: "+err.Error()))
		return
	}

	relation, err := h.relationSvc.AddRelation(id, req.RelatedType, req.RelatedID, req.RelationType)
	if err != nil {
		h.logger.Error("添加关联失败", zap.Error(err))
		c.JSON(500, dto.InternalError("添加关联失败"))
		return
	}

	c.JSON(201, dto.Success(dto.RelationResponse{
		ID:            relation.ID,
		RequirementID: relation.RequirementID,
		RelatedType:   relation.RelatedType,
		RelatedID:     relation.RelatedID,
		RelationType:  relation.RelationType,
	}))
}

// GetRelations 获取关联列表
func (h *RequirementHandler) GetRelations(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效的需求ID"))
		return
	}

	relations, err := h.relationSvc.GetRelations(id)
	if err != nil {
		h.logger.Error("获取关联列表失败", zap.Error(err))
		c.JSON(500, dto.InternalError("获取关联列表失败"))
		return
	}

	var responses []dto.RelationResponse
	for _, r := range relations {
		responses = append(responses, dto.RelationResponse{
			ID:            r.ID,
			RequirementID: r.RequirementID,
			RelatedType:   r.RelatedType,
			RelatedID:     r.RelatedID,
			RelationType:  r.RelationType,
		})
	}

	c.JSON(200, dto.Success(responses))
}

// GetMatrix 需求跟踪矩�?
func (h *RequirementHandler) GetMatrix(c *gin.Context) {
	// 返回需求跟踪矩阵数据（按状态、优先级等维度统计）
	list, _, err := h.reqService.ListRequirements(&repository.RequirementListParams{Page: 1, PageSize: 1000})
	if err != nil {
		c.JSON(500, dto.InternalError("获取需求数据失�?))
		return
	}

	matrix := gin.H{
		"by_status":   countByField(list, "Status"),
		"by_type":     countByField(list, "Type"),
		"by_priority": countByPriority(list),
		"by_stage":    countByField(list, "Stage"),
		"total":       len(list),
	}

	c.JSON(200, dto.Success(matrix))
}

// toRequirementResponse 将实体转换为响应DTO
func toRequirementResponse(req *entity.Requirement) dto.RequirementResponse {
	return dto.RequirementResponse{
		ID:             req.ID,
		Code:           req.Code,
		Title:          req.Title,
		Description:    req.Description,
		Type:           req.Type,
		ParentID:       req.ParentID,
		Level:          req.Level,
		ProductID:      req.ProductID,
		ProjectID:      req.ProjectID,
		Status:         req.Status,
		Priority:       req.Priority,
		Source:         req.Source,
		Category:       req.Category,
		OwnerID:        req.OwnerID,
		ReviewerID:     req.ReviewerID,
		StoryPoints:    req.StoryPoints,
		EstimatedHours: req.EstimatedHours,
		ReleaseVersion: req.ReleaseVersion,
		Stage:          req.Stage,
		SourceURL:      req.SourceURL,
		Tags:           req.Tags,
		CreatedAt:      req.CreatedAt,
		UpdatedAt:      req.UpdatedAt,
	}
}

// toRequirementResponseWithChildren 递归转换带子需求的响应
func toRequirementResponseWithChildren(req *entity.Requirement) dto.RequirementResponse {
	resp := toRequirementResponse(req)
	for i := range req.Children {
		resp.Children = append(resp.Children, toRequirementResponseWithChildren(&req.Children[i]))
	}
	return resp
}

// toReviewResponse 将评审实体转换为响应DTO
func toReviewResponse(review *entity.RequirementReview) dto.ReviewResponse {
	resp := dto.ReviewResponse{
		ID:            review.ID,
		RequirementID: review.RequirementID,
		Title:         review.Title,
		MeetingDate:   review.MeetingDate,
		Status:        review.Status,
		Conclusion:    review.Conclusion,
		Decision:      review.Decision,
		CreatorID:     review.CreatorID,
		CreatedAt:     review.CreatedAt,
	}
	for _, p := range review.Participants {
		resp.Participants = append(resp.Participants, dto.ParticipantResponse{
			ID:      p.ID,
			UserID:  p.UserID,
			Opinion: p.Opinion,
			Comment: p.Comment,
		})
	}
	return resp
}

// 辅助统计函数
func countByField(list []*entity.Requirement, field string) map[string]int {
	result := make(map[string]int)
	for _, req := range list {
		var val string
		switch field {
		case "Status":
			val = req.Status
		case "Type":
			val = req.Type
		case "Stage":
			val = req.Stage
		}
		result[val]++
	}
	return result
}

func countByPriority(list []*entity.Requirement) map[int]int {
	result := make(map[int]int)
	for _, req := range list {
		result[req.Priority]++
	}
	return result
}
