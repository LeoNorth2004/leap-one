package application

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"leap-one/service-task/internal/domain/entity"
	"leap-one/service-task/internal/domain/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// 任务服务相关错误定义
var (
	ErrTaskNotFound      = errors.New("任务不存在")
	ErrInvalidTaskStatus = errors.New("无效的任务状态转换")
	ErrTaskAlreadyClosed = errors.New("任务已关闭，无法操作")
)

// 有效的任务状态流转规则
var validTaskTransitions = map[string][]string{
	"waiting":     {"in_progress", "cancelled"},
	"in_progress": {"done", "paused", "cancelled"},
	"paused":      {"in_progress", "cancelled"},
	"done":        {"in_progress", "closed"},
	"cancelled":   {"waiting"},
}

// TaskService 任务应用服务 - 协调任务相关的业务流程
type TaskService struct {
	taskRepo       repository.TaskRepository
	assignmentRepo repository.TaskAssignmentRepository
	commentRepo    repository.TaskCommentRepository
	attachmentRepo repository.TaskAttachmentRepository
	linkRepo       repository.TaskLinkRepository
	worklogRepo    repository.TaskWorklogRepository
	logger         *zap.Logger
}

// NewTaskService 创建任务应用服务实例
func NewTaskService(
	taskRepo repository.TaskRepository,
	assignmentRepo repository.TaskAssignmentRepository,
	commentRepo repository.TaskCommentRepository,
	attachmentRepo repository.TaskAttachmentRepository,
	linkRepo repository.TaskLinkRepository,
	worklogRepo repository.TaskWorklogRepository,
	logger *zap.Logger,
) *TaskService {
	return &TaskService{
		taskRepo:       taskRepo,
		assignmentRepo: assignmentRepo,
		commentRepo:    commentRepo,
		attachmentRepo: attachmentRepo,
		linkRepo:       linkRepo,
		worklogRepo:    worklogRepo,
		logger:         logger,
	}
}

// CreateTask 创建任务用例
func (s *TaskService) CreateTask(ctx context.Context, req *CreateTaskInput) (*entity.Task, error) {
	tagsJSON, _ := json.Marshal(req.Tags)

	task := &entity.Task{
		Title:          req.Title,
		Description:    req.Description,
		Type:           req.Type,
		ProjectID:      req.ProjectID,
		IterationID:    req.IterationID,
		RequirementID:  req.RequirementID,
		ParentID:       req.ParentID,
		AssigneeID:     req.AssigneeID,
		CreatorID:      req.CreatorID,
		Priority:       req.Priority,
		Severity:       req.Severity,
		StoryPoints:    req.StoryPoints,
		EstimatedHours: req.EstimatedHours,
		StartDate:      req.StartDate,
		DueDate:        req.DueDate,
		KanbanColumn:   req.KanbanColumn,
		Tags:           string(tagsJSON),
	}

	if task.Type == "" {
		task.Type = "task"
	}
	if task.Status == "" {
		task.Status = "waiting"
	}
	if task.Priority == 0 {
		task.Priority = 3
	}
	if task.Severity == 0 {
		task.Severity = 1
	}
	if task.KanbanColumn == "" {
		task.KanbanColumn = "todo"
	}
	if task.ParentID != nil {
		task.Type = "sub_task"
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		s.logger.Error("创建任务失败", zap.Error(err), zap.String("title", req.Title))
		return nil, errors.New("创建任务失败")
	}

	s.logger.Info("任务创建成功",
		zap.String("task_id", task.ID.String()),
		zap.String("title", task.Title),
	)
	return task, nil
}

// UpdateTask 更新任务用例
func (s *TaskService) UpdateTask(ctx context.Context, id uuid.UUID, req *UpdateTaskInput) (*entity.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil || task == nil {
		return nil, ErrTaskNotFound
	}

	if task.Status == "closed" || task.Status == "cancelled" {
		return nil, ErrTaskAlreadyClosed
	}

	applyTaskUpdate(task, req)

	tagsJSON, _ := json.Marshal(req.Tags)
	task.Tags = string(tagsJSON)

	if err := s.taskRepo.Update(ctx, task); err != nil {
		s.logger.Error("更新任务失败", zap.Error(err), zap.String("task_id", id.String()))
		return nil, errors.New("更新任务失败")
	}

	return task, nil
}

// DeleteTask 删除任务用例
func (s *TaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil || task == nil {
		return ErrTaskNotFound
	}

	if err := s.taskRepo.Delete(ctx, id); err != nil {
		s.logger.Error("删除任务失败", zap.Error(err), zap.String("task_id", id.String()))
		return errors.New("删除任务失败")
	}

	s.logger.Info("任务已删除", zap.String("task_id", id.String()), zap.String("title", task.Title))
	return nil
}

// ChangeTaskStatus 状态流转用例（含状态机校验）
func (s *TaskService) ChangeTaskStatus(ctx context.Context, id uuid.UUID, newStatus string) (*entity.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil || task == nil {
		return nil, ErrTaskNotFound
	}

	if !isValidTransition(task.Status, newStatus) {
		return nil, ErrInvalidTaskStatus
	}

	if err := s.taskRepo.UpdateStatus(ctx, id, newStatus); err != nil {
		s.logger.Error("更新任务状态失败", zap.Error(err),
			zap.String("task_id", id.String()),
			zap.String("from", task.Status),
			zap.String("to", newStatus),
		)
		return nil, errors.New("状态更新失败")
	}

	task.Status = newStatus
	s.logger.Info("任务状态变更",
		zap.String("task_id", id.String()),
		zap.String("old_status", task.Status),
		zap.String("new_status", newStatus),
	)

	return task, nil
}

// StartTask 开始任务用例
func (s *TaskService) StartTask(ctx context.Context, id uuid.UUID) (*entity.Task, error) {
	return s.ChangeTaskStatus(ctx, id, "in_progress")
}

// CompleteTask 完成任务用例
func (s *TaskService) CompleteTask(ctx context.Context, id uuid.UUID) (*entity.Task, error) {
	return s.ChangeTaskStatus(ctx, id, "done")
}

// AssignTask 分配任务用例（支持多人）
func (s *TaskService) AssignTask(ctx context.Context, taskID uuid.UUID, userIDs []string, role string) ([]*entity.TaskAssignment, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil || task == nil {
		return nil, ErrTaskNotFound
	}

	if role == "" {
		role = "assignee"
	}

	now := time.Now()
	var assignments []*entity.TaskAssignment

	for _, uidStr := range userIDs {
		uid, parseErr := uuid.Parse(uidStr)
		if parseErr != nil {
			continue
		}

		assignment := &entity.TaskAssignment{
			TaskID:     taskID,
			UserID:     uid,
			Role:       role,
			AssignedAt: now,
		}

		if assignErr := s.assignmentRepo.Create(ctx, assignment); assignErr != nil {
			s.logger.Warn("分配任务用户失败",
				zap.Error(assignErr),
				zap.String("user_id", uidStr),
			)
			continue
		}
		assignments = append(assignments, assignment)
	}

	s.logger.Info("任务分配完成",
		zap.String("task_id", taskID.String()),
		zap.Int("assigned_count", len(assignments)),
	)

	return assignments, nil
}

// RemoveAssignment 移除任务分配
func (s *TaskService) RemoveAssignment(ctx context.Context, taskID, userID uuid.UUID) error {
	return s.assignmentRepo.Delete(ctx, taskID, userID)
}

// AddComment 添加评论用例
func (s *TaskService) AddComment(ctx context.Context, taskID, userID uuid.UUID, content string, parentID *uuid.UUID) (*entity.TaskComment, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil || task == nil {
		return nil, ErrTaskNotFound
	}

	comment := &entity.TaskComment{
		TaskID:   taskID,
		UserID:   userID,
		Content:  content,
		ParentID: parentID,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, errors.New("添加评论失败")
	}

	return comment, nil
}

// ListComments 获取评论列表
func (s *TaskService) ListComments(ctx context.Context, taskID uuid.UUID) ([]*entity.TaskComment, error) {
	return s.commentRepo.ListByTaskID(ctx, taskID)
}

// DeleteComment 删除评论
func (s *TaskService) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	return s.commentRepo.Delete(ctx, commentID)
}

// AddAttachment 添加附件用例
func (s *TaskService) AddAttachment(ctx context.Context, taskID, uploadedBy uuid.UUID, fileName string, fileSize int64, fileType, fileURL string) (*entity.TaskAttachment, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil || task == nil {
		return nil, ErrTaskNotFound
	}

	attachment := &entity.TaskAttachment{
		TaskID:     taskID,
		FileName:   fileName,
		FileSize:   fileSize,
		FileType:   fileType,
		FileURL:    fileURL,
		UploadedBy: uploadedBy,
	}

	if err := s.attachmentRepo.Create(ctx, attachment); err != nil {
		return nil, errors.New("添加附件失败")
	}

	return attachment, nil
}

// ListAttachments 获取附件列表
func (s *TaskService) ListAttachments(ctx context.Context, taskID uuid.UUID) ([]*entity.TaskAttachment, error) {
	return s.attachmentRepo.ListByTaskID(ctx, taskID)
}

// DeleteAttachment 删除附件
func (s *TaskService) DeleteAttachment(ctx context.Context, attachmentID uuid.UUID) error {
	return s.attachmentRepo.Delete(ctx, attachmentID)
}

// AddWorklog 添加工作日志用例
func (s *TaskService) AddWorklog(ctx context.Context, taskID, userID uuid.UUID, spentHours float64, workDate time.Time, summary string) (*entity.TaskWorklog, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil || task == nil {
		return nil, ErrTaskNotFound
	}

	worklog := &entity.TaskWorklog{
		TaskID:     taskID,
		UserID:     userID,
		SpentHours: spentHours,
		WorkDate:   workDate,
		Summary:    summary,
	}

	if logErr := s.worklogRepo.Create(ctx, worklog); logErr != nil {
		return nil, errors.New("添加工作日志失败")
	}

	totalHours, _ := s.worklogRepo.GetTotalSpentHours(ctx, taskID)
	s.logger.Info("工作日志已添加",
		zap.String("task_id", taskID.String()),
		zap.Float64("spent_hours", spentHours),
		zap.Float64("total_hours", totalHours),
	)

	return worklog, nil
}

// ListWorklogs 获取工作日志列表
func (s *TaskService) ListWorklogs(ctx context.Context, taskID uuid.UUID) ([]*entity.TaskWorklog, error) {
	return s.worklogRepo.ListByTaskID(ctx, taskID)
}

// CreateSubTask 创建子任务用例
func (s *TaskService) CreateSubTask(ctx context.Context, parentID, creatorID uuid.UUID, title, description string, assigneeID *uuid.UUID, priority int) (*entity.Task, error) {
	parent, err := s.taskRepo.GetByID(ctx, parentID)
	if err != nil || parent == nil {
		return nil, ErrTaskNotFound
	}

	subTask := &entity.Task{
		Title:        title,
		Description:  description,
		Type:         "sub_task",
		ParentID:     &parentID,
		CreatorID:    creatorID,
		AssigneeID:   assigneeID,
		Priority:     priority,
		Status:       "waiting",
		KanbanColumn: "todo",
	}

	if priority == 0 {
		subTask.Priority = parent.Priority
	}

	if createErr := s.taskRepo.Create(ctx, subTask); createErr != nil {
		return nil, errors.New("创建子任务失败")
	}

	s.logger.Info("子任务已创建",
		zap.String("parent_id", parentID.String()),
		zap.String("subtask_id", subTask.ID.String()),
	)

	return subTask, nil
}

// ListSubTasks 获取子任务列表
func (s *TaskService) ListSubTasks(ctx context.Context, parentID uuid.UUID) ([]*entity.Task, int64, error) {
	filter := &repository.TaskFilter{
		Page:      1,
		PageSize:  100,
		ParentID:  &parentID,
		IsSubTask: true,
	}
	return s.taskRepo.List(ctx, filter)
}

// AddTaskLink 创建任务关联用例
func (s *TaskService) AddTaskLink(ctx context.Context, sourceTaskID, targetTaskID uuid.UUID, linkType, name string) (*entity.TaskLink, error) {
	source, srcErr := s.taskRepo.GetByID(ctx, sourceTaskID)
	if srcErr != nil || source == nil {
		return nil, ErrTaskNotFound
	}

	target, tgtErr := s.taskRepo.GetByID(ctx, targetTaskID)
	if tgtErr != nil || target == nil {
		return nil, errors.New("目标任务不存在")
	}

	if sourceTaskID == targetTaskID {
		return nil, errors.New("不能关联自己")
	}

	link := &entity.TaskLink{
		SourceTaskID: sourceTaskID,
		TargetTaskID: targetTaskID,
		LinkType:     linkType,
		Name:         name,
	}

	if err := s.linkRepo.Create(ctx, link); err != nil {
		return nil, errors.New("创建任务关联失败")
	}

	return link, nil
}

// ListTaskLinks 获取任务关联列表
func (s *TaskService) ListTaskLinks(ctx context.Context, taskID uuid.UUID) ([]*entity.TaskLink, error) {
	return s.linkRepo.ListByTaskID(ctx, taskID)
}

// DeleteTaskLink 删除任务关联
func (s *TaskService) DeleteTaskLink(ctx context.Context, linkID uuid.UUID) error {
	return s.linkRepo.Delete(ctx, linkID)
}

// ListMyTasks 获取我的任务列表
func (s *TaskService) ListMyTasks(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*entity.Task, int64, error) {
	return s.taskRepo.ListByAssigneeID(ctx, userID, page, pageSize)
}

// GetTask 获取任务详情
func (s *TaskService) GetTask(ctx context.Context, id uuid.UUID) (*entity.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil || task == nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

// ListTasks 分页查询任务列表
func (s *TaskService) ListTasks(ctx context.Context, filter *repository.TaskFilter) ([]*entity.Task, int64, error) {
	return s.taskRepo.List(ctx, filter)
}

// isValidTransition 校验状态流转是否合法
func isValidTransition(fromStatus, toStatus string) bool {
	validTos, ok := validTaskTransitions[fromStatus]
	if !ok {
		return false
	}
	for _, v := range validTos {
		if v == toStatus {
			return true
		}
	}
	return false
}

// ==================== 输入结构体 ====================

// CreateTaskInput 创建任务输入
type CreateTaskInput struct {
	Title          string
	Description    string
	Type           string
	ProjectID      *uuid.UUID
	IterationID    *uuid.UUID
	RequirementID  *uuid.UUID
	ParentID       *uuid.UUID
	AssigneeID     *uuid.UUID
	CreatorID      uuid.UUID
	Priority       int
	Severity       int
	StoryPoints    *float64
	EstimatedHours *float64
	StartDate      *time.Time
	DueDate        *time.Time
	KanbanColumn   string
	Tags           []string
}

// UpdateTaskInput 更新任务输入
type UpdateTaskInput struct {
	Title          *string
	Description    *string
	Type           *string
	ProjectID      *uuid.UUID
	IterationID    *uuid.UUID
	RequirementID  *uuid.UUID
	AssigneeID     *uuid.UUID
	Priority       *int
	Severity       *int
	StoryPoints    *float64
	EstimatedHours *float64
	ActualHours    *float64
	RemainingHours *float64
	StartDate      *time.Time
	DueDate        *time.Time
	KanbanColumn   *string
	Tags           []string
}

// applyTaskUpdate 将更新输入应用到实体
func applyTaskUpdate(task *entity.Task, input *UpdateTaskInput) {
	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Type != nil {
		task.Type = *input.Type
	}
	if input.ProjectID != nil {
		task.ProjectID = input.ProjectID
	}
	if input.IterationID != nil {
		task.IterationID = input.IterationID
	}
	if input.RequirementID != nil {
		task.RequirementID = input.RequirementID
	}
	if input.AssigneeID != nil {
		task.AssigneeID = input.AssigneeID
	}
	if input.Priority != nil {
		task.Priority = *input.Priority
	}
	if input.Severity != nil {
		task.Severity = *input.Severity
	}
	if input.StoryPoints != nil {
		task.StoryPoints = input.StoryPoints
	}
	if input.EstimatedHours != nil {
		task.EstimatedHours = input.EstimatedHours
	}
	if input.ActualHours != nil {
		task.ActualHours = input.ActualHours
	}
	if input.RemainingHours != nil {
		task.RemainingHours = input.RemainingHours
	}
	if input.StartDate != nil {
		task.StartDate = input.StartDate
	}
	if input.DueDate != nil {
		task.DueDate = input.DueDate
	}
	if input.KanbanColumn != nil {
		task.KanbanColumn = *input.KanbanColumn
	}
}
