package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-kanban/internal/domain/entity"
	"leap-one/service-kanban/internal/domain/repository"
)

// BoardService 看板应用服务
type BoardService struct {
	boardRepo repository.BoardRepository
	logger    *zap.Logger
}

func NewBoardService(repo repository.BoardRepository, logger *zap.Logger) *BoardService {
	return &BoardService{boardRepo: repo, logger: logger}
}

func (s *BoardService) Create(b *entity.KanbanBoard) (*entity.KanbanBoard, error) {
	if b.Type == "" {
		b.Type = "project"
	}
	if err := s.boardRepo.Create(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *BoardService) GetByID(id uuid.UUID) (*entity.KanbanBoard, error) {
	return s.boardRepo.GetByID(id)
}

func (s *BoardService) Update(b *entity.KanbanBoard) error { return s.boardRepo.Update(b) }

func (s *BoardService) Delete(id uuid.UUID) error { return s.boardRepo.Delete(id) }

func (s *BoardService) List(ownerID uuid.UUID, boardType string) ([]*entity.KanbanBoard, error) {
	return s.boardRepo.List(ownerID, boardType)
}

// ColumnService 列管理服�?
type ColumnService struct {
	repo     repository.ColumnRepository
	cardRepo repository.CardRepository
	logger   *zap.Logger
}

func NewColumnService(repo repository.ColumnRepository, cardRepo repository.CardRepository, logger *zap.Logger) *ColumnService {
	return &ColumnService{repo: repo, cardRepo: cardRepo, logger: logger}
}

func (s *ColumnService) Create(c *entity.KanbanColumn) (*entity.KanbanColumn, error) {
	if c.Type == "" {
		c.Type = "normal"
	}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *ColumnService) Update(c *entity.KanbanColumn) error { return s.repo.Update(c) }

func (s *ColumnService) Delete(id uuid.UUID) error { return s.repo.Delete(id) }

func (s *ColumnService) List(boardID uuid.UUID) ([]*entity.KanbanColumn, error) {
	return s.repo.ListByBoardID(boardID)
}

func (s *ColumnService) Reorder(boardID uuid.UUID, columnIDs []uuid.UUID) error {
	return s.repo.UpdateSortOrder(boardID, columnIDs)
}

// CardService 卡片操作服务
type CardService struct {
	repo        repository.CardRepository
	historyRepo repository.MoveHistoryRepository
	columnRepo  repository.ColumnRepository
	logger      *zap.Logger
}

func NewCardService(
	repo repository.CardRepository,
	historyRepo repository.MoveHistoryRepository,
	columnRepo repository.ColumnRepository,
	logger *zap.Logger,
) *CardService {
	return &CardService{repo: repo, historyRepo: historyRepo, columnRepo: columnRepo, logger: logger}
}

func (s *CardService) Create(c *entity.KanbanCard) (*entity.KanbanCard, error) {
	if c.CardType == "" {
		c.CardType = "task"
	}
	c.MovedAt = time.Now()
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CardService) GetByID(id uuid.UUID) (*entity.KanbanCard, error) { return s.repo.GetByID(id) }

func (s *CardService) Update(c *entity.KanbanCard) error { return s.repo.Update(c) }

func (s *CardService) Delete(id uuid.UUID) error { return s.repo.Delete(id) }

func (s *CardService) ListByBoard(boardID uuid.UUID) ([]*entity.KanbanCard, error) {
	return s.repo.ListByBoardID(boardID)
}

func (s *CardService) Move(cardID, toColumnID, movedBy uuid.UUID) error {
	// 检查WIP限制
	col, err := s.columnRepo.GetByID(toColumnID)
	if err != nil {
		return err
	}
	if col.WIPLimit != nil && *col.WIPLimit > 0 {
		existingCards, _ := s.repo.ListByColumnID(toColumnID)
		if len(existingCards) >= *col.WIPLimit {
			return fmt.Errorf("列[%s]已达WIP限制(%d)，无法移动卡�?, col.Name, *col.WIPLimit)
		}
	}
	return s.repo.Move(cardID, toColumnID, movedBy)
}

func (s *CardService) GetMoveHistory(cardID uuid.UUID) ([]*entity.KanbanCardMoveHistory, error) {
	return s.historyRepo.ListByCardID(cardID)
}

// SwimlaneService 泳道服务
type SwimlaneService struct{ repo repository.SwimlaneRepository }

func NewSwimlaneService(repo repository.SwimlaneRepository) *SwimlaneService {
	return &SwimlaneService{repo: repo}
}

func (s *SwimlaneService) Create(sw *entity.KanbanSwimlane) error { return s.repo.Create(sw) }
func (s *SwimlaneService) Update(sw *entity.KanbanSwimlane) error { return s.repo.Update(sw) }
func (s *SwimlaneService) Delete(id uuid.UUID) error              { return s.repo.Delete(id) }
func (s *SwimlaneService) List(boardID uuid.UUID) ([]*entity.KanbanSwimlane, error) {
	return s.repo.ListByBoardID(boardID)
}

// StatisticsService 统计服务
type StatisticsService struct {
	cardRepo   repository.CardRepository
	columnRepo repository.ColumnRepository
}

func NewStatisticsService(cardRepo repository.CardRepository, columnRepo repository.ColumnRepository) *StatisticsService {
	return &StatisticsService{cardRepo: cardRepo, columnRepo: columnRepo}
}

func (s *StatisticsService) GetBoardStats(boardID uuid.UUID) (map[string]interface{}, error) {
	cards, _ := s.cardRepo.ListByBoardID(boardID)
	columns, _ := s.columnRepo.ListByBoardID(boardID)

	// 按列统计卡片�?
	cardCountByCol := make(map[uuid.UUID]int)
	for _, c := range cards {
		cardCountByCol[c.ColumnID]++
	}

	// 按优先级统计
	priorityCount := make(map[int]int)
	for _, c := range cards {
		priorityCount[c.Priority]++
	}

	// 按类型统�?
	typeCount := make(map[string]int)
	for _, c := range cards {
		typeCount[c.CardType]++
	}

	// 阻塞卡片数量
	blockedCount := 0
	for _, c := range cards {
		if c.BlockReason != "" {
			blockedCount++
		}
	}

	return map[string]interface{}{
		"total_cards":     len(cards),
		"total_columns":   len(columns),
		"cards_by_column": cardCountByCol,
		"by_priority":     priorityCount,
		"by_type":         typeCount,
		"blocked_count":   blockedCount,
	}, nil
}

func (s *StatisticsService) GetCFDData(boardID uuid.UUID) ([]map[string]interface{}, error) {
	// 简化版CFD数据：返回各列的累积流量数据
	cards, _ := s.cardRepo.ListByBoardID(boardID)
	columns, _ := s.columnRepo.ListByBoardID(boardID)

	result := []map[string]interface{}{}
	for _, col := range columns {
		count := 0
		for _, c := range cards {
			if c.ColumnID == col.ID {
				count++
			}
		}
		result = append(result, map[string]interface{}{
			"column_name": col.Name,
			"column_key":  col.Key,
			"count":       count,
		})
	}
	return result, nil
}
