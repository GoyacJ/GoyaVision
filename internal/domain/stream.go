package domain

import (
	"time"

	"github.com/google/uuid"
)

type Stream struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	URL       string    `gorm:"not null"`
	Name      string    `gorm:"not null"`
	Enabled   bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Stream) TableName() string { return "streams" }
