package repository_impl

import (
	"context"

	"leap-one/service-project/internal/domain/entity"
	"leap-one/service-project/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProjectRepositoryImpl 项目仓库实现
type ProjectRepositoryImpl struct {
	db *gorm.DB
}

// NewProjectRepository 创建项目仓库实例
func NewProjectRepository(db *gorm.DB) repository.ProjectRepository {
	return &ProjectRepositoryImpl{db: db}
}

// Create 创建项目
func (r *ProjectRepositoryImpl) Create(ctx context.Context, project *entity.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

// GetByID 根据ID获取项目
func (r *ProjectRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Project, error) {
	var project entity.Project
	err := r.db.WithContext(ctx).First(&project, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

// GetByCode 根据项目编号获取项目
func (r *ProjectRepositoryImpl) GetByCode(ctx context.Context, code string) (*entity.Project, error) {
	var project entity.Project
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&project).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

// Update 更新项目信息
func (r *ProjectRepositoryImpl) Update(ctx context.Context, project *entity.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

// Delete 软删除项�?
func (r *ProjectRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Project{}, "id = ?", id).Error
}

// List 分页查询项目列表（支持筛选和搜索）
func (r *ProjectRepositoryImpl) List(
	ctx context.Context,
	page, pageSize int,
	keyword, status, programID, pmID, sortBy, sortOrder string,
) ([]*entity.Project, int64, error) {
	var projects []*entity.Project
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Project{})

	// 关键词搜索：匹配项目名称、编号、描�?
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where(
			"name LIKE ? OR code LIKE ? OR description LIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// 状态筛�?
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 项目集筛�?
	if programID != "" {
		pid, parseErr := uuid.Parse(programID)
		if parseErr == nil {
			query = query.Where("program_id = ?", pid)
		}
	}

	// 项目经理筛�?
	if pmID != "" {
		pmid, parseErr := uuid.Parse(pmID)
		if parseErr == nil {
			query = query.Where("pm_id = ?", pmid)
		}
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序处理
	orderClause := buildOrderClause(sortBy, sortOrder, "created_at", "DESC")
	query = query.Order(orderClause)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&projects).Error
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

// UpdateStatus 更新项目状�?
func (r *ProjectRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).
		Model(&entity.Project{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// ListByPMID 查询某项目经理下的所有项目
func (r *ProjectRepositoryImpl) ListByPMID(ctx context.Context, pmID uuid.UUID) ([]*entity.Project, error) {
	var projects []*entity.Project
	err := r.db.WithContext(ctx).
		Where("pm_id = ?", pmID).
		Order("created_at DESC").
		Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// CountByStatus 按状态统计项目数量
func (r *ProjectRepositoryImpl) CountByStatus(ctx context.Context) (map[string]int64, error) {
	type StatusCount struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var results []StatusCount

	err := r.db.WithContext(ctx).
		Model(&entity.Project{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	resultMap := make(map[string]int64)
	for _, item := range results {
		resultMap[item.Status] = item.Count
	}
	return resultMap, nil
}
