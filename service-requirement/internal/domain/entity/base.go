package entity

import (
	"time"

	"github.com/google/uuid"
)

// BaseModel 基础审计模型
type BaseModel struct {
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
}
