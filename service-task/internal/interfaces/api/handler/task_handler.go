package handler

import (
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

// TaskHandler 任务管理Handler
type TaskHandler struct {
	taskSvc *application.TaskService
	logger  *zap.Logger
}

// NewTaskHandler 创建任务管理Handler实例
func NewTaskHandler(taskSvc *application.TaskService, logger *zap.Logger) *TaskHandler {
	return &TaskHandler{
		taskSvc: taskSvc,
		logger:  logger,
	}
}

// CreateTask 创建任务（POST /api/v1/tasks）
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	userIDVal, _ := c.Get("userID")
	var creatorID uuid.UUID
	if uidStr, ok := userIDVal.(string); ok {
		if parsed, err := uuid.Parse(uidStr); err == nil {
			creatorID = parsed
		} else {
			creatorID = uuid.New()
		}
	} else {
		creatorID = uuid.New()
	}

	input := &application.CreateTaskInput{
		Title:          req.Title,
		Description:    req.Description,
		Type:           req.Type,
		ProjectID:      req.ProjectID,
		IterationID:    req.IterationID,
		RequirementID:  req.RequirementID,
		ParentID:       req.ParentID,
		AssigneeID:     req.AssigneeID,
		CreatorID:      creatorID,
		Priority:       req.Priority,
		Severity:       req.Severity,
		StoryPoints:    req.StoryPoints,
		EstimatedHours: req.EstimatedHours,
		StartDate:      req.StartDate,
		DueDate:        req.DueDate,
		KanbanColumn:   req.KanbanColumn,
		Tags:           req.Tags,
	}
	if input.Priority == 0 {
		input.Priority = 3
	}
	if input.Severity == 0 {
		input.Severity = 1
	}

	task, err := h.taskSvc.CreateTask(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "任务创建成功",
		"task_id": task.ID.String(),
	})
}

// GetTask 获取任务详情（GET /api/v1/tasks/:id）
func (h *TaskHandler) GetTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}

	ctx := c.Request.Context()
	task, err := h.taskSvc.GetTask(ctx, id)
	if err != nil || task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	resp := h.buildTaskDetailResponse(task)
	c.JSON(http.StatusOK, resp)
}

// UpdateTask 更新任务（PUT /api/v1/tasks/:id）
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	input := &application.UpdateTaskInput{
		Title:          req.Title,
		Description:    req.Description,
		Type:           req.Type,
		ProjectID:      req.ProjectID,
		IterationID:    req.IterationID,
		RequirementID:  req.RequirementID,
		AssigneeID:     req.AssigneeID,
		Priority:       req.Priority,
		Severity:       req.Severity,
		StoryPoints:    req.StoryPoints,
		EstimatedHours: req.EstimatedHours,
		ActualHours:    req.ActualHours,
		RemainingHours: req.RemainingHours,
		StartDate:      req.StartDate,
		DueDate:        req.DueDate,
		KanbanColumn:   req.KanbanColumn,
		Tags:           req.Tags,
	}

	task, svcErr := h.taskSvc.UpdateTask(ctx, id, input)
	if svcErr != nil {
		if svcErr == application.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务更新成功", "task_id": task.ID.String()})
}

// DeleteTask 删除任务（DELETE /api/v1/tasks/:id）
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.taskSvc.DeleteTask(ctx, id); err != nil {
		if err == application.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务删除成功"})
}

// ListTasks 任务列表（GET /api/v1/tasks）
func (h *TaskHandler) ListTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	filter := &repository.TaskFilter{
		Page:      page,
		PageSize:  size,
		Keyword:   c.Query("keyword"),
		Status:    c.Query("status"),
		Type:      c.Query("type"),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "DESC"),
	}

	if pStr := c.Query("priority"); pStr != "" {
		var p int
		if _, err := strconv.Atoi(pStr); err == nil {
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
	tasks, total, err := h.taskSvc.ListTasks(ctx, filter)
	if err != nil {
		h.logger.Error("查询任务列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询任务列表失败"})
		return
	}

	list := make([]dto.TaskInfo, len(tasks))
	for i, t := range tasks {
		list[i] = buildTaskInfo(t)
	}

	c.JSON(http.StatusOK, dto.TaskListResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// AssignTask 分配任务（POST /api/v1/tasks/:id/assign）
func (h *TaskHandler) AssignTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}

	var req dto.AssignTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	assignments, svcErr := h.taskSvc.AssignTask(ctx, id, req.UserIDs, req.Role)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "任务分配成功",
		"assigned_count": len(assignments),
	})
}

// RemoveAssignment 移除分配（DELETE /api/v1/tasks/:id/assignments/:uid）
func (h *TaskHandler) RemoveAssignment(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}
	userID, err := uuid.Parse(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.taskSvc.RemoveAssignment(ctx, taskID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "移除分配失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "分配已移除"})
}

// AddComment 添加评论（POST /api/v1/tasks/:id/comments）
func (h *TaskHandler) AddComment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}

	var req dto.CreateCommentRequest
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

	comment, svcErr := h.taskSvc.AddComment(ctx, id, userID, req.Content, req.ParentID)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "评论添加成功", "comment_id": comment.ID.String()})
}

// ListComments 评论列表（GET /api/v1/tasks/:id/comments）
func (h *TaskHandler) ListComments(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}

	ctx := c.Request.Context()
	comments, svcErr := h.taskSvc.ListComments(ctx, id)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	list := make([]dto.CommentInfo, len(comments))
	for i, cm := range comments {
		list[i] = dto.CommentInfo{
			ID:        cm.ID.String(),
			UserID:    cm.UserID.String(),
			Content:   cm.Content,
			CreatedAt: cm.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if cm.ParentID != nil {
			list[i].ParentID = cm.ParentID.String()
		}
	}

	c.JSON(http.StatusOK, gin.H{"list": list, "count": len(list)})
}

// DeleteComment 删除评论（DELETE /api/v1/tasks/:id/comments/:cid）
func (h *TaskHandler) DeleteComment(c *gin.Context) {
	commentID, err := uuid.Parse(c.Param("cid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.taskSvc.DeleteComment(ctx, commentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除评论失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评论已删除"})
}

// AddAttachment 上传附件（POST /api/v1/tasks/:id/attachments）
func (h *TaskHandler) AddAttachment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}

	var req dto.UploadAttachmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	userIDVal, _ := c.Get("userID")
	var uploadedBy uuid.UUID
	if uidStr, ok := userIDVal.(string); ok {
		uploadedBy, _ = uuid.Parse(uidStr)
	} else {
		uploadedBy = uuid.New()
	}

	attachment, svcErr := h.taskSvc.AddAttachment(ctx, id, uploadedBy, req.FileName, req.FileSize, req.FileType, req.FileURL)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "附件上传成功", "attachment_id": attachment.ID.String()})
}

// ListAttachments 附件列表（GET /api/v1/tasks/:id/attachments）
func (h *TaskHandler) ListAttachments(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}

	ctx := c.Request.Context()
	attachments, svcErr := h.taskSvc.ListAttachments(ctx, id)
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

// DeleteAttachment 删除附件（DELETE /api/v1/tasks/:id/attachments/:aid）
func (h *TaskHandler) DeleteAttachment(c *gin.Context) {
	attachmentID, err := uuid.Parse(c.Param("aid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的附件ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.taskSvc.DeleteAttachment(ctx, attachmentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除附件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "附件已删除"})
}

// AddWorklog 添加工作日志（POST /api/v1/tasks/:id/worklogs）
func (h *TaskHandler) AddWorklog(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}

	var req dto.CreateWorklogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	workDate, parseErr := time.Parse(time.RFC3339, req.WorkDate)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "工作日期格式无效，请使用RFC3339格式"})
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

	worklog, svcErr := h.taskSvc.AddWorklog(ctx, id, userID, req.SpentHours, workDate, req.Summary)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "工作日志添加成功", "worklog_id": worklog.ID.String()})
}

// ListWorklogs 工作日志列表（GET /api/v1/tasks/:id/worklogs）
func (h *TaskHandler) ListWorklogs(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}

	ctx := c.Request.Context()
	worklogs, svcErr := h.taskSvc.ListWorklogs(ctx, id)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	list := make([]dto.WorklogInfo, len(worklogs))
	for i, w := range worklogs {
		list[i] = dto.WorklogInfo{
			ID:         w.ID.String(),
			UserID:     w.UserID.String(),
			SpentHours: w.SpentHours,
			WorkDate:   w.WorkDate.Format("2006-01-02"),
			Summary:    w.Summary,
			CreatedAt:  w.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"list": list, "count": len(list)})
}

// CreateSubTask 创建子任务（POST /api/v1/tasks/:id/subtasks）
func (h *TaskHandler) CreateSubTask(c *gin.Context) {
	parentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}

	var req dto.CreateSubTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	userIDVal, _ := c.Get("userID")
	var creatorID uuid.UUID
	if uidStr, ok := userIDVal.(string); ok {
		creatorID, _ = uuid.Parse(uidStr)
	} else {
		creatorID = uuid.New()
	}

	subTask, svcErr := h.taskSvc.CreateSubTask(ctx, parentID, creatorID, req.Title, req.Description, req.AssigneeID, req.Priority)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "子任务创建成功", "subtask_id": subTask.ID.String()})
}

// ListSubTasks 子任务列表（GET /api/v1/tasks/:id/subtasks）
func (h *TaskHandler) ListSubTasks(c *gin.Context) {
	parentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}

	ctx := c.Request.Context()
	subTasks, total, svcErr := h.taskSvc.ListSubTasks(ctx, parentID)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	list := make([]dto.TaskInfo, len(subTasks))
	for i, t := range subTasks {
		list[i] = buildTaskInfo(t)
	}

	c.JSON(http.StatusOK, gin.H{"list": list, "total": total})
}

// AddTaskLink 创建任务关联（POST /api/v1/tasks/:id/links）
func (h *TaskHandler) AddTaskLink(c *gin.Context) {
	sourceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}

	var req dto.CreateTaskLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	link, svcErr := h.taskSvc.AddTaskLink(ctx, sourceID, req.TargetTaskID, req.LinkType, req.Name)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "任务关联创建成功", "link_id": link.ID.String()})
}

// ListTaskLinks 任务关联列表（GET /api/v1/tasks/:id/links）
func (h *TaskHandler) ListTaskLinks(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}

	ctx := c.Request.Context()
	links, svcErr := h.taskSvc.ListTaskLinks(ctx, taskID)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	list := make([]dto.TaskLinkInfo, len(links))
	for i, l := range links {
		list[i] = dto.TaskLinkInfo{
			ID:           l.ID.String(),
			SourceTaskID: l.SourceTaskID.String(),
			TargetTaskID: l.TargetTaskID.String(),
			LinkType:     l.LinkType,
			Name:         l.Name,
			CreatedAt:    l.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"list": list, "count": len(list)})
}

// DeleteTaskLink 删除任务关联（DELETE /api/v1/tasks/:id/links/:lid）
func (h *TaskHandler) DeleteTaskLink(c *gin.Context) {
	linkID, err := uuid.Parse(c.Param("lid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的关联ID格式"})
		return
	}

	ctx := c.Request.Context()
	if err := h.taskSvc.DeleteTaskLink(ctx, linkID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除关联失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "关联已删除"})
}

// UpdateStatus 更改状态（PUT /api/v1/tasks/:id/status）
func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}

	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	ctx := c.Request.Context()
	task, svcErr := h.taskSvc.ChangeTaskStatus(ctx, id, req.Status)
	if svcErr != nil {
		if svcErr == application.ErrInvalidTaskStatus {
			c.JSON(http.StatusBadRequest, gin.H{"error": "非法的状态转换: " + svcErr.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "状态更新成功", "status": task.Status})
}

// StartTask 开始任务（POST /api/v1/tasks/:id/start）
func (h *TaskHandler) StartTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}

	ctx := c.Request.Context()
	task, svcErr := h.taskSvc.StartTask(ctx, id)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务已开始", "status": task.Status})
}

// CompleteTask 完成任务（POST /api/v1/tasks/:id/complete）
func (h *TaskHandler) CompleteTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID格式"})
		return
	}

	ctx := c.Request.Context()
	task, svcErr := h.taskSvc.CompleteTask(ctx, id)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": svcErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务已完成", "status": task.Status})
}

// MyTasks 我的任务（GET /api/v1/tasks/my）
func (h *TaskHandler) MyTasks(c *gin.Context) {
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
	tasks, total, err := h.taskSvc.ListMyTasks(ctx, userID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询我的任务失败"})
		return
	}

	list := make([]dto.TaskInfo, len(tasks))
	for i, t := range tasks {
		list[i] = buildTaskInfo(t)
	}

	c.JSON(http.StatusOK, dto.TaskListResponse{List: list, Total: total, Page: page, Size: size})
}

// ==================== 辅助方法 ====================

// buildTaskDetailResponse 构建任务详情响应
func (h *TaskHandler) buildTaskDetailResponse(task *entity.Task) dto.TaskDetailResponse {
	info := buildTaskInfo(task)

	var assignees []dto.AssignmentInfo
	if task.Assignees != nil {
		assignees = make([]dto.AssignmentInfo, len(task.Assignees))
		for i, a := range task.Assignees {
			assignees[i] = dto.AssignmentInfo{
				ID:         a.ID.String(),
				UserID:     a.UserID.String(),
				Role:       a.Role,
				AssignedAt: a.AssignedAt.Format("2006-01-02 15:04:05"),
			}
		}
	}

	var comments []dto.CommentInfo
	if task.Comments != nil {
		comments = make([]dto.CommentInfo, len(task.Comments))
		for i, cm := range task.Comments {
			comments[i] = dto.CommentInfo{
				ID:        cm.ID.String(),
				UserID:    cm.UserID.String(),
				Content:   cm.Content,
				CreatedAt: cm.CreatedAt.Format("2006-01-02 15:04:05"),
			}
			if cm.ParentID != nil {
				comments[i].ParentID = cm.ParentID.String()
			}
		}
	}

	var attachments []dto.AttachmentInfo
	if task.Attachments != nil {
		attachments = make([]dto.AttachmentInfo, len(task.Attachments))
		for i, a := range task.Attachments {
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

	return dto.TaskDetailResponse{
		TaskInfo:    info,
		Assignees:   assignees,
		Comments:    comments,
		Attachments: attachments,
	}
}

// buildTaskInfo 构建任务基本信息
func buildTaskInfo(t *entity.Task) dto.TaskInfo {
	info := dto.TaskInfo{
		ID:             t.ID.String(),
		Title:          t.Title,
		Description:    t.Description,
		Type:           t.Type,
		CreatorID:      t.CreatorID.String(),
		Status:         t.Status,
		Priority:       t.Priority,
		Severity:       t.Severity,
		KanbanColumn:   t.KanbanColumn,
		StoryPoints:    t.StoryPoints,
		EstimatedHours: t.EstimatedHours,
		ActualHours:    t.ActualHours,
		RemainingHours: t.RemainingHours,
		CreatedAt:      t.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      t.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if t.ProjectID != nil {
		info.ProjectID = t.ProjectID.String()
	}
	if t.IterationID != nil {
		info.IterationID = t.IterationID.String()
	}
	if t.RequirementID != nil {
		info.RequirementID = t.RequirementID.String()
	}
	if t.ParentID != nil {
		info.ParentID = t.ParentID.String()
	}
	if t.AssigneeID != nil {
		info.AssigneeID = t.AssigneeID.String()
	}
	if t.StartDate != nil {
		info.StartDate = t.StartDate.Format("2006-01-02")
	}
	if t.DueDate != nil {
		info.DueDate = t.DueDate.Format("2006-01-02")
	}
	if t.FinishedDate != nil {
		info.FinishedDate = t.FinishedDate.Format("2006-01-02 15:04:05")
	}

	return info
}
