package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-quality/internal/domain/entity"
	"leap-one/service-quality/internal/domain/repository"
	"gorm.io/gorm"
)

// EnvironmentRepositoryImpl 测试环境仓库实现
type EnvironmentRepositoryImpl struct {
	db *gorm.DB
}

// NewEnvironmentRepository 创建测试环境仓库实例
func NewEnvironmentRepository(db *gorm.DB) repository.EnvironmentRepository {
	return &EnvironmentRepositoryImpl{db: db}
}

// Create 创建测试环境
func (r *EnvironmentRepositoryImpl) Create(ctx context.Context, env *entity.TestEnvironment) error {
	return r.db.WithContext(ctx).Create(env).Error
}

// GetByID 根据ID获取测试环境
func (r *EnvironmentRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.TestEnvironment, error) {
	var env entity.TestEnvironment
	err := r.db.WithContext(ctx).First(&env, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &env, nil
}

// Update 更新测试环境
func (r *EnvironmentRepositoryImpl) Update(ctx context.Context, env *entity.TestEnvironment) error {
	return r.db.WithContext(ctx).Save(env).Error
}

// Delete 删除测试环境（软删除�?func (r *EnvironmentRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.TestEnvironment{}, "id = ?", id).Error
}

// List 查询所有测试环境列�?func (r *EnvironmentRepositoryImpl) List(ctx context.Context, includeInactive bool) ([]*entity.TestEnvironment, error) {
	var envs []*entity.TestEnvironment

	query := r.db.WithContext(ctx).Model(&entity.TestEnvironment{})
	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}

	err := query.Order("type ASC, created_at ASC").Find(&envs).Error
	if err != nil {
		return nil, err
	}
	return envs, nil
}
