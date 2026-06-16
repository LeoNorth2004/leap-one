package service

import (
	"github.com/google/uuid"

	"leap-one/service-document/internal/domain/entity"
	"leap-one/service-document/internal/domain/repository"
)

// VersionService 版本管理服务
type VersionService struct {
	repo repository.DocumentVersionRepository
}

func NewVersionService(repo repository.DocumentVersionRepository) *VersionService {
	return &VersionService{repo: repo}
}

func (s *VersionService) ListVersions(documentID uuid.UUID) ([]*entity.DocumentVersion, error) {
	return s.repo.ListByDocumentID(documentID)
}

func (s *VersionService) GetVersion(documentID uuid.UUID, ver int) (*entity.DocumentVersion, error) {
	return s.repo.GetByDocumentAndVersion(documentID, ver)
}

func (s *VersionService) RestoreToVersion(documentID uuid.UUID, targetVer int, currentDoc *entity.Document, docSvc *DocumentService) error {
	version, err := s.repo.GetByDocumentAndVersion(documentID, targetVer)
	if err != nil {
		return err
	}
	currentDoc.Content = version.Content
	currentDoc.Title = version.Title
	currentDoc.Version++
	_, err = docSvc.Update(currentDoc.ID, map[string]interface{}{"content": currentDoc.Content, "title": currentDoc.Title})
	return err
}

// CommentService 评论服务
type CommentService struct {
	repo repository.DocumentCommentRepository
}

func NewCommentService(repo repository.DocumentCommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) Add(comment *entity.DocumentComment) error { return s.repo.Create(comment) }
func (s *CommentService) List(documentID uuid.UUID) ([]*entity.DocumentComment, error) {
	return s.repo.ListByDocumentID(documentID)
}
func (s *CommentService) Delete(id uuid.UUID) error { return s.repo.Delete(id) }

// CategoryService 分类服务
type CategoryService struct {
	repo repository.DocumentCategoryRepository
}

func NewCategoryService(repo repository.DocumentCategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(c *entity.DocumentCategory) error   { return s.repo.Create(c) }
func (s *CategoryService) List() ([]*entity.DocumentCategory, error) { return s.repo.List() }
func (s *CategoryService) Update(c *entity.DocumentCategory) error   { return s.repo.Update(c) }
func (s *CategoryService) Delete(id uuid.UUID) error                 { return s.repo.Delete(id) }

// KnowledgeBaseService 知识库服�?
type KnowledgeBaseService struct {
	repo repository.KnowledgeBaseRepository
}

func NewKnowledgeBaseService(repo repository.KnowledgeBaseRepository) *KnowledgeBaseService {
	return &KnowledgeBaseService{repo: repo}
}

func (s *KnowledgeBaseService) Create(kb *entity.KnowledgeBase) error { return s.repo.Create(kb) }
func (s *KnowledgeBaseService) GetByID(id uuid.UUID) (*entity.KnowledgeBase, error) {
	return s.repo.GetByID(id)
}
func (s *KnowledgeBaseService) Update(kb *entity.KnowledgeBase) error { return s.repo.Update(kb) }
func (s *KnowledgeBaseService) Delete(id uuid.UUID) error             { return s.repo.Delete(id) }
func (s *KnowledgeBaseService) List(ownerID uuid.UUID) ([]*entity.KnowledgeBase, error) {
	return s.repo.List(ownerID)
}

// AttachmentService 附件服务
type AttachmentService struct {
	repo repository.DocumentAttachmentRepository
}

func NewAttachmentService(repo repository.DocumentAttachmentRepository) *AttachmentService {
	return &AttachmentService{repo: repo}
}

func (s *AttachmentService) Upload(a *entity.DocumentAttachment) error { return s.repo.Create(a) }
func (s *AttachmentService) List(documentID uuid.UUID) ([]*entity.DocumentAttachment, error) {
	return s.repo.ListByDocumentID(documentID)
}
func (s *AttachmentService) Delete(id uuid.UUID) error { return s.repo.Delete(id) }

// TagService 标签服务
type TagService struct {
	repo repository.DocumentTagRepository
}

func NewTagService(repo repository.DocumentTagRepository) *TagService { return &TagService{repo: repo} }

func (s *TagService) Create(t *entity.DocumentTag) error   { return s.repo.Create(t) }
func (s *TagService) List() ([]*entity.DocumentTag, error) { return s.repo.List() }
func (s *TagService) Delete(id uuid.UUID) error            { return s.repo.Delete(id) }
