package dto

import (
	"time"

	"goyavision/internal/domain/agent"

	"github.com/google/uuid"
)

type AgentSessionCreateReq struct {
	TaskID uuid.UUID              `json:"task_id" validate:"required"`
	Budget map[string]interface{} `json:"budget,omitempty"`
}

type AgentSessionStopReq struct {
	Status string `json:"status,omitempty"`
}

type AgentSessionRunReq struct {
	MaxActions int `json:"max_actions,omitempty"`
}

type AgentSessionListQuery struct {
	TaskID *uuid.UUID `query:"task_id"`
	Status *string    `query:"status"`
	Limit  int        `query:"limit"`
	Offset int        `query:"offset"`
}

type AgentSessionEventListQuery struct {
	Source  *string `query:"source"`
	NodeKey *string `query:"node_key"`
	Limit   int     `query:"limit"`
	Offset  int     `query:"offset"`
}

type AgentSessionResponse struct {
	ID        uuid.UUID              `json:"id"`
	TaskID    uuid.UUID              `json:"task_id"`
	Status    string                 `json:"status"`
	Budget    map[string]interface{} `json:"budget,omitempty"`
	StepCount int                    `json:"step_count"`
	StartedAt time.Time              `json:"started_at"`
	EndedAt   *time.Time             `json:"ended_at,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type AgentSessionListResponse struct {
	Items []*AgentSessionResponse `json:"items"`
	Total int64                   `json:"total"`
}

type AgentSessionEventResponse struct {
	ID        uuid.UUID              `json:"id"`
	TaskID    uuid.UUID              `json:"task_id"`
	SessionID *uuid.UUID             `json:"session_id,omitempty"`
	EventType string                 `json:"event_type"`
	Source    string                 `json:"source"`
	NodeKey   string                 `json:"node_key,omitempty"`
	ToolName  string                 `json:"tool_name,omitempty"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

type AgentSessionEventListResponse struct {
	Items []*AgentSessionEventResponse `json:"items"`
	Total int64                        `json:"total"`
}

func AgentSessionToResponse(session *agent.Session) *AgentSessionResponse {
	if session == nil {
		return nil
	}
	return &AgentSessionResponse{
		ID:        session.ID,
		TaskID:    session.TaskID,
		Status:    string(session.Status),
		Budget:    session.Budget,
		StepCount: session.StepCount,
		StartedAt: session.StartedAt,
		EndedAt:   session.EndedAt,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
	}
}

func AgentSessionsToResponse(items []*agent.Session) []*AgentSessionResponse {
	out := make([]*AgentSessionResponse, len(items))
	for i := range items {
		out[i] = AgentSessionToResponse(items[i])
	}
	return out
}

func AgentSessionEventToResponse(event *agent.RunEvent) *AgentSessionEventResponse {
	if event == nil {
		return nil
	}
	return &AgentSessionEventResponse{
		ID:        event.ID,
		TaskID:    event.TaskID,
		SessionID: event.SessionID,
		EventType: string(event.EventType),
		Source:    event.Source,
		NodeKey:   event.NodeKey,
		ToolName:  event.ToolName,
		Payload:   event.Payload,
		CreatedAt: event.CreatedAt,
	}
}

func AgentSessionEventsToResponse(items []*agent.RunEvent) []*AgentSessionEventResponse {
	out := make([]*AgentSessionEventResponse, len(items))
	for i := range items {
		out[i] = AgentSessionEventToResponse(items[i])
	}
	return out
}
