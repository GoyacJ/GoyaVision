package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// OperatorCategory 算子分类
type OperatorCategory string

const (
	OperatorCategoryAnalysis   OperatorCategory = "analysis"
	OperatorCategoryProcessing OperatorCategory = "processing"
	OperatorCategoryGeneration OperatorCategory = "generation"
	OperatorCategoryUtility    OperatorCategory = "utility"
)

// OperatorType 算子类型
type OperatorType string

const (
	OperatorTypeFrameExtract    OperatorType = "frame_extract"
	OperatorTypeObjectDetection OperatorType = "object_detection"
	OperatorTypeOCR             OperatorType = "ocr"
	OperatorTypeTracking        OperatorType = "tracking"
	OperatorTypeSegmentation    OperatorType = "segmentation"
	OperatorTypeClassification  OperatorType = "classification"
	OperatorTypeEnhancement     OperatorType = "enhancement"
	OperatorTypeASR             OperatorType = "asr"
	OperatorTypeTTS             OperatorType = "tts"
	OperatorTypeClip            OperatorType = "clip"
	OperatorTypeRedact          OperatorType = "redact"
	OperatorTypeSubtitle        OperatorType = "subtitle"
	OperatorTypeShotDetect      OperatorType = "shot_detect"
	OperatorTypeHighlight       OperatorType = "highlight"
	OperatorTypeTranscode       OperatorType = "transcode"
)

// OperatorStatus 算子状态
type OperatorStatus string

const (
	OperatorStatusEnabled  OperatorStatus = "enabled"
	OperatorStatusDisabled OperatorStatus = "disabled"
	OperatorStatusDraft    OperatorStatus = "draft"
)

// Operator 算子实体（重构自 Algorithm）
type Operator struct {
	ID          uuid.UUID        `gorm:"type:uuid;primaryKey"`
	Code        string           `gorm:"type:varchar(100);not null;uniqueIndex"`
	Name        string           `gorm:"type:varchar(255);not null"`
	Description string           `gorm:"type:text"`
	Category    OperatorCategory `gorm:"type:varchar(50);not null;index:idx_operators_category"`
	Type        OperatorType     `gorm:"type:varchar(50);not null;index:idx_operators_type"`
	Version     string           `gorm:"type:varchar(50);not null;default:'1.0.0'"`
	Endpoint    string           `gorm:"type:varchar(1024);not null"`
	Method      string           `gorm:"type:varchar(10);not null;default:'POST'"`
	InputSchema datatypes.JSON   `gorm:"type:jsonb"`
	OutputSpec  datatypes.JSON   `gorm:"type:jsonb"`
	Config      datatypes.JSON   `gorm:"type:jsonb"`
	Status      OperatorStatus   `gorm:"type:varchar(20);not null;default:'draft';index:idx_operators_status"`
	IsBuiltin   bool             `gorm:"not null;default:false;index:idx_operators_builtin"`
	Tags        datatypes.JSON   `gorm:"type:jsonb"`
	CreatedAt   time.Time        `gorm:"autoCreateTime;index:idx_operators_created_at"`
	UpdatedAt   time.Time        `gorm:"autoUpdateTime"`
}

func (Operator) TableName() string { return "operators" }

// IsEnabled 判断算子是否启用
func (o *Operator) IsEnabled() bool {
	return o.Status == OperatorStatusEnabled
}

// IsDisabled 判断算子是否禁用
func (o *Operator) IsDisabled() bool {
	return o.Status == OperatorStatusDisabled
}

// IsDraft 判断算子是否草稿
func (o *Operator) IsDraft() bool {
	return o.Status == OperatorStatusDraft
}

// IsAnalysis 判断是否为分析类算子
func (o *Operator) IsAnalysis() bool {
	return o.Category == OperatorCategoryAnalysis
}

// IsProcessing 判断是否为处理类算子
func (o *Operator) IsProcessing() bool {
	return o.Category == OperatorCategoryProcessing
}

// IsGeneration 判断是否为生成类算子
func (o *Operator) IsGeneration() bool {
	return o.Category == OperatorCategoryGeneration
}

// OperatorFilter 算子过滤器
type OperatorFilter struct {
	Category  *OperatorCategory
	Type      *OperatorType
	Status    *OperatorStatus
	IsBuiltin *bool
	Tags      []string
	Keyword   string
	Limit     int
	Offset    int
}

// OperatorInput 算子标准输入
type OperatorInput struct {
	AssetID uuid.UUID              `json:"asset_id"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

// OperatorOutput 算子标准输出
type OperatorOutput struct {
	OutputAssets []OutputAsset      `json:"output_assets,omitempty"`
	Results      []OperatorResult   `json:"results,omitempty"`
	Timeline     []TimelineEvent    `json:"timeline,omitempty"`
	Diagnostics  map[string]interface{} `json:"diagnostics,omitempty"`
}

// OutputAsset 输出资产
type OutputAsset struct {
	Type     AssetType              `json:"type"`
	Path     string                 `json:"path"`
	Format   string                 `json:"format,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// OperatorResult 算子结果
type OperatorResult struct {
	Type       string                 `json:"type"`
	Data       map[string]interface{} `json:"data"`
	Confidence float64                `json:"confidence,omitempty"`
}

// TimelineEvent 时间轴事件
type TimelineEvent struct {
	Start      float64                `json:"start"`
	End        float64                `json:"end"`
	EventType  string                 `json:"event_type"`
	Confidence float64                `json:"confidence,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}
