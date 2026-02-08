package mapper

import (
	"encoding/json"

	"goyavision/internal/domain/agent"
	"goyavision/internal/infra/persistence/model"

	"gorm.io/datatypes"
)

func AgentSessionToModel(s *agent.Session) *model.AgentSessionModel {
	m := &model.AgentSessionModel{
		ID:        s.ID,
		TaskID:    s.TaskID,
		Status:    string(s.Status),
		StepCount: s.StepCount,
		StartedAt: s.StartedAt,
		EndedAt:   s.EndedAt,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
	if s.Budget != nil {
		b, _ := json.Marshal(s.Budget)
		m.Budget = datatypes.JSON(b)
	}
	return m
}

func AgentSessionToDomain(m *model.AgentSessionModel) *agent.Session {
	s := &agent.Session{
		ID:        m.ID,
		TaskID:    m.TaskID,
		Status:    agent.SessionStatus(m.Status),
		StepCount: m.StepCount,
		StartedAt: m.StartedAt,
		EndedAt:   m.EndedAt,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.Budget != nil {
		_ = json.Unmarshal(m.Budget, &s.Budget)
	}
	return s
}

func RunEventToModel(e *agent.RunEvent) *model.RunEventModel {
	m := &model.RunEventModel{
		ID:        e.ID,
		TaskID:    e.TaskID,
		SessionID: e.SessionID,
		EventType: string(e.EventType),
		Source:    e.Source,
		NodeKey:   e.NodeKey,
		ToolName:  e.ToolName,
		CreatedAt: e.CreatedAt,
	}
	if e.Payload != nil {
		b, _ := json.Marshal(e.Payload)
		m.Payload = datatypes.JSON(b)
	}
	return m
}

func RunEventToDomain(m *model.RunEventModel) *agent.RunEvent {
	e := &agent.RunEvent{
		ID:        m.ID,
		TaskID:    m.TaskID,
		SessionID: m.SessionID,
		EventType: agent.EventType(m.EventType),
		Source:    m.Source,
		NodeKey:   m.NodeKey,
		ToolName:  m.ToolName,
		CreatedAt: m.CreatedAt,
	}
	if m.Payload != nil {
		_ = json.Unmarshal(m.Payload, &e.Payload)
	}
	return e
}

func ToolPolicyToModel(p *agent.ToolPolicy) *model.ToolPolicyModel {
	m := &model.ToolPolicyModel{
		ID:          p.ID,
		ToolName:    p.ToolName,
		RiskLevel:   p.RiskLevel,
		Determinism: string(p.Determinism),
		Enabled:     p.Enabled,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	if p.Permissions != nil {
		b, _ := json.Marshal(p.Permissions)
		m.Permissions = datatypes.JSON(b)
	}
	{
		b, _ := json.Marshal(p.DataAccess)
		m.DataAccess = datatypes.JSON(b)
	}
	if p.Limits != nil {
		b, _ := json.Marshal(p.Limits)
		m.Limits = datatypes.JSON(b)
	}
	return m
}

func ToolPolicyToDomain(m *model.ToolPolicyModel) *agent.ToolPolicy {
	p := &agent.ToolPolicy{
		ID:          m.ID,
		ToolName:    m.ToolName,
		RiskLevel:   m.RiskLevel,
		Determinism: agent.Determinism(m.Determinism),
		Enabled:     m.Enabled,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	if m.Permissions != nil {
		_ = json.Unmarshal(m.Permissions, &p.Permissions)
	}
	if m.DataAccess != nil {
		_ = json.Unmarshal(m.DataAccess, &p.DataAccess)
	}
	if m.Limits != nil {
		_ = json.Unmarshal(m.Limits, &p.Limits)
	}
	return p
}
