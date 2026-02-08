package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type WorkflowModel struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey"`
	TenantID       uuid.UUID      `gorm:"type:uuid;not null;index:idx_workflows_tenant_id"`
	OwnerID        uuid.UUID      `gorm:"type:uuid;index:idx_workflows_owner_id"`
	Visibility     int            `gorm:"default:0;index:idx_workflows_visibility"`
	VisibleRoleIDs datatypes.JSON `gorm:"serializer:json"`
	Code           string         `gorm:"type:varchar(100);not null;uniqueIndex"`
	Name           string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Version     string         `gorm:"type:varchar(50);not null;default:'1.0.0'"`
	TriggerType string         `gorm:"type:varchar(50);not null;index:idx_workflows_trigger_type"`
	TriggerConf datatypes.JSON `gorm:"serializer:json"`
	Status      string         `gorm:"type:varchar(20);not null;default:'draft';index:idx_workflows_status"`
	Tags        datatypes.JSON `gorm:"serializer:json"`
	CreatedAt   time.Time      `gorm:"autoCreateTime;index:idx_workflows_created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`

	Nodes []WorkflowNodeModel `gorm:"foreignKey:WorkflowID;constraint:OnDelete:CASCADE"`
	Edges []WorkflowEdgeModel `gorm:"foreignKey:WorkflowID;constraint:OnDelete:CASCADE"`
}

func (WorkflowModel) TableName() string { return "workflows" }

type WorkflowNodeModel struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey"`
	WorkflowID uuid.UUID      `gorm:"type:uuid;not null;index:idx_workflow_nodes_workflow_id"`
	NodeKey    string         `gorm:"type:varchar(100);not null"`
	NodeType   string         `gorm:"type:varchar(50);not null"`
	OperatorID *uuid.UUID     `gorm:"type:uuid;index:idx_workflow_nodes_operator_id"`
	Config     datatypes.JSON `gorm:"serializer:json"`
	Position   datatypes.JSON `gorm:"serializer:json"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`

	Operator *OperatorModel `gorm:"foreignKey:OperatorID"`
}

func (WorkflowNodeModel) TableName() string { return "workflow_nodes" }

type WorkflowEdgeModel struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey"`
	WorkflowID uuid.UUID      `gorm:"type:uuid;not null;index:idx_workflow_edges_workflow_id"`
	SourceKey  string         `gorm:"type:varchar(100);not null"`
	TargetKey  string         `gorm:"type:varchar(100);not null"`
	Condition  datatypes.JSON `gorm:"serializer:json"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
}

func (WorkflowEdgeModel) TableName() string { return "workflow_edges" }
