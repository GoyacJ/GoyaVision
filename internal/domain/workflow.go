package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// WorkflowStatus 工作流状态
type WorkflowStatus string

const (
	WorkflowStatusEnabled  WorkflowStatus = "enabled"
	WorkflowStatusDisabled WorkflowStatus = "disabled"
	WorkflowStatusDraft    WorkflowStatus = "draft"
)

// TriggerType 触发器类型
type TriggerType string

const (
	TriggerTypeManual    TriggerType = "manual"
	TriggerTypeSchedule  TriggerType = "schedule"
	TriggerTypeEvent     TriggerType = "event"
	TriggerTypeAssetNew  TriggerType = "asset_new"
	TriggerTypeAssetDone TriggerType = "asset_done"
)

// Workflow 工作流实体（替代 AlgorithmBinding）
type Workflow struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Code        string         `gorm:"type:varchar(100);not null;uniqueIndex"`
	Name        string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Version     string         `gorm:"type:varchar(50);not null;default:'1.0.0'"`
	TriggerType TriggerType    `gorm:"type:varchar(50);not null;index:idx_workflows_trigger_type"`
	TriggerConf datatypes.JSON `gorm:"type:jsonb"`
	Status      WorkflowStatus `gorm:"type:varchar(20);not null;default:'draft';index:idx_workflows_status"`
	Tags        datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt   time.Time      `gorm:"autoCreateTime;index:idx_workflows_created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`

	Nodes []WorkflowNode `gorm:"foreignKey:WorkflowID;constraint:OnDelete:CASCADE"`
	Edges []WorkflowEdge `gorm:"foreignKey:WorkflowID;constraint:OnDelete:CASCADE"`
}

func (Workflow) TableName() string { return "workflows" }

// IsEnabled 判断工作流是否启用
func (w *Workflow) IsEnabled() bool {
	return w.Status == WorkflowStatusEnabled
}

// IsDisabled 判断工作流是否禁用
func (w *Workflow) IsDisabled() bool {
	return w.Status == WorkflowStatusDisabled
}

// IsDraft 判断工作流是否草稿
func (w *Workflow) IsDraft() bool {
	return w.Status == WorkflowStatusDraft
}

// IsManualTrigger 判断是否为手动触发
func (w *Workflow) IsManualTrigger() bool {
	return w.TriggerType == TriggerTypeManual
}

// IsScheduleTrigger 判断是否为定时触发
func (w *Workflow) IsScheduleTrigger() bool {
	return w.TriggerType == TriggerTypeSchedule
}

// IsEventTrigger 判断是否为事件触发
func (w *Workflow) IsEventTrigger() bool {
	return w.TriggerType == TriggerTypeEvent
}

// WorkflowNode 工作流节点
type WorkflowNode struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey"`
	WorkflowID uuid.UUID      `gorm:"type:uuid;not null;index:idx_workflow_nodes_workflow_id"`
	NodeKey    string         `gorm:"type:varchar(100);not null"`
	NodeType   string         `gorm:"type:varchar(50);not null"`
	OperatorID *uuid.UUID     `gorm:"type:uuid;index:idx_workflow_nodes_operator_id"`
	Config     datatypes.JSON `gorm:"type:jsonb"`
	Position   datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`

	Operator *Operator `gorm:"foreignKey:OperatorID"`
}

func (WorkflowNode) TableName() string { return "workflow_nodes" }

// WorkflowEdge 工作流边
type WorkflowEdge struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey"`
	WorkflowID uuid.UUID      `gorm:"type:uuid;not null;index:idx_workflow_edges_workflow_id"`
	SourceKey  string         `gorm:"type:varchar(100);not null"`
	TargetKey  string         `gorm:"type:varchar(100);not null"`
	Condition  datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
}

func (WorkflowEdge) TableName() string { return "workflow_edges" }

// WorkflowFilter 工作流过滤器
type WorkflowFilter struct {
	Status      *WorkflowStatus
	TriggerType *TriggerType
	Tags        []string
	Keyword     string
	Limit       int
	Offset      int
}

// TriggerConfig 触发器配置
type TriggerConfig struct {
	Schedule      string                 `json:"schedule,omitempty"`
	IntervalSec   int                    `json:"interval_sec,omitempty"`
	EventType     string                 `json:"event_type,omitempty"`
	EventFilter   map[string]interface{} `json:"event_filter,omitempty"`
	AssetSourceID *uuid.UUID             `json:"asset_source_id,omitempty"`
	AssetType     *AssetType             `json:"asset_type,omitempty"`
}

// NodeConfig 节点配置
type NodeConfig struct {
	Params         map[string]interface{} `json:"params,omitempty"`
	RetryCount     int                    `json:"retry_count,omitempty"`
	TimeoutSeconds int                    `json:"timeout_seconds,omitempty"`
}

// EdgeCondition 边条件
type EdgeCondition struct {
	Type       string                 `json:"type,omitempty"`
	Expression string                 `json:"expression,omitempty"`
	Value      map[string]interface{} `json:"value,omitempty"`
}

// NodePosition 节点位置
type NodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
