package entity

import ("time"; "github.com/google/uuid"; "gorm.io/gorm")

type SystemConfig struct{
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Category string `gorm:"size:100;not null;index" json:"category"` // 分类
	Key string `gorm:"size:200;not null" json:"key"` // 配置�?	Value string `gorm:"type:text" json:"value"` // 配置�?	ValueType string `gorm:"size:20;default:'string'" json:"value_type"` // string/int/bool/json
	IsEncrypted bool `gorm:"default:false" json:"is_encrypted"`
	IsPublic bool `gorm:"default:false" json:"is_public"` // 是否公开可读
	Description string `gorm:"size:500" json:"description"`
	SortOrder int `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
func(SystemConfig)TableName()string{return "system_configs"}
func(s*SystemConfig)BeforeCreate(tx*gorm.DB)error{if s.ID==uuid.Nil{s.ID=uuid.New()};return nil}

type FeatureFlag struct{
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Key string `gorm:"size:200;uniqueIndex;not null" json:"key"`
	Name string `gorm:"size:200;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Enabled bool `gorm:"default:false" json:"enabled"`
	Rules string `gorm:"type:text" json:"rules"` // JSON启用规则
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
func(FeatureFlag)TableName()string{return "feature_flags"}
func(f*FeatureFlag)BeforeCreate(tx*gorm.DB)error{if f.ID==uuid.Nil{f.ID=uuid.New()};return nil}

type AuditLog struct{
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID uuid.UUID `gorm:"index" json:"user_id"`
	Action string `gorm:"size:100;not null" json:"action"` // 动作
	Resource string `gorm:"size:100;not null" json:"resource"` // 资源类型
	ResourceID uuid.UUID `gorm:"index" json:"resource_id"`
	Detail string `gorm:"type:text" json:"detail"` // 详情JSON
	IPAddress string `gorm:"size:50" json:"ip_address"`
	UserAgent string `gorm:"size:500" json:"user_agent"`
	CreatedAt time.Time `gorm:"default:NOW()" json:"created_at"`
}
func(AuditLog)TableName()string{return "audit_logs"}
func(a*AuditLog)BeforeCreate(tx*gorm.DB)error{if a.ID==uuid.Nil{a.ID=uuid.New()};return nil}
