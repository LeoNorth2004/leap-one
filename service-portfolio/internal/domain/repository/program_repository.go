package repository

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-portfolio/internal/domain/entity"
)

// ProgramRepository 项目集仓库接口定义
type ProgramRepository interface {
	// Create 创建项目集
	Create(ctx context.Context, program *entity.Program) error

	// GetByID 根据ID获取项目集
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Program, error)

	// GetByCode 根据编号获取项目集
	GetByCode(ctx context.Context, code string) (*entity.Program, error)

	// Update 更新项目集信息
	Update(ctx context.Context, program *entity.Program) error

	// Delete 软删除项目集
	Delete(ctx context.Context, id uuid.UUID) error

	// List 分页查询项目集列表（支持关键词搜索和状态筛选）
	List(ctx context.Context, page, pageSize int, keyword, status string) ([]*entity.Program, int64, error)

	// GetTree 获取完整项目集树形结构
	GetTree(ctx context.Context) ([]*entity.Program, error)

	// GetChildren 获取子项目集列表
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entity.Program, error)

	// GetByOwnerID 根据负责人ID获取项目集列表
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*entity.Program, error)
}

// MilestoneRepository 里程碑仓库接口定义
type MilestoneRepository interface {
	// Create 创建里程碑
	Create(ctx context.Context, milestone *entity.Milestone) error

	// GetByID 根据ID获取里程碑
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Milestone, error)

	// Update 更新里程碑
	Update(ctx context.Context, milestone *entity.Milestone) error

	// Delete 删除里程碑
	Delete(ctx context.Context, id uuid.UUID) error

	// ListByProgramID 根据项目集ID查询里程碑列表
	ListByProgramID(ctx context.Context, programID uuid.UUID) ([]*entity.Milestone, error)
}

// RiskRepository 风险仓库接口定义
type RiskRepository interface {
	// Create 创建风险项
	Create(ctx context.Context, risk *entity.Risk) error

	// GetByID 根据ID获取风险项
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Risk, error)

	// Update 更新风险项
	Update(ctx context.Context, risk *entity.Risk) error

	// Delete 删除风险项
	Delete(ctx context.Context, id uuid.UUID) error

	// ListByProgramID 根据项目集ID查询风险列表
	ListByProgramID(ctx context.Context, programID uuid.UUID) ([]*entity.Risk, error)
}
