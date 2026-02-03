package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// AssetType 资产类型
type AssetType string

const (
	AssetTypeVideo  AssetType = "video"
	AssetTypeImage  AssetType = "image"
	AssetTypeAudio  AssetType = "audio"
	AssetTypeStream AssetType = "stream" // 流媒体类型
)

// AssetSourceType 资产来源类型
type AssetSourceType string

const (
	AssetSourceLive      AssetSourceType = "live"
	AssetSourceVOD       AssetSourceType = "vod"
	AssetSourceUpload    AssetSourceType = "upload"
	AssetSourceGenerated AssetSourceType = "generated"
)

// AssetStatus 资产状态
type AssetStatus string

const (
	AssetStatusPending AssetStatus = "pending"
	AssetStatusReady   AssetStatus = "ready"
	AssetStatusFailed  AssetStatus = "failed"
)

// MediaAsset 媒体资产实体
type MediaAsset struct {
	ID         uuid.UUID       `gorm:"type:uuid;primaryKey"`
	Type       AssetType       `gorm:"type:varchar(20);not null;index:idx_assets_type"`
	SourceType AssetSourceType `gorm:"type:varchar(20);not null;index:idx_assets_source_type"`
	SourceID   *uuid.UUID      `gorm:"type:uuid;index:idx_assets_source_id"`
	ParentID   *uuid.UUID      `gorm:"type:uuid;index:idx_assets_parent_id"`
	Name       string          `gorm:"type:varchar(255);not null"`
	Path       string          `gorm:"type:varchar(1024);not null"`
	Duration   *float64        `gorm:"type:float8"`
	Size       int64           `gorm:"not null;default:0"`
	Format     string          `gorm:"type:varchar(50)"`
	Metadata   datatypes.JSON  `gorm:"type:jsonb"`
	Status     AssetStatus     `gorm:"type:varchar(20);not null;default:'pending';index:idx_assets_status"`
	Tags       datatypes.JSON  `gorm:"type:jsonb"`
	CreatedAt  time.Time       `gorm:"autoCreateTime;index:idx_assets_created_at"`
	UpdatedAt  time.Time       `gorm:"autoUpdateTime"`
}

func (MediaAsset) TableName() string { return "media_assets" }

// IsVideo 判断是否为视频资产
func (a *MediaAsset) IsVideo() bool {
	return a.Type == AssetTypeVideo
}

// IsImage 判断是否为图片资产
func (a *MediaAsset) IsImage() bool {
	return a.Type == AssetTypeImage
}

// IsAudio 判断是否为音频资产
func (a *MediaAsset) IsAudio() bool {
	return a.Type == AssetTypeAudio
}

// IsStream 判断是否为流媒体资产
func (a *MediaAsset) IsStream() bool {
	return a.Type == AssetTypeStream
}

// HasParent 判断是否有父资产（派生资产）
func (a *MediaAsset) HasParent() bool {
	return a.ParentID != nil
}

// HasSource 判断是否关联媒体源
func (a *MediaAsset) HasSource() bool {
	return a.SourceID != nil
}

// IsReady 判断资产是否就绪
func (a *MediaAsset) IsReady() bool {
	return a.Status == AssetStatusReady
}

// IsFailed 判断资产是否失败
func (a *MediaAsset) IsFailed() bool {
	return a.Status == AssetStatusFailed
}

// MediaAssetFilter 媒体资产过滤器
type MediaAssetFilter struct {
	Type       *AssetType
	SourceType *AssetSourceType
	SourceID   *uuid.UUID
	ParentID   *uuid.UUID
	Status     *AssetStatus
	Tags       []string
	From       *time.Time
	To         *time.Time
	Limit      int
	Offset     int
}
