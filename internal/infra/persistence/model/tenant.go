package model

import (
	"time"

	"github.com/google/uuid"
)

type TenantModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Code      string    `gorm:"type:varchar(64);uniqueIndex;not null"`
	Status    int       `gorm:"default:1"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (TenantModel) TableName() string { return "tenants" }
