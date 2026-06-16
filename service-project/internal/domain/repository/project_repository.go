package repository

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-project/internal/domain/entity"
)

// ProjectRepository 项目仓库接口定义
type ProjectRepository interface {
	// Create 创建项目
	Create(ctx context.Context, project *entity.Project) error

	// GetByID 根据ID获取项目
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Project, error)

	// GetByCode 根据项目编号获取项目
	GetByCode(ctx context.Context, code string) (*entity.Project, error)

	// Update 更新项目信息
	Update(ctx context.Context, project *entity.Project) error

	// Delete 软删除项目
	Delete(ctx context.Context, id uuid.UUID) error

	// List 分页查询项目列表（支持筛选和搜索）
	List(ctx context.Context, page, pageSize int, keyword, status, programID, pmID, sortBy, sortOrder string) ([]*entity.Project, int64, error)

	// UpdateStatus 更新项目状态（planning→executing→paused/completed/cancelled/archived）
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error

	// ListByPMID 查询某项目经理下的所有项目
	ListByPMID(ctx context.Context, pmID uuid.UUID) ([]*entity.Project, error)

	// CountByStatus 按状态统计项目数量
	CountByStatus(ctx context.Context) (map[string]int64, error)
}

// ProjectMemberRepository 项目成员仓库接口定义
type ProjectMemberRepository interface {
	// Add 添加项目成员
	Add(ctx context.Context, member *entity.ProjectMember) error

	// Remove 移除项目成员
	Remove(ctx context.Context, projectID, userID uuid.UUID) error

	// Get 获取单个成员
	Get(ctx context.Context, projectID, userID uuid.UUID) (*entity.ProjectMember, error)

	// ListByProjectID 获取项目的所有成员
	ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entity.ProjectMember, error)

	// UpdateRole 更新成员角色
	UpdateRole(ctx context.Context, projectID, userID uuid.UUID, role string) error

	// CountByProjectID 统计项目成员数量
	CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error)
}

// ProjectTemplateRepository 项目模板仓库接口定义
type ProjectTemplateRepository interface {
	// Create 创建模板
	Create(ctx context.Context, template *entity.ProjectTemplate) error

	// GetByID 根据ID获取模板
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ProjectTemplate, error)

	// List 列出所有模板（分页）
	List(ctx context.Context, page, pageSize int, templateType string) ([]*entity.ProjectTemplate, int64, error)

	// ListSystemTemplates 列出系统预置模板
	ListSystemTemplates(ctx context.Context) ([]*entity.ProjectTemplate, error)

	// Update 更新模板
	Update(ctx context.Context, template *entity.ProjectTemplate) error

	// Delete 删除模板（仅非系统预置）
	Delete(ctx context.Context, id uuid.UUID) error
}

// ProjectMilestoneRepository 项目里程碑仓库接口定义
type ProjectMilestoneRepository interface {
	// Create 创建里程碑
	Create(ctx context.Context, milestone *entity.ProjectMilestone) error

	// GetByID 根据ID获取里程碑
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ProjectMilestone, error)

	// ListByProjectID 获取项目的所有里程碑
	ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entity.ProjectMilestone, error)

	// Update 更新里程碑
	Update(ctx context.Context, milestone *entity.ProjectMilestone) error

	// Delete 删除里程碑
	Delete(ctx context.Context, id uuid.UUID) error

	// Complete 完成里程碑
	Complete(ctx context.Context, id uuid.UUID, completedBy uuid.UUID) error

	// CountByProjectID 统计项目里程碑数量
	CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error)
}

// ProjectRiskRepository 项目风险仓库接口定义
type ProjectRiskRepository interface {
	// Create 创建风险
	Create(ctx context.Context, risk *entity.ProjectRisk) error

	// GetByID 根据ID获取风险
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ProjectRisk, error)

	// ListByProjectID 获取项目的所有风险
	ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entity.ProjectRisk, error)

	// Update 更新风险
	Update(ctx context.Context, risk *entity.ProjectRisk) error

	// Delete 删除风险
	Delete(ctx context.Context, id uuid.UUID) error

	// CountHighRisk 统计高风险数量
	CountHighRisk(ctx context.Context, projectID uuid.UUID) (int64, error)
}

// CustomFieldRepository 自定义字段仓库接口定义
type CustomFieldRepository interface {
	// Create 创建自定义字段
	Create(ctx context.Context, field *entity.CustomField) error

	// GetByID 根据ID获取字段
	GetByID(ctx context.Context, id uuid.UUID) (*entity.CustomField, error)

	// ListByProjectID 获取项目的所有自定义字段
	ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entity.CustomField, error)

	// Update 更新字段
	Update(ctx context.Context, field *entity.CustomField) error

	// Delete 删除字段
	Delete(ctx context.Context, id uuid.UUID) error
}

// IterationRepository 迭代仓库接口定义
type IterationRepository interface {
	// Create 创建迭代
	Create(ctx context.Context, iteration *entity.Iteration) error

	// GetByID 根据ID获取迭代
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Iteration, error)

	// List 分页查询迭代列表
	List(ctx context.Context, page, pageSize int, projectID uuid.UUID, status string) ([]*entity.Iteration, int64, error)

	// ListByProjectID 获取项目的所有迭代
	ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entity.Iteration, error)

	// Update 更新迭代
	Update(ctx context.Context, iteration *entity.Iteration) error

	// Delete 删除迭代
	Delete(ctx context.Context, id uuid.UUID) error

	// UpdateStatus 更新迭代状态
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error

	// GetActiveIteration 获取项目当前活跃的迭代
	GetActiveIteration(ctx context.Context, projectID uuid.UUID) (*entity.Iteration, error)

	// ListCompleted 获取已完成的迭代列表（用于计算平均速度）
	ListCompleted(ctx context.Context, projectID uuid.UUID) ([]*entity.Iteration, error)
}
