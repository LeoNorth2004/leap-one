package repository

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-quality/internal/domain/entity"
)

// TestSuiteRepository 测试套件仓库接口定义
type TestSuiteRepository interface {
	// Create 创建测试套件
	Create(ctx context.Context, suite *entity.TestSuite) error

	// GetByID 根据ID获取测试套件（含关联用例列表�?
	GetByID(ctx context.Context, id uuid.UUID) (*entity.TestSuite, error)

	// Update 更新测试套件
	Update(ctx context.Context, suite *entity.TestSuite) error

	// Delete 删除测试套件（软删除，同时清理关联关系）
	Delete(ctx context.Context, id uuid.UUID) error

	// List 分页查询测试套件列表
	List(ctx context.Context, page, pageSize int, productID, projectID *uuid.UUID) ([]*entity.TestSuite, int64, error)

	// AddCases 添加用例到套�?
	AddCases(ctx context.Context, suiteID uuid.UUID, caseIDs []uuid.UUID) error

	// RemoveCase 从套件中移除指定用例
	RemoveCase(ctx context.Context, suiteID, caseID uuid.UUID) error

	// GetSuiteCases 获取套件中的用例列表（含排序�?
	GetSuiteCases(ctx context.Context, suiteID uuid.UUID) ([]*entity.TestCaseSuiteRelation, error)
}
