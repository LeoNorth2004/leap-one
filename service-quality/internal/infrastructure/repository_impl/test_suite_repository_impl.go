package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-quality/internal/domain/entity"
	"leap-one/service-quality/internal/domain/repository"
	"gorm.io/gorm"
)

// TestSuiteRepositoryImpl 测试套件仓库实现
type TestSuiteRepositoryImpl struct {
	db *gorm.DB
}

// NewTestSuiteRepository 创建测试套件仓库实例
func NewTestSuiteRepository(db *gorm.DB) repository.TestSuiteRepository {
	return &TestSuiteRepositoryImpl{db: db}
}

// Create 创建测试套件
func (r *TestSuiteRepositoryImpl) Create(ctx context.Context, suite *entity.TestSuite) error {
	return r.db.WithContext(ctx).Create(suite).Error
}

// GetByID 根据ID获取测试套件（预加载关联用例列表�?
func (r *TestSuiteRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.TestSuite, error) {
	var suite entity.TestSuite
	err := r.db.WithContext(ctx).
		Preload("Cases").
		First(&suite, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &suite, nil
}

// Update 更新测试套件
func (r *TestSuiteRepositoryImpl) Update(ctx context.Context, suite *entity.TestSuite) error {
	return r.db.WithContext(ctx).Save(suite).Error
}

// Delete 删除测试套件（软删除，同时清理关联关系）
func (r *TestSuiteRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	r.db.WithContext(ctx).Where("suite_id = ?", id).Delete(&entity.TestCaseSuiteRelation{})
	return r.db.WithContext(ctx).Delete(&entity.TestSuite{}, "id = ?", id).Error
}

// List 分页查询测试套件列表
func (r *TestSuiteRepositoryImpl) List(ctx context.Context, page, pageSize int, productID, projectID *uuid.UUID) ([]*entity.TestSuite, int64, error) {
	var suites []*entity.TestSuite
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.TestSuite{})

	if productID != nil {
		query = query.Where("product_id = ?", *productID)
	}
	if projectID != nil {
		query = query.Where("project_id = ?", *projectID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&suites).Error

	if err != nil {
		return nil, 0, err
	}

	return suites, total, nil
}

// AddCases 添加用例到套件（批量插入关联关系�?
func (r *TestSuiteRepositoryImpl) AddCases(ctx context.Context, suiteID uuid.UUID, caseIDs []uuid.UUID) error {
	var maxOrder int
	r.db.WithContext(ctx).Model(&entity.TestCaseSuiteRelation{}).
		Where("suite_id = ?", suiteID).
		Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder)

	relations := make([]entity.TestCaseSuiteRelation, len(caseIDs))
	for i, caseID := range caseIDs {
		relations[i] = entity.TestCaseSuiteRelation{
			SuiteID:   suiteID,
			CaseID:    caseID,
			SortOrder: maxOrder + i + 1,
		}
	}

	return r.db.WithContext(ctx).Create(&relations).Error
}

// RemoveCase 从套件中移除指定用例
func (r *TestSuiteRepositoryImpl) RemoveCase(ctx context.Context, suiteID, caseID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("suite_id = ? AND case_id = ?", suiteID, caseID).
		Delete(&entity.TestCaseSuiteRelation{}).Error
}

// GetSuiteCases 获取套件中的用例列表（含排序信息�?
func (r *TestSuiteRepositoryImpl) GetSuiteCases(ctx context.Context, suiteID uuid.UUID) ([]*entity.TestCaseSuiteRelation, error) {
	var relations []*entity.TestCaseSuiteRelation
	err := r.db.WithContext(ctx).
		Where("suite_id = ?", suiteID).
		Order("sort_order ASC").
		Find(&relations).Error
	if err != nil {
		return nil, err
	}
	return relations, nil
}
