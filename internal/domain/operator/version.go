package operator

import (
	"time"

	"github.com/google/uuid"
)

type VersionStatus string

const (
	VersionStatusDraft    VersionStatus = "draft"
	VersionStatusTesting  VersionStatus = "testing"
	VersionStatusActive   VersionStatus = "active"
	VersionStatusArchived VersionStatus = "archived"
)

type OperatorVersion struct {
	ID          uuid.UUID
	OperatorID  uuid.UUID
	Version     string
	ExecMode    ExecMode
	ExecConfig  *ExecConfig
	InputSchema map[string]interface{}
	OutputSpec  map[string]interface{}
	Config      map[string]interface{}
	Changelog   string
	Status      VersionStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
