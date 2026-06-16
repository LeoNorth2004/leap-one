package repository_impl

import (
	"context"
	"time"

	"github.com/google/uuid"
	"leap-one/service-quality/internal/domain/entity"
	"leap-one/service-quality/internal/domain/repository"
	"gorm.io/gorm"
)

// TestCaseRepositoryImpl 测试用例仓库实现
type TestCaseRepositoryImpl struct {
	db *gorm.DB
}

// NewTestCaseRepository 创建测试用例仓库实例
func NewTestCaseRepository(db *gorm.DB) repository.TestCaseRepository {
	return &TestCaseRepositoryImpl{db: db}
}

// Create 创建测试用例
func (r *TestCaseRepositoryImpl) Create(ctx context.Context, testCase *entity.TestCase) error {
	return r.db.WithContext(ctx).Create(testCase).Error
}

// GetByID 根据ID获取测试用例
func (r *TestCaseRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.TestCase, error) {
	var tc entity.TestCase
	err := r.db.WithContext(ctx).First(&tc, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &tc, nil
}

// Update 更新测试用例
func (r *TestCaseRepositoryImpl) Update(ctx context.Context, testCase *entity.TestCase) error {
	return r.db.WithContext(ctx).Save(testCase).Error
}

// Delete 删除测试用例（软删除�?
func (r *TestCaseRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.TestCase{}, "id = ?", id).Error
}

// List 分页查询测试用例列表（支持筛选条件）
func (r *TestCaseRepositoryImpl) List(ctx context.Context, page, pageSize int, filter *repository.TestCaseFilter) ([]*entity.TestCase, int64, error) {
	var cases []*entity.TestCase
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.TestCase{})

	// 应用筛选条�?
	query = r.applyTestCaseFilter(query, filter)

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按创建时间倒序
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&cases).Error

	if err != nil {
		return nil, 0, err
	}

	return cases, total, nil
}

// BatchCreate 批量创建测试用例（导入场景）
func (r *TestCaseRepositoryImpl) BatchCreate(ctx context.Context, cases []*entity.TestCase) error {
	return r.db.WithContext(ctx).Create(&cases).Error
}

// BatchDelete 批量删除测试用例
func (r *TestCaseRepositoryImpl) BatchDelete(ctx context.Context, ids []uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&entity.TestCase{}).Error
}

// Review 评审用例（设置评审人和评审时间）
func (r *TestCaseRepositoryImpl) Review(ctx context.Context, id uuid.UUID, reviewerID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.TestCase{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"reviewer_id": reviewerID,
			"reviewed_at": now,
			"status":      "active",
		}).Error
}

// applyTestCaseFilter 应用测试用例筛选条件到查询构建�?
func (r *TestCaseRepositoryImpl) applyTestCaseFilter(query *gorm.DB, filter *repository.TestCaseFilter) *gorm.DB {
	if filter == nil {
		return query
	}

	// 关键词搜索：匹配标题和模�?
	if filter.Keyword != "" {
		searchPattern := "%" + filter.Keyword + "%"
		query = query.Where("title LIKE ? OR module LIKE ?", searchPattern, searchPattern)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}
	if filter.ProductID != nil {
		query = query.Where("product_id = ?", *filter.ProductID)
	}
	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.CreatorID != nil {
		query = query.Where("creator_id = ?", *filter.CreatorID)
	}
	if filter.Automation != nil {
		query = query.Where("automation = ?", *filter.Automation)
	}
	if filter.RequirementID != nil {
		query = query.Where("requirement_id = ?", *filter.RequirementID)
	}

	return query
}
