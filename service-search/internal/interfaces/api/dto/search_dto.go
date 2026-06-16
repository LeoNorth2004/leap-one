package dto

import "github.com/google/uuid"

// SearchRequest 全局搜索请求
type SearchRequest struct {
	Query    string `form:"q" binding:"required,max=500"`
	Type     string `form:"type" binding:"omitempty,max=30"` // doc_type过滤
	Page     int    `form:"page"`
	PageSize int    `form:"size"`
}

// AdvancedSearchRequest 高级搜索请求
type AdvancedSearchRequest struct {
	Query    string                 `json:"query" binding:"max=500"`
	Filters  map[string]interface{} `json:"filters"` // 高级筛选条�?
	Sort     string                 `json:"sort,omitempty"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
}

// SaveSearchRequest 保存搜索条件请求
type SaveSearchRequest struct {
	UserID  uuid.UUID `json:"user_id" binding:"required"`
	Name    string    `json:"name" binding:"required,max=200"`
	Scope   string    `json:"scope" binding:"required,max=30"`
	Filters string    `json:"filters"` // JSON筛选条�?
	Sort    string    `json:"sort"`
}

// SearchResult 搜索结果�?
type SearchResult struct {
	ID        string `json:"id"`
	DocType   string `json:"doc_type"`
	RefID     string `json:"ref_id"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Tags      string `json:"tags"`
	MetaData  string `json:"metadata"`
	IndexedAt string `json:"indexed_at"`
	Highlight string `json:"highlight,omitempty"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	List       []SearchResult `json:"list"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Size       int            `json:"size"`
	Query      string         `json:"query"`
	DurationMs int64          `json:"duration_ms"`
}

// SavedSearchInfo 保存的搜索信�?
type SavedSearchInfo struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Scope     string `json:"scope"`
	Filters   string `json:"filters"`
	Sort      string `json:"sort"`
	CreatedAt string `json:"created_at"`
}

// SearchHistoryItem 搜索历史�?
type SearchHistoryItem struct {
	ID          string `json:"id"`
	Query       string `json:"query"`
	Scope       string `json:"scope"`
	ResultCount int    `json:"result_count"`
	SearchedAt  string `json:"searched_at"`
}

// SuggestionResponse 搜索建议响应
type SuggestionResponse struct {
	Suggestions []string `json:"suggestions"`
	Query       string   `json:"query"`
}

// IndexStatusResponse 索引状态响�?
type IndexStatusResponse struct {
	TotalDocuments int64            `json:"total_documents"`
	LastIndexedAt  string           `json:"last_indexed_at"`
	DocTypes       map[string]int64 `json:"doc_types"`
}
