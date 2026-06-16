package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"leap-one/service-document/internal/domain/entity"
	repo "leap-one/service-document/internal/domain/repository"
)

// documentVersionRepository 文档版本仓储实现
type documentVersionRepository struct{ db *gorm.DB }

func NewDocumentVersionRepository(db *gorm.DB) repo.DocumentVersionRepository {
	return &documentVersionRepository{db: db}
}

func (r *documentVersionRepository) Create(v *entity.DocumentVersion) error {
	return r.db.Create(v).Error
}

func (r *documentVersionRepository) GetByDocumentAndVersion(docID uuid.UUID, ver int) (*entity.DocumentVersion, error) {
	var v entity.DocumentVersion
	err := r.db.Where("document_id = ? AND version_no = ?", docID, ver).First(&v).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *documentVersionRepository) ListByDocumentID(docID uuid.UUID) ([]*entity.DocumentVersion, error) {
	var vs []*entity.DocumentVersion
	err := r.db.Where("document_id = ?", docID).Order("version_no DESC").Find(&vs).Error
	return vs, err
}

func (r *documentVersionRepository) GetLatest(docID uuid.UUID) (*entity.DocumentVersion, error) {
	var v entity.DocumentVersion
	err := r.db.Where("document_id = ?", docID).Order("version_no DESC").First(&v).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// ==================== 其他仓储实现 ====================

// categoryRepoImpl
type categoryRepoImpl struct{ db *gorm.DB }

func NewCategoryRepository(db *gorm.DB) repo.DocumentCategoryRepository {
	return &categoryRepoImpl{db: db}
}
func (r *categoryRepoImpl) Create(c *entity.DocumentCategory) error { return r.db.Create(c).Error }
func (r *categoryRepoImpl) GetByID(id uuid.UUID) (*entity.DocumentCategory, error) {
	var c entity.DocumentCategory
	err := r.db.First(&c, "id = ?", id).Error
	return &c, err
}
func (r *categoryRepoImpl) Update(c *entity.DocumentCategory) error { return r.db.Save(c).Error }
func (r *categoryRepoImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.DocumentCategory{}, "id = ?", id).Error
}
func (r *categoryRepoImpl) List() ([]*entity.DocumentCategory, error) {
	var list []*entity.DocumentCategory
	err := r.db.Order("sort_order ASC").Find(&list).Error
	return list, err
}
func (r *categoryRepoImpl) GetTree() ([]*entity.DocumentCategory, error) {
	var all []*entity.DocumentCategory
	_ = r.db.Find(&all)
	return buildCategoryTree(all), nil
}

// attachmentRepoImpl
type attachmentRepoImpl struct{ db *gorm.DB }

func NewAttachmentRepository(db *gorm.DB) repo.DocumentAttachmentRepository {
	return &attachmentRepoImpl{db: db}
}
func (r *attachmentRepoImpl) Create(a *entity.DocumentAttachment) error { return r.db.Create(a).Error }
func (r *attachmentRepoImpl) GetByID(id uuid.UUID) (*entity.DocumentAttachment, error) {
	var a entity.DocumentAttachment
	err := r.db.First(&a, "id = ?", id).Error
	return &a, err
}
func (r *attachmentRepoImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.DocumentAttachment{}, "id = ?", id).Error
}
func (r *attachmentRepoImpl) ListByDocumentID(docID uuid.UUID) ([]*entity.DocumentAttachment, error) {
	var list []*entity.DocumentAttachment
	err := r.db.Where("document_id = ?", docID).Find(&list).Error
	return list, err
}

// commentRepoImpl
type commentRepoImpl struct{ db *gorm.DB }

func NewCommentRepository(db *gorm.DB) repo.DocumentCommentRepository {
	return &commentRepoImpl{db: db}
}
func (r *commentRepoImpl) Create(c *entity.DocumentComment) error { return r.db.Create(c).Error }
func (r *commentRepoImpl) GetByID(id uuid.UUID) (*entity.DocumentComment, error) {
	var c entity.DocumentComment
	err := r.db.First(&c, "id = ?", id).Error
	return &c, err
}
func (r *commentRepoImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.DocumentComment{}, "id = ?", id).Error
}
func (r *commentRepoImpl) ListByDocumentID(docID uuid.UUID) ([]*entity.DocumentComment, error) {
	var list []*entity.DocumentComment
	err := r.db.Where("document_id = ?", docID).Order("created_at ASC").Find(&list).Error
	return list, err
}

// knowledgeBaseRepoImpl
type kbRepoImpl struct{ db *gorm.DB }

func NewKnowledgeBaseRepository(db *gorm.DB) repo.KnowledgeBaseRepository {
	return &kbRepoImpl{db: db}
}
func (r *kbRepoImpl) Create(kb *entity.KnowledgeBase) error { return r.db.Create(kb).Error }
func (r *kbRepoImpl) GetByID(id uuid.UUID) (*entity.KnowledgeBase, error) {
	var kb entity.KnowledgeBase
	err := r.db.First(&kb, "id = ?", id).Error
	return &kb, err
}
func (r *kbRepoImpl) Update(kb *entity.KnowledgeBase) error { return r.db.Save(kb).Error }
func (r *kbRepoImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.KnowledgeBase{}, "id = ?", id).Error
}
func (r *kbRepoImpl) List(ownerID uuid.UUID) ([]*entity.KnowledgeBase, error) {
	var list []*entity.KnowledgeBase
	err := r.db.Where("owner_id = ? OR is_public = true", ownerID).Find(&list).Error
	return list, err
}

// favoriteRepoImpl
type favoriteRepoImpl struct{ db *gorm.DB }

func NewFavoriteRepository(db *gorm.DB) repo.DocumentFavoriteRepository {
	return &favoriteRepoImpl{db: db}
}
func (r *favoriteRepoImpl) Add(uid, did uuid.UUID) error {
	fav := &entity.DocumentFavorite{ID: uuid.New(), UserID: uid, DocumentID: did}
	return r.db.Create(fav).Error
}
func (r *favoriteRepoImpl) Remove(uid, did uuid.UUID) error {
	return r.db.Where("user_id = ? AND document_id = ?", uid, did).Delete(&entity.DocumentFavorite{}).Error
}
func (r *favoriteRepoImpl) IsFavorited(uid, did uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&entity.DocumentFavorite{}).
		Where("user_id = ? AND document_id = ?", uid, did).Count(&count).Error
	return count > 0, err
}
func (r *favoriteRepoImpl) ListByUserID(uid uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := r.db.Model(&entity.DocumentFavorite{}).
		Where("user_id = ?", uid).Pluck("document_id", &ids).Error
	return ids, err
}

// tagRepoImpl
type tagRepoImpl struct{ db *gorm.DB }

func NewTagRepository(db *gorm.DB) repo.DocumentTagRepository {
	return &tagRepoImpl{db: db}
}
func (r *tagRepoImpl) Create(t *entity.DocumentTag) error { return r.db.Create(t).Error }
func (r *tagRepoImpl) GetByID(id uuid.UUID) (*entity.DocumentTag, error) {
	var t entity.DocumentTag
	err := r.db.First(&t, "id = ?", id).Error
	return &t, err
}
func (r *tagRepoImpl) Update(t *entity.DocumentTag) error { return r.db.Save(t).Error }
func (r *tagRepoImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.DocumentTag{}, "id = ?", id).Error
}
func (r *tagRepoImpl) List() ([]*entity.DocumentTag, error) {
	var list []*entity.DocumentTag
	err := r.db.Order("count DESC").Find(&list).Error
	return list, err
}
func (r *tagRepoImpl) GetByName(name string) (*entity.DocumentTag, error) {
	var t entity.DocumentTag
	err := r.db.Where("name = ?", name).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}
func (r *tagRepoImpl) IncrementCount(id uuid.UUID) error {
	return r.db.Model(&entity.DocumentTag{}).Where("id = ?", id).UpdateColumn("count", gorm.Expr("count + 1")).Error
}

func buildCategoryTree(categories []*entity.DocumentCategory) []*entity.DocumentCategory {
	m := make(map[uuid.UUID]*entity.DocumentCategory)
	var roots []*entity.DocumentCategory
	for _, c := range categories {
		m[c.ID] = c
	}
	for _, c := range categories {
		if c.ParentID != nil {
			if parent, ok := m[*c.ParentID]; ok {
				_ = append([]*entity.DocumentCategory{}, parent) // 占位，实际需要Children字段
			}
		} else {
			roots = append(roots, c)
		}
	}
	return roots
}
