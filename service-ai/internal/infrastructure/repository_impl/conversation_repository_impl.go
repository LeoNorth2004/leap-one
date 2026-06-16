package repository_impl

import (
	"context"

	"github.com/google/uuid"
	"leap-one/service-ai/internal/domain/entity"
	"leap-one/service-ai/internal/domain/repository"
	"gorm.io/gorm"
)

// ConversationRepositoryImpl 对话仓库实现
type ConversationRepositoryImpl struct{ db *gorm.DB }

func NewConversationRepository(db *gorm.DB) repository.ConversationRepository {
	return &ConversationRepositoryImpl{db: db}
}

func (r *ConversationRepositoryImpl) Create(ctx context.Context, conv *entity.AIConversation) error {
	return r.db.WithContext(ctx).Create(conv).Error
}
func (r *ConversationRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.AIConversation, error) {
	var c entity.AIConversation
	err := r.db.WithContext(ctx).First(&c, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
func (r *ConversationRepositoryImpl) ListByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*entity.AIConversation, int64, error) {
	var list []*entity.AIConversation
	var total int64
	query := r.db.WithContext(ctx).Model(&entity.AIConversation{}).Where("user_id = ?", userID)
	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
func (r *ConversationRepositoryImpl) Update(ctx context.Context, conv *entity.AIConversation) error {
	return r.db.WithContext(ctx).Save(conv).Error
}
func (r *ConversationRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.AIConversation{}, "id = ?", id).Error
}
func (r *ConversationRepositoryImpl) GetWithMessages(ctx context.Context, id uuid.UUID) (*entity.AIConversation, error) {
	var conv entity.AIConversation
	err := r.db.WithContext(ctx).Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).First(&conv, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

// MessageRepositoryImpl 消息仓库实现
type MessageRepositoryImpl struct{ db *gorm.DB }

func NewMessageRepository(db *gorm.DB) repository.MessageRepository {
	return &MessageRepositoryImpl{db: db}
}

func (r *MessageRepositoryImpl) Create(ctx context.Context, msg *entity.AIMessage) error {
	return r.db.WithContext(ctx).Create(msg).Error
}
func (r *MessageRepositoryImpl) BatchCreate(ctx context.Context, msgs []*entity.AIMessage) error {
	if len(msgs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&msgs).Error
}
func (r *MessageRepositoryImpl) ListByConversation(ctx context.Context, conversationID uuid.UUID, page, pageSize int) ([]*entity.AIMessage, int64, error) {
	var list []*entity.AIMessage
	var total int64
	query := r.db.WithContext(ctx).Model(&entity.AIMessage{}).Where("conversation_id = ?", conversationID)
	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
func (r *MessageRepositoryImpl) GetRecentByConversation(ctx context.Context, conversationID uuid.UUID, limit int) ([]*entity.AIMessage, error) {
	var list []*entity.AIMessage
	err := r.db.WithContext(ctx).Where("conversation_id = ?", conversationID).Order("created_at DESC").Limit(limit).Find(&list).Error
	return list, err
}
