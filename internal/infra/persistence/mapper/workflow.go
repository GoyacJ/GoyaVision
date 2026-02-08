package mapper

import (
	"encoding/json"

	"goyavision/internal/domain/workflow"
	"goyavision/internal/infra/persistence/model"

	"gorm.io/datatypes"
)

func WorkflowToModel(w *workflow.Workflow) *model.WorkflowModel {
	m := &model.WorkflowModel{
		ID:             w.ID,
		TenantID:       w.TenantID,
		OwnerID:        w.OwnerID,
		Visibility:     int(w.Visibility),
		Code:           w.Code,
		Name:           w.Name,
		Description: w.Description,
		Version:     w.Version,
		TriggerType: string(w.TriggerType),
		Status:      string(w.Status),
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
	}
	if w.TriggerConf != nil {
		data, _ := json.Marshal(w.TriggerConf)
		m.TriggerConf = datatypes.JSON(data)
	}
	if w.Tags != nil {
		data, _ := json.Marshal(w.Tags)
		m.Tags = datatypes.JSON(data)
	}
	if w.VisibleRoleIDs != nil {
		data, _ := json.Marshal(w.VisibleRoleIDs)
		m.VisibleRoleIDs = datatypes.JSON(data)
	}
	for _, n := range w.Nodes {
		m.Nodes = append(m.Nodes, *NodeToModel(&n))
	}
	for _, e := range w.Edges {
		m.Edges = append(m.Edges, *EdgeToModel(&e))
	}
	return m
}

func WorkflowToDomain(m *model.WorkflowModel) *workflow.Workflow {
	w := &workflow.Workflow{
		ID:             m.ID,
		TenantID:       m.TenantID,
		OwnerID:        m.OwnerID,
		Visibility:     workflow.Visibility(m.Visibility),
		Code:           m.Code,
		Name:           m.Name,
		Description: m.Description,
		Version:     m.Version,
		TriggerType: workflow.TriggerType(m.TriggerType),
		Status:      workflow.Status(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	if m.TriggerConf != nil {
		var tc workflow.TriggerConfig
		if err := json.Unmarshal(m.TriggerConf, &tc); err == nil {
			w.TriggerConf = &tc
		}
	}
	if m.Tags != nil {
		_ = json.Unmarshal(m.Tags, &w.Tags)
	}
	if m.VisibleRoleIDs != nil {
		_ = json.Unmarshal(m.VisibleRoleIDs, &w.VisibleRoleIDs)
	}
	for _, n := range m.Nodes {
		w.Nodes = append(w.Nodes, *NodeToDomain(&n))
	}
	for _, e := range m.Edges {
		w.Edges = append(w.Edges, *EdgeToDomain(&e))
	}
	return w
}

func NodeToModel(n *workflow.Node) *model.WorkflowNodeModel {
	m := &model.WorkflowNodeModel{
		ID:         n.ID,
		WorkflowID: n.WorkflowID,
		NodeKey:    n.NodeKey,
		NodeType:   n.NodeType,
		OperatorID: n.OperatorID,
		CreatedAt:  n.CreatedAt,
		UpdatedAt:  n.UpdatedAt,
	}
	if n.Config != nil {
		data, _ := json.Marshal(n.Config)
		m.Config = datatypes.JSON(data)
	}
	if n.Position != nil {
		data, _ := json.Marshal(n.Position)
		m.Position = datatypes.JSON(data)
	}
	return m
}

func NodeToDomain(m *model.WorkflowNodeModel) *workflow.Node {
	n := &workflow.Node{
		ID:         m.ID,
		WorkflowID: m.WorkflowID,
		NodeKey:    m.NodeKey,
		NodeType:   m.NodeType,
		OperatorID: m.OperatorID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
	if m.Config != nil {
		var nc workflow.NodeConfig
		if err := json.Unmarshal(m.Config, &nc); err == nil {
			n.Config = &nc
		}
	}
	if m.Position != nil {
		var np workflow.NodePosition
		if err := json.Unmarshal(m.Position, &np); err == nil {
			n.Position = &np
		}
	}
	return n
}

func EdgeToModel(e *workflow.Edge) *model.WorkflowEdgeModel {
	m := &model.WorkflowEdgeModel{
		ID:         e.ID,
		WorkflowID: e.WorkflowID,
		SourceKey:  e.SourceKey,
		TargetKey:  e.TargetKey,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
	if e.Condition != nil {
		data, _ := json.Marshal(e.Condition)
		m.Condition = datatypes.JSON(data)
	}
	return m
}

func EdgeToDomain(m *model.WorkflowEdgeModel) *workflow.Edge {
	e := &workflow.Edge{
		ID:         m.ID,
		WorkflowID: m.WorkflowID,
		SourceKey:  m.SourceKey,
		TargetKey:  m.TargetKey,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
	if m.Condition != nil {
		var ec workflow.EdgeCondition
		if err := json.Unmarshal(m.Condition, &ec); err == nil {
			e.Condition = &ec
		}
	}
	return e
}

func TaskToModel(t *workflow.Task) *model.TaskModel {
	m := &model.TaskModel{
		ID:                t.ID,
		TenantID:          t.TenantID,
		TriggeredByUserID: t.TriggeredByUserID,
		WorkflowID:        t.WorkflowID,
		AssetID:           t.AssetID,
		Status:      string(t.Status),
		Progress:    t.Progress,
		CurrentNode: t.CurrentNode,
		Error:       t.Error,
		StartedAt:   t.StartedAt,
		CompletedAt: t.CompletedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
	if t.InputParams != nil {
		data, _ := json.Marshal(t.InputParams)
		m.InputParams = datatypes.JSON(data)
	}
	if t.NodeExecutions != nil {
		data, _ := json.Marshal(t.NodeExecutions)
		m.NodeExecutions = datatypes.JSON(data)
	}
	return m
}

func TaskToDomain(m *model.TaskModel) *workflow.Task {
	t := &workflow.Task{
		ID:                m.ID,
		TenantID:          m.TenantID,
		TriggeredByUserID: m.TriggeredByUserID,
		WorkflowID:        m.WorkflowID,
		AssetID:           m.AssetID,
		Status:      workflow.TaskStatus(m.Status),
		Progress:    m.Progress,
		CurrentNode: m.CurrentNode,
		Error:       m.Error,
		StartedAt:   m.StartedAt,
		CompletedAt: m.CompletedAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	if m.InputParams != nil {
		_ = json.Unmarshal(m.InputParams, &t.InputParams)
	}
	if m.NodeExecutions != nil {
		_ = json.Unmarshal(m.NodeExecutions, &t.NodeExecutions)
	}
	return t
}

func ArtifactToModel(a *workflow.Artifact) *model.ArtifactModel {
	m := &model.ArtifactModel{
		ID:        a.ID,
		TenantID:  a.TenantID,
		TaskID:    a.TaskID,
		Type:      string(a.Type),
		AssetID:   a.AssetID,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
	if a.Data != nil {
		data, _ := json.Marshal(a.Data)
		m.Data = datatypes.JSON(data)
	}
	return m
}

func ArtifactToDomain(m *model.ArtifactModel) *workflow.Artifact {
	a := &workflow.Artifact{
		ID:        m.ID,
		TenantID:  m.TenantID,
		TaskID:    m.TaskID,
		Type:      workflow.ArtifactType(m.Type),
		AssetID:   m.AssetID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.Data != nil {
		var ad workflow.ArtifactData
		if err := json.Unmarshal(m.Data, &ad); err == nil {
			a.Data = &ad
		}
	}
	return a
}
