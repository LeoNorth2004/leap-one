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

// 工单服务相关错误定义
var (
	ErrIssueNotFound       = errors.New("工单不存在")
	ErrInvalidIssueStatus  = errors.New("无效的工单状态转换")
	ErrIssueAlreadyClosed  = errors.New("工单已关闭，无法操作")
	ErrInvalidSatisfaction = errors.New("无效的满意度评分")
	ErrWorkflowNotFound    = errors.New("工作流不存在")
)

// IssueService 工单应用服务 - 协调工单相关的业务流程（含SLA计算、状态机等）
type IssueService struct {
	issueRepo      repository.IssueRepository
	commentRepo    repository.IssueCommentRepository
	attachmentRepo repository.IssueAttachmentRepository
	templateRepo   repository.IssueTemplateRepository
	workflowRepo   repository.IssueWorkflowRepository
	slaConfigRepo  repository.IssueSLAConfigRepository
	logger         *zap.Logger
}

// NewIssueService 创建工单应用服务实例
func NewIssueService(
	issueRepo repository.IssueRepository,
	commentRepo repository.IssueCommentRepository,
	attachmentRepo repository.IssueAttachmentRepository,
	templateRepo repository.IssueTemplateRepository,
	workflowRepo repository.IssueWorkflowRepository,
	slaConfigRepo repository.IssueSLAConfigRepository,
	logger *zap.Logger,
) *IssueService {
	return &IssueService{
		issueRepo:      issueRepo,
		commentRepo:    commentRepo,
		attachmentRepo: attachmentRepo,
		templateRepo:   templateRepo,
		workflowRepo:   workflowRepo,
		slaConfigRepo:  slaConfigRepo,
		logger:         logger,
	}
}

// CreateIssue 创建工单用例
func (s *IssueService) CreateIssue(ctx context.Context, input *CreateIssueInput) (*entity.Issue, error) {
	tagsJSON, _ := json.Marshal(input.Tags)

	issue := &entity.Issue{
		Title:       input.Title,
		Description: input.Description,
		Type:        input.Type,
		ProjectID:   input.ProjectID,
		ProductID:   input.ProductID,
		ReporterID:  input.ReporterID,
		AssigneeID:  input.AssigneeID,
		Priority:    input.Priority,
		Severity:    input.Severity,
		Source:      input.Source,
		TemplateID:  input.TemplateID,
		Tags:        string(tagsJSON),
	}

	if issue.Type == "" {
		issue.Type = "bug"
	}
	if issue.Status == "" {
		issue.Status = "new"
	}
	if issue.Priority == 0 {
		issue.Priority = 3
	}
	if issue.Severity == 0 {
		issue.Severity = 2
	}
	if issue.Source == "" {
		issue.Source = "manual"
	}

	if err := s.issueRepo.Create(ctx, issue); err != nil {
		s.logger.Error("创建工单失败", zap.Error(err), zap.String("title", input.Title))
		return nil, errors.New("创建工单失败")
	}

	// 计算并设置SLA截止时间
	s.calculateAndSetSLADates(ctx, issue)

	s.logger.Info("工单创建成功",
		zap.String("issue_id", issue.ID.String()),
		zap.String("title", issue.Title),
	)
	return issue, nil
}

// GetIssue 获取工单详情
func (s *IssueService) GetIssue(ctx context.Context, id uuid.UUID) (*entity.Issue, error) {
	issue, err := s.issueRepo.GetByID(ctx, id)
	if err != nil || issue == nil {
		return nil, ErrIssueNotFound
	}
	return issue, nil
}

// UpdateIssue 更新工单用例
func (s *IssueService) UpdateIssue(ctx context.Context, id uuid.UUID, input *UpdateIssueInput) (*entity.Issue, error) {
	issue, err := s.issueRepo.GetByID(ctx, id)
	if err != nil || issue == nil {
		return nil, ErrIssueNotFound
	}

	if issue.Status == "closed" || issue.Status == "cancelled" {
		return nil, ErrIssueAlreadyClosed
	}

	applyIssueUpdate(issue, input)

	tagsJSON, _ := json.Marshal(input.Tags)
	issue.Tags = string(tagsJSON)

	if err := s.issueRepo.Update(ctx, issue); err != nil {
		s.logger.Error("更新工单失败", zap.Error(err), zap.String("issue_id", id.String()))
		return nil, errors.New("更新工单失败")
	}

	return issue, nil
}

// DeleteIssue 删除工单用例
func (s *IssueService) DeleteIssue(ctx context.Context, id uuid.UUID) error {
	issue, err := s.issueRepo.GetByID(ctx, id)
	if err != nil || issue == nil {
		return ErrIssueNotFound
	}

	if err := s.issueRepo.Delete(ctx, id); err != nil {
		s.logger.Error("删除工单失败", zap.Error(err), zap.String("issue_id", id.String()))
		return errors.New("删除工单失败")
	}

	s.logger.Info("工单已删除", zap.String("issue_id", id.String()), zap.String("title", issue.Title))
	return nil
}

// ListIssues 分页查询工单列表
func (s *IssueService) ListIssues(ctx context.Context, filter *repository.IssueFilter) ([]*entity.Issue, int64, error) {
	return s.issueRepo.List(ctx, filter)
}

// TransitionIssue 状态流转用例（含状态机校验）
func (s *IssueService) TransitionIssue(ctx context.Context, id uuid.UUID, newStatus string, resolvedBy *uuid.UUID) (*entity.Issue, error) {
	issue, err := s.issueRepo.GetByID(ctx, id)
	if err != nil || issue == nil {
		return nil, ErrIssueNotFound
	}

	if !isValidIssueTransition(issue.Status, newStatus) {
		return nil, ErrInvalidIssueStatus
	}

	now := time.Now()
	if newStatus == "resolved" && resolvedBy != nil {
		issue.ResolvedBy = resolvedBy
		issue.ResolvedAt = &now
		issue.Status = newStatus
	} else if newStatus == "closed" {
		if resolvedBy != nil {
			issue.ClosedBy = resolvedBy
		}
		issue.ClosedAt = &now
		issue.Status = newStatus
	} else {
		issue.Status = newStatus
	}

	if err := s.issueRepo.UpdateStatus(ctx, id, newStatus); err != nil {
		s.logger.Error("更新工单状态失败", zap.Error(err),
			zap.String("issue_id", id.String()),
			zap.String("from", issue.Status),
			zap.String("to", newStatus),
		)
		return nil, errors.New("状态更新失败")
	}

	// 如果是已解决或已关闭，更新解决信息
	if newStatus == "resolved" || newStatus == "closed" {
		s.issueRepo.Update(ctx, issue)
	}

	s.logger.Info("工单状态变更",
		zap.String("issue_id", id.String()),
		zap.String("old_status", issue.Status),
		zap.String("new_status", newStatus),
	)

	return issue, nil
}

// AddComment 添加评论用例
func (s *IssueService) AddComment(ctx context.Context, issueID, userID uuid.UUID, content string, isInternal bool, parentID *uuid.UUID) (*entity.IssueComment, error) {
	issue, err := s.issueRepo.GetByID(ctx, issueID)
	if err != nil || issue == nil {
		return nil, ErrIssueNotFound
	}

	comment := &entity.IssueComment{
		IssueID:    issueID,
		UserID:     userID,
		Content:    content,
		IsInternal: isInternal,
		ParentID:   parentID,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, errors.New("添加评论失败")
	}

	return comment, nil
}

// ListComments 获取评论列表
func (s *IssueService) ListComments(ctx context.Context, issueID uuid.UUID) ([]*entity.IssueComment, error) {
	return s.commentRepo.ListByIssueID(ctx, issueID)
}

// AddAttachment 添加附件用例
func (s *IssueService) AddAttachment(ctx context.Context, issueID uuid.UUID, uploadedBy *uuid.UUID, fileName string, fileSize int64, fileType, fileURL string) (*entity.IssueAttachment, error) {
	issue, err := s.issueRepo.GetByID(ctx, issueID)
	if err != nil || issue == nil {
		return nil, ErrIssueNotFound
	}

	var uploaderID uuid.UUID
	if uploadedBy != nil {
		uploaderID = *uploadedBy
	} else {
		uploaderID = uuid.Nil
	}

	attachment := &entity.IssueAttachment{
		IssueID:    issueID,
		FileName:   fileName,
		FileSize:   fileSize,
		FileType:   fileType,
		FileURL:    fileURL,
		UploadedBy: uploaderID,
	}

	if err := s.attachmentRepo.Create(ctx, attachment); err != nil {
		return nil, errors.New("添加附件失败")
	}

	return attachment, nil
}

// ListAttachments 获取附件列表
func (s *IssueService) ListAttachments(ctx context.Context, issueID uuid.UUID) ([]*entity.IssueAttachment, error) {
	return s.attachmentRepo.ListByIssueID(ctx, issueID)
}

// SLAResult SLA计算结果
type SLAResult struct {
	ResponseSLA       int
	ResolveSLA        int
	SLADueDate        time.Time
	ResponseDueDate   time.Time
	IsOverdue         bool
	ResponseOverdue   bool
	RemainingMinutes  *int64
	BusinessHoursOnly bool
}

// GetSLAInfo 获取SLA信息
func (s *IssueService) GetSLAInfo(ctx context.Context, issue *entity.Issue) (*SLAResult, error) {
	result := &SLAResult{}

	// 根据工单类型和优先级查找SLA配置
	slaCfg, err := s.slaConfigRepo.GetByTypeAndPriority(ctx, issue.Type, issue.Priority)
	if err != nil || slaCfg == nil {
		// 使用默认值
		result.ResponseSLA = 60
		result.ResolveSLA = 480
		result.BusinessHoursOnly = false
	} else {
		result.ResponseSLA = slaCfg.ResponseSLA
		result.ResolveSLA = slaCfg.ResolveSLA
		result.BusinessHoursOnly = slaCfg.BusinessHoursOnly
	}

	// 计算截止时间
	now := time.Now()
	result.ResponseDueDate = now.Add(time.Duration(result.ResponseSLA) * time.Minute)
	result.SLADueDate = now.Add(time.Duration(result.ResolveSLA) * time.Minute)

	// 判断是否超期
	if issue.ResponseDueDate != nil && now.After(*issue.ResponseDueDate) {
		result.ResponseOverdue = true
	}
	if issue.SLADueDate != nil && now.After(*issue.SLADueDate) {
		result.IsOverdue = true
	}

	// 计算剩余时间
	if issue.SLADueDate != nil && now.Before(*issue.SLADueDate) {
		remaining := issue.SLADueDate.Sub(now).Minutes()
		remInt64 := int64(remaining)
		result.RemainingMinutes = &remInt64
	}

	return result, nil
}

// RateSatisfaction 满意度评价用例
func (s *IssueService) RateSatisfaction(ctx context.Context, id uuid.UUID, score int) error {
	issue, err := s.issueRepo.GetByID(ctx, id)
	if err != nil || issue == nil {
		return ErrIssueNotFound
	}

	if score < 1 || score > 5 {
		return ErrInvalidSatisfaction
	}

	issue.Satisfaction = &score
	if err := s.issueRepo.Update(ctx, issue); err != nil {
		s.logger.Error("提交满意度评价失败", zap.Error(err), zap.String("issue_id", id.String()))
		return errors.New("提交满意度评价失败")
	}

	s.logger.Info("满意度评价已提交",
		zap.String("issue_id", id.String()),
		zap.Int("score", score),
	)
	return nil
}

// ListMyIssues 我的工单列表（我提报的 + 分配给我的）
func (s *IssueService) ListMyIssues(ctx context.Context, reporterID, assigneeID uuid.UUID, page, pageSize int) ([]*entity.Issue, int64, error) {
	// 同时查询提报的和分配的，合并去重
	reported, reportedTotal, _ := s.issueRepo.ListByReporterID(ctx, reporterID, page, pageSize)
	assigned, assignedTotal, _ := s.issueRepo.ListByAssigneeID(ctx, assigneeID, page, pageSize)

	// 合并结果（简单合并，实际生产中可能需要更复杂的去重逻辑）
	seen := make(map[string]bool)
	var result []*entity.Issue
	for _, iss := range reported {
		key := iss.ID.String()
		if !seen[key] {
			seen[key] = true
			result = append(result, iss)
		}
	}
	for _, iss := range assigned {
		key := iss.ID.String()
		if !seen[key] {
			seen[key] = true
			result = append(result, iss)
		}
	}

	total := reportedTotal
	if assignedTotal > total {
		total = assignedTotal
	}

	return result, total, nil
}

// calculateAndSetSLADates 计算并设置SLA截止时间
func (s *IssueService) calculateAndSetSLADates(ctx context.Context, issue *entity.Issue) {
	slaCfg, err := s.slaConfigRepo.GetByTypeAndPriority(ctx, issue.Type, issue.Priority)
	if err != nil || slaCfg == nil {
		return
	}

	now := time.Now()

	responseDue := now.Add(time.Duration(slaCfg.ResponseSLA) * time.Minute)
	slaDue := now.Add(time.Duration(slaCfg.ResolveSLA) * time.Minute)

	s.issueRepo.UpdateSLADates(ctx, issue.ID, &slaDue, &responseDue)
}

// isValidIssueTransition 校验工单状态流转是否合法
var validIssueTransitions = map[string][]string{
	"new":         {"in_progress", "waiting", "cancelled"},
	"in_progress": {"resolved", "waiting", "cancelled"},
	"waiting":     {"in_progress", "cancelled"},
	"resolved":    {"closed", "in_progress", "cancelled"},
	"closed":      {},
	"cancelled":   {"new"},
}

func isValidIssueTransition(fromStatus, toStatus string) bool {
	validTos, ok := validIssueTransitions[fromStatus]
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

// CreateIssueInput 创建工单输入
type CreateIssueInput struct {
	Title       string
	Description string
	Type        string
	ProjectID   *uuid.UUID
	ProductID   *uuid.UUID
	ReporterID  uuid.UUID
	AssigneeID  *uuid.UUID
	Priority    int
	Severity    int
	Source      string
	TemplateID  *uuid.UUID
	Tags        []string
}

// UpdateIssueInput 更新工单输入
type UpdateIssueInput struct {
	Title       *string
	Description *string
	Type        *string
	ProjectID   *uuid.UUID
	ProductID   *uuid.UUID
	AssigneeID  *uuid.UUID
	Priority    *int
	Severity    *int
	Resolution  *string
	Tags        []string
}

// applyIssueUpdate 将更新输入应用到实体
func applyIssueUpdate(issue *entity.Issue, input *UpdateIssueInput) {
	if input.Title != nil {
		issue.Title = *input.Title
	}
	if input.Description != nil {
		issue.Description = *input.Description
	}
	if input.Type != nil {
		issue.Type = *input.Type
	}
	if input.ProjectID != nil {
		issue.ProjectID = input.ProjectID
	}
	if input.ProductID != nil {
		issue.ProductID = input.ProductID
	}
	if input.AssigneeID != nil {
		issue.AssigneeID = input.AssigneeID
	}
	if input.Priority != nil {
		issue.Priority = *input.Priority
	}
	if input.Severity != nil {
		issue.Severity = *input.Severity
	}
	if input.Resolution != nil {
		issue.Resolution = *input.Resolution
	}
}
