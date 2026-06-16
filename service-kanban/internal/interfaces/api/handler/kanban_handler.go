package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"leap-one/service-kanban/internal/application/service"
	"leap-one/service-kanban/internal/domain/entity"
	"leap-one/service-kanban/internal/interfaces/api/dto"
)

// KanbanHandler 看板HTTP处理�?type KanbanHandler struct {
	boardSvc      *service.BoardService
	columnSvc     *service.ColumnService
	cardSvc       *service.CardService
	swimlaneSvc   *service.SwimlaneService
	statisticsSvc *service.StatisticsService
	logger        *zap.Logger
}

func NewKanbanHandler(
	boardSvc *service.BoardService,
	columnSvc *service.ColumnService,
	cardSvc *service.CardService,
	swimlaneSvc *service.SwimlaneService,
	statisticsSvc *service.StatisticsService,
	logger *zap.Logger,
) *KanbanHandler {
	return &KanbanHandler{boardSvc, columnSvc, cardSvc, swimlaneSvc, statisticsSvc, logger}
}

// ==================== 看板 CRUD ====================

func (h *KanbanHandler) CreateBoard(c *gin.Context) {
	var req dto.CreateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest(err.Error()))
		return
	}
	ownerID, _ := uuid.Parse(c.GetHeader("X-User-ID"))
	board := &entity.KanbanBoard{ID: uuid.New(), Name: req.Name, Type: req.Type, RefID: req.RefID, OwnerID: ownerID, Description: req.Description, IsDefault: req.IsDefault}
	result, err := h.boardSvc.Create(board)
	if err != nil {
		c.JSON(500, dto.InternalError("创建看板失败"))
		return
	}
	c.JSON(201, dto.Success(toBoardResp(result)))
}

func (h *KanbanHandler) ListBoards(c *gin.Context) {
	ownerID, _ := uuid.Parse(c.DefaultQuery("owner_id", "00000000-0000-0000-0000-000000000000"))
	boards, err := h.boardSvc.List(ownerID, c.Query("type"))
	if err != nil {
		c.JSON(500, dto.InternalError("获取看板列表失败"))
		return
	}
	var resps []dto.BoardResponse
	for _, b := range boards {
		resps = append(resps, toBoardResp(b))
	}
	c.JSON(200, dto.Success(resps))
}

func (h *KanbanHandler) GetBoard(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	board, err := h.boardSvc.GetByID(id)
	if err != nil {
		c.JSON(404, dto.NotFound("看板不存�?))
		return
	}
	c.JSON(200, dto.Success(toBoardFullResp(board)))
}

func (h *KanbanHandler) UpdateBoard(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	var req dto.CreateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest(err.Error()))
		return
	}
	board, _ := h.boardSvc.GetByID(id)
	board.Name = req.Name
	board.Description = req.Description
	board.Type = req.Type
	board.RefID = req.RefID
	board.IsDefault = req.IsDefault
	if err := h.boardSvc.Update(board); err != nil {
		c.JSON(500, dto.InternalError("更新看板失败"))
		return
	}
	c.JSON(200, dto.Success(toBoardResp(board)))
}

func (h *KanbanHandler) DeleteBoard(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效ID"))
		return
	}
	if err := h.boardSvc.Delete(id); err != nil {
		c.JSON(500, dto.InternalError("删除看板失败"))
		return
	}
	c.JSON(200, dto.Success(nil))
}

// ==================== 列操�?====================

func (h *KanbanHandler) CreateColumn(c *gin.Context) {
	boardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效看板ID"))
		return
	}
	var req dto.CreateColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest(err.Error()))
		return
	}
	col := &entity.KanbanColumn{ID: uuid.New(), BoardID: boardID, Name: req.Name, Key: req.Key, WIPLimit: req.WIPLimit, Color: req.Color, Type: req.Type}
	result, err := h.columnSvc.Create(col)
	if err != nil {
		c.JSON(500, dto.InternalError("创建列失�?))
		return
	}
	c.JSON(201, dto.Success(toColResp(result)))
}

func (h *KanbanHandler) UpdateColumn(c *gin.Context) {
	colID, err := uuid.Parse(c.Param("cid"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效列ID"))
		return
	}
	var req dto.UpdateColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest(err.Error()))
		return
	}
	col := &entity.KanbanColumn{ID: colID}
	if req.Name != nil {
		col.Name = *req.Name
	}
	if req.Key != nil {
		col.Key = *req.Key
	}
	if req.WIPLimit != nil {
		col.WIPLimit = req.WIPLimit
	}
	if req.Color != nil {
		col.Color = *req.Color
	}
	if req.Type != nil {
		col.Type = *req.Type
	}
	if req.SortOrder != nil {
		col.SortOrder = *req.SortOrder
	}
	if err := h.columnSvc.Update(col); err != nil {
		c.JSON(500, dto.InternalError("更新列失�?))
		return
	}
	c.JSON(200, dto.Success(nil))
}

func (h *KanbanHandler) DeleteColumn(c *gin.Context) {
	colID, err := uuid.Parse(c.Param("cid"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效列ID"))
		return
	}
	if err := h.columnSvc.Delete(colID); err != nil {
		c.JSON(500, dto.InternalError("删除列失�?))
		return
	}
	c.JSON(200, dto.Success(nil))
}

func (h *KanbanHandler) ReorderColumns(c *gin.Context) {
	boardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效看板ID"))
		return
	}
	var req struct {
		ColumnIDs []uuid.UUID `json:"column_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest(err.Error()))
		return
	}
	if err := h.columnSvc.Reorder(boardID, req.ColumnIDs); err != nil {
		c.JSON(500, dto.InternalError("排序列失�?))
		return
	}
	c.JSON(200, dto.Success(nil))
}

// ==================== 泳道操作 ====================

func (h *KanbanHandler) CreateSwimlane(c *gin.Context) {
	boardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效看板ID"))
		return
	}
	var req dto.CreateSwimlaneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest(err.Error()))
		return
	}
	sw := &entity.KanbanSwimlane{ID: uuid.New(), BoardID: boardID, Name: req.Name, Key: req.Key, Color: req.Color, SortOrder: req.SortOrder}
	if err := h.swimlaneSvc.Create(sw); err != nil {
		c.JSON(500, dto.InternalError("创建泳道失败"))
		return
	}
	c.JSON(201, dto.Success(toSwimResp(sw)))
}

func (h *KanbanHandler) UpdateSwimlane(c *gin.Context) {
	sid, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效泳道ID"))
		return
	}
	var req dto.UpdateSwimlaneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest(err.Error()))
		return
	}
	sw := &entity.KanbanSwimlane{ID: sid}
	if req.Name != nil {
		sw.Name = *req.Name
	}
	if req.Key != nil {
		sw.Key = *req.Key
	}
	if req.Color != nil {
		sw.Color = *req.Color
	}
	if req.SortOrder != nil {
		sw.SortOrder = *req.SortOrder
	}
	if err := h.swimlaneSvc.Update(sw); err != nil {
		c.JSON(500, dto.InternalError("更新泳道失败"))
		return
	}
	c.JSON(200, dto.Success(nil))
}

func (h *KanbanHandler) DeleteSwimlane(c *gin.Context) {
	sid, err := uuid.Parse(c.Param("sid"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效泳道ID"))
		return
	}
	if err := h.swimlaneSvc.Delete(sid); err != nil {
		c.JSON(500, dto.InternalError("删除泳道失败"))
		return
	}
	c.JSON(200, dto.Success(nil))
}

// ==================== 卡片操作 ====================

func (h *KanbanHandler) CreateCard(c *gin.Context) {
	boardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效看板ID"))
		return
	}
	var req dto.CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest(err.Error()))
		return
	}
	card := &entity.KanbanCard{
		ID: uuid.New(), BoardID: boardID, ColumnID: req.ColumnID,
		SwimlaneID: req.SwimlaneID, CardType: req.CardType,
		RefID: req.RefID, Title: req.Title, Priority: req.Priority,
		AssigneeID: req.AssigneeID, DueDate: req.DueDate,
		Tags: req.Tags, BlockReason: req.BlockReason,
	}
	result, err := h.cardSvc.Create(card)
	if err != nil {
		c.JSON(500, dto.InternalError("创建卡片失败"))
		return
	}
	c.JSON(201, dto.Success(toCardResp(result)))
}

func (h *KanbanHandler) UpdateCard(c *gin.Context) {
	cardID, err := uuid.Parse(c.Param("cardId"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效卡片ID"))
		return
	}
	var req dto.UpdateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest(err.Error()))
		return
	}
	card, err := h.cardSvc.GetByID(cardID)
	if err != nil {
		c.JSON(404, dto.NotFound("卡片不存�?))
		return
	}
	if req.Title != nil {
		card.Title = *req.Title
	}
	if req.Priority != nil {
		card.Priority = *req.Priority
	}
	if req.AssigneeID != nil {
		card.AssigneeID = req.AssigneeID
	}
	if req.DueDate != nil {
		card.DueDate = req.DueDate
	}
	if req.Tags != nil {
		card.Tags = *req.Tags
	}
	if req.BlockReason != nil {
		card.BlockReason = *req.BlockReason
	}
	if req.SortOrder != nil {
		card.SortOrder = *req.SortOrder
	}
	if err := h.cardSvc.Update(card); err != nil {
		c.JSON(500, dto.InternalError("更新卡片失败"))
		return
	}
	c.JSON(200, dto.Success(toCardResp(card)))
}

func (h *KanbanHandler) DeleteCard(c *gin.Context) {
	cardID, err := uuid.Parse(c.Param("cardId"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效卡片ID"))
		return
	}
	if err := h.cardSvc.Delete(cardID); err != nil {
		c.JSON(500, dto.InternalError("删除卡片失败"))
		return
	}
	c.JSON(200, dto.Success(nil))
}

func (h *KanbanHandler) MoveCard(c *gin.Context) {
	cardID, err := uuid.Parse(c.Param("cardId"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效卡片ID"))
		return
	}
	var req dto.MoveCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.BadRequest(err.Error()))
		return
	}
	movedBy, _ := uuid.Parse(c.GetHeader("X-User-ID"))
	if err := h.cardSvc.Move(cardID, req.ToColumnID, movedBy); err != nil {
		c.JSON(400, dto.BadRequest(err.Error()))
		return
	}
	c.JSON(200, dto.Success(gin.H{"moved": true}))
}

func (h *KanbanHandler) GetMoveHistory(c *gin.Context) {
	cardID, err := uuid.Parse(c.Param("cardId"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效卡片ID"))
		return
	}
	histories, err := h.cardSvc.GetMoveHistory(cardID)
	if err != nil {
		c.JSON(500, dto.InternalError("获取移动历史失败"))
		return
	}
	var resps []dto.MoveHistoryResponse
	for _, h := range histories {
		resps = append(resps, dto.MoveHistoryResponse{
			ID: h.ID, CardID: h.CardID, FromColID: h.FromColID,
			ToColID: h.ToColID, MovedBy: h.MovedBy, MoveTime: h.MoveTime,
		})
	}
	c.JSON(200, dto.Success(resps))
}

// ==================== 统计 ====================

func (h *KanbanHandler) GetStatistics(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效看板ID"))
		return
	}
	stats, err := h.statisticsSvc.GetBoardStats(id)
	if err != nil {
		c.JSON(500, dto.InternalError("获取统计数据失败"))
		return
	}
	c.JSON(200, dto.Success(stats))
}

func (h *KanbanHandler) GetCFD(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, dto.BadRequest("无效看板ID"))
		return
	}
	data, err := h.statisticsSvc.GetCFDData(id)
	if err != nil {
		c.JSON(500, dto.InternalError("获取CFD数据失败"))
		return
	}
	c.JSON(200, dto.Success(data))
}

// ==================== 转换函数 ====================

func toBoardResp(b *entity.KanbanBoard) dto.BoardResponse {
	return dto.BoardResponse{
		ID: b.ID, Name: b.Name, Type: b.Type, RefID: b.RefID,
		OwnerID: b.OwnerID, Description: b.Description, IsDefault: b.IsDefault,
		CreatedAt: b.CreatedAt,
	}
}

func toBoardFullResp(b *entity.KanbanBoard) dto.BoardResponse {
	resp := toBoardResp(b)
	for _, col := range b.Columns {
		resp.Columns = append(resp.Columns, toColResp(&col))
	}
	for _, sw := range b.Swimlanes {
		resp.Swimlanes = append(resp.Swimlanes, toSwimResp(&sw))
	}
	for _, card := range b.Cards {
		resp.Cards = append(resp.Cards, toCardResp(&card))
	}
	return resp
}

func toColResp(c *entity.KanbanColumn) dto.ColumnResponse {
	return dto.ColumnResponse{
		ID: c.ID, BoardID: c.BoardID, Name: c.Name, Key: c.Key,
		WIPLimit: c.WIPLimit, Color: c.Color, SortOrder: c.SortOrder, Type: c.Type,
	}
}

func toCardResp(cd *entity.KanbanCard) dto.CardResponse {
	return dto.CardResponse{
		ID: cd.ID, BoardID: cd.BoardID, ColumnID: cd.ColumnID,
		SwimlaneID: cd.SwimlaneID, CardType: cd.CardType, RefID: cd.RefID,
		Title: cd.Title, Priority: cd.Priority, AssigneeID: cd.AssigneeID,
		DueDate: cd.DueDate, Tags: cd.Tags, BlockReason: cd.BlockReason,
		SortOrder: cd.SortOrder, MovedAt: cd.MovedAt, CreatedAt: cd.CreatedAt,
	}
}

func toSwimResp(s *entity.KanbanSwimlane) dto.SwimlaneResponse {
	return dto.SwimlaneResponse{
		ID: s.ID, BoardID: s.BoardID, Name: s.Name, Key: s.Key,
		SortOrder: s.SortOrder, Color: s.Color,
	}
}
