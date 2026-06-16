package repository

import (
	"github.com/google/uuid"
	"leap-one/service-kanban/internal/domain/entity"
)

// BoardRepository 看板仓储接口
type BoardRepository interface {
	Create(board *entity.KanbanBoard) error
	GetByID(id uuid.UUID) (*entity.KanbanBoard, error)
	Update(board *entity.KanbanBoard) error
	Delete(id uuid.UUID) error
	List(ownerID uuid.UUID, boardType string) ([]*entity.KanbanBoard, error)
	GetByRefID(refID uuid.UUID) (*entity.KanbanBoard, error)
}

// ColumnRepository 看板列仓储接�?type ColumnRepository interface {
	Create(column *entity.KanbanColumn) error
	GetByID(id uuid.UUID) (*entity.KanbanColumn, error)
	Update(column *entity.KanbanColumn) error
	Delete(id uuid.UUID) error
	ListByBoardID(boardID uuid.UUID) ([]*entity.KanbanColumn, error)
	UpdateSortOrder(boardID uuid.UUID, columnIDs []uuid.UUID) error
}

// CardRepository 卡片仓储接口
type CardRepository interface {
	Create(card *entity.KanbanCard) error
	GetByID(id uuid.UUID) (*entity.KanbanCard) error
	Update(card *entity.KanbanCard) error
	Delete(id uuid.UUID) error
	ListByBoardID(boardID uuid.UUID) ([]*entity.KanbanCard, error)
	ListByColumnID(columnID uuid.UUID) ([]*entity.KanbanCard, error)
	Move(cardID, toColumnID, movedBy uuid.UUID) error
}

// SwimlaneRepository 泳道仓储接口
type SwimlaneRepository interface {
	Create(swimlane *entity.KanbanSwimlane) error
	GetByID(id uuid.UUID) (*entity.KanbanSwimlane, error)
	Update(swimlane *entity.KanbanSwimlane) error
	Delete(id uuid.UUID) error
	ListByBoardID(boardID uuid.UUID) ([]*entity.KanbanSwimlane, error)
}

// MoveHistoryRepository 移动历史仓储接口
type MoveHistoryRepository interface {
	Create(history *entity.KanbanCardMoveHistory) error
	ListByCardID(cardID uuid.UUID) ([]*entity.KanbanCardMoveHistory, error)
}
