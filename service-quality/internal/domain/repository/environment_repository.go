package repository

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-quality/internal/domain/entity"
)

// EnvironmentRepository 测试环境仓库接口定义
type EnvironmentRepository interface {
	// Create 创建测试环境
	Create(ctx context.Context, env *entity.TestEnvironment) error

	// GetByID 根据ID获取测试环境
	GetByID(ctx context.Context, id uuid.UUID) (*entity.TestEnvironment, error)

	// Update 更新测试环境
	Update(ctx context.Context, env *entity.TestEnvironment) error

	// Delete 删除测试环境（软删除�?
	Delete(ctx context.Context, id uuid.UUID) error

	// List 查询所有启用的测试环境列表
	List(ctx context.Context, includeInactive bool) ([]*entity.TestEnvironment, error)
}
