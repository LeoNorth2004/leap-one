package repository

import (
	"github.com/google/uuid"
	"leap-one/service-document/internal/domain/entity"
)

// DocumentCategoryRepository 文档分类仓储接口
type DocumentCategoryRepository interface {
	Create(category *entity.DocumentCategory) error
	GetByID(id uuid.UUID) (*entity.DocumentCategory, error)
	Update(category *entity.DocumentCategory) error
	Delete(id uuid.UUID) error
	List() ([]*entity.DocumentCategory, error)
	GetTree() ([]*entity.DocumentCategory, error)
}

// DocumentAttachmentRepository 文档附件仓储接口
type DocumentAttachmentRepository interface {
	Create(attachment *entity.DocumentAttachment) error
	GetByID(id uuid.UUID) (*entity.DocumentAttachment, error)
	Delete(id uuid.UUID) error
	ListByDocumentID(documentID uuid.UUID) ([]*entity.DocumentAttachment, error)
}

// DocumentCommentRepository 文档评论仓储接口
type DocumentCommentRepository interface {
	Create(comment *entity.DocumentComment) error
	GetByID(id uuid.UUID) (*entity.DocumentComment, error)
	Delete(id uuid.UUID) error
	ListByDocumentID(documentID uuid.UUID) ([]*entity.DocumentComment, error)
}

// KnowledgeBaseRepository 知识库仓储接�?type KnowledgeBaseRepository interface {
	Create(kb *entity.KnowledgeBase) error
	GetByID(id uuid.UUID) (*entity.KnowledgeBase, error)
	Update(kb *entity.KnowledgeBase) error
	Delete(id uuid.UUID) error
	List(ownerID uuid.UUID) ([]*entity.KnowledgeBase, error)
}

// DocumentFavoriteRepository 文档收藏仓储接口
type DocumentFavoriteRepository interface {
	Add(userID, documentID uuid.UUID) error
	Remove(userID, documentID uuid.UUID) error
	IsFavorited(userID, documentID uuid.UUID) (bool, error)
	ListByUserID(userID uuid.UUID) ([]uuid.UUID, error)
}

// DocumentTagRepository 文档标签仓储接口
type DocumentTagRepository interface {
	Create(tag *entity.DocumentTag) error
	GetByID(id uuid.UUID) (*entity.DocumentTag, error)
	Update(tag *entity.DocumentTag) error
	Delete(id uuid.UUID) error
	List() ([]*entity.DocumentTag, error)
	GetByName(name string) (*entity.DocumentTag, error)
	IncrementCount(id uuid.UUID) error
}
