package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"leap-one/service-kanban/internal/domain/entity"
)

// ==================== Board Repository ====================

type boardRepoImpl struct{ db *gorm.DB }

func NewBoardRepository(db *gorm.DB) BoardRepository { return &boardRepoImpl{db: db} }

func (r *boardRepoImpl) Create(b *entity.KanbanBoard) error { return r.db.Create(b).Error }
func (r *boardRepoImpl) GetByID(id uuid.UUID) (*entity.KanbanBoard, error) {
	var b entity.KanbanBoard
	err := r.db.Preload("Columns").Preload("Cards").Preload("Swimlanes").
		Where("id = ? AND deleted_at IS NULL", id).First(&b).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}
func (r *boardRepoImpl) Update(b *entity.KanbanBoard) error { return r.db.Save(b).Error }
func (r *boardRepoImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.KanbanBoard{}, "id = ?", id).Error
}

func (r *boardRepoImpl) List(ownerID uuid.UUID, boardType string) ([]*entity.KanbanBoard, error) {
	var list []*entity.KanbanBoard
	query := r.db.Where("(owner_id = ? OR is_default = true) AND deleted_at IS NULL", ownerID)
	if boardType != "" {
		query = query.Where("type = ?", boardType)
	}
	err := query.Order("is_default DESC, created_at DESC").Find(&list).Error
	return list, err
}

func (r *boardRepoImpl) GetByRefID(refID uuid.UUID) (*entity.KanbanBoard, error) {
	var b entity.KanbanBoard
	err := r.db.Where("ref_id = ? AND deleted_at IS NULL", refID).First(&b).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// ==================== Column Repository ====================

type columnRepoImpl struct{ db *gorm.DB }

func NewColumnRepository(db *gorm.DB) ColumnRepository { return &columnRepoImpl{db: db} }

func (r *columnRepoImpl) Create(c *entity.KanbanColumn) error { return r.db.Create(c).Error }
func (r *columnRepoImpl) GetByID(id uuid.UUID) (*entity.KanbanColumn, error) {
	var c entity.KanbanColumn
	err := r.db.First(&c, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}
func (r *columnRepoImpl) Update(c *entity.KanbanColumn) error { return r.db.Save(c).Error }
func (r *columnRepoImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.KanbanColumn{}, "id = ?", id).Error
}

func (r *columnRepoImpl) ListByBoardID(boardID uuid.UUID) ([]*entity.KanbanColumn, error) {
	var list []*entity.KanbanColumn
	err := r.db.Where("board_id = ? AND deleted_at IS NULL", boardID).Order("sort_order ASC").Find(&list).Error
	return list, err
}

func (r *columnRepoImpl) UpdateSortOrder(boardID uuid.UUID, columnIDs []uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range columnIDs {
			if err := tx.Model(&entity.KanbanColumn{}).Where("id = ?", id).Update("sort_order", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ==================== Card Repository ====================

type cardRepoImpl struct{ db *gorm.DB }

func NewCardRepository(db *gorm.DB) CardRepository { return &cardRepoImpl{db: db} }

func (r *cardRepoImpl) Create(c *entity.KanbanCard) error { return r.db.Create(c).Error }
func (r *cardRepoImpl) GetByID(id uuid.UUID) (*entity.KanbanCard, error) {
	var c entity.KanbanCard
	err := r.db.First(&c, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}
func (r *cardRepoImpl) Update(c *entity.KanbanCard) error { return r.db.Save(c).Error }
func (r *cardRepoImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.KanbanCard{}, "id = ?", id).Error
}

func (r *cardRepoImpl) ListByBoardID(boardID uuid.UUID) ([]*entity.KanbanCard, error) {
	var list []*entity.KanbanCard
	err := r.db.Where("board_id = ? AND deleted_at IS NULL", boardID).
		Order("column_id ASC, sort_order ASC").Find(&list).Error
	return list, err
}

func (r *cardRepoImpl) ListByColumnID(columnID uuid.UUID) ([]*entity.KanbanCard, error) {
	var list []*entity.KanbanCard
	err := r.db.Where("column_id = ? AND deleted_at IS NULL", columnID).
		Order("sort_order ASC").Find(&list).Error
	return list, err
}

func (r *cardRepoImpl) Move(cardID, toColumnID, movedBy uuid.UUID) error {
	now := time.Now()
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 获取当前卡片信息
		var card entity.KanbanCard
		if err := tx.First(&card, "id = ?", cardID).Error; err != nil {
			return err
		}

		// 记录移动历史
		history := &entity.KanbanCardMoveHistory{
			ID: uuid.New(), CardID: cardID, FromColID: card.ColumnID,
			ToColID: toColumnID, MovedBy: movedBy, MoveTime: now,
		}
		if err := tx.Table("kanban_card_move_histories").Create(history).Error; err != nil {
			return err
		}

		// 更新卡片位置
		return tx.Model(&entity.KanbanCard{}).Where("id = ?", cardID).Updates(map[string]interface{}{
			"column_id": toColumnID, "moved_at": now, "moved_by": movedBy,
		}).Error
	})
}

// ==================== Swimlane Repository ====================

type swimlaneRepoImpl struct{ db *gorm.DB }

func NewSwimlaneRepository(db *gorm.DB) SwimlaneRepository { return &swimlaneRepoImpl{db: db} }

func (r *swimlaneRepoImpl) Create(s *entity.KanbanSwimlane) error { return r.db.Create(s).Error }
func (r *swimlaneRepoImpl) GetByID(id uuid.UUID) (*entity.KanbanSwimlane, error) {
	var s entity.KanbanSwimlane
	err := r.db.First(&s, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}
func (r *swimlaneRepoImpl) Update(s *entity.KanbanSwimlane) error { return r.db.Save(s).Error }
func (r *swimlaneRepoImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.KanbanSwimlane{}, "id = ?", id).Error
}
func (r *swimlaneRepoImpl) ListByBoardID(boardID uuid.UUID) ([]*entity.KanbanSwimlane, error) {
	var list []*entity.KanbanSwimlane
	err := r.db.Where("board_id = ? AND deleted_at IS NULL", boardID).Order("sort_order ASC").Find(&list).Error
	return list, err
}

// ==================== Move History Repository ====================

type moveHistoryRepoImpl struct{ db *gorm.DB }

func NewMoveHistoryRepository(db *gorm.DB) MoveHistoryRepository { return &moveHistoryRepoImpl{db: db} }

func (r *moveHistoryRepoImpl) Create(h *entity.KanbanCardMoveHistory) error {
	return r.db.Create(h).Error
}

func (r *moveHistoryRepoImpl) ListByCardID(cardID uuid.UUID) ([]*entity.KanbanCardMoveHistory, error) {
	var list []*entity.KanbanCardMoveHistory
	err := r.db.Where("card_id = ?", cardID).Order("move_time DESC").Find(&list).Error
	return list, err
}
