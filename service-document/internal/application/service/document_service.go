package service

import (
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-document/internal/domain/entity"
	"leap-one/service-document/internal/domain/repository"
)

// DocumentService 文档应用服务
type DocumentService struct {
	docRepo    repository.DocumentRepository
	versionRepo repository.DocumentVersionRepository
	favRepo    repository.DocumentFavoriteRepository
	logger     *zap.Logger
}

func NewDocumentService(docRepo repository.DocumentRepository, versionRepo repository.DocumentVersionRepository, favRepo repository.DocumentFavoriteRepository, logger *zap.Logger) *DocumentService {
	return &DocumentService{docRepo: docRepo, versionRepo: versionRepo, favRepo: favRepo, logger: logger}
}

func (s *DocumentService) Create(doc *entity.Document) (*entity.Document, error) {
	if doc.Type == "" { doc.Type = "markdown" }
	if doc.Status == "" { doc.Status = "draft" }
	if doc.Visibility == "" { doc.Visibility = "public" }
	if err := s.docRepo.Create(doc); err != nil { return nil, err }

	// 创建初始版本
	version := &entity.DocumentVersion{
		ID: uuid.New(), DocumentID: doc.ID, VersionNo: 1,
		Title: doc.Title, Content: doc.Content, ChangeNote: "初始版本", CreatedBy: doc.OwnerID,
	}
	s.versionRepo.Create(version)
	return doc, nil
}

func (s *DocumentService) GetByID(id uuid.UUID) (*entity.Document, error) { return s.docRepo.GetByID(id) }

func (s *DocumentService) Update(id uuid.UUID, updates map[string]interface{}) (*entity.Document, error) {
	doc, err := s.docRepo.GetByID(id)
	if err != nil { return nil, err }
	for field, val := range updates {
		switch field {
		case "title": if v, ok := val.(string); ok { doc.Title = v; doc.Content = v } // 简化处�?		case "content": if v, ok := val.(string); ok { doc.Content = v }
		case "status": if v, ok := val.(string); ok { doc.Status = v }
		case "visibility": if v, ok := val.(string); ok { doc.Visibility = v }
		case "tags": if v, ok := val.(string); ok { doc.Tags = v }
		case "category_id": if v, ok := val.(uuid.UUID); ok { doc.CategoryID = &v }
		}
	}
	doc.Version++
	if err := s.docRepo.Update(doc); err != nil { return nil, err }

	// 保存版本快照
	version := &entity.DocumentVersion{
		ID: uuid.New(), DocumentID: id, VersionNo: doc.Version,
		Title: doc.Title, Content: doc.Content, ChangeNote: "更新保存", CreatedBy: doc.OwnerID,
	}
	s.versionRepo.Create(version)
	return doc, nil
}

func (s *DocumentService) Delete(id uuid.UUID) error { return s.docRepo.Delete(id) }
func (s *DocumentService) List(params *repository.DocumentListParams) ([]*entity.Document, int64, error) {
	return s.docRepo.List(params)
}
func (s *DocumentService) GetTree(projectID uuid.UUID) ([]*entity.Document, error) { return s.docRepo.GetTree(projectID) }
func (s *DocumentService) Publish(id uuid.UUID) error { return s.docRepo.Publish(id) }
func (s *DocumentService) Search(keyword string) ([]*entity.Document, error) { return s.docRepo.Search(keyword) }
func (s *DocumentService) AddFavorite(userID, docID uuid.UUID) error { return s.favRepo.Add(userID, docID) }
func (s *DocumentService) RemoveFavorite(userID, docID uuid.UUID) error { return s.favRepo.Remove(userID, docID) }
