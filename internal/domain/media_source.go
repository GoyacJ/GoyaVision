package domain

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SourceType 媒体源类型
type SourceType string

const (
	SourceTypePull SourceType = "pull"
	SourceTypePush SourceType = "push"
)

// MediaSource 媒体源实体（与 MediaMTX path 一一对应）
type MediaSource struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Name          string     `gorm:"type:varchar(255);not null"`
	PathName      string     `gorm:"type:varchar(255);not null;uniqueIndex:idx_media_sources_path_name"`
	Type          SourceType `gorm:"type:varchar(20);not null;index:idx_media_sources_type"`
	URL           string     `gorm:"type:varchar(1024)"`
	Protocol      string     `gorm:"type:varchar(20)"`
	Enabled       bool       `gorm:"not null;default:true"`
	RecordEnabled bool       `gorm:"not null;default:false"`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime"`
}

func (MediaSource) TableName() string { return "media_sources" }

// IsPull 是否为拉流
func (s *MediaSource) IsPull() bool {
	return s.Type == SourceTypePull
}

// IsPush 是否为推流
func (s *MediaSource) IsPush() bool {
	return s.Type == SourceTypePush
}

// MediaSourceFilter 媒体源列表过滤
type MediaSourceFilter struct {
	Type   *SourceType
	Limit  int
	Offset int
}

// GeneratePathName 生成 MediaMTX path name：live/{slug}-{short_uuid}
// 过滤非法字符（仅保留字母、数字、连字符、下划线），保证全局唯一。name 为空时 slug 为 "stream"
func GeneratePathName(name string) string {
	slug := slugFromName(name)
	if slug == "" {
		slug = "stream"
	}
	short := uuid.New().String()[:8]
	return "live/" + slug + "-" + short
}

var reSlugInvalid = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
var reSlugDashes = regexp.MustCompile(`-+`)

func slugFromName(name string) string {
	s := strings.TrimSpace(name)
	s = reSlugInvalid.ReplaceAllString(s, "-")
	s = reSlugDashes.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 64 {
		s = s[:64]
		s = strings.TrimRight(s, "-")
	}
	return s
}
