package operator

import (
	"time"

	"goyavision/internal/domain/media"

	"github.com/google/uuid"
)

type Category string

const (
	CategoryAnalysis   Category = "analysis"
	CategoryProcessing Category = "processing"
	CategoryGeneration Category = "generation"
	CategoryUtility    Category = "utility"
)

type Type string

const (
	TypeFrameExtract    Type = "frame_extract"
	TypeObjectDetection Type = "object_detection"
	TypeOCR             Type = "ocr"
	TypeTracking        Type = "tracking"
	TypeSegmentation    Type = "segmentation"
	TypeClassification  Type = "classification"
	TypeEnhancement     Type = "enhancement"
	TypeASR             Type = "asr"
	TypeTTS             Type = "tts"
	TypeClip            Type = "clip"
	TypeRedact          Type = "redact"
	TypeSubtitle        Type = "subtitle"
	TypeShotDetect      Type = "shot_detect"
	TypeHighlight       Type = "highlight"
	TypeTranscode       Type = "transcode"
)

type Status string

const (
	StatusEnabled  Status = "enabled"
	StatusDisabled Status = "disabled"
	StatusDraft    Status = "draft"
)

type Operator struct {
	ID          uuid.UUID
	Code        string
	Name        string
	Description string
	Category    Category
	Type        Type
	Version     string
	Endpoint    string
	Method      string
	InputSchema map[string]interface{}
	OutputSpec  map[string]interface{}
	Config      map[string]interface{}
	Status      Status
	IsBuiltin   bool
	Tags        []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (o *Operator) IsEnabled() bool {
	return o.Status == StatusEnabled
}

func (o *Operator) IsDisabled() bool {
	return o.Status == StatusDisabled
}

func (o *Operator) IsDraft() bool {
	return o.Status == StatusDraft
}

func (o *Operator) IsAnalysis() bool {
	return o.Category == CategoryAnalysis
}

func (o *Operator) IsProcessing() bool {
	return o.Category == CategoryProcessing
}

func (o *Operator) IsGeneration() bool {
	return o.Category == CategoryGeneration
}

type Filter struct {
	Category  *Category
	Type      *Type
	Status    *Status
	IsBuiltin *bool
	Tags      []string
	Keyword   string
	Limit     int
	Offset    int
}

type Input struct {
	AssetID uuid.UUID              `json:"asset_id"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

type Output struct {
	OutputAssets []OutputAsset          `json:"output_assets,omitempty"`
	Results      []Result               `json:"results,omitempty"`
	Timeline     []TimelineEvent        `json:"timeline,omitempty"`
	Diagnostics  map[string]interface{} `json:"diagnostics,omitempty"`
}

type OutputAsset struct {
	Type     media.AssetType        `json:"type"`
	Path     string                 `json:"path"`
	Format   string                 `json:"format,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type Result struct {
	Type       string                 `json:"type"`
	Data       map[string]interface{} `json:"data"`
	Confidence float64                `json:"confidence,omitempty"`
}

type TimelineEvent struct {
	Start      float64                `json:"start"`
	End        float64                `json:"end"`
	EventType  string                 `json:"event_type"`
	Confidence float64                `json:"confidence,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

// Type aliases for backward compatibility with port interfaces
type OperatorInput = Input
type OperatorOutput = Output
