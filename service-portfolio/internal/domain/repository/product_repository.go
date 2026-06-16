package repository

import (
	"context"

	"leap-one/service-portfolio/internal/domain/entity"

	"github.com/google/uuid"
)

// ProductRepository 产品仓库接口定义
type ProductRepository interface {
	// Create 创建产品
	Create(ctx context.Context, product *entity.Product) error

	// GetByID 根据ID获取产品
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)

	// GetByCode 根据编码获取产品
	GetByCode(ctx context.Context, code string) (*entity.Product, error)

	// Update 更新产品信息
	Update(ctx context.Context, product *entity.Product) error

	// Delete 软删除产品
	Delete(ctx context.Context, id uuid.UUID) error

	// List 分页查询产品列表（支持关键词搜索、状态筛选、产品线筛选）
	List(ctx context.Context, page, pageSize int, keyword, status string, productLineID *uuid.UUID) ([]*entity.Product, int64, error)

	// ListByProgramID 根据项目集ID查询关联的产品列表
	ListByProgramID(ctx context.Context, programID uuid.UUID) ([]*entity.Product, error)

	// ListByProductLineID 根据产品线ID查询产品列表
	ListByProductLineID(ctx context.Context, productLineID uuid.UUID) ([]*entity.Product, error)

	// GetByOwnerID 根据负责人ID获取产品列表
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*entity.Product, error)
}

// ProductLineRepository 产品线仓库接口定义
type ProductLineRepository interface {
	// Create 创建产品线
	Create(ctx context.Context, line *entity.ProductLine) error

	// GetByID 根据ID获取产品线
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ProductLine, error)

	// Update 更新产品线
	Update(ctx context.Context, line *entity.ProductLine) error

	// Delete 软删除产品线
	Delete(ctx context.Context, id uuid.UUID) error

	// List 查询全部产品线列表（按排序字段）
	List(ctx context.Context) ([]*entity.ProductLine, error)
}
