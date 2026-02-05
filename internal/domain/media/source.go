package media

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SourceType string

const (
	SourceTypePull SourceType = "pull"
	SourceTypePush SourceType = "push"
)

type Source struct {
	ID            uuid.UUID
	Name          string
	PathName      string
	Type          SourceType
	URL           string
	Protocol      string
	Enabled       bool
	RecordEnabled bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewSource(name string, sourceType SourceType, url string, protocol string) *Source {
	return &Source{
		ID:       uuid.New(),
		Name:     name,
		PathName: GeneratePathName(name),
		Type:     sourceType,
		URL:      url,
		Protocol: protocol,
		Enabled:  true,
	}
}

func (s *Source) IsPull() bool {
	return s.Type == SourceTypePull
}

func (s *Source) IsPush() bool {
	return s.Type == SourceTypePush
}

func (s *Source) Enable() {
	s.Enabled = true
}

func (s *Source) Disable() {
	s.Enabled = false
}

func (s *Source) Validate() error {
	if s.Name == "" {
		return ErrSourceNameRequired
	}
	if s.Type != SourceTypePull && s.Type != SourceTypePush {
		return ErrInvalidSourceType
	}
	if s.Type == SourceTypePull && s.URL == "" {
		return ErrPullSourceRequiresURL
	}
	return nil
}

type SourceFilter struct {
	Type   *SourceType
	Limit  int
	Offset int
}

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
