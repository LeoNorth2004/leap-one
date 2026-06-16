package repository

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-task/internal/domain/entity"
)

// IssueTemplateRepository 工单模板仓库接口定义
type IssueTemplateRepository interface {
	// Create 创建模板
	Create(ctx context.Context, template *entity.IssueTemplate) error

	// GetByID 根据ID获取模板
	GetByID(ctx context.Context, id uuid.UUID) (*entity.IssueTemplate, error)

	// Update 更新模板
	Update(ctx context.Context, template *entity.IssueTemplate) error

	// Delete 软删除模板
	Delete(ctx context.Context, id uuid.UUID) error

	// List 分页查询模板列表
	List(ctx context.Context, page, pageSize int, tmplType string) ([]*entity.IssueTemplate, int64, error)

	// ListByType 按类型查询模板
	ListByType(ctx context.Context, tmplType string) ([]*entity.IssueTemplate, error)
}
