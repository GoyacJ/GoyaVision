package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ArtifactType 产物类型
type ArtifactType string

const (
	ArtifactTypeAsset    ArtifactType = "asset"
	ArtifactTypeResult   ArtifactType = "result"
	ArtifactTypeTimeline ArtifactType = "timeline"
	ArtifactTypeReport   ArtifactType = "report"
)

// Artifact 产物实体（任务执行结果）
type Artifact struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	TaskID    uuid.UUID      `gorm:"type:uuid;not null;index:idx_artifacts_task_id"`
	Type      ArtifactType   `gorm:"type:varchar(50);not null;index:idx_artifacts_type"`
	AssetID   *uuid.UUID     `gorm:"type:uuid;index:idx_artifacts_asset_id"`
	Data      datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt time.Time      `gorm:"autoCreateTime;index:idx_artifacts_created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`

	Task  *Task       `gorm:"foreignKey:TaskID"`
	Asset *MediaAsset `gorm:"foreignKey:AssetID"`
}

func (Artifact) TableName() string { return "artifacts" }

// IsAsset 判断是否为资产产物
func (a *Artifact) IsAsset() bool {
	return a.Type == ArtifactTypeAsset
}

// IsResult 判断是否为结果产物
func (a *Artifact) IsResult() bool {
	return a.Type == ArtifactTypeResult
}

// IsTimeline 判断是否为时间轴产物
func (a *Artifact) IsTimeline() bool {
	return a.Type == ArtifactTypeTimeline
}

// IsReport 判断是否为报告产物
func (a *Artifact) IsReport() bool {
	return a.Type == ArtifactTypeReport
}

// ArtifactFilter 产物过滤器
type ArtifactFilter struct {
	TaskID  *uuid.UUID
	Type    *ArtifactType
	AssetID *uuid.UUID
	From    *time.Time
	To      *time.Time
	Limit   int
	Offset  int
}

// ArtifactData 产物数据结构
type ArtifactData struct {
	AssetInfo      *AssetInfo              `json:"asset_info,omitempty"`
	Results        []map[string]interface{} `json:"results,omitempty"`
	Timeline       []TimelineSegment        `json:"timeline,omitempty"`
	Diagnostics    map[string]interface{}   `json:"diagnostics,omitempty"`
	Summary        string                   `json:"summary,omitempty"`
	Metadata       map[string]interface{}   `json:"metadata,omitempty"`
}

// AssetInfo 资产信息
type AssetInfo struct {
	AssetID  uuid.UUID              `json:"asset_id"`
	Type     AssetType              `json:"type"`
	Path     string                 `json:"path"`
	Format   string                 `json:"format,omitempty"`
	Duration *float64               `json:"duration,omitempty"`
	Size     int64                  `json:"size,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// TimelineSegment 时间轴片段
type TimelineSegment struct {
	Start      float64                `json:"start"`
	End        float64                `json:"end"`
	EventType  string                 `json:"event_type"`
	Confidence float64                `json:"confidence,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

// AnalysisResult 分析结果
type AnalysisResult struct {
	Type       string                 `json:"type"`
	Confidence float64                `json:"confidence,omitempty"`
	Data       map[string]interface{} `json:"data"`
	Timestamp  *float64               `json:"timestamp,omitempty"`
}
