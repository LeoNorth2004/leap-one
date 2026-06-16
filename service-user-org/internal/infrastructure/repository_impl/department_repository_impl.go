package repository_impl

import (
	"context"

	"leap-one/service-user-org/internal/domain/entity"
	"leap-one/service-user-org/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DepartmentRepositoryImpl 部门仓库实现（支持树形结构管理）
type DepartmentRepositoryImpl struct {
	db *gorm.DB
}

// NewDepartmentRepository 创建部门仓库实例
func NewDepartmentRepository(db *gorm.DB) repository.DepartmentRepository {
	return &DepartmentRepositoryImpl{db: db}
}

// Create 创建部门
func (r *DepartmentRepositoryImpl) Create(ctx context.Context, dept *entity.Department) error {
	return r.db.WithContext(ctx).Create(dept).Error
}

// GetByID 根据ID获取部门
func (r *DepartmentRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Department, error) {
	var dept entity.Department
	err := r.db.WithContext(ctx).First(&dept, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dept, nil
}

// GetByCode 根据编码获取部门
func (r *DepartmentRepositoryImpl) GetByCode(ctx context.Context, code string) (*entity.Department, error) {
	var dept entity.Department
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&dept).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dept, nil
}

// Update 更新部门信息
func (r *DepartmentRepositoryImpl) Update(ctx context.Context, dept *entity.Department) error {
	return r.db.WithContext(ctx).Save(dept).Error
}

// Delete 软删除部门
func (r *DepartmentRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Department{}, "id = ?", id).Error
}

// List 分页查询部门列表
func (r *DepartmentRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*entity.Department, int64, error) {
	var depts []*entity.Department
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Department{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("sort_order ASC, created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&depts).Error

	if err != nil {
		return nil, 0, err
	}

	return depts, total, nil
}

// GetChildren 获取直接子部门列表
func (r *DepartmentRepositoryImpl) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entity.Department, error) {
	var children []*entity.Department
	err := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("sort_order ASC").
		Find(&children).Error
	return children, err
}

// GetTree 获取完整部门树形结构（先查出全部再在内存中构建树）
func (r *DepartmentRepositoryImpl) GetTree(ctx context.Context) ([]*entity.Department, error) {
	var allDepts []*entity.Department
	err := r.db.WithContext(ctx).
		Where("status = 1").
		Order("sort_order ASC, level ASC, created_at ASC").
		Find(&allDepts).Error
	if err != nil {
		return nil, err
	}
	return buildDeptTree(allDepts, uuid.Nil), nil
}

// Move 移动部门（变更父级部门并更新层级）
func (r *DepartmentRepositoryImpl) Move(ctx context.Context, deptID uuid.UUID, newParentID uuid.UUID) error {
	dept, err := r.GetByID(ctx, deptID)
	if err != nil || dept == nil {
		return err
	}

	newLevel := 1
	if newParentID != uuid.Nil {
		parent, pErr := r.GetByID(ctx, newParentID)
		if pErr != nil || parent == nil {
			return pErr
		}
		newLevel = parent.Level + 1
	}

	var parentIDPtr *uuid.UUID
	if newParentID != uuid.Nil {
		parentIDPtr = &newParentID
	}

	return r.db.WithContext(ctx).
		Model(&entity.Department{}).
		Where("id = ?", deptID).
		Updates(map[string]interface{}{
			"parent_id": parentIDPtr,
			"level":     newLevel,
		}).Error
}

// HasChildren 检查部门是否有子部门
func (r *DepartmentRepositoryImpl) HasChildren(ctx context.Context, parentID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Department{}).
		Where("parent_id = ?", parentID).
		Count(&count).Error
	return count > 0, err
}

// CountMembers 统计部门下的成员数量
func (r *DepartmentRepositoryImpl) CountMembers(ctx context.Context, deptID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("department_id = ? AND deleted_at IS NULL", deptID).
		Count(&count).Error
	return count, err
}

// GetDepartmentMembers 获取部门成员列表（分页）
func (r *DepartmentRepositoryImpl) GetDepartmentMembers(ctx context.Context, deptID uuid.UUID, page, pageSize int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("department_id = ? AND deleted_at IS NULL", deptID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// buildDeptTree 在内存中递归构建部门树
// rootID为Nil时从顶层开始构建
func buildDeptTree(depts []*entity.Department, rootID uuid.UUID) []*entity.Department {
	tree := make([]*entity.Department, 0)
	// 使用map建立索引加速查找
	deptMap := make(map[uuid.UUID]*entity.Department, len(depts))
	childrenMap := make(map[uuid.UUID][]entity.Department)

	for _, d := range depts {
		copy := *d // 浅拷贝避免修改原数据
		deptMap[d.ID] = &copy
		if d.ParentID != nil {
			childrenMap[*d.ParentID] = append(childrenMap[*d.ParentID], copy)
		} else {
			childrenMap[uuid.Nil] = append(childrenMap[uuid.Nil], copy)
		}
	}

	// 从rootID的子节点开始构建
	for i := range childrenMap[rootID] {
		buildSubTree(&childrenMap[rootID][i], childrenMap)
		tree = append(tree, &childrenMap[rootID][i])
	}

	return tree
}

// buildSubTree 递归构建子树
func buildSubTree(node *entity.Department, childrenMap map[uuid.UUID][]entity.Department) {
	node.Children = childrenMap[node.ID]
	for i := range node.Children {
		buildSubTree(&node.Children[i], childrenMap)
	}
}
