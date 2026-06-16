package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-task/internal/application"
	"leap-one/service-task/internal/domain/entity"
	"leap-one/service-task/internal/domain/repository"
	"leap-one/service-task/internal/interfaces/api/dto"
)

// IssueHandler 工单管理Handler
type IssueHandler struct {
	issueSvc *application.IssueService
	logger   *zap.Logger
}

// NewIssueHandler 创建工单管理Handler实例
func NewIssueHandler(issueSvc *application.IssueService, logger *zap.Logger) *IssueHandler {
	return &IssueHandler{
		issueSvc: issueSvc,
		logger:   logger,
	}
}

// CreateIssue 创建工单（POST /api/v1/issues）
func (h *IssueHandler) CreateIssue(c *gin.Context) {
	var req dto.CreateIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	userIDVal, _ := c.Get("userID")
	var reporterID uuid.UUID
	if uidStr, ok := userIDVal.(string); ok {
		if parsed, err := uuid.Parse(uidStr); err == nil {
			reporterID = parsed
		} else {
			reporterID = uuid.New()
		}
	} else {
		reporterID = uuid.New()
	}

	input := &application.CreateIssueInput{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		ProjectID:   req.ProjectID,
		ProductID:   req.ProductID,
		ReporterID:  reporterID,
		AssigneeID:  req.AssigneeID,
		Priority:    req.Priority,
		Severity:    req.Severity,
		Source:      req.Source,
		TemplateID:  req.TemplateID,
		Tags:        req.Tags,
	}
	if input.Priority == 0 {
		input.Priority = 3
	}
	if input.Severity == 0 {
		input.Severity = 2
	}

	issue, err := h.issueSvc.CreateIssue(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "工单创建成功",
		"issue_id": issue.ID.String(),
	})
}

// GetIssue 获取工单详情（GET /api/v1/issues/:id）
func (h *IssueHandler) GetIssue(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工单ID格式"})
		return
	}

	ctx := c.Request.Context()
	issue, err := h.issueSvc.GetIssue(ctx, id)
	if err != nil || issue == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "工单不存在"})
		return
	}

	resp := h.buildIssueDetailResponse(issue)
	c.JSON(http.StatusOK, resp)
}

// UpdateIssue 更新工单（PUT /api/v1/issues/:id）
func (h *IssueHandler) UpdateIssue(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工单ID格式"})
		return
	}

	var req dto.UpdateIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	input := &application.UpdateIssueInput{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		ProjectID:   req.ProjectID,
		ProductID:   req.ProductID,
		AssigneeID:  req.AssigneeID,
		Priority:    req.Priority,
		Severity:    req.Severity,
		Resolution:  req.Resolution,
		Tags:        req.Tags,
	}

	issue, svcErr := h.issueSvc.UpdateIssue(ctx, id, input)
	if svcErr != nil {
		if svcErr == application.ErrIssueNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "工单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "工单更新成功", "issue_id": issue.ID.String()})
}

// DeleteIssue 删除工单（DELETE /api/v1/issues/:id）
func (h *IssueHandler) DeleteIssue(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工单ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.issueSvc.DeleteIssue(ctx, id); err != nil {
		if err == application.ErrIssueNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "工单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "工单删除成功"})
}

// ListIssues 工单列表（GET /api/v1/issues）
func (h *IssueHandler) ListIssues(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	filter := &repository.IssueFilter{
		Page:      page,
		PageSize:  size,
		Keyword:   c.Query("keyword"),
		Status:    c.Query("status"),
		Type:      c.Query("type"),
		Source:    c.Query("source"),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "DESC"),
	}

	if pStr := c.Query("priority"); pStr != "" {
		var p int
		if _, e := strconv.Atoi(pStr); e == nil {
			p, _ = strconv.Atoi(pStr)
			filter.Priority = &p
		}
	}
	if pidStr := c.Query("project_id"); pidStr != "" {
		if pid, parseErr := uuid.Parse(pidStr); parseErr == nil {
			filter.ProjectID = &pid
		}
	}
	if aidStr := c.Query("assignee_id"); aidStr != "" {
		if aid, parseErr := uuid.Parse(aidStr); parseErr == nil {
			filter.AssigneeID = &aid
		}
	}

	ctx := c.Request.Context()
	issues, total, err := h.issueSvc.ListIssues(ctx, filter)
	if err != nil {
		h.logger.Error("查询工单列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询工单列表失败"})
		return
	}

	list := make([]dto.IssueInfo, len(issues))
	for i, iss := range issues {
		list[i] = buildIssueInfo(iss)
	}

	c.JSON(http.StatusOK, dto.IssueListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// TransitionIssue 状态流转（POST /api/v1/issues/:id/transition）
func (h *IssueHandler) TransitionIssue(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工单ID格式"})
		return
	}

	var req dto.TransitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	issue, svcErr := h.issueSvc.TransitionIssue(ctx, id, req.Status, req.ResolvedBy)
	if svcErr != nil {
		if svcErr == application.ErrInvalidIssueStatus {
			c.JSON(http.StatusBadRequest, gin.H{"error": "非法的状态转换: " + svcErr.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "状态流转成功", "status": issue.Status})
}

// AddComment 添加评论（POST /api/v1/issues/:id/comments）
func (h *IssueHandler) AddComment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工单ID格式"})
		return
	}

	var req dto.CreateIssueCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	userIDVal, _ := c.Get("userID")
	var userID uuid.UUID
	if uidStr, ok := userIDVal.(string); ok {
		userID, _ = uuid.Parse(uidStr)
	} else {
		userID = uuid.New()
	}

	comment, svcErr := h.issueSvc.AddComment(ctx, id, userID, req.Content, req.IsInternal, req.ParentID)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "评论添加成功", "comment_id": comment.ID.String()})
}

// ListComments 评论列表（GET /api/v1/issues/:id/comments）
func (h *IssueHandler) ListComments(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工单ID格式"})
		return
	}

	ctx := c.Request.Context()
	comments, svcErr := h.issueSvc.ListComments(ctx, id)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	list := make([]dto.IssueCommentInfo, len(comments))
	for i, cm := range comments {
		list[i] = dto.IssueCommentInfo{
			ID:         cm.ID.String(),
			UserID:     cm.UserID.String(),
			Content:    cm.Content,
			IsInternal: cm.IsInternal,
			CreatedAt:  cm.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if cm.ParentID != nil {
			list[i].ParentID = cm.ParentID.String()
		}
	}

	c.JSON(http.StatusOK, gin.H{"list": list, "count": len(list)})
}

// AddAttachment 上传附件（POST /api/v1/issues/:id/attachments）
func (h *IssueHandler) AddAttachment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工单ID格式"})
		return
	}

	var req dto.UploadAttachmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	attachment, svcErr := h.issueSvc.AddAttachment(ctx, id, nil, req.FileName, req.FileSize, req.FileType, req.FileURL)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "附件上传成功", "attachment_id": attachment.ID.String()})
}

// ListAttachments 附件列表（GET /api/v1/issues/:id/attachments）
func (h *IssueHandler) ListAttachments(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工单ID格式"})
		return
	}

	ctx := c.Request.Context()
	attachments, svcErr := h.issueSvc.ListAttachments(ctx, id)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	list := make([]dto.AttachmentInfo, len(attachments))
	for i, a := range attachments {
		list[i] = dto.AttachmentInfo{
			ID:         a.ID.String(),
			FileName:   a.FileName,
			FileSize:   a.FileSize,
			FileType:   a.FileType,
			FileURL:    a.FileURL,
			UploadedBy: a.UploadedBy.String(),
			CreatedAt:  a.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"list": list, "count": len(list)})
}

// GetSLAInfo SLA信息（GET /api/v1/issues/:id/sla）
func (h *IssueHandler) GetSLAInfo(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工单ID格式"})
		return
	}

	ctx := c.Request.Context()
	issue, getErr := h.issueSvc.GetIssue(ctx, id)
	if getErr != nil || issue == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "工单不存在"})
		return
	}

	slaResult, slaErr := h.issueSvc.GetSLAInfo(ctx, issue)
	if slaErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取SLA信息失败"})
		return
	}

	resp := dto.SLAInfo{
		ResponseSLA:       slaResult.ResponseSLA,
		ResolveSLA:        slaResult.ResolveSLA,
		BusinessHoursOnly: slaResult.BusinessHoursOnly,
	}

	if !slaResult.ResponseDueDate.IsZero() {
		resp.ResponseDueDate = slaResult.ResponseDueDate.Format(time.RFC3339)
	}
	if !slaResult.SLADueDate.IsZero() {
		resp.SLADueDate = slaResult.SLADueDate.Format(time.RFC3339)
	}
	resp.IsOverdue = slaResult.IsOverdue
	resp.ResponseOverdue = slaResult.ResponseOverdue
	resp.RemainingMinutes = slaResult.RemainingMinutes

	c.JSON(http.StatusOK, resp)
}

// RateSatisfaction 满意度评价（POST /api/v1/issues/:id/satisfaction）
func (h *IssueHandler) RateSatisfaction(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的工单ID格式"})
		return
	}

	var req dto.SatisfactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	svcErr := h.issueSvc.RateSatisfaction(ctx, id, req.Score)
	if svcErr != nil {
		if svcErr == application.ErrInvalidSatisfaction {
			c.JSON(http.StatusBadRequest, gin.H{"error": svcErr.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "满意度评价成功"})
}

// MyIssues 我的工单（GET /api/v1/issues/my）
func (h *IssueHandler) MyIssues(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	userIDVal, _ := c.Get("userID")
	var userID uuid.UUID
	if uidStr, ok := userIDVal.(string); ok {
		userID, _ = uuid.Parse(uidStr)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	ctx := c.Request.Context()
	issues, total, err := h.issueSvc.ListMyIssues(ctx, userID, userID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询我的工单失败"})
		return
	}

	list := make([]dto.IssueInfo, len(issues))
	for i, iss := range issues {
		list[i] = buildIssueInfo(iss)
	}

	c.JSON(http.StatusOK, dto.IssueListResponse{List: list, Total: total, Page: page, Size: size})
}

// ==================== 辅助方法 ====================

func (h *IssueHandler) buildIssueDetailResponse(issue *entity.Issue) dto.IssueDetailResponse {
	info := buildIssueInfo(issue)

	var comments []dto.IssueCommentInfo
	if issue.Comments != nil {
		comments = make([]dto.IssueCommentInfo, len(issue.Comments))
		for i, cm := range issue.Comments {
			comments[i] = dto.IssueCommentInfo{
				ID:         cm.ID.String(),
				UserID:     cm.UserID.String(),
				Content:    cm.Content,
				IsInternal: cm.IsInternal,
				CreatedAt:  cm.CreatedAt.Format("2006-01-02 15:04:05"),
			}
			if cm.ParentID != nil {
				comments[i].ParentID = cm.ParentID.String()
			}
		}
	}

	var attachments []dto.AttachmentInfo
	if issue.Attachments != nil {
		attachments = make([]dto.AttachmentInfo, len(issue.Attachments))
		for i, a := range issue.Attachments {
			attachments[i] = dto.AttachmentInfo{
				ID:         a.ID.String(),
				FileName:   a.FileName,
				FileSize:   a.FileSize,
				FileType:   a.FileType,
				FileURL:    a.FileURL,
				UploadedBy: a.UploadedBy.String(),
				CreatedAt:  a.CreatedAt.Format("2006-01-02 15:04:05"),
			}
		}
	}

	return dto.IssueDetailResponse{
		IssueInfo:   info,
		Comments:    comments,
		Attachments: attachments,
	}
}

func buildIssueInfo(iss *entity.Issue) dto.IssueInfo {
	info := dto.IssueInfo{
		ID:           iss.ID.String(),
		Title:        iss.Title,
		Description:  iss.Description,
		Type:         iss.Type,
		ReporterID:   iss.ReporterID.String(),
		Status:       iss.Status,
		Priority:     iss.Priority,
		Severity:     iss.Severity,
		Source:       iss.Source,
		Satisfaction: iss.Satisfaction,
		Resolution:   iss.Resolution,
		CreatedAt:    iss.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    iss.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if iss.ProjectID != nil {
		info.ProjectID = iss.ProjectID.String()
	}
	if iss.ProductID != nil {
		info.ProductID = iss.ProductID.String()
	}
	if iss.AssigneeID != nil {
		info.AssigneeID = iss.AssigneeID.String()
	}
	if iss.TemplateID != nil {
		info.TemplateID = iss.TemplateID.String()
	}
	if iss.SLADueDate != nil {
		info.SLADueDate = iss.SLADueDate.Format(time.RFC3339)
	}
	if iss.ResponseDueDate != nil {
		info.ResponseDueDate = iss.ResponseDueDate.Format(time.RFC3339)
	}
	if iss.ResolvedAt != nil {
		info.ResolvedAt = iss.ResolvedAt.Format("2006-01-02 15:04:05")
	}
	if iss.ResolvedBy != nil {
		info.ResolvedBy = iss.ResolvedBy.String()
	}
	if iss.ClosedAt != nil {
		info.ClosedAt = iss.ClosedAt.Format("2006-01-02 15:04:05")
	}
	if iss.ClosedBy != nil {
		info.ClosedBy = iss.ClosedBy.String()
	}

	if iss.Tags != "" {
		var tags []string
		json.Unmarshal([]byte(iss.Tags), &tags)
		info.Tags = tags
	}

	return info
}
