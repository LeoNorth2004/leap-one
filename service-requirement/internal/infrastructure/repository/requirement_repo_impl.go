package repository

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"leap-one/service-requirement/internal/domain/entity"
	"leap-one/service-requirement/internal/domain/repository"
)

// requirementRepository 需求仓储实�?type requirementRepository struct {
	db *gorm.DB
}

// NewRequirementRepository 创建需求仓储实�?func NewRequirementRepository(db *gorm.DB) repository.RequirementRepository {
	return &requirementRepository{db: db}
}

func (r *requirementRepository) Create(req *entity.Requirement) error {
	return r.db.Create(req).Error
}

func (r *requirementRepository) GetByID(id uuid.UUID) (*entity.Requirement, error) {
	var req entity.Requirement
	err := r.db.Preload("Children").Preload("Relations").Preload("Reviews").Preload("ChangeLogs").
		Where("id = ? AND deleted_at IS NULL", id).First(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *requirementRepository) Update(req *entity.Requirement) error {
	return r.db.Save(req).Error
}

func (r *requirementRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Requirement{}, "id = ?", id).Error
}

func (r *requirementRepository) List(params *repository.RequirementListParams) ([]*entity.Requirement, int64, error) {
	var requirements []*entity.Requirement
	var total int64

	query := r.db.Model(&entity.Requirement{}).Where("deleted_at IS NULL")

	query = applyRequirementFilters(query, params)

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("title LIKE ? OR code LIKE ? OR description LIKE ?", keyword, keyword, keyword)
	}

	// 排序
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
 sortOrder := strings.ToUpper(params.SortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}
	query = query.Order(sortBy + " " + sortOrder)

	// 总数统计
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	page := params.Page
	if page <= 0 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	err := query.Offset(offset).Limit(pageSize).Find(&requirements).Error
	if err != nil {
		return nil, 0, err
	}
	return requirements, total, nil
}

func (r *requirementRepository) GetTree(productID uuid.UUID) ([]*entity.Requirement, error) {
	var requirements []*entity.Requirement
	err := r.db.Where("product_id = ? AND deleted_at IS NULL", productID).
		Order("level ASC, created_at ASC").Find(&requirements).Error
	if err != nil {
		return nil, err
	}
	return buildRequirementTree(requirements), nil
}

func (r *requirementRepository) GetChildren(parentID uuid.UUID) ([]*entity.Requirement, error) {
	var children []*entity.Requirement
	err := r.db.Where("parent_id = ? AND deleted_at IS NULL", parentID).
		Order("level ASC, created_at ASC").Find(&children).Error
	return children, err
}

func (r *requirementRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&entity.Requirement{}).Where("id = ?", id).Update("status", status).Error
}

func (r *requirementRepository) GenerateCode() (string, error) {
	var maxCode string
	err := r.db.Model(&entity.Requirement{}).
		Select("COALESCE(MAX(code), 'REQ-000')").
		Row().Scan(&maxCode)
	if err != nil {
		return "", fmt.Errorf("生成需求编号失�? %w", err)
	}
	// 解析最大编号并递增
	num := 0
	fmt.Sscanf(maxCode, "REQ-%d", &num)
	return fmt.Sprintf("REQ-%03d", num+1), nil
}

// applyRequirementFilters 应用查询过滤条件
func applyRequirementFilters(query *gorm.DB, params *repository.RequirementListParams) *gorm.DB {
	if params.ProductID != nil {
		query = query.Where("product_id = ?", *params.ProductID)
	}
	if params.ProjectID != nil {
		query = query.Where("project_id = ?", *params.ProjectID)
	}
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Priority != nil {
		query = query.Where("priority = ?", *params.Priority)
	}
	if params.OwnerID != nil {
		query = query.Where("owner_id = ?", *params.OwnerID)
	}
	if params.Category != "" {
		query = query.Where("category = ?", params.Category)
	}
	if params.Stage != "" {
		query = query.Where("stage = ?", params.Stage)
	}
	return query
}

// buildRequirementTree 构建需求树形结�?func buildRequirementTree(requirements []*entity.Requirement) []*entity.Requirement {
	reqMap := make(map[uuid.UUID]*entity.Requirement)
	var roots []*entity.Requirement

	for _, req := range requirements {
		reqMap[req.ID] = req
	}

	for _, req := range requirements {
		if req.ParentID != nil {
			if parent, ok := reqMap[*req.ParentID]; ok {
				parent.Children = append(parent.Children, *req)
			}
		} else {
			roots = append(roots, req)
		}
	}
	return roots
}
