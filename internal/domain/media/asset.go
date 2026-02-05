package media

import (
	"time"

	"github.com/google/uuid"
)

type AssetType string

const (
	AssetTypeVideo  AssetType = "video"
	AssetTypeImage  AssetType = "image"
	AssetTypeAudio  AssetType = "audio"
	AssetTypeStream AssetType = "stream"
)

type AssetSourceType string

const (
	AssetSourceLive      AssetSourceType = "live"
	AssetSourceVOD       AssetSourceType = "vod"
	AssetSourceUpload    AssetSourceType = "upload"
	AssetSourceGenerated AssetSourceType = "generated"
)

type AssetStatus string

const (
	AssetStatusPending AssetStatus = "pending"
	AssetStatusReady   AssetStatus = "ready"
	AssetStatusFailed  AssetStatus = "failed"
)

type Asset struct {
	ID         uuid.UUID
	Type       AssetType
	SourceType AssetSourceType
	SourceID   *uuid.UUID
	ParentID   *uuid.UUID
	Name       string
	Path       string
	Duration   *float64
	Size       int64
	Format     string
	Metadata   map[string]interface{}
	Status     AssetStatus
	Tags       []string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (a *Asset) IsVideo() bool {
	return a.Type == AssetTypeVideo
}

func (a *Asset) IsImage() bool {
	return a.Type == AssetTypeImage
}

func (a *Asset) IsAudio() bool {
	return a.Type == AssetTypeAudio
}

func (a *Asset) IsStream() bool {
	return a.Type == AssetTypeStream
}

func (a *Asset) HasParent() bool {
	return a.ParentID != nil
}

func (a *Asset) HasSource() bool {
	return a.SourceID != nil
}

func (a *Asset) IsReady() bool {
	return a.Status == AssetStatusReady
}

func (a *Asset) IsFailed() bool {
	return a.Status == AssetStatusFailed
}

type AssetFilter struct {
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
