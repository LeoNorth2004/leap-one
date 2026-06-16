package repository

import (
	"github.com/google/uuid"
	"leap-one/service-requirement/internal/domain/entity"
)

// RequirementRepository 需求仓储接�?
type RequirementRepository interface {
	// Create 创建需�?
	Create(req *entity.Requirement) error
	// GetByID 根据ID获取需�?
	GetByID(id uuid.UUID) (*entity.Requirement, error)
	// Update 更新需�?
	Update(req *entity.Requirement) error
	// Delete 删除需求（软删除）
	Delete(id uuid.UUID) error
	// List 分页查询需求列�?
	List(params *RequirementListParams) ([]*entity.Requirement, int64, error)
	// GetTree 获取需求树形结构（按产品维度）
	GetTree(productID uuid.UUID) ([]*entity.Requirement, error)
	// GetChildren 获取子需求列�?
	GetChildren(parentID uuid.UUID) ([]*entity.Requirement, error)
	// UpdateStatus 更新需求状�?
	UpdateStatus(id uuid.UUID, status string) error
	// GenerateCode 生成下一个需求编�?
	GenerateCode() (string, error)
}

// RequirementListParams 需求列表查询参�?
type RequirementListParams struct {
	Page      int
	PageSize  int
	ProductID *uuid.UUID
	ProjectID *uuid.UUID
	Type      string
	Status    string
	Priority  *int
	OwnerID   *uuid.UUID
	Category  string
	Stage     string
	Keyword   string // 模糊搜索标题/描述/编号
	SortBy    string // 排序字段
	SortOrder string // asc/desc
}
