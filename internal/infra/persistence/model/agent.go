package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AgentSessionModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	TaskID    uuid.UUID      `gorm:"type:uuid;not null;index:idx_agent_sessions_task_id"`
	Status    string         `gorm:"type:varchar(20);not null;index:idx_agent_sessions_status"`
	Budget    datatypes.JSON `gorm:"serializer:json"`
	StepCount int            `gorm:"not null;default:0"`
	StartedAt time.Time      `gorm:"not null;index:idx_agent_sessions_started_at"`
	EndedAt   *time.Time
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_agent_sessions_created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (AgentSessionModel) TableName() string { return "agent_sessions" }

type RunEventModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	TaskID    uuid.UUID      `gorm:"type:uuid;not null;index:idx_run_events_task_id"`
	SessionID *uuid.UUID     `gorm:"type:uuid;index:idx_run_events_session_id"`
	EventType string         `gorm:"type:varchar(50);not null;index:idx_run_events_event_type"`
	Source    string         `gorm:"type:varchar(50);not null;index:idx_run_events_source"`
	NodeKey   string         `gorm:"type:varchar(100);index:idx_run_events_node_key"`
	ToolName  string         `gorm:"type:varchar(120);index:idx_run_events_tool_name"`
	Payload   datatypes.JSON `gorm:"serializer:json"`
	CreatedAt time.Time      `gorm:"autoCreateTime;index:idx_run_events_created_at"`
}

func (RunEventModel) TableName() string { return "run_events" }

type ToolPolicyModel struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	ToolName    string         `gorm:"type:varchar(120);not null;uniqueIndex"`
	RiskLevel   string         `gorm:"type:varchar(20);not null;default:'medium'"`
	Permissions datatypes.JSON `gorm:"serializer:json"`
	DataAccess  datatypes.JSON `gorm:"serializer:json"`
	Determinism string         `gorm:"type:varchar(20);not null;default:'deterministic'"`
	Limits      datatypes.JSON `gorm:"serializer:json"`
	Enabled     bool           `gorm:"not null;default:true;index:idx_tool_policies_enabled"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
}

func (ToolPolicyModel) TableName() string { return "tool_policies" }
