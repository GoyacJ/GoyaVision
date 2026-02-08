package workflow

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	SharedConflictReject    = "reject"
	SharedConflictOverwrite = "overwrite"
	SharedConflictMerge     = "merge"
	SharedConflictAppend    = "append"
)

var ErrContextVersionConflict = errors.New("task context version conflict")

type ContextSpec struct {
	Vars       map[string]ContextVarSpec `json:"vars,omitempty"`
	SharedKeys map[string]SharedKeySpec  `json:"shared_keys,omitempty"`
	Meta       map[string]interface{}    `json:"meta,omitempty"`
}

type ContextVarSpec struct {
	Type        string      `json:"type,omitempty"`
	Required    bool        `json:"required,omitempty"`
	ReadOnly    bool        `json:"readonly,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description,omitempty"`
}

type SharedKeySpec struct {
	Type           string `json:"type,omitempty"`
	ConflictPolicy string `json:"conflict_policy,omitempty"`
	CAS            bool   `json:"cas,omitempty"`
	Description    string `json:"description,omitempty"`
}

type TaskContextState struct {
	TaskID    uuid.UUID
	Version   int64
	Data      map[string]interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ContextDiff struct {
	Set   map[string]interface{} `json:"set,omitempty"`
	Unset []string               `json:"unset,omitempty"`
}

type TaskContextPatch struct {
	ID            uuid.UUID
	TaskID        uuid.UUID
	WriterNodeKey string
	BeforeVersion int64
	AfterVersion  int64
	Diff          ContextDiff
	CreatedAt     time.Time
}

type TaskContextSnapshot struct {
	ID        uuid.UUID
	TaskID    uuid.UUID
	Version   int64
	Data      map[string]interface{}
	Trigger   string
	CreatedAt time.Time
}
