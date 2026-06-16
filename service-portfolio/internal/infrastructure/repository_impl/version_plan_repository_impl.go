package repository_impl

import (
	"context"

	"leap-one/service-portfolio/internal/domain/entity"
	"leap-one/service-portfolio/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductVersionRepositoryImpl 产品版本仓库实现
type ProductVersionRepositoryImpl struct {
	db *gorm.DB
}

// NewProductVersionRepository 创建产品版本仓库实例
func NewProductVersionRepository(db *gorm.DB) repository.ProductVersionRepository {
	return &ProductVersionRepositoryImpl{db: db}
}

func (r *ProductVersionRepositoryImpl) Create(ctx context.Context, version *entity.ProductVersion) error {
	return r.db.WithContext(ctx).Create(version).Error
}

func (r *ProductVersionRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.ProductVersion, error) {
	var v entity.ProductVersion
	err := r.db.WithContext(ctx).First(&v, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &v, nil
}

func (r *ProductVersionRepositoryImpl) Update(ctx context.Context, version *entity.ProductVersion) error {
	return r.db.WithContext(ctx).Save(version).Error
}

func (r *ProductVersionRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ProductVersion{}, "id = ?", id).Error
}

func (r *ProductVersionRepositoryImpl) ListByProductID(ctx context.Context, productID uuid.UUID) ([]*entity.ProductVersion, error) {
	var versions []*entity.ProductVersion
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&versions).Error
	return versions, err
}

// ProductRoadmapRepositoryImpl 产品路线图仓库实现
type ProductRoadmapRepositoryImpl struct {
	db *gorm.DB
}

// NewProductRoadmapRepository 创建路线图仓库实例
func NewProductRoadmapRepository(db *gorm.DB) repository.ProductRoadmapRepository {
	return &ProductRoadmapRepositoryImpl{db: db}
}

func (r *ProductRoadmapRepositoryImpl) Create(ctx context.Context, item *entity.ProductRoadmapItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ProductRoadmapRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.ProductRoadmapItem, error) {
	var item entity.ProductRoadmapItem
	err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *ProductRoadmapRepositoryImpl) Update(ctx context.Context, item *entity.ProductRoadmapItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *ProductRoadmapRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ProductRoadmapItem{}, "id = ?", id).Error
}

func (r *ProductRoadmapRepositoryImpl) ListByProductID(ctx context.Context, productID uuid.UUID) ([]*entity.ProductRoadmapItem, error) {
	var items []*entity.ProductRoadmapItem
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("sort_order ASC, priority ASC, created_at ASC").
		Find(&items).Error
	return items, err
}

// UpdateSortOrder 批量更新路线图项的排序顺序
func (r *ProductRoadmapRepositoryImpl) UpdateSortOrder(ctx context.Context, productID uuid.UUID, itemIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, id := range itemIDs {
			if err := tx.Model(&entity.ProductRoadmapItem{}).
				Where("id = ? AND product_id = ?", id, productID).
				Update("sort_order", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ProductPlanRepositoryImpl 产品计划仓库实现
type ProductPlanRepositoryImpl struct {
	db *gorm.DB
}

// NewProductPlanRepository 创建产品计划仓库实例
func NewProductPlanRepository(db *gorm.DB) repository.ProductPlanRepository {
	return &ProductPlanRepositoryImpl{db: db}
}

func (r *ProductPlanRepositoryImpl) Create(ctx context.Context, plan *entity.ProductPlan) error {
	return r.db.WithContext(ctx).Create(plan).Error
}

func (r *ProductPlanRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.ProductPlan, error) {
	var plan entity.ProductPlan
	err := r.db.WithContext(ctx).First(&plan, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &plan, nil
}

func (r *ProductPlanRepositoryImpl) Update(ctx context.Context, plan *entity.ProductPlan) error {
	return r.db.WithContext(ctx).Save(plan).Error
}

func (r *ProductPlanRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ProductPlan{}, "id = ?", id).Error
}

func (r *ProductPlanRepositoryImpl) ListByProductID(ctx context.Context, productID uuid.UUID) ([]*entity.ProductPlan, error) {
	var plans []*entity.ProductPlan
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("start_date ASC, created_at ASC").
		Find(&plans).Error
	return plans, err
}

// List 分页查询全部计划列表
func (r *ProductPlanRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*entity.ProductPlan, int64, error) {
	var plans []*entity.ProductPlan
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.ProductPlan{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&plans).Error

	if err != nil {
		return nil, 0, err
	}

	return plans, total, nil
}
