package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AlgorithmModel struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey"`
	TenantID       uuid.UUID      `gorm:"type:uuid;not null;index:idx_algorithms_tenant_id"`
	OwnerID        uuid.UUID      `gorm:"type:uuid;index:idx_algorithms_owner_id"`
	Visibility     int            `gorm:"default:0;index:idx_algorithms_visibility"`
	VisibleRoleIDs datatypes.JSON `gorm:"serializer:json"`
	Code           string         `gorm:"type:varchar(100);not null;uniqueIndex"`
	Name           string         `gorm:"type:varchar(255);not null"`
	Description    string         `gorm:"type:text"`
	Scenario       string         `gorm:"type:varchar(100);index:idx_algorithms_scenario"`
	Status         string         `gorm:"type:varchar(20);not null;default:'draft';index:idx_algorithms_status"`
	Tags           datatypes.JSON `gorm:"serializer:json"`
	CreatedAt      time.Time      `gorm:"autoCreateTime;index:idx_algorithms_created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`

	Versions []AlgorithmVersionModel `gorm:"foreignKey:AlgorithmID;constraint:OnDelete:CASCADE"`
}

func (AlgorithmModel) TableName() string { return "algorithms" }

type AlgorithmVersionModel struct {
	ID                      uuid.UUID  `gorm:"type:uuid;primaryKey"`
	AlgorithmID             uuid.UUID  `gorm:"type:uuid;not null;index:idx_algorithm_versions_algorithm_id"`
	Version                 string     `gorm:"type:varchar(50);not null"`
	Status                  string     `gorm:"type:varchar(20);not null;default:'draft';index:idx_algorithm_versions_status"`
	SelectionPolicy         string     `gorm:"type:varchar(30);not null;default:'stable'"`
	DefaultImplementationID *uuid.UUID `gorm:"type:uuid;index:idx_algorithm_versions_default_impl_id"`
	CreatedAt               time.Time  `gorm:"autoCreateTime;index:idx_algorithm_versions_created_at"`
	UpdatedAt               time.Time  `gorm:"autoUpdateTime"`

	Implementations []AlgorithmImplementationModel `gorm:"foreignKey:VersionID;constraint:OnDelete:CASCADE"`
	Evaluations     []AlgorithmEvaluationModel     `gorm:"foreignKey:VersionID;constraint:OnDelete:CASCADE"`
}

func (AlgorithmVersionModel) TableName() string { return "algorithm_versions" }

type AlgorithmImplementationModel struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey"`
	VersionID    uuid.UUID      `gorm:"type:uuid;not null;index:idx_algorithm_impls_version_id"`
	Name         string         `gorm:"type:varchar(120);not null"`
	Type         string         `gorm:"type:varchar(30);not null;index:idx_algorithm_impls_type"`
	BindingRef   string         `gorm:"type:varchar(255);not null"`
	Config       datatypes.JSON `gorm:"serializer:json"`
	LatencyMS    int            `gorm:"not null;default:0"`
	CostScore    float64        `gorm:"not null;default:0"`
	QualityScore float64        `gorm:"not null;default:0"`
	Tier         string         `gorm:"type:varchar(30);not null;default:'stable';index:idx_algorithm_impls_tier"`
	IsDefault    bool           `gorm:"not null;default:false;index:idx_algorithm_impls_is_default"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
}

func (AlgorithmImplementationModel) TableName() string { return "algorithm_implementations" }

type AlgorithmEvaluationModel struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey"`
	VersionID        uuid.UUID      `gorm:"type:uuid;not null;index:idx_algorithm_eval_version_id"`
	DatasetRef       string         `gorm:"type:varchar(255);not null"`
	Metrics          datatypes.JSON `gorm:"serializer:json"`
	ReportArtifactID *uuid.UUID     `gorm:"type:uuid;index:idx_algorithm_eval_report_artifact_id"`
	Summary          string         `gorm:"type:text"`
	CreatedAt        time.Time      `gorm:"autoCreateTime;index:idx_algorithm_eval_created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime"`
}

func (AlgorithmEvaluationModel) TableName() string { return "algorithm_evaluation_profiles" }
