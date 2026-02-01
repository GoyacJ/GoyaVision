package domain

import (
	"time"

	"github.com/google/uuid"
)

// StreamType 流类型
type StreamType string

const (
	StreamTypePull StreamType = "pull"
	StreamTypePush StreamType = "push"
)

// Stream 视频流实体
type Stream struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	URL       string     `gorm:"not null"`
	Name      string     `gorm:"not null;uniqueIndex"`
	Type      StreamType `gorm:"type:varchar(20);default:pull"`
	Enabled   bool       `gorm:"default:true"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
}

func (Stream) TableName() string { return "streams" }

// StreamStatus 流状态（非持久化，从 MediaMTX 获取）
type StreamStatus struct {
	StreamID      uuid.UUID
	PathName      string
	Ready         bool
	Online        bool
	Tracks        []string
	BytesReceived uint64
	BytesSent     uint64
	ReaderCount   int

	RTSPUrl   string
	RTMPUrl   string
	HLSUrl    string
	WebRTCUrl string
}

// StreamWithStatus 流及其实时状态
type StreamWithStatus struct {
	Stream *Stream
	Status *StreamStatus
}
