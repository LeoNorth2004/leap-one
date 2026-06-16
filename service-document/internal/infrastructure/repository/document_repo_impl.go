package repository

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"leap-one/service-document/internal/domain/entity"
	"leap-one/service-document/internal/domain/repository"
)

// documentRepository 文档仓储实现
type documentRepository struct {
	db *gorm.DB
}

// NewDocumentRepository 创建文档仓储实例
func NewDocumentRepository(db *gorm.DB) repository.DocumentRepository {
	return &documentRepository{db: db}
}

func (r *documentRepository) Create(doc *entity.Document) error { return r.db.Create(doc).Error }

func (r *documentRepository) GetByID(id uuid.UUID) (*entity.Document, error) {
	var doc entity.Document
	err := r.db.Preload("Children").Preload("Versions").Preload("Attachments").
		Where("id = ? AND deleted_at IS NULL", id).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *documentRepository) Update(doc *entity.Document) error { return r.db.Save(doc).Error }

func (r *documentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Document{}, "id = ?", id).Error
}

func (r *documentRepository) List(params *repository.DocumentListParams) ([]*entity.Document, int64, error) {
	var docs []*entity.Document
	var total int64

	query := r.db.Model(&entity.Document{}).Where("deleted_at IS NULL")
	query = applyDocFilters(query, params)

	if params.Keyword != "" {
		kw := "%" + params.Keyword + "%"
		query = query.Where("title LIKE ? OR content LIKE ?", kw, kw)
	}

	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "updated_at"
	}
	sortOrder := strings.ToUpper(params.SortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}
	query = query.Order(sortBy + " " + sortOrder)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := params.Page
	if page <= 0 {
		page = 1
	}
	size := params.PageSize
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	err := query.Offset((page - 1) * size).Limit(size).Find(&docs).Error
	return docs, total, err
}

func (r *documentRepository) GetTree(projectID uuid.UUID) ([]*entity.Document, error) {
	var docs []*entity.Document
	err := r.db.Where("project_id = ? AND deleted_at IS NULL", projectID).
		Order("order_index ASC, created_at ASC").Find(&docs).Error
	if err != nil {
		return nil, err
	}
	return buildDocTree(docs), nil
}

func (r *documentRepository) GetChildren(parentID uuid.UUID) ([]*entity.Document, error) {
	var children []*entity.Document
	err := r.db.Where("parent_id = ? AND deleted_at IS NULL", parentID).
		Order("order_index ASC").Find(&children).Error
	return children, err
}

func (r *documentRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&entity.Document{}).Where("id = ?", id).Update("status", status).Error
}

func (r *documentRepository) Publish(id uuid.UUID) error {
	return r.db.Model(&entity.Document{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": "published", "version": gorm.Expr("version + 1")}).Error
}

func (r *documentRepository) Search(keyword string) ([]*entity.Document, error) {
	var docs []*entity.Document
	kw := "%" + keyword + "%"
	err := r.db.Where("(title LIKE ? OR content LIKE ? OR tags LIKE ?) AND deleted_at IS NULL AND status = ?", kw, kw, kw, "published").
		Order("updated_at DESC").Limit(50).Find(&docs).Error
	return docs, err
}

func applyDocFilters(query *gorm.DB, p *repository.DocumentListParams) *gorm.DB {
	if p.ProductID != nil {
		query = query.Where("product_id = ?", *p.ProductID)
	}
	if p.ProjectID != nil {
		query = query.Where("project_id = ?", *p.ProjectID)
	}
	if p.CategoryID != nil {
		query = query.Where("category_id = ?", *p.CategoryID)
	}
	if p.Type != "" {
		query = query.Where("type = ?", p.Type)
	}
	if p.Status != "" {
		query = query.Where("status = ?", p.Status)
	}
	if p.Visibility != "" {
		query = query.Where("visibility = ?", p.Visibility)
	}
	if p.OwnerID != nil {
		query = query.Where("owner_id = ?", *p.OwnerID)
	}
	if p.IsTemplate != nil {
		query = query.Where("is_template = ?", *p.IsTemplate)
	}
	return query
}

func buildDocTree(docs []*entity.Document) []*entity.Document {
	m := make(map[uuid.UUID]*entity.Document)
	var roots []*entity.Document
	for _, d := range docs {
		m[d.ID] = d
	}
	for _, d := range docs {
		if d.ParentID != nil {
			if parent, ok := m[*d.ParentID]; ok {
				parent.Children = append(parent.Children, *d)
			}
		} else {
			roots = append(roots, d)
		}
	}
	return roots
}
