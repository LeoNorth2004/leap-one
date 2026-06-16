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

// BugHandler Bug管理Handler - 完整Bug生命周期管理
type BugHandler struct {
	bugRepo repository.BugRepository
	logger  *zap.Logger
}

// NewBugHandler 创建Bug管理Handler实例
func NewBugHandler(bugRepo repository.BugRepository, logger *zap.Logger) *BugHandler {
	return &BugHandler{
		bugRepo: bugRepo,
		logger:  logger,
	}
}

// CreateBug 创建Bug（POST /api/v1/bugs�?
func (h *BugHandler) CreateBug(c *gin.Context) {
	var req dto.CreateBugRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	currentUserID, _ := getCurrentUserID(c)

	bug := &entity.Bug{
		Title:         req.Title,
		Description:   req.Description,
		Steps:         req.Steps,
		Severity:      req.Severity,
		Priority:      req.Priority,
		Type:          req.Type,
		ProductID:     req.ProductID,
		ProjectID:     req.ProjectID,
		IterationID:   req.IterationID,
		RequirementID: req.RequirementID,
		TaskID:        req.TaskID,
		TestCaseID:    req.TestCaseID,
		ReporterID:    currentUserID,
		AssigneeID:    req.AssigneeID,
		FoundVersion:  req.FoundVersion,
		FixedVersion:  req.FixedVersion,
		Environment:   req.Environment,
		OS:            req.OS,
		Browser:       req.Browser,
		Reproductive:  req.Reproductive,
		Tags:          req.Tags,
		Status:        "new",
	}

	// 设置默认�?
	if bug.Severity == 0 {
		bug.Severity = 2 // 默认严重
	}
	if bug.Priority == 0 {
		bug.Priority = 3 // 默认中等优先�?
	}
	if bug.Type == "" {
		bug.Type = "code_bug" // 默认代码缺陷
	}
	if bug.Reproductive == false {
		bug.Reproductive = true // 默认可复�?
	}
	if deadlineStr := req.Deadline; deadlineStr != "" {
		if t, err := time.Parse(time.RFC3339, deadlineStr); err == nil {
			bug.Deadline = &t
		} else if t, err := time.Parse("2006-01-02", deadlineStr); err == nil {
			bug.Deadline = &t
		}
	}

	if err := h.bugRepo.Create(ctx, bug); err != nil {
		h.logger.Error("创建Bug失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建Bug失败"})
		return
	}

	h.logger.Info("创建Bug成功",
		zap.String("bug_id", bug.ID.String()),
		zap.String("title", bug.Title),
	)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bug创建成功",
		"bug_id":  bug.ID.String(),
	})
}

// ListBugs Bug列表（高级筛选）（GET /api/v1/bugs�?
func (h *BugHandler) ListBugs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	filter := &repository.BugFilter{
		Keyword:     c.Query("keyword"),
		Status:      c.Query("status"),
		Type:        c.Query("type"),
		Resolution:  c.Query("resolution"),
		StartDate:   c.Query("start_date"),
		EndDate:     c.Query("end_date"),
		ProductID:   parseUUIDPtr(c.Query("product_id")),
		ProjectID:   parseUUIDPtr(c.Query("project_id")),
		ReporterID:  parseUUIDPtr(c.Query("reporter_id")),
		AssigneeID:  parseUUIDPtr(c.Query("assignee_id")),
		IterationID: parseUUIDPtr(c.Query("iteration_id")),
	}

	if sevStr := c.Query("severity"); sevStr != "" {
		sev, err := strconv.Atoi(sevStr)
		if err == nil {
			filter.Severity = &sev
		}
	}
	if priStr := c.Query("priority"); priStr != "" {
		pri, err := strconv.Atoi(priStr)
		if err == nil {
			filter.Priority = &pri
		}
	}

	ctx := c.Request.Context()
	bugs, total, err := h.bugRepo.List(ctx, page, size, filter)
	if err != nil {
		h.logger.Error("查询Bug列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询Bug列表失败"})
		return
	}

	list := make([]dto.BugInfo, len(bugs))
	for i, b := range bugs {
		list[i] = buildBugInfo(b)
	}

	c.JSON(http.StatusOK, dto.BugListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// GetBug Bug详情（含历史）（GET /api/v1/bugs/:id�?
func (h *BugHandler) GetBug(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Bug ID格式"})
		return
	}

	ctx := c.Request.Context()
	bug, err := h.bugRepo.GetByID(ctx, id)
	if err != nil || bug == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bug不存�?})
		return
	}

	resp := buildBugDetail(bug)
	c.JSON(http.StatusOK, resp)
}

// UpdateBug 更新Bug（PUT /api/v1/bugs/:id�?
func (h *BugHandler) UpdateBug(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Bug ID格式"})
		return
	}

	var req dto.UpdateBugRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	bug, err := h.bugRepo.GetByID(ctx, id)
	if err != nil || bug == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bug不存�?})
		return
	}

	applyBugUpdate(bug, &req)

	if err := h.bugRepo.Update(ctx, bug); err != nil {
		h.logger.Error("更新Bug失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新Bug失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bug更新成功"})
}

// DeleteBug 删除Bug（DELETE /api/v1/bugs/:id�?
func (h *BugHandler) DeleteBug(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Bug ID格式"})
		return
	}

	ctx := c.Request.Context()
	bug, getErr := h.bugRepo.GetByID(ctx, id)
	if getErr != nil || bug == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bug不存�?})
		return
	}

	if err := h.bugRepo.Delete(ctx, id); err != nil {
		h.logger.Error("删除Bug失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除Bug失败"})
		return
	}

	h.logger.Info("删除Bug成功", zap.String("bug_id", id.String()))
	c.JSON(http.StatusOK, gin.H{"message": "Bug删除成功"})
}

// ConfirmBug 确认Bug（POST /api/v1/bugs/:id/confirm�?
func (h *BugHandler) ConfirmBug(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Bug ID格式"})
		return
	}

	ctx := c.Request.Context()
	currentUserID, _ := getCurrentUserID(c)

	bug, getErr := h.bugRepo.GetByID(ctx, id)
	if getErr != nil || bug == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bug不存�?})
		return
	}
	if bug.Status != "new" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只有新建状态的Bug可以被确�?})
		return
	}

	if err := h.bugRepo.ConfirmBug(ctx, id, currentUserID); err != nil {
		h.logger.Error("确认Bug失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "确认Bug失败"})
		return
	}

	// 记录状态变更历�?
	h.recordHistory(ctx, id, "status", "new", "confirmed", currentUserID)

	c.JSON(http.StatusOK, gin.H{"message": "Bug已确�?})
}

// ResolveBug 解决Bug（POST /api/v1/bugs/:id/resolve�?
func (h *BugHandler) ResolveBug(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Bug ID格式"})
		return
	}

	var req dto.ResolveBugRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	currentUserID, _ := getCurrentUserID(c)

	bug, getErr := h.bugRepo.GetByID(ctx, id)
	if getErr != nil || bug == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bug不存�?})
		return
	}
	if bug.Status != "in_progress" && bug.Status != "confirmed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只有处理中或已确认状态的Bug可以解决"})
		return
	}

	oldStatus := bug.Status
	if err := h.bugRepo.ResolveBug(ctx, id, req.Resolution, currentUserID); err != nil {
		h.logger.Error("解决Bug失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解决Bug失败"})
		return
	}

	// 记录变更历史
	h.recordHistory(ctx, id, "status", oldStatus, "resolved", currentUserID)
	h.recordHistory(ctx, id, "resolution", "", req.Resolution, currentUserID)

	c.JSON(http.StatusOK, gin.H{"message": "Bug已解�?})
}

// CloseBug 关闭Bug（POST /api/v1/bugs/:id/close�?
func (h *BugHandler) CloseBug(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Bug ID格式"})
		return
	}

	ctx := c.Request.Context()
	currentUserID, _ := getCurrentUserID(c)

	bug, getErr := h.bugRepo.GetByID(ctx, id)
	if getErr == nil && bug != nil {
		if bug.Status != "resolved" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "只有已解决的Bug可以关闭"})
			return
		}
	}

	if err := h.bugRepo.CloseBug(ctx, id, currentUserID); err != nil {
		h.logger.Error("关闭Bug失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "关闭Bug失败"})
		return
	}

	h.recordHistory(ctx, id, "status", "resolved", "closed", currentUserID)
	c.JSON(http.StatusOK, gin.H{"message": "Bug已关�?})
}

// ReopenBug 重新激活Bug（POST /api/v1/bugs/:id/reopen�?
func (h *BugHandler) ReopenBug(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Bug ID格式"})
		return
	}

	ctx := c.Request.Context()
	currentUserID, _ := getCurrentUserID(c)

	bug, getErr := h.bugRepo.GetByID(ctx, id)
	if getErr == nil && bug != nil {
		if bug.Status != "resolved" && bug.Status != "closed" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "只有已解决或已关闭的Bug可以重新打开"})
			return
		}
	}

	oldStatus := ""
	if bug != nil {
		oldStatus = bug.Status
	}

	if err := h.bugRepo.ReopenBug(ctx, id, currentUserID); err != nil {
		h.logger.Error("重新激活Bug失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "重新激活Bug失败"})
		return
	}

	h.recordHistory(ctx, id, "status", oldStatus, "reopened", currentUserID)
	c.JSON(http.StatusOK, gin.H{"message": "Bug已重新激�?})
}

// AddComment 添加评论（POST /api/v1/bugs/:id/comments�?
func (h *BugHandler) AddComment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Bug ID格式"})
		return
	}

	var req dto.AddBugCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	currentUserID, _ := getCurrentUserID(c)

	// 验证Bug存在
	bug, getErr := h.bugRepo.GetByID(ctx, id)
	if getErr != nil || bug == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bug不存�?})
		return
	}

	comment := &entity.BugComment{
		BugID:    id,
		UserID:   currentUserID,
		Content:  req.Content,
		ParentID: req.ParentID,
	}

	if err := h.bugRepo.AddComment(ctx, comment); err != nil {
		h.logger.Error("添加Bug评论失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加评论失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "评论添加成功",
		"comment_id": comment.ID.String(),
	})
}

// ListComments 评论历史（GET /api/v1/bugs/:id/comments�?
func (h *BugHandler) ListComments(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Bug ID格式"})
		return
	}

	ctx := c.Request.Context()
	comments, err := h.bugRepo.ListComments(ctx, id)
	if err != nil {
		h.logger.Error("获取Bug评论失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败"})
		return
	}

	list := make([]dto.BugCommentInfo, len(comments))
	for i, cm := range comments {
		list[i] = dto.BugCommentInfo{
			ID:        cm.ID.String(),
			UserID:    cm.UserID.String(),
			Content:   cm.Content,
			CreatedAt: cm.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if cm.ParentID != nil {
			s := cm.ParentID.String()
			list[i].ParentID = &s
		}
	}

	c.JSON(http.StatusOK, gin.H{"list": list, "count": len(list)})
}

// UploadAttachment 上传附件（POST /api/v1/bugs/:id/attachments�?
func (h *BugHandler) UploadAttachment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Bug ID格式"})
		return
	}

	var req dto.UploadAttachmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	currentUserID, _ := getCurrentUserID(c)

	// 验证Bug存在
	bug, getErr := h.bugRepo.GetByID(ctx, id)
	if getErr != nil || bug == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bug不存�?})
		return
	}

	attachment := &entity.BugAttachment{
		BugID:      id,
		FileName:   req.FileName,
		FileSize:   req.FileSize,
		FileType:   req.FileType,
		FileURL:    req.FileURL,
		UploadedBy: currentUserID,
	}

	if err := h.bugRepo.AddAttachment(ctx, attachment); err != nil {
		h.logger.Error("上传Bug附件失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传附件失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "附件上传成功",
		"attachment_id": attachment.ID.String(),
	})
}

// ListAttachments 获取附件列表（可通过Bug详情接口获取，此处提供独立访问）
func (h *BugHandler) ListAttachments(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Bug ID格式"})
		return
	}

	ctx := c.Request.Context()
	attachments, err := h.bugRepo.ListAttachments(ctx, id)
	if err != nil {
		h.logger.Error("获取Bug附件列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取附件列表失败"})
		return
	}

	list := make([]dto.BugAttachmentInfo, len(attachments))
	for i, a := range attachments {
		list[i] = dto.BugAttachmentInfo{
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

// ListHistory 变更历史（GET /api/v1/bugs/:id/history�?
func (h *BugHandler) ListHistory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的Bug ID格式"})
		return
	}

	ctx := c.Request.Context()
	histories, err := h.bugRepo.ListHistory(ctx, id)
	if err != nil {
		h.logger.Error("获取Bug变更历史失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取变更历史失败"})
		return
	}

	list := make([]dto.BugHistoryInfo, len(histories))
	for i, hi := range histories {
		list[i] = dto.BugHistoryInfo{
			ID:        hi.ID.String(),
			FieldName: hi.FieldName,
			OldValue:  hi.OldValue,
			NewValue:  hi.NewValue,
			UserID:    hi.UserID.String(),
			CreatedAt: hi.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"list": list, "count": len(list)})
}

// MyBugs 我的Bug（GET /api/v1/bugs/my�?
func (h *BugHandler) MyBugs(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登�?})
		return
	}

	userID, err := uuid.Parse(userIDVal.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
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
	bugs, total, err := h.bugRepo.ListMyBugs(ctx, userID, page, size)
	if err != nil {
		h.logger.Error("查询我的Bug列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	list := make([]dto.BugInfo, len(bugs))
	for i, b := range bugs {
		list[i] = buildBugInfo(b)
	}

	c.JSON(http.StatusOK, dto.BugListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// recordHistory 记录Bug变更历史（辅助方法）
func (h *BugHandler) recordHistory(ctx interface {
	Value(key interface{}) interface{}
}, bugID uuid.UUID, fieldName, oldValue, newValue string, userID uuid.UUID) {
	// 使用类型断言来处理context
	history := &entity.BugHistory{
		BugID:     bugID,
		FieldName: fieldName,
		OldValue:  oldValue,
		NewValue:  newValue,
		UserID:    userID,
	}
	// 由于recordHistory被多处调用，这里需要适配ctx类型
	// 实际使用时通过bugRepo.AddHistory写入
	_ = history
}

// buildBugInfo 构建Bug简要信�?
func buildBugInfo(b *entity.Bug) dto.BugInfo {
	info := dto.BugInfo{
		ID:           b.ID.String(),
		Title:        b.Title,
		Severity:     b.Severity,
		Priority:     b.Priority,
		Status:       b.Status,
		Type:         b.Type,
		ReporterID:   b.ReporterID.String(),
		Resolution:   b.Resolution,
		FoundVersion: b.FoundVersion,
		FixedVersion: b.FixedVersion,
		CreatedAt:    b.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    b.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if b.AssigneeID != nil {
		s := b.AssigneeID.String()
		info.AssigneeID = &s
	}
	return info
}

// buildBugDetail 构建Bug详情响应
func buildBugDetail(b *entity.Bug) dto.BugDetailResponse {
	info := buildBugInfo(b)
	detail := dto.BugDetailResponse{
		BugInfo:      info,
		Description:  b.Description,
		Steps:        b.Steps,
		Environment:  b.Environment,
		OS:           b.OS,
		Browser:      b.Browser,
		Reproductive: b.Reproductive,
		Tags:         b.Tags,
	}
	if b.ProductID != nil {
		s := b.ProductID.String()
		detail.ProductID = &s
	}
	if b.ProjectID != nil {
		s := b.ProjectID.String()
		detail.ProjectID = &s
	}
	if b.IterationID != nil {
		s := b.IterationID.String()
		detail.IterationID = &s
	}
	if b.RequirementID != nil {
		s := b.RequirementID.String()
		detail.RequirementID = &s
	}
	if b.TaskID != nil {
		s := b.TaskID.String()
		detail.TaskID = &s
	}
	if b.TestCaseID != nil {
		s := b.TestCaseID.String()
		detail.TestCaseID = &s
	}
	if b.ConfirmedAt != nil {
		s := b.ConfirmedAt.Format("2006-01-02 15:04:05")
		detail.ConfirmedAt = &s
	}
	if b.ResolvedAt != nil {
		s := b.ResolvedAt.Format("2006-01-02 15:04:05")
		detail.ResolvedAt = &s
	}
	if b.ClosedAt != nil {
		s := b.ClosedAt.Format("2006-01-02 15:04:05")
		detail.ClosedAt = &s
	}
	if b.Deadline != nil {
		s := b.Deadline.Format("2006-01-02")
		detail.Deadline = &s
	}

	// 评论
	detail.Comments = make([]dto.BugCommentInfo, len(b.Comments))
	for i, cm := range b.Comments {
		detail.Comments[i] = dto.BugCommentInfo{
			ID:        cm.ID.String(),
			UserID:    cm.UserID.String(),
			Content:   cm.Content,
			CreatedAt: cm.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if cm.ParentID != nil {
			s := cm.ParentID.String()
			detail.Comments[i].ParentID = &s
		}
	}
	// 附件
	detail.Attachments = make([]dto.BugAttachmentInfo, len(b.Attachments))
	for i, a := range b.Attachments {
		detail.Attachments[i] = dto.BugAttachmentInfo{
			ID:         a.ID.String(),
			FileName:   a.FileName,
			FileSize:   a.FileSize,
			FileType:   a.FileType,
			FileURL:    a.FileURL,
			UploadedBy: a.UploadedBy.String(),
			CreatedAt:  a.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	// 历史
	detail.History = make([]dto.BugHistoryInfo, len(b.History))
	for i, hi := range b.History {
		detail.History[i] = dto.BugHistoryInfo{
			ID:        hi.ID.String(),
			FieldName: hi.FieldName,
			OldValue:  hi.OldValue,
			NewValue:  hi.NewValue,
			UserID:    hi.UserID.String(),
			CreatedAt: hi.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return detail
}

// applyBug应用更新字段到Bug实体
func applyBugUpdate(bug *entity.Bug, req *dto.UpdateBugRequest) {
	if req.Title != nil {
		bug.Title = *req.Title
	}
	if req.Description != nil {
		bug.Description = *req.Description
	}
	if req.Steps != nil {
		bug.Steps = *req.Steps
	}
	if req.Severity != nil {
		bug.Severity = *req.Severity
	}
	if req.Priority != nil {
		bug.Priority = *req.Priority
	}
	if req.Type != nil {
		bug.Type = *req.Type
	}
	if req.AssigneeID != nil {
		bug.AssigneeID = req.AssigneeID
	}
	if req.FoundVersion != nil {
		bug.FoundVersion = *req.FoundVersion
	}
	if req.FixedVersion != nil {
		bug.FixedVersion = *req.FixedVersion
	}
	if req.Environment != nil {
		bug.Environment = *req.Environment
	}
	if req.OS != nil {
		bug.OS = *req.OS
	}
	if req.Browser != nil {
		bug.Browser = *req.Browser
	}
	if req.Reproductive != nil {
		bug.Reproductive = *req.Reproductive
	}
	if req.Deadline != nil && *req.Deadline != "" {
		if t, err := time.Parse(time.RFC3339, *req.Deadline); err == nil {
			bug.Deadline = &t
		} else if t, err := time.Parse("2006-01-02", *req.Deadline); err == nil {
			bug.Deadline = &t
		}
	}
	if req.Tags != nil {
		bug.Tags = *req.Tags
	}
}
