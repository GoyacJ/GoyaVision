package workflow

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusEnabled  Status = "enabled"
	StatusDisabled Status = "disabled"
	StatusDraft    Status = "draft"
)

type TriggerType string

const (
	TriggerTypeManual    TriggerType = "manual"
	TriggerTypeSchedule  TriggerType = "schedule"
	TriggerTypeEvent     TriggerType = "event"
	TriggerTypeAssetNew  TriggerType = "asset_new"
	TriggerTypeAssetDone TriggerType = "asset_done"
)

type Visibility int

const (
	VisibilityPrivate Visibility = 0
	VisibilityRole    Visibility = 1
	VisibilityPublic  Visibility = 2
)

type Workflow struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	OwnerID        uuid.UUID
	Visibility     Visibility
	VisibleRoleIDs []string
	Code           string
	Name           string
	Description string
	Version     string
	TriggerType TriggerType
	TriggerConf *TriggerConfig
	Status      Status
	Tags        []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Nodes       []Node
	Edges       []Edge
}

func (w *Workflow) IsEnabled() bool {
	return w.Status == StatusEnabled
}

func (w *Workflow) IsDisabled() bool {
	return w.Status == StatusDisabled
}

func (w *Workflow) IsDraft() bool {
	return w.Status == StatusDraft
}

func (w *Workflow) IsManualTrigger() bool {
	return w.TriggerType == TriggerTypeManual
}

func (w *Workflow) IsScheduleTrigger() bool {
	return w.TriggerType == TriggerTypeSchedule
}

func (w *Workflow) IsEventTrigger() bool {
	return w.TriggerType == TriggerTypeEvent
}

func (w *Workflow) Activate() error {
	if len(w.Nodes) == 0 {
		return errors.New("workflow must have at least one node to activate")
	}
	w.Status = StatusEnabled
	return nil
}

func (w *Workflow) Pause() {
	w.Status = StatusDisabled
}

func (w *Workflow) AddNode(node Node) {
	w.Nodes = append(w.Nodes, node)
}

func (w *Workflow) AddEdge(edge Edge) {
	w.Edges = append(w.Edges, edge)
}

func (w *Workflow) Validate() error {
	if w.Code == "" {
		return errors.New("workflow code is required")
	}
	if w.Name == "" {
		return errors.New("workflow name is required")
	}
	return nil
}

type Node struct {
	ID         uuid.UUID
	WorkflowID uuid.UUID
	NodeKey    string
	NodeType   string
	OperatorID *uuid.UUID
	Config     *NodeConfig
	Position   *NodePosition
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Edge struct {
	ID         uuid.UUID
	WorkflowID uuid.UUID
	SourceKey  string
	TargetKey  string
	Condition  *EdgeCondition
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Filter struct {
	Status      *Status
	TriggerType *TriggerType
	Tags        []string
	Keyword     string
	Limit       int
	Offset      int
}

type TriggerConfig struct {
	Schedule    string                 `json:"schedule,omitempty"`
	IntervalSec int                    `json:"interval_sec,omitempty"`
	EventType   string                 `json:"event_type,omitempty"`
	EventFilter map[string]interface{} `json:"event_filter,omitempty"`
}

type NodeConfig struct {
	Params         map[string]interface{} `json:"params,omitempty"`
	RetryCount     int                    `json:"retry_count,omitempty"`
	TimeoutSeconds int                    `json:"timeout_seconds,omitempty"`
}

type EdgeCondition struct {
	Type       string                 `json:"type,omitempty"`
	Expression string                 `json:"expression,omitempty"`
	Value      map[string]interface{} `json:"value,omitempty"`
}

type NodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
