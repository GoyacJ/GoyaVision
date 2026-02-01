package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"goyavision/config"
	"goyavision/internal/adapter/mediamtx"
	"goyavision/internal/domain"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PlaybackURLs 点播 URL 集合
type PlaybackURLs struct {
	HLS string `json:"hls"`
	MP4 string `json:"mp4"`
}

// PlaybackSegment 点播段信息
type PlaybackSegment struct {
	Start       time.Time `json:"start"`
	PlaybackURL string    `json:"playback_url"`
}

type PlaybackService struct {
	repo   port.Repository
	mtxCli *mediamtx.Client
	mtxCfg config.MediaMTX
}

func NewPlaybackService(repo port.Repository, mtxCli *mediamtx.Client, mtxCfg config.MediaMTX) *PlaybackService {
	return &PlaybackService{
		repo:   repo,
		mtxCli: mtxCli,
		mtxCfg: mtxCfg,
	}
}

// GetPlaybackURL 获取录制文件的点播 URL
func (s *PlaybackService) GetPlaybackURL(ctx context.Context, streamID uuid.UUID, start time.Time) (*PlaybackURLs, error) {
	stream, err := s.repo.GetStream(ctx, streamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("stream not found")
		}
		return nil, err
	}

	pathName := s.pathName(stream)
	startStr := start.Format(time.RFC3339)

	urls := &PlaybackURLs{
		HLS: fmt.Sprintf("%s/get?path=%s&start=%s", s.mtxCfg.PlaybackAddress, pathName, startStr),
		MP4: fmt.Sprintf("%s/get?path=%s&start=%s&format=mp4", s.mtxCfg.PlaybackAddress, pathName, startStr),
	}

	return urls, nil
}

// ListRecordingSegments 列出流的录制段
func (s *PlaybackService) ListRecordingSegments(ctx context.Context, streamID uuid.UUID) ([]PlaybackSegment, error) {
	stream, err := s.repo.GetStream(ctx, streamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("stream not found")
		}
		return nil, err
	}

	pathName := s.pathName(stream)

	recording, err := s.mtxCli.GetRecordings(ctx, pathName)
	if err != nil {
		return nil, fmt.Errorf("get recordings: %w", err)
	}

	segments := make([]PlaybackSegment, len(recording.Segments))
	for i, seg := range recording.Segments {
		startStr := seg.Start.Format(time.RFC3339)
		segments[i] = PlaybackSegment{
			Start:       seg.Start,
			PlaybackURL: fmt.Sprintf("%s/get?path=%s&start=%s", s.mtxCfg.PlaybackAddress, pathName, startStr),
		}
	}

	return segments, nil
}

// pathName 生成 MediaMTX 路径名
func (s *PlaybackService) pathName(stream *domain.Stream) string {
	name := strings.ReplaceAll(stream.Name, " ", "_")
	name = strings.ToLower(name)
	return name
}
