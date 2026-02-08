package mapper

import (
	"encoding/json"

	"goyavision/internal/domain/workflow"
	"goyavision/internal/infra/persistence/model"

	"gorm.io/datatypes"
)

func TaskContextStateToModel(s *workflow.TaskContextState) *model.TaskContextStateModel {
	m := &model.TaskContextStateModel{
		TaskID:    s.TaskID,
		Version:   s.Version,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
	if s.Data != nil {
		data, _ := json.Marshal(s.Data)
		m.Data = datatypes.JSON(data)
	}
	return m
}

func TaskContextStateToDomain(m *model.TaskContextStateModel) *workflow.TaskContextState {
	s := &workflow.TaskContextState{
		TaskID:    m.TaskID,
		Version:   m.Version,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.Data != nil {
		_ = json.Unmarshal(m.Data, &s.Data)
	}
	return s
}

func TaskContextPatchToModel(p *workflow.TaskContextPatch) *model.TaskContextPatchModel {
	m := &model.TaskContextPatchModel{
		ID:            p.ID,
		TaskID:        p.TaskID,
		WriterNodeKey: p.WriterNodeKey,
		BeforeVersion: p.BeforeVersion,
		AfterVersion:  p.AfterVersion,
		CreatedAt:     p.CreatedAt,
	}
	data, _ := json.Marshal(p.Diff)
	m.Diff = datatypes.JSON(data)
	return m
}

func TaskContextPatchToDomain(m *model.TaskContextPatchModel) *workflow.TaskContextPatch {
	p := &workflow.TaskContextPatch{
		ID:            m.ID,
		TaskID:        m.TaskID,
		WriterNodeKey: m.WriterNodeKey,
		BeforeVersion: m.BeforeVersion,
		AfterVersion:  m.AfterVersion,
		CreatedAt:     m.CreatedAt,
	}
	if m.Diff != nil {
		_ = json.Unmarshal(m.Diff, &p.Diff)
	}
	return p
}

func TaskContextSnapshotToModel(s *workflow.TaskContextSnapshot) *model.TaskContextSnapshotModel {
	m := &model.TaskContextSnapshotModel{
		ID:        s.ID,
		TaskID:    s.TaskID,
		Version:   s.Version,
		Trigger:   s.Trigger,
		CreatedAt: s.CreatedAt,
	}
	if s.Data != nil {
		data, _ := json.Marshal(s.Data)
		m.Data = datatypes.JSON(data)
	}
	return m
}
