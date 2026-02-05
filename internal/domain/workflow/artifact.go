package workflow

import (
	"time"

	"github.com/google/uuid"
)

type ArtifactType string

const (
	ArtifactTypeAsset    ArtifactType = "asset"
	ArtifactTypeResult   ArtifactType = "result"
	ArtifactTypeTimeline ArtifactType = "timeline"
	ArtifactTypeReport   ArtifactType = "report"
)

type Artifact struct {
	ID        uuid.UUID
	TaskID    uuid.UUID
	Type      ArtifactType
	AssetID   *uuid.UUID
	Data      *ArtifactData
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *Artifact) IsAsset() bool {
	return a.Type == ArtifactTypeAsset
}

func (a *Artifact) IsResult() bool {
	return a.Type == ArtifactTypeResult
}

func (a *Artifact) IsTimeline() bool {
	return a.Type == ArtifactTypeTimeline
}

func (a *Artifact) IsReport() bool {
	return a.Type == ArtifactTypeReport
}

type ArtifactFilter struct {
	TaskID  *uuid.UUID
	Type    *ArtifactType
	AssetID *uuid.UUID
	From    *time.Time
	To      *time.Time
	Limit   int
	Offset  int
}

type ArtifactData struct {
	AssetInfo   *AssetInfo               `json:"asset_info,omitempty"`
	Results     []map[string]interface{} `json:"results,omitempty"`
	Timeline    []TimelineSegment        `json:"timeline,omitempty"`
	Diagnostics map[string]interface{}   `json:"diagnostics,omitempty"`
	Summary     string                   `json:"summary,omitempty"`
	Metadata    map[string]interface{}   `json:"metadata,omitempty"`
}

type AssetInfo struct {
	AssetID  uuid.UUID              `json:"asset_id"`
	Type     string                 `json:"type"`
	Path     string                 `json:"path"`
	Format   string                 `json:"format,omitempty"`
	Duration *float64               `json:"duration,omitempty"`
	Size     int64                  `json:"size,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type TimelineSegment struct {
	Start      float64                `json:"start"`
	End        float64                `json:"end"`
	EventType  string                 `json:"event_type"`
	Confidence float64                `json:"confidence,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

type AnalysisResult struct {
	Type       string                 `json:"type"`
	Confidence float64                `json:"confidence,omitempty"`
	Data       map[string]interface{} `json:"data"`
	Timestamp  *float64               `json:"timestamp,omitempty"`
}
