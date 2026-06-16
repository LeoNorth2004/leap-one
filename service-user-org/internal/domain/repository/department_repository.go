package repository

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-user-org/internal/domain/entity"
)

// DepartmentRepository 部门仓库接口定义
type DepartmentRepository interface {
	// Create 创建部门
	Create(ctx context.Context, dept *entity.Department) error

	// GetByID 根据ID获取部门
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Department, error)

	// GetByCode 根据编码获取部门
	GetByCode(ctx context.Context, code string) (*entity.Department, error)

	// Update 更新部门信息
	Update(ctx context.Context, dept *entity.Department) error

	// Delete 软删除部门
	Delete(ctx context.Context, id uuid.UUID) error

	// List 分页查询部门列表
	List(ctx context.Context, page, pageSize int) ([]*entity.Department, int64, error)

	// GetChildren 获取子部门列表
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entity.Department, error)

	// GetTree 获取完整部门树形结构
	GetTree(ctx context.Context) ([]*entity.Department, error)

	// CountMembers 统计部门成员数量
	CountMembers(ctx context.Context, deptID uuid.UUID) (int64, error)

	// HasChildren 检查部门是否有子部门
	HasChildren(ctx context.Context, deptID uuid.UUID) (bool, error)

	// GetDepartmentMembers 获取部门成员列表（分页）
	GetDepartmentMembers(ctx context.Context, deptID uuid.UUID, page, pageSize int) ([]*entity.User, int64, error)

	// Move 移动部门到新的父部门下
	Move(ctx context.Context, id, newParentID uuid.UUID) error
}
