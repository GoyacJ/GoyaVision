package agent

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Determinism string

const (
	DeterminismDeterministic Determinism = "deterministic"
	DeterminismStochastic    Determinism = "stochastic"
)

type SessionStatus string

const (
	SessionStatusRunning   SessionStatus = "running"
	SessionStatusSucceeded SessionStatus = "succeeded"
	SessionStatusFailed    SessionStatus = "failed"
	SessionStatusCancelled SessionStatus = "cancelled"
)

type EventType string

const (
	EventTypeNodeStarted     EventType = "node_started"
	EventTypeNodeSucceeded   EventType = "node_succeeded"
	EventTypeNodeFailed      EventType = "node_failed"
	EventTypeToolCalled      EventType = "tool_called"
	EventTypeToolFailed      EventType = "tool_failed"
	EventTypeRecover         EventType = "recover"
	EventTypeAgentDecision   EventType = "agent_decision"
	EventTypeAgentAction     EventType = "agent_action"
	EventTypeAgentEscalation EventType = "agent_escalation"
)

type ErrorCategory string

const (
	ErrorCategoryTransient      ErrorCategory = "TRANSIENT"
	ErrorCategoryTimeout        ErrorCategory = "TIMEOUT"
	ErrorCategoryResourceLimit  ErrorCategory = "RESOURCE_LIMIT"
	ErrorCategoryValidation     ErrorCategory = "VALIDATION"
	ErrorCategoryPolicyDeny     ErrorCategory = "POLICY_DENY"
	ErrorCategoryDependency     ErrorCategory = "DEPENDENCY"
	ErrorCategoryToolBug        ErrorCategory = "TOOL_BUG"
	ErrorCategoryModelReasoning ErrorCategory = "MODEL_REASONING"
	ErrorCategoryUnknown        ErrorCategory = "UNKNOWN"
)

type DataAccess struct {
	ReadScopes       []string `json:"read_scopes,omitempty"`
	WriteScopes      []string `json:"write_scopes,omitempty"`
	NetworkAllowlist []string `json:"network_allowlist,omitempty"`
}

type ToolSpec struct {
	Name                string                 `json:"name"`
	Description         string                 `json:"description,omitempty"`
	InputSchema         map[string]interface{} `json:"input_schema,omitempty"`
	OutputSchema        map[string]interface{} `json:"output_schema,omitempty"`
	Determinism         Determinism            `json:"determinism"`
	DataAccess          DataAccess             `json:"data_access"`
	TimeoutSeconds      int                    `json:"timeout_seconds,omitempty"`
	RetryCount          int                    `json:"retry_count,omitempty"`
	Idempotent          bool                   `json:"idempotent,omitempty"`
	RiskLevel           string                 `json:"risk_level,omitempty"`
	RequiredPermissions []string               `json:"required_permissions,omitempty"`
}

type RecoveryPolicy struct {
	Retryable    bool   `json:"retryable"`
	Backoff      string `json:"backoff,omitempty"`
	Fallback     string `json:"fallback,omitempty"`
	RequireHuman bool   `json:"require_human"`
	Severity     string `json:"severity"`
}

var DefaultRecoveryPolicies = map[ErrorCategory]RecoveryPolicy{
	ErrorCategoryTransient:      {Retryable: true, Backoff: "exponential", Fallback: "enabled", RequireHuman: false, Severity: "medium"},
	ErrorCategoryTimeout:        {Retryable: true, Backoff: "exponential", Fallback: "enabled", RequireHuman: false, Severity: "medium"},
	ErrorCategoryResourceLimit:  {Retryable: true, Backoff: "linear_throttle", Fallback: "enabled", RequireHuman: false, Severity: "high"},
	ErrorCategoryValidation:     {Retryable: false, Backoff: "none", Fallback: "param_repair", RequireHuman: true, Severity: "medium"},
	ErrorCategoryPolicyDeny:     {Retryable: false, Backoff: "none", Fallback: "none", RequireHuman: true, Severity: "high"},
	ErrorCategoryDependency:     {Retryable: true, Backoff: "exponential", Fallback: "enabled", RequireHuman: false, Severity: "high"},
	ErrorCategoryToolBug:        {Retryable: false, Backoff: "none", Fallback: "switch_impl", RequireHuman: true, Severity: "high"},
	ErrorCategoryModelReasoning: {Retryable: true, Backoff: "short_retry", Fallback: "switch_model", RequireHuman: false, Severity: "medium"},
	ErrorCategoryUnknown:        {Retryable: false, Backoff: "none", Fallback: "none", RequireHuman: true, Severity: "high"},
}

type ToolError struct {
	Category     ErrorCategory `json:"category"`
	RootCause    string        `json:"root_cause,omitempty"`
	ActionHint   string        `json:"action_hint,omitempty"`
	Retryable    bool          `json:"retryable"`
	ProviderCode string        `json:"provider_code,omitempty"`
	Message      string        `json:"message"`
}

func (e *ToolError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("%s: %s", e.Category, e.Message)
}

type Session struct {
	ID        uuid.UUID
	TaskID    uuid.UUID
	Status    SessionStatus
	Budget    map[string]interface{}
	StepCount int
	StartedAt time.Time
	EndedAt   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RunEvent struct {
	ID        uuid.UUID
	TaskID    uuid.UUID
	SessionID *uuid.UUID
	EventType EventType
	Source    string
	NodeKey   string
	ToolName  string
	Payload   map[string]interface{}
	CreatedAt time.Time
}

type ToolPolicy struct {
	ID          uuid.UUID
	ToolName    string
	RiskLevel   string
	Permissions []string
	DataAccess  DataAccess
	Determinism Determinism
	Limits      map[string]interface{}
	Enabled     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
