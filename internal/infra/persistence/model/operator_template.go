package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OperatorTemplateModel struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Code        string         `gorm:"type:varchar(100);not null;uniqueIndex"`
	Name        string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Category    string         `gorm:"type:varchar(50);not null;index:idx_operator_templates_category"`
	Type        string         `gorm:"type:varchar(50);not null;index:idx_operator_templates_type"`
	ExecMode    string         `gorm:"type:varchar(20);not null;index:idx_operator_templates_exec_mode"`
	ExecConfig  datatypes.JSON `gorm:"type:jsonb"`
	InputSchema datatypes.JSON `gorm:"type:jsonb"`
	OutputSpec  datatypes.JSON `gorm:"type:jsonb"`
	Config      datatypes.JSON `gorm:"type:jsonb"`
	Author      string         `gorm:"type:varchar(255)"`
	Tags        datatypes.JSON `gorm:"type:jsonb"`
	IconURL     string         `gorm:"type:varchar(1024)"`
	Downloads   int64          `gorm:"not null;default:0"`
	CreatedAt   time.Time      `gorm:"autoCreateTime;index:idx_operator_templates_created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
}

func (OperatorTemplateModel) TableName() string { return "operator_templates" }
