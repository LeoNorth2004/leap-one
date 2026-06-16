package repository_impl

import (
	"context"

	"leap-one/service-portfolio/internal/domain/entity"
	"leap-one/service-portfolio/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProgramRepositoryImpl 项目集仓库实现（支持树形结构管理）
type ProgramRepositoryImpl struct {
	db *gorm.DB
}

// NewProgramRepository 创建项目集仓库实例
func NewProgramRepository(db *gorm.DB) repository.ProgramRepository {
	return &ProgramRepositoryImpl{db: db}
}

// Create 创建项目集
func (r *ProgramRepositoryImpl) Create(ctx context.Context, program *entity.Program) error {
	return r.db.WithContext(ctx).Create(program).Error
}

// GetByID 根据ID获取项目集
func (r *ProgramRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Program, error) {
	var program entity.Program
	err := r.db.WithContext(ctx).First(&program, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &program, nil
}

// GetByCode 根据编号获取项目集
func (r *ProgramRepositoryImpl) GetByCode(ctx context.Context, code string) (*entity.Program, error) {
	var program entity.Program
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&program).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &program, nil
}

// Update 更新项目集信息
func (r *ProgramRepositoryImpl) Update(ctx context.Context, program *entity.Program) error {
	return r.db.WithContext(ctx).Save(program).Error
}

// Delete 软删除项目集
func (r *ProgramRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Program{}, "id = ?", id).Error
}

// List 分页查询项目集列表（支持关键词搜索和状态筛选）
func (r *ProgramRepositoryImpl) List(ctx context.Context, page, pageSize int, keyword, status string) ([]*entity.Program, int64, error) {
	var programs []*entity.Program
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Program{})

	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR description LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("priority ASC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&programs).Error

	if err != nil {
		return nil, 0, err
	}

	return programs, total, nil
}

// GetTree 获取完整项目集树形结构（先查出全部再在内存中构建树）
func (r *ProgramRepositoryImpl) GetTree(ctx context.Context) ([]*entity.Program, error) {
	var allPrograms []*entity.Program
	err := r.db.WithContext(ctx).
		Order("priority ASC, created_at ASC").
		Find(&allPrograms).Error
	if err != nil {
		return nil, err
	}
	return buildProgramTree(allPrograms, uuid.Nil), nil
}

// GetChildren 获取直接子项目集列表
func (r *ProgramRepositoryImpl) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entity.Program, error) {
	var children []*entity.Program
	err := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("priority ASC, created_at ASC").
		Find(&children).Error
	return children, err
}

// GetByOwnerID 根据负责人ID获取项目集列表
func (r *ProgramRepositoryImpl) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*entity.Program, error) {
	var programs []*entity.Program
	err := r.db.WithContext(ctx).
		Where("owner_id = ?", ownerID).
		Order("priority ASC, created_at DESC").
		Find(&programs).Error
	return programs, err
}

// buildProgramTree 在内存中递归构建项目集树
func buildProgramTree(programs []*entity.Program, rootID uuid.UUID) []*entity.Program {
	tree := make([]*entity.Program, 0)
	programMap := make(map[uuid.UUID]*entity.Program)
	childrenMap := make(map[uuid.UUID][]entity.Program)

	for _, p := range programs {
		copy := *p
		programMap[p.ID] = &copy
		if p.ParentID != nil {
			childrenMap[*p.ParentID] = append(childrenMap[*p.ParentID], copy)
		} else {
			childrenMap[uuid.Nil] = append(childrenMap[uuid.Nil], copy)
		}
	}

	for i := range childrenMap[rootID] {
		buildProgramSubTree(&childrenMap[rootID][i], childrenMap)
		tree = append(tree, &childrenMap[rootID][i])
	}

	return tree
}

// buildProgramSubTree 递归构建子树
func buildProgramSubTree(node *entity.Program, childrenMap map[uuid.UUID][]entity.Program) {
	node.Children = childrenMap[node.ID]
	for i := range node.Children {
		buildProgramSubTree(&node.Children[i], childrenMap)
	}
}
