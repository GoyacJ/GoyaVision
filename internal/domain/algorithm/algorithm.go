package algorithm

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusDraft      Status = "draft"
	StatusPublished  Status = "published"
	StatusDeprecated Status = "deprecated"
)

type VersionStatus string

const (
	VersionStatusDraft     VersionStatus = "draft"
	VersionStatusTested    VersionStatus = "tested"
	VersionStatusPublished VersionStatus = "published"
	VersionStatusArchived  VersionStatus = "archived"
)

type ImplementationType string

const (
	ImplementationOperatorVersion ImplementationType = "operator_version"
	ImplementationMCPTool         ImplementationType = "mcp_tool"
	ImplementationAIChain         ImplementationType = "ai_chain"
)

type SelectionPolicy string

const (
	SelectionPolicyStable      SelectionPolicy = "stable"
	SelectionPolicyHighQuality SelectionPolicy = "high_quality"
	SelectionPolicyLowCost     SelectionPolicy = "low_cost"
)

type Visibility int

const (
	VisibilityPrivate Visibility = 0
	VisibilityRole    Visibility = 1
	VisibilityPublic  Visibility = 2
)

type Algorithm struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	OwnerID        uuid.UUID
	Visibility     Visibility
	VisibleRoleIDs []string
	Code           string
	Name           string
	Description    string
	Scenario       string
	Status         Status
	Tags           []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Versions       []Version
}

type Version struct {
	ID                    uuid.UUID
	AlgorithmID           uuid.UUID
	Version               string
	Status                VersionStatus
	SelectionPolicy       SelectionPolicy
	DefaultImplementation *uuid.UUID
	CreatedAt             time.Time
	UpdatedAt             time.Time
	Implementations       []Implementation
	Evaluations           []EvaluationProfile
}

type Implementation struct {
	ID           uuid.UUID
	VersionID    uuid.UUID
	Name         string
	Type         ImplementationType
	BindingRef   string
	Config       map[string]interface{}
	LatencyMS    int
	CostScore    float64
	QualityScore float64
	Tier         string
	IsDefault    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type EvaluationProfile struct {
	ID               uuid.UUID
	VersionID        uuid.UUID
	DatasetRef       string
	Metrics          map[string]float64
	ReportArtifactID *uuid.UUID
	Summary          string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Filter struct {
	Status   *Status
	Scenario string
	Tags     []string
	Keyword  string
	Limit    int
	Offset   int
}
