package dto

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Response 统一响应
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

// Document DTO
type CreateDocRequest struct {
	Title      string     `json:"title" binding:"required,max=500"`
	Content    string     `json:"content"`
	Type       string     `json:"type"`
	CategoryID *uuid.UUID `json:"category_id"`
	ParentID   *uuid.UUID `json:"parent_id"`
	ProductID  *uuid.UUID `json:"product_id"`
	ProjectID  *uuid.UUID `json:"project_id"`
	OwnerID    uuid.UUID  `json:"owner_id" binding:"required"`
	Visibility string     `json:"visibility"`
	Tags       string     `json:"tags"`
	IsTemplate bool       `json:"is_template"`
}

type UpdateDocRequest struct {
	Title      *string    `json:"title" binding:"omitempty,max=500"`
	Content    *string    `json:"content"`
	Type       *string    `json:"type"`
	Status     *string    `json:"status"`
	Visibility *string    `json:"visibility"`
	CategoryID *uuid.UUID `json:"category_id"`
	Tags       *string    `json:"tags"`
}

type DocumentResponse struct {
	ID         uuid.UUID          `json:"id"`
	Title      string             `json:"title"`
	Content    string             `json:"content"`
	Type       string             `json:"type"`
	CategoryID *uuid.UUID         `json:"category_id"`
	ParentID   *uuid.UUID         `json:"parent_id"`
	ProductID  *uuid.UUID         `json:"product_id"`
	ProjectID  *uuid.UUID         `json:"project_id"`
	OwnerID    uuid.UUID          `json:"owner_id"`
	Status     string             `json:"status"`
	Visibility string             `json:"visibility"`
	Version    int                `json:"version"`
	Tags       string             `json:"tags"`
	IsTemplate bool               `json:"is_template"`
	OrderIndex int                `json:"order_index"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	Children   []DocumentResponse `json:"children,omitempty"`
}

// Version DTO
type VersionResponse struct {
	ID         uuid.UUID `json:"id"`
	DocumentID uuid.UUID `json:"document_id"`
	VersionNo  int       `json:"version_no"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	ChangeNote string    `json:"change_note"`
	CreatedBy  uuid.UUID `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
}

// Category DTO
type CategoryRequest struct {
	Name      string     `json:"name" binding:"required,max=200"`
	ParentID  *uuid.UUID `json:"parent_id"`
	SortOrder int        `json:"sort_order"`
}

type CategoryResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	ParentID  *uuid.UUID `json:"parent_id"`
	SortOrder int        `json:"sort_order"`
}

// Comment DTO
type CommentRequest struct {
	UserID   uuid.UUID  `json:"user_id" binding:"required"`
	Content  string     `json:"content" binding:"required"`
	Position string     `json:"position"`
	ParentID *uuid.UUID `json:"parent_id"`
}

type CommentResponse struct {
	ID         uuid.UUID `json:"id"`
	DocumentID uuid.UUID `json:"document_id"`
	UserID     uuid.UUID `json:"user_id"`
	Content    string    `json:"content"`
	Position   string    `json:"position"`
	CreatedAt  time.Time `json:"created_at"`
}

// KnowledgeBase DTO
type KBRequest struct {
	Name        string `json:"name" binding:"required,max=200"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type KBResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     uuid.UUID `json:"owner_id"`
	IsPublic    bool      `json:"is_public"`
}

// Attachment DTO
type AttachRequest struct {
	FileName string `json:"file_name" binding:"required,max=255"`
	FileSize int64  `json:"file_size"`
	FileType string `json:"file_type"`
	FileURL  string `json:"file_url" binding:"required"`
}

type AttachResponse struct {
	ID         uuid.UUID `json:"id"`
	DocumentID uuid.UUID `json:"document_id"`
	FileName   string    `json:"fileName"`
	FileSize   int64     `json:"file_size"`
	FileType   string    `json:"file_type"`
	FileURL    string    `json:"file_url"`
}

// Tag DTO
type TagRequest struct {
	Name  string `json:"name" binding:"required,max=50"`
	Color string `json:"color" binding:"omitempty,len=7"`
}

type TagResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Color string    `json:"color"`
	Count int       `json:"count"`
}
