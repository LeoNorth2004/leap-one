package repository

import (
	"github.com/google/uuid"
	"leap-one/service-document/internal/domain/entity"
)

// DocumentRepository 文档仓储接口
type DocumentRepository interface {
	Create(doc *entity.Document) error
	GetByID(id uuid.UUID) (*entity.Document, error)
	Update(doc *entity.Document) error
	Delete(id uuid.UUID) error
	List(params *DocumentListParams) ([]*entity.Document, int64, error)
	GetTree(projectID uuid.UUID) ([]*entity.Document, error)
	GetChildren(parentID uuid.UUID) ([]*entity.Document, error)
	UpdateStatus(id uuid.UUID, status string) error
	Publish(id uuid.UUID) error
	Search(keyword string) ([]*entity.Document, error)
}

// DocumentListParams 文档列表查询参数
type DocumentListParams struct {
	Page       int
	PageSize   int
	ProductID  *uuid.UUID
	ProjectID  *uuid.UUID
	CategoryID *uuid.UUID
	Type       string
	Status     string
	Visibility string
	OwnerID    *uuid.UUID
	IsTemplate *bool
	Keyword    string
	SortBy     string
	SortOrder  string
}
