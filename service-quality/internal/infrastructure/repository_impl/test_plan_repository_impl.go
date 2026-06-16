package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-quality/internal/domain/entity"
	"leap-one/service-quality/internal/domain/repository"
	"gorm.io/gorm"
)

// TestPlanRepositoryImpl 测试计划仓库实现
type TestPlanRepositoryImpl struct {
	db *gorm.DB
}

// NewTestPlanRepository 创建测试计划仓库实例
func NewTestPlanRepository(db *gorm.DB) repository.TestPlanRepository {
	return &TestPlanRepositoryImpl{db: db}
}

// Create 创建测试计划
func (r *TestPlanRepositoryImpl) Create(ctx context.Context, plan *entity.TestPlan) error {
	return r.db.WithContext(ctx).Create(plan).Error
}

// GetByID 根据ID获取测试计划（预加载关联的用例执行记录）
func (r *TestPlanRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.TestPlan, error) {
	var plan entity.TestPlan
	err := r.db.WithContext(ctx).
		Preload("Cases").
		First(&plan, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &plan, nil
}

// Update 更新测试计划
func (r *TestPlanRepositoryImpl) Update(ctx context.Context, plan *entity.TestPlan) error {
	return r.db.WithContext(ctx).Save(plan).Error
}

// Delete 删除测试计划（软删除�?
func (r *TestPlanRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	// 清理关联的执行记�?
	r.db.WithContext(ctx).Where("plan_id = ?", id).Delete(&entity.TestPlanCase{})
	return r.db.WithContext(ctx).Delete(&entity.TestPlan{}, "id = ?", id).Error
}

// List 分页查询测试计划列表
func (r *TestPlanRepositoryImpl) List(ctx context.Context, page, pageSize int, filter *repository.TestPlanFilter) ([]*entity.TestPlan, int64, error) {
	var plans []*entity.TestPlan
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.TestPlan{})

	if filter != nil {
		if filter.Status != "" {
			query = query.Where("status = ?", filter.Status)
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
	}

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

// AddCases 添加用例到测试计�?
func (r *TestPlanRepositoryImpl) AddCases(ctx context.Context, planID uuid.UUID, caseIDs []uuid.UUID) error {
	// 获取当前最大排序号
	var maxOrder int
	r.db.WithContext(ctx).Model(&entity.TestPlanCase{}).
		Where("plan_id = ?", planID).
		Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder)

	planCases := make([]entity.TestPlanCase, len(caseIDs))
	for i, caseID := range caseIDs {
		planCases[i] = entity.TestPlanCase{
			PlanID:    planID,
			CaseID:    caseID,
			Result:    "not_run",
			SortOrder: maxOrder + i + 1,
		}
	}

	return r.db.WithContext(ctx).Create(&planCases).Error
}

// ExecuteCase 执行测试计划中的某个用例
func (r *TestPlanRepositoryImpl) ExecuteCase(ctx context.Context, planCaseID uuid.UUID, result *entity.TestPlanCase) error {
	return r.db.WithContext(ctx).Model(result).Where("id = ?", planCaseID).Updates(map[string]interface{}{
		"assignee_id":   result.AssigneeID,
		"result":        result.Result,
		"execute_time":  result.ExecuteTime,
		"actual_result": result.ActualResult,
		"bug_ids":       result.BugIDs,
		"comment":       result.Comment,
	}).Error
}

// GetPlanCase 获取测试计划中的单个用例记录
func (r *TestPlanRepositoryImpl) GetPlanCase(ctx context.Context, planCaseID uuid.UUID) (*entity.TestPlanCase, error) {
	var pc entity.TestPlanCase
	err := r.db.WithContext(ctx).First(&pc, "id = ?", planCaseID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &pc, nil
}

// UpdatePlanStatus 更新测试计划状�?
func (r *TestPlanRepositoryImpl) UpdatePlanStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).
		Model(&entity.TestPlan{}).
		Where("id = ?", id).
		Update("status", status).Error
}
