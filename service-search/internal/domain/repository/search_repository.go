package repository

import ("context"; "github.com/google/uuid"; "leap-one/service-search/internal/domain/entity")

type SearchDocumentRepository interface{
	Create(ctx context.Context, doc *entity.SearchDocument) error
	BatchCreate(ctx context.Context, docs []*entity.SearchDocument) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.SearchDocument, error)
	GetByRefID(ctx context.Context, docType string, refID uuid.UUID) (*entity.SearchDocument, error)
	Search(ctx context.Context, query string, docTypes []string, page, pageSize int) ([]*entity.SearchDocument, int64, error)
	AdvancedSearch(ctx context.Context, query string, filters map[string]interface{}, page, pageSize int) ([]*entity.SearchDocument, int64, error)
	Update(ctx context.Context, doc *entity.SearchDocument) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByRefID(ctx context.Context, docType string, refID uuid.UUID) error
	GetSuggestions(ctx context.Context, prefix string, limit int) ([]string, error)
}

type SavedSearchRepository interface{
	Create(ctx context.Context, s *entity.SavedSearch) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.SavedSearch, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.SavedSearch, error)
	Update(ctx context.Context, s *entity.SavedSearch) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SearchHistoryRepository interface{
	Create(ctx context.Context, h *entity.SearchHistory) error
	ListByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.SearchHistory, error)
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
