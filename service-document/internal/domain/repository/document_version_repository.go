package repository

import (
	"github.com/google/uuid"
	"leap-one/service-document/internal/domain/entity"
)

// DocumentVersionRepository 文档版本仓储接口
type DocumentVersionRepository interface {
	Create(version *entity.DocumentVersion) error
	GetByDocumentAndVersion(documentID uuid.UUID, versionNo int) (*entity.DocumentVersion, error)
	ListByDocumentID(documentID uuid.UUID) ([]*entity.DocumentVersion, error)
	GetLatest(documentID uuid.UUID) (*entity.DocumentVersion, error)
}
