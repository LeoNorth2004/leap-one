package repository

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-quality/internal/domain/entity"
)

// TestPlanRepository 测试计划仓库接口定义
type TestPlanRepository interface {
	// Create 创建测试计划
	Create(ctx context.Context, plan *entity.TestPlan) error

	// GetByID 根据ID获取测试计划（含关联的用例执行记录）
	GetByID(ctx context.Context, id uuid.UUID) (*entity.TestPlan, error)

	// Update 更新测试计划
	Update(ctx context.Context, plan *entity.TestPlan) error

	// Delete 删除测试计划（软删除�?
	Delete(ctx context.Context, id uuid.UUID) error

	// List 分页查询测试计划列表
	List(ctx context.Context, page, pageSize int, filter *TestPlanFilter) ([]*entity.TestPlan, int64, error)

	// AddCases 添加用例到测试计�?
	AddCases(ctx context.Context, planID uuid.UUID, caseIDs []uuid.UUID) error

	// ExecuteCase 执行测试计划中的某个用例
	ExecuteCase(ctx context.Context, planCaseID uuid.UUID, result *entity.TestPlanCase) error

	// GetPlanCase 获取测试计划中的单个用例记录
	GetPlanCase(ctx context.Context, planCaseID uuid.UUID) (*entity.TestPlanCase, error)

	// UpdatePlanStatus 更新测试计划状�?
	UpdatePlanStatus(ctx context.Context, id uuid.UUID, status string) error
}

// TestPlanFilter 测试计划查询筛选条�?
type TestPlanFilter struct {
	Status    string     // 状�?planning/executing/completed/cancelled
	ProductID *uuid.UUID // 产品ID
	ProjectID *uuid.UUID // 项目ID
	CreatorID *uuid.UUID // 创建人ID
}
