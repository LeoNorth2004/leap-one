package repository

import (
	"context"

	"leap-one/service-portfolio/internal/domain/entity"

	"github.com/google/uuid"
)

// ProductVersionRepository 产品版本仓库接口定义
type ProductVersionRepository interface {
	// Create 创建版本
	Create(ctx context.Context, version *entity.ProductVersion) error

	// GetByID 根据ID获取版本
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ProductVersion, error)

	// Update 更新版本信息
	Update(ctx context.Context, version *entity.ProductVersion) error

	// Delete 软删除版本
	Delete(ctx context.Context, id uuid.UUID) error

	// ListByProductID 根据产品ID查询版本列表
	ListByProductID(ctx context.Context, productID uuid.UUID) ([]*entity.ProductVersion, error)
}

// ProductRoadmapRepository 产品路线图仓库接口定义
type ProductRoadmapRepository interface {
	// Create 创建路线图项
	Create(ctx context.Context, item *entity.ProductRoadmapItem) error

	// GetByID 根据ID获取路线图项
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ProductRoadmapItem, error)

	// Update 更新路线图项
	Update(ctx context.Context, item *entity.ProductRoadmapItem) error

	// Delete 删除路线图项
	Delete(ctx context.Context, id uuid.UUID) error

	// ListByProductID 根据产品ID查询路线图列表（按排序字段）
	ListByProductID(ctx context.Context, productID uuid.UUID) ([]*entity.ProductRoadmapItem, error)

	// UpdateSortOrder 批量更新路线图项的排序顺序
	UpdateSortOrder(ctx context.Context, productID uuid.UUID, itemIDs []uuid.UUID) error
}

// ProductPlanRepository 产品计划仓库接口定义
type ProductPlanRepository interface {
	// Create 创建计划
	Create(ctx context.Context, plan *entity.ProductPlan) error

	// GetByID 根据ID获取计划
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ProductPlan, error)

	// Update 更新计划
	Update(ctx context.Context, plan *entity.ProductPlan) error

	// Delete 软删除计划
	Delete(ctx context.Context, id uuid.UUID) error

	// ListByProductID 根据产品ID查询计划列表
	ListByProductID(ctx context.Context, productID uuid.UUID) ([]*entity.ProductPlan, error)

	// List 分页查询全部计划列表
	List(ctx context.Context, page, pageSize int) ([]*entity.ProductPlan, int64, error)
}
