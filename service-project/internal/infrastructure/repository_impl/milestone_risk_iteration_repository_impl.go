package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-project/internal/domain/entity"
	"leap-one/service-project/internal/domain/repository"
	"gorm.io/gorm"
)

// ProjectMilestoneRepositoryImpl 项目里程碑仓库实现
type ProjectMilestoneRepositoryImpl struct {
	db *gorm.DB
}

// NewProjectMilestoneRepository 创建里程碑仓库实例
func NewProjectMilestoneRepository(db *gorm.DB) repository.ProjectMilestoneRepository {
	return &ProjectMilestoneRepositoryImpl{db: db}
}

// Create 创建里程碑
func (r *ProjectMilestoneRepositoryImpl) Create(ctx context.Context, milestone *entity.ProjectMilestone) error {
	return r.db.WithContext(ctx).Create(milestone).Error
}

// GetByID 根据ID获取里程碑
func (r *ProjectMilestoneRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.ProjectMilestone, error) {
	var milestone entity.ProjectMilestone
	err := r.db.WithContext(ctx).First(&milestone, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &milestone, nil
}

// ListByProjectID 获取项目的所有里程碑
func (r *ProjectMilestoneRepositoryImpl) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entity.ProjectMilestone, error) {
	var milestones []*entity.ProjectMilestone
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("sort_order ASC, due_date ASC").
		Find(&milestones).Error
	if err != nil {
		return nil, err
	}
	return milestones, nil
}

// Update 更新里程碑
func (r *ProjectMilestoneRepositoryImpl) Update(ctx context.Context, milestone *entity.ProjectMilestone) error {
	return r.db.WithContext(ctx).Save(milestone).Error
}

// Delete 删除里程碑
func (r *ProjectMilestoneRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ProjectMilestone{}, "id = ?", id).Error
}

// Complete 完成里程碑
func (r *ProjectMilestoneRepositoryImpl) Complete(ctx context.Context, id uuid.UUID, completedBy uuid.UUID) error {
	now := NowFunc()
	return r.db.WithContext(ctx).
		Model(&entity.ProjectMilestone{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       "completed",
			"completed_at": now,
			"completed_by": completedBy,
		}).Error
}

// CountByProjectID 统计项目里程碑数
func (r *ProjectMilestoneRepositoryImpl) CountByProjectID(ctx context.Context, projectID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.ProjectMilestone{}).
		Where("project_id = ?", projectID).
		Count(&count).Error
	return count, err
}

// ==================== 风险仓库实现 ====================

// ProjectRiskRepositoryImpl 项目风险仓库实现
type ProjectRiskRepositoryImpl struct {
	db *gorm.DB
}

// NewProjectRiskRepository 创建风险仓库实例
func NewProjectRiskRepository(db *gorm.DB) repository.ProjectRiskRepository {
	return &ProjectRiskRepositoryImpl{db: db}
}

// Create 创建风险
func (r *ProjectRiskRepositoryImpl) Create(ctx context.Context, risk *entity.ProjectRisk) error {
	return r.db.WithContext(ctx).Create(risk).Error
}

// GetByID 根据ID获取风险
func (r *ProjectRiskRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.ProjectRisk, error) {
	var risk entity.ProjectRisk
	err := r.db.WithContext(ctx).First(&risk, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &risk, nil
}

// ListByProjectID 获取项目的所有风险
func (r *ProjectRiskRepositoryImpl) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entity.ProjectRisk, error) {
	var risks []*entity.ProjectRisk
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("severity DESC, created_at DESC").
		Find(&risks).Error
	if err != nil {
		return nil, err
	}
	return risks, nil
}

// Update 更新风险
func (r *ProjectRiskRepositoryImpl) Update(ctx context.Context, risk *entity.ProjectRisk) error {
	return r.db.WithContext(ctx).Save(risk).Error
}

// Delete 删除风险
func (r *ProjectRiskRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ProjectRisk{}, "id = ?", id).Error
}

// CountHighRisk 统计高风险数量（严重程度>=12）
func (r *ProjectRiskRepositoryImpl) CountHighRisk(ctx context.Context, projectID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.ProjectRisk{}).
		Where("project_id = ? AND severity >= ?", projectID, 12).
		Count(&count).Error
	return count, err
}

// ==================== 迭代仓库实现 ====================

// IterationRepositoryImpl 迭代仓库实现
type IterationRepositoryImpl struct {
	db *gorm.DB
}

// NewIterationRepository 创建迭代仓库实例
func NewIterationRepository(db *gorm.DB) repository.IterationRepository {
	return &IterationRepositoryImpl{db: db}
}

// Create 创建迭代
func (r *IterationRepositoryImpl) Create(ctx context.Context, iteration *entity.Iteration) error {
	return r.db.WithContext(ctx).Create(iteration).Error
}

// GetByID 根据ID获取迭代
func (r *IterationRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Iteration, error) {
	var iteration entity.Iteration
	err := r.db.WithContext(ctx).First(&iteration, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &iteration, nil
}

// List 分页查询迭代列表
func (r *IterationRepositoryImpl) List(ctx context.Context, page, pageSize int, projectID uuid.UUID, status string) ([]*entity.Iteration, int64, error) {
	var iterations []*entity.Iteration
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Iteration{})

	if projectID != uuid.Nil {
		query = query.Where("project_id = ?", projectID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("sort_order ASC, start_date DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&iterations).Error
	if err != nil {
		return nil, 0, err
	}

	return iterations, total, nil
}

// ListByProjectID 获取项目的所有迭代
func (r *IterationRepositoryImpl) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entity.Iteration, error) {
	var iterations []*entity.Iteration
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("sort_order ASC, start_date DESC").
		Find(&iterations).Error
	if err != nil {
		return nil, err
	}
	return iterations, nil
}

// Update 更新迭代
func (r *IterationRepositoryImpl) Update(ctx context.Context, iteration *entity.Iteration) error {
	return r.db.WithContext(ctx).Save(iteration).Error
}

// Delete 删除迭代
func (r *IterationRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Iteration{}, "id = ?", id).Error
}

// UpdateStatus 更新迭代状态
func (r *IterationRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).
		Model(&entity.Iteration{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// GetActiveIteration 获取项目当前活跃的迭代
func (r *IterationRepositoryImpl) GetActiveIteration(ctx context.Context, projectID uuid.UUID) (*entity.Iteration, error) {
	var iteration entity.Iteration
	err := r.db.WithContext(ctx).
		Where("project_id = ? AND status = ?", projectID, "active").
		First(&iteration).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &iteration, nil
}

// ListCompleted 获取已完成的迭代列表
func (r *IterationRepositoryImpl) ListCompleted(ctx context.Context, projectID uuid.UUID) ([]*entity.Iteration, error) {
	var iterations []*entity.Iteration
	err := r.db.WithContext(ctx).
		Where("project_id = ? AND status = ?", projectID, "completed").
		Order("end_date DESC").
		Find(&iterations).Error
	if err != nil {
		return nil, err
	}
	return iterations, nil
}
