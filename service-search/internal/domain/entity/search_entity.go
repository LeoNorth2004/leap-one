package entity

import ("time"; "github.com/google/uuid"; "gorm.io/gorm")

type SearchDocument struct{
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	DocType string `gorm:"size:30;not null" json:"doc_type"` // product/project/requirement/task/bug/document/user/issue
	RefID uuid.UUID `gorm:"index;not null" json:"ref_id"`
	Title string `gorm:"size:500;not null" json:"title"`
	Content string `gorm:"type:text" json:"content"`
	Summary string `gorm:"size:500" json:"summary"`
	Tags string `gorm:"type:text" json:"tags"`
	MetaData string `gorm:"type:text" json:"metadata"` // JSON扩展元数�?	Permission string `gorm:"type:text" json:"permission"`
	IndexedAt time.Time `gorm:"default:NOW()" json:"indexed_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
func (SearchDocument) TableName() string { return "search_documents" }
func (s *SearchDocument) BeforeCreate(tx *gorm.DB) error { if s.ID==uuid.Nil{s.ID=uuid.New()}; return nil }

type SavedSearch struct{
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID uuid.UUID `gorm:"index;not null" json:"user_id"`
	Name string `gorm:"size:200;not null" json:"name"`
	Scope string `gorm:"size:30;not null" json:"scope"` // 搜索范围
	Filters string `gorm:"type;text" json:"filters"` // JSON筛选条�?	Sort string `gorm:"size:100" json:"sort"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
func (SavedSearch) TableName() string { return "saved_searches" }
func (s *SavedSearch) BeforeCreate(tx *gorm.DB) error { if s.ID==uuid.Nil{s.ID=uuid.New()}; return nil }

type SearchHistory struct{
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID uuid.UUID `gorm:"index;not null" json:"user_id"`
	Query string `gorm:"size:500;not null" json:"query"`
	Scope string `gorm:"size:30;default:'all'" json:"scope"`
	ResultCount int `json:"result_count"`
	SearchedAt time.Time `gorm:"default:NOW()" json:"searched_at"`
}
func (SearchHistory) TableName() string { return "search_histories" }
func (s *SearchHistory) BeforeCreate(tx *gorm.DB) error { if s.ID==uuid.Nil{s.ID=uuid.New()}; return nil }
