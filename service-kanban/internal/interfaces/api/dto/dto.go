package dto

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ==================== 统一响应 ====================

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(data interface{}) Response {
	return Response{Code: http.StatusOK, Message: "success", Data: data}
}
func PageSuccess(list interface{}, total int64, page, size int) Response {
	return Response{Code: http.StatusOK, Message: "success",
		Data: gin.H{"list": list, "total": total, "page": page, "size": size}}
}
func Error(code int, msg string) Response { return Response{Code: code, Message: msg} }
func BadRequest(msg string) Response      { return Error(http.StatusBadRequest, msg) }
func NotFound(msg string) Response        { return Error(http.StatusNotFound, msg) }
func InternalError(msg string) Response   { return Error(http.StatusInternalServerError, msg) }

// ==================== 看板 DTO ====================

type CreateBoardRequest struct {
	Name        string     `json:"name" binding:"required,max=200"`
	Type        string     `json:"type"` // project/product/personal
	RefID       *uuid.UUID `json:"ref_id"`
	Description string     `json:"description"`
	IsDefault   bool       `json:"is_default"`
}

type BoardResponse struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	Type        string             `json:"type"`
	RefID       *uuid.UUID         `json:"ref_id"`
	OwnerID     uuid.UUID          `json:"owner_id"`
	Description string             `json:"description"`
	IsDefault   bool               `json:"is_default"`
	CreatedAt   time.Time          `json:"created_at"`
	Columns     []ColumnResponse   `json:"columns,omitempty"`
	Swimlanes   []SwimlaneResponse `json:"swimlanes,omitempty"`
	Cards       []CardResponse     `json:"cards,omitempty"`
}

// ==================== �?DTO ====================

type CreateColumnRequest struct {
	Name     string `json:"name" binding:"required,max=100"`
	Key      string `json:"key" binding:"required,max=50"`
	WIPLimit *int   `json:"wip_limit"`
	Color    string `json:"color" binding:"omitempty,len=7"`
	Type     string `json:"type"` // normal/backlog/done
}

type UpdateColumnRequest struct {
	Name      *string `json:"name" binding:"omitempty,max=100"`
	Key       *string `json:"key" binding:"omitempty,max=50"`
	WIPLimit  *int    `json:"wip_limit"`
	Color     *string `json:"color" binding:"omitempty,len=7"`
	Type      *string `json:"type"`
	SortOrder *int    `json:"sort_order"`
}

type ColumnResponse struct {
	ID        uuid.UUID `json:"id"`
	BoardID   uuid.UUID `json:"board_id"`
	Name      string    `json:"name"`
	Key       string    `json:"key"`
	WIPLimit  *int      `json:"wip_limit"`
	Color     string    `json:"color"`
	SortOrder int       `json:"sort_order"`
	Type      string    `json:"type"`
}

// ==================== 卡片 DTO ====================

type CreateCardRequest struct {
	ColumnID    uuid.UUID  `json:"column_id" binding:"required"`
	SwimlaneID  *uuid.UUID `json:"swimlane_id"`
	CardType    string     `json:"card_type"` // task/requirement/bug
	RefID       uuid.UUID  `json:"ref_id" binding:"required"`
	Title       string     `json:"title" binding:"required,max=500"`
	Priority    int        `json:"priority"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
	Tags        string     `json:"tags"`
	BlockReason string     `json:"block_reason"`
}

type UpdateCardRequest struct {
	Title       *string    `json:"title" binding:"omitempty,max=500"`
	Priority    *int       `json:"priority"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
	Tags        *string    `json:"tags"`
	BlockReason *string    `json:"block_reason"`
	SortOrder   *int       `json:"sort_order"`
}

type MoveCardRequest struct {
	ToColumnID uuid.UUID `json:"to_column_id" binding:"required"`
}

type CardResponse struct {
	ID          uuid.UUID  `json:"id"`
	BoardID     uuid.UUID  `json:"board_id"`
	ColumnID    uuid.UUID  `json:"column_id"`
	SwimlaneID  *uuid.UUID `json:"swimlane_id"`
	CardType    string     `json:"card_type"`
	RefID       uuid.UUID  `json:"ref_id"`
	Title       string     `json:"title"`
	Priority    int        `json:"priority"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
	Tags        string     `json:"tags"`
	BlockReason string     `json:"block_reason"`
	SortOrder   int        `json:"sort_order"`
	MovedAt     time.Time  `json:"moved_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ==================== 泳道 DTO ====================

type CreateSwimlaneRequest struct {
	Name      string `json:"name" binding:"required,max=100"`
	Key       string `json:"key" binding:"required,max=50"`
	Color     string `json:"color" binding:"omitempty,len=7"`
	SortOrder int    `json:"sort_order"`
}

type UpdateSwimlaneRequest struct {
	Name      *string `json:"name" binding:"omitempty,max=100"`
	Key       *string `json:"key" binding:"omitempty,max=50"`
	Color     *string `json:"color" binding:"omitempty,len=7"`
	SortOrder *int    `json:"sort_order"`
}

type SwimlaneResponse struct {
	ID        uuid.UUID `json:"id"`
	BoardID   uuid.UUID `json:"board_id"`
	Name      string    `json:"name"`
	Key       string    `json:"key"`
	SortOrder int       `json:"sort_order"`
	Color     string    `json:"color"`
}

// ==================== 移动历史 DTO ====================

type MoveHistoryResponse struct {
	ID        uuid.UUID `json:"id"`
	CardID    uuid.UUID `json:"card_id"`
	FromColID uuid.UUID `json:"from_column_id"`
	ToColID   uuid.UUID `json:"to_column_id"`
	MovedBy   uuid.UUID `json:"moved_by"`
	MoveTime  time.Time `json:"move_time"`
}

// ==================== 统计 DTO ====================

type StatisticsResponse struct {
	TotalCards    int            `json:"total_cards"`
	TotalColumns  int            `json:"total_columns"`
	CardsByColumn map[string]int `json:"cards_by_column"`
	ByPriority    map[int]int    `json:"by_priority"`
	ByType        map[string]int `json:"by_type"`
	BlockedCount  int            `json:"blocked_count"`
}

type CFDItemResponse struct {
	ColumnName string `json:"column_name"`
	ColumnKey  string `json:"column_key"`
	Count      int    `json:"count"`
}
