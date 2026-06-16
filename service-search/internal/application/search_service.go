package application

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-search/internal/domain/entity"
	"leap-one/service-search/internal/domain/repository"
	"go.uber.org/zap"
)

type SearchService struct {
	docRepo     repository.SearchDocumentRepository
	savedRepo   repository.SavedSearchRepository
	historyRepo repository.SearchHistoryRepository
	logger      *zap.Logger
}

func NewSearchService(docRepo repository.SearchDocumentRepository, savedRepo repository.SavedSearchRepository, historyRepo repository.SearchHistoryRepository, logger *zap.Logger) *SearchService {
	return &SearchService{docRepo: docRepo, savedRepo: savedRepo, historyRepo: historyRepo, logger: logger}
}
func (s *SearchService) IndexDocumentUseCase(ctx context.Context, docType string, refID uuid.UUID, title, content, summary, tags, metaData string) error {
	doc := &entity.SearchDocument{DocType: docType, RefID: refID, Title: title, Content: content, Summary: summary, Tags: tags, MetaData: metaData}
	existing, _ := s.docRepo.GetByRefID(ctx, docType, refID)
	if existing != nil {
		doc.ID = existing.ID
		return s.docRepo.Update(ctx, doc)
	}
	return s.docRepo.Create(ctx, doc)
}
func (s *SearchService) RemoveDocumentUseCase(ctx context.Context, docType string, refID uuid.UUID) error {
	return s.docRepo.DeleteByRefID(ctx, docType, refID)
}
