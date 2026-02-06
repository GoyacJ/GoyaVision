package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OperatorVersionModel struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey"`
	OperatorID uuid.UUID      `gorm:"type:uuid;not null;index:idx_operator_versions_operator_id;uniqueIndex:uk_operator_versions_operator_version"`
	Version    string         `gorm:"type:varchar(50);not null;uniqueIndex:uk_operator_versions_operator_version"`
	ExecMode   string         `gorm:"type:varchar(20);not null;index:idx_operator_versions_exec_mode"`
	ExecConfig datatypes.JSON `gorm:"type:jsonb"`
	InputSchema datatypes.JSON `gorm:"type:jsonb"`
	OutputSpec datatypes.JSON `gorm:"type:jsonb"`
	Config     datatypes.JSON `gorm:"type:jsonb"`
	Changelog  string         `gorm:"type:text"`
	Status     string         `gorm:"type:varchar(20);not null;default:'draft';index:idx_operator_versions_status"`
	CreatedAt  time.Time      `gorm:"autoCreateTime;index:idx_operator_versions_created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
}

func (OperatorVersionModel) TableName() string { return "operator_versions" }
