package repository

import (
	"context"

	"leap-one/service-task/internal/domain/entity"

	"github.com/google/uuid"
)

// IssueWorkflowRepository 工作流仓库接口定义
type IssueWorkflowRepository interface {
	// Create 创建工作流
	Create(ctx context.Context, workflow *entity.IssueWorkflow) error

	// GetByID 根据ID获取工作流（含转换规则）
	GetByID(ctx context.Context, id uuid.UUID) (*entity.IssueWorkflow, error)

	// Update 更新工作流
	Update(ctx context.Context, workflow *entity.IssueWorkflow) error

	// Delete 软删除工作流
	Delete(ctx context.Context, id uuid.UUID) error

	// List 分页查询工作流列表
	List(ctx context.Context, page, pageSize int, wfType string) ([]*entity.IssueWorkflow, int64, error)

	// AddTransition 添加状态转换规则
	AddTransition(ctx context.Context, transition *entity.IssueWorkflowTransition) error

	// RemoveTransition 删除状态转换规则
	RemoveTransition(ctx context.Context, id uuid.UUID) error

	// ListTransitions 查询工作流的转换规则列表
	ListTransitions(ctx context.Context, workflowID uuid.UUID) ([]*entity.IssueWorkflowTransition, error)

	// GetByTypeAndStatus 根据类型和初始状态查找工作流
	GetByTypeAndStatus(ctx context.Context, wfType, initialStatus string) (*entity.IssueWorkflow, error)
}

// IssueSLAConfigRepository SLA配置仓库接口定义
type IssueSLAConfigRepository interface {
	// Create 创建SLA配置
	Create(ctx context.Context, config *entity.IssueSLAConfig) error

	// GetByID 根据ID获取配置
	GetByID(ctx context.Context, id uuid.UUID) (*entity.IssueSLAConfig, error)

	// Update 更新配置
	Update(ctx context.Context, config *entity.IssueSLAConfig) error

	// Delete 删除配置
	Delete(ctx context.Context, id uuid.UUID) error

	// List 分页查询SLA配置列表
	List(ctx context.Context, page, pageSize int) ([]*entity.IssueSLAConfig, int64, error)

	// GetByTypeAndPriority 根据类型和优先级获取SLA配置
	GetByTypeAndPriority(ctx context.Context, slaType string, priority int) (*entity.IssueSLAConfig, error)
}
