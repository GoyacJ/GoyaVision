package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UserIdentityModel struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index:idx_user_identities_user_id"`
	IdentityType string         `gorm:"type:varchar(20);not null;index:idx_user_identities_type_identifier"`
	Identifier   string         `gorm:"type:varchar(255);not null;index:idx_user_identities_type_identifier"`
	Credential   string         `gorm:"type:varchar(255)"`
	Meta         datatypes.JSON `gorm:"serializer:json"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
}

func (UserIdentityModel) TableName() string { return "user_identities" }
