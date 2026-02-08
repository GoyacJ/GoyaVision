package workflow

import (
	"time"

	"github.com/google/uuid"
)

type RevisionStatus string

const (
	RevisionStatusDraft    RevisionStatus = "draft"
	RevisionStatusActive   RevisionStatus = "active"
	RevisionStatusArchived RevisionStatus = "archived"
)

type Definition struct {
	Version     string         `json:"version,omitempty"`
	TriggerType TriggerType    `json:"trigger_type"`
	TriggerConf *TriggerConfig `json:"trigger_conf,omitempty"`
	ContextSpec *ContextSpec   `json:"context_spec,omitempty"`
	Nodes       []Node         `json:"nodes,omitempty"`
	Edges       []Edge         `json:"edges,omitempty"`
}

type WorkflowRevision struct {
	ID         uuid.UUID
	WorkflowID uuid.UUID
	Revision   int64
	Status     RevisionStatus
	Definition Definition
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type RevisionFilter struct {
	WorkflowID uuid.UUID
	Limit      int
	Offset     int
}

func BuildDefinitionFromWorkflow(wf *Workflow) Definition {
	if wf == nil {
		return Definition{}
	}
	nodes := make([]Node, len(wf.Nodes))
	copy(nodes, wf.Nodes)
	edges := make([]Edge, len(wf.Edges))
	copy(edges, wf.Edges)
	return Definition{
		Version:     wf.Version,
		TriggerType: wf.TriggerType,
		TriggerConf: wf.TriggerConf,
		ContextSpec: wf.ContextSpec,
		Nodes:       nodes,
		Edges:       edges,
	}
}
