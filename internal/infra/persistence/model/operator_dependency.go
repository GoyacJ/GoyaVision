package model

import (
	"time"

	"github.com/google/uuid"
)

type OperatorDependencyModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	OperatorID  uuid.UUID `gorm:"type:uuid;not null;index:idx_operator_dependencies_operator_id;uniqueIndex:uk_operator_dependencies_pair"`
	DependsOnID uuid.UUID `gorm:"type:uuid;not null;index:idx_operator_dependencies_depends_on_id;uniqueIndex:uk_operator_dependencies_pair"`
	MinVersion  string    `gorm:"type:varchar(50)"`
	IsOptional  bool      `gorm:"not null;default:false"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index:idx_operator_dependencies_created_at"`
}

func (OperatorDependencyModel) TableName() string { return "operator_dependencies" }
