package agent

import (
	"context"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	Get(ctx context.Context, id uuid.UUID) (*Session, error)
	Update(ctx context.Context, session *Session) error
	List(ctx context.Context, filter SessionFilter) ([]*Session, int64, error)
}

type SessionFilter struct {
	TaskID *uuid.UUID
	Status *SessionStatus
	Limit  int
	Offset int
}

type EventFilter struct {
	TaskID    *uuid.UUID
	SessionID *uuid.UUID
	Limit     int
	Offset    int
	Source    string
	NodeKey   string
}

type RunEventRepository interface {
	Create(ctx context.Context, event *RunEvent) error
	List(ctx context.Context, filter EventFilter) ([]*RunEvent, int64, error)
}

type ToolPolicyRepository interface {
	Upsert(ctx context.Context, policy *ToolPolicy) error
	GetByToolName(ctx context.Context, toolName string) (*ToolPolicy, error)
	ListEnabled(ctx context.Context) ([]*ToolPolicy, error)
}
