package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OperatorModel struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Code        string         `gorm:"type:varchar(100);not null;uniqueIndex"`
	Name        string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Category    string         `gorm:"type:varchar(50);not null;index:idx_operators_category"`
	Type        string         `gorm:"type:varchar(50);not null;index:idx_operators_type"`
	Version     string         `gorm:"type:varchar(50);not null;default:'1.0.0'"`
	Endpoint    string         `gorm:"type:varchar(1024);not null"`
	Method      string         `gorm:"type:varchar(10);not null;default:'POST'"`
	InputSchema datatypes.JSON `gorm:"type:jsonb"`
	OutputSpec  datatypes.JSON `gorm:"type:jsonb"`
	Config      datatypes.JSON `gorm:"type:jsonb"`
	Status      string         `gorm:"type:varchar(20);not null;default:'draft';index:idx_operators_status"`
	IsBuiltin   bool           `gorm:"not null;default:false;index:idx_operators_builtin"`
	Tags        datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt   time.Time      `gorm:"autoCreateTime;index:idx_operators_created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
}

func (OperatorModel) TableName() string { return "operators" }
