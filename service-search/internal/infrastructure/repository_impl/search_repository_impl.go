package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-search/internal/domain/entity"
	"leap-one/service-search/internal/domain/repository"
	"gorm.io/gorm"
)

type SearchDocumentRepositoryImpl struct{ db *gorm.DB }

func NewSearchDocumentRepository(db *gorm.DB) repository.SearchDocumentRepository {
	return &SearchDocumentRepositoryImpl{db: db}
}
func (r *SearchDocumentRepositoryImpl) Create(ctx context.Context, d *entity.SearchDocument) error {
	return r.db.WithContext(ctx).Create(d).Error
}
func (r *SearchDocumentRepositoryImpl) BatchCreate(ctx context.Context, docs []*entity.SearchDocument) error {
	if len(docs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&docs).Error
}
func (r *SearchDocumentRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.SearchDocument, error) {
	var d entity.SearchDocument
	err := r.db.WithContext(ctx).First(&d, "id=?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}
func (r *SearchDocumentRepositoryImpl) GetByRefID(ctx context.Context, docType string, refID uuid.UUID) (*entity.SearchDocument, error) {
	var d entity.SearchDocument
	err := r.db.WithContext(ctx).Where("doc_type=? AND ref_id=?", docType, refID).First(&d).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}
func (r *SearchDocumentRepositoryImpl) Search(ctx context.Context, query string, docTypes []string, page, pageSize int) ([]*entity.SearchDocument, int64, error) {
	var list []*entity.SearchDocument
	var total int64
	q := r.db.WithContext(ctx).Model(&entity.SearchDocument{})
	pattern := "%" + query + "%"
	q = q.Where("title ILIKE ? OR content ILIKE ? OR summary ILIKE ?", pattern, pattern, pattern)
	if len(docTypes) > 0 {
		q = q.Where("doc_type IN ?", docTypes)
	}
	q.Count(&total)
	offset := (page - 1) * pageSize
	if err := q.Order("indexed_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
func (r *SearchDocumentRepositoryImpl) AdvancedSearch(ctx context.Context, query string, filters map[string]interface{}, page, pageSize int) ([]*entity.SearchDocument, int64, error) {
	var list []*entity.SearchDocument
	var total int64
	q := r.db.WithContext(ctx).Model(&entity.SearchDocument{})
	if query != "" {
		p := "%" + query + "%"
		q = q.Where("title ILIKE ? OR content ILIKE ?", p, p)
	}
	for k, v := range filters {
		q = q.Where(k+" = ?", v)
	}
	q.Count(&total)
	offset := (page - 1) * pageSize
	if err := q.Order("indexed_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
func (r *SearchDocumentRepositoryImpl) Update(ctx context.Context, d *entity.SearchDocument) error {
	return r.db.WithContext(ctx).Save(d).Error
}
func (r *SearchDocumentRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.SearchDocument{}, "id=?", id).Error
}
func (r *SearchDocumentRepositoryImpl) DeleteByRefID(ctx context.Context, docType string, refID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("doc_type=? AND ref_id=?", docType, refID).Delete(&entity.SearchDocument{}).Error
}
func (r *SearchDocumentRepositoryImpl) GetSuggestions(ctx context.Context, prefix string, limit int) ([]string, error) {
	var titles []string
	err := r.db.WithContext(ctx).Model(&entity.SearchDocument{}).Where("title ILIKE ?", prefix+"%").Distinct("title").Limit(limit).Pluck("title", &titles).Error
	return titles, err
}

type SavedSearchRepositoryImpl struct{ db *gorm.DB }

func NewSavedSearchRepository(db *gorm.DB) repository.SavedSearchRepository {
	return &SavedSearchRepositoryImpl{db: db}
}
func (r *SavedSearchRepositoryImpl) Create(ctx context.Context, s *entity.SavedSearch) error {
	return r.db.WithContext(ctx).Create(s).Error
}
func (r *SavedSearchRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.SavedSearch, error) {
	var s entity.SavedSearch
	err := r.db.WithContext(ctx).First(&s, "id=?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}
func (r *SavedSearchRepositoryImpl) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.SavedSearch, error) {
	var list []*entity.SavedSearch
	err := r.db.WithContext(ctx).Where("user_id=?", userID).Order("created_at DESC").Find(&list).Error
	return list, err
}
func (r *SavedSearchRepositoryImpl) Update(ctx context.Context, s *entity.SavedSearch) error {
	return r.db.WithContext(ctx).Save(s).Error
}
func (r *SavedSearchRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.SavedSearch{}, "id=?", id).Error
}

type SearchHistoryRepositoryImpl struct{ db *gorm.DB }

func NewSearchHistoryRepository(db *gorm.DB) repository.SearchHistoryRepository {
	return &SearchHistoryRepositoryImpl{db: db}
}
func (r *SearchHistoryRepositoryImpl) Create(ctx context.Context, h *entity.SearchHistory) error {
	return r.db.WithContext(ctx).Create(h).Error
}
func (r *SearchHistoryRepositoryImpl) ListByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.SearchHistory, error) {
	var list []*entity.SearchHistory
	err := r.db.WithContext(ctx).Where("user_id=?", userID).Order("searched_at DESC").Limit(limit).Find(&list).Error
	return list, err
}
func (r *SearchHistoryRepositoryImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id=?", userID).Delete(&entity.SearchHistory{}).Error
}
