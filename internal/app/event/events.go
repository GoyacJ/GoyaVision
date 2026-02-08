package event

import (
	"time"

	"goyavision/internal/app/port"

	"github.com/google/uuid"
)

const (
	EventTypeAssetNew  = "asset_new"
	EventTypeAssetDone = "asset_done"
)

// AssetCreatedEvent 资产创建事件（新资产上传/生成后发布）
type AssetCreatedEvent struct {
	AssetID uuid.UUID
	At      int64
}

func (e *AssetCreatedEvent) EventType() string { return EventTypeAssetNew }
func (e *AssetCreatedEvent) OccurredAt() int64  { return e.At }

var _ port.Event = (*AssetCreatedEvent)(nil)

// AssetDoneEvent 资产就绪事件（如录制完成、任务产出资产就绪后发布）
type AssetDoneEvent struct {
	AssetID uuid.UUID
	At      int64
}

func (e *AssetDoneEvent) EventType() string { return EventTypeAssetDone }
func (e *AssetDoneEvent) OccurredAt() int64  { return e.At }

var _ port.Event = (*AssetDoneEvent)(nil)

// NewAssetCreatedEvent 构造资产创建事件
func NewAssetCreatedEvent(assetID uuid.UUID) *AssetCreatedEvent {
	return &AssetCreatedEvent{AssetID: assetID, At: time.Now().Unix()}
}

// NewAssetDoneEvent 构造资产就绪事件
func NewAssetDoneEvent(assetID uuid.UUID) *AssetDoneEvent {
	return &AssetDoneEvent{AssetID: assetID, At: time.Now().Unix()}
}
