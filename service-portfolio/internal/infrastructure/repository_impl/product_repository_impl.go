package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-portfolio/internal/domain/entity"
	"leap-one/service-portfolio/internal/domain/repository"
	"gorm.io/gorm"
)

// ProductRepositoryImpl 产品仓库实现
type ProductRepositoryImpl struct {
	db *gorm.DB
}

// NewProductRepository 创建产品仓库实例
func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &ProductRepositoryImpl{db: db}
}

// Create 创建产品
func (r *ProductRepositoryImpl) Create(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// GetByID 根据ID获取产品
func (r *ProductRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).First(&product, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// GetByCode 根据编码获取产品
func (r *ProductRepositoryImpl) GetByCode(ctx context.Context, code string) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// Update 更新产品信息
func (r *ProductRepositoryImpl) Update(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// Delete 软删除产品
func (r *ProductRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Product{}, "id = ?", id).Error
}

// List 分页查询产品列表（支持关键词搜索、状态筛选、产品线筛选）
func (r *ProductRepositoryImpl) List(ctx context.Context, page, pageSize int, keyword, status string, productLineID *uuid.UUID) ([]*entity.Product, int64, error) {
	var products []*entity.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Product{})

	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR description LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if productLineID != nil && *productLineID != uuid.Nil {
		query = query.Where("product_line_id = ?", *productLineID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// ListByProgramID 根据项目集ID查询关联的产品列表
func (r *ProductRepositoryImpl) ListByProgramID(ctx context.Context, programID uuid.UUID) ([]*entity.Product, error) {
	var products []*entity.Product
	err := r.db.WithContext(ctx).
		Where("program_id = ?", programID).
		Order("created_at DESC").
		Find(&products).Error
	return products, err
}

// ListByProductLineID 根据产品线ID查询产品列表
func (r *ProductRepositoryImpl) ListByProductLineID(ctx context.Context, productLineID uuid.UUID) ([]*entity.Product, error) {
	var products []*entity.Product
	err := r.db.WithContext(ctx).
		Where("product_line_id = ?", productLineID).
		Order("created_at DESC").
		Find(&products).Error
	return products, err
}

// GetByOwnerID 根据负责人ID获取产品列表
func (r *ProductRepositoryImpl) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*entity.Product, error) {
	var products []*entity.Product
	err := r.db.WithContext(ctx).
		Where("owner_id = ?", ownerID).
		Order("created_at DESC").
		Find(&products).Error
	return products, err
}

// ProductLineRepositoryImpl 产品线仓库实现
type ProductLineRepositoryImpl struct {
	db *gorm.DB
}

// NewProductLineRepository 创建产品线仓库实例
func NewProductLineRepository(db *gorm.DB) repository.ProductLineRepository {
	return &ProductLineRepositoryImpl{db: db}
}

// Create 创建产品线
func (r *ProductLineRepositoryImpl) Create(ctx context.Context, line *entity.ProductLine) error {
	return r.db.WithContext(ctx).Create(line).Error
}

// GetByID 根据ID获取产品线
func (r *ProductLineRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.ProductLine, error) {
	var line entity.ProductLine
	err := r.db.WithContext(ctx).First(&line, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &line, nil
}

// Update 更新产品线
func (r *ProductLineRepositoryImpl) Update(ctx context.Context, line *entity.ProductLine) error {
	return r.db.WithContext(ctx).Save(line).Error
}

// Delete 软删除产品线
func (r *ProductLineRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ProductLine{}, "id = ?", id).Error
}

// List 查询全部产品线列表（按排序字段）
func (r *ProductLineRepositoryImpl) List(ctx context.Context) ([]*entity.ProductLine, error) {
	var lines []*entity.ProductLine
	err := r.db.WithContext(ctx).
		Order("sort_order ASC, created_at ASC").
		Find(&lines).Error
	return lines, err
}
