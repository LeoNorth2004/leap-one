package repository

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-quality/internal/domain/entity"
)

// TestCaseRepository 测试用例仓库接口定义
type TestCaseRepository interface {
	// Create 创建测试用例
	Create(ctx context.Context, testCase *entity.TestCase) error

	// GetByID 根据ID获取测试用例
	GetByID(ctx context.Context, id uuid.UUID) (*entity.TestCase, error)

	// Update 更新测试用例
	Update(ctx context.Context, testCase *entity.TestCase) error

	// Delete 删除测试用例（软删除�?
	Delete(ctx context.Context, id uuid.UUID) error

	// List 分页查询测试用例列表（支持筛选条件）
	List(ctx context.Context, page, pageSize int, filter *TestCaseFilter) ([]*entity.TestCase, int64, error)

	// BatchCreate 批量创建测试用例（导入场景）
	BatchCreate(ctx context.Context, cases []*entity.TestCase) error

	// BatchDelete 批量删除测试用例
	BatchDelete(ctx context.Context, ids []uuid.UUID) error

	// Review 评审用例（设置评审人和评审时间）
	Review(ctx context.Context, id uuid.UUID, reviewerID uuid.UUID) error
}

// TestCaseFilter 测试用例查询筛选条�?
type TestCaseFilter struct {
	Keyword       string     // 关键词搜索（标题、模块）
	Type          string     // 用例类型 manual/automated/smoke/regression
	Status        string     // 状�?draft/active/archived
	Priority      *int       // 优先�?
	ProductID     *uuid.UUID // 产品ID
	ProjectID     *uuid.UUID // 项目ID
	CreatorID     *uuid.UUID // 创建人ID
	Automation    *bool      // 是否自动�?
	RequirementID *uuid.UUID // 需求ID
}
