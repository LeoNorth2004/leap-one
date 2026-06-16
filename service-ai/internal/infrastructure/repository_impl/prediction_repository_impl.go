package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-ai/internal/domain/entity"
	"leap-one/service-ai/internal/domain/repository"
	"gorm.io/gorm"
)

// PredictionRepositoryImpl 预测记录仓库实现
type PredictionRepositoryImpl struct{ db *gorm.DB }

func NewPredictionRepository(db *gorm.DB) repository.PredictionRepository { return &PredictionRepositoryImpl{db: db} }

func (r *PredictionRepositoryImpl) Create(ctx context.Context, pred *entity.AIPrediction) error {
	return r.db.WithContext(ctx).Create(pred).Error
}
func (r *PredictionRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.AIPrediction, error) {
	var p entity.AIPrediction
	err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound { return nil, nil }
	if err != nil { return nil, err }
	return &p, nil
}
func (r *PredictionRepositoryImpl) List(ctx context.Context, page, pageSize int, predType string, targetID uuid.UUID) ([]*entity.AIPrediction, int64, error) {
	var list []*entity.AIPrediction; var total int64
	query := r.db.WithContext(ctx).Model(&entity.AIPrediction{})
	if predType != "" { query = query.Where("type = ?", predType) }
	if targetID != uuid.Nil { query = query.Where("target_id = ?", targetID) }
	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil { return nil, 0, err }
	return list, total, nil
}
func (r *PredictionRepositoryImpl) ListByTarget(ctx context.Context, targetID uuid.UUID) ([]*entity.AIPrediction, error) {
	var list []*entity.AIPrediction
	err := r.db.WithContext(ctx).Where("target_id = ?", targetID).Order("created_at DESC").Find(&list).Error
	return list, err
}
func (r *PredictionRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.AIPrediction{}, "id = ?", id).Error
}

// AIConfigRepositoryImpl AI配置仓库实现
type AIConfigRepositoryImpl struct{ db *gorm.DB }

func NewAIConfigRepository(db *gorm.DB) repository.AIConfigRepository { return &AIConfigRepositoryImpl{db: db} }

func (r *AIConfigRepositoryImpl) GetActive(ctx context.Context) (*entity.AIConfig, error) {
	var cfg entity.AIConfig
	err := r.db.WithContext(ctx).Where("is_active = true").First(&cfg).Error
	if err == gorm.ErrRecordNotFound { return nil, nil }
	if err != nil { return nil, err }
	return &cfg, nil
}
func (r *AIConfigRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.AIConfig, error) {
	var cfg entity.AIConfig
	err := r.db.WithContext(ctx).First(&cfg, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound { return nil, nil }
	if err != nil { return nil, err }
	return &cfg, nil
}
func (r *AIConfigRepositoryImpl) Update(ctx context.Context, cfg *entity.AIConfig) error {
	return r.db.WithContext(ctx).Save(cfg).Error
}
func (r *AIConfigRepositoryImpl) Create(ctx context.Context, cfg *entity.AIConfig) error {
	return r.db.WithContext(ctx).Create(cfg).Error
}
func (r *AIConfigRepositoryImpl) ListAll(ctx context.Context) ([]*entity.AIConfig, error) {
	var list []*entity.AIConfig
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&list).Error
	return list, err
}
