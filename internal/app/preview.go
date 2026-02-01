package app

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"goyavision/config"
	"goyavision/internal/adapter/mediamtx"
	"goyavision/internal/domain"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PreviewURLs 预览 URL 集合
type PreviewURLs struct {
	HLS    string `json:"hls"`
	RTSP   string `json:"rtsp"`
	RTMP   string `json:"rtmp"`
	WebRTC string `json:"webrtc"`
}

type PreviewService struct {
	repo   port.Repository
	mtxCli *mediamtx.Client
	mtxCfg config.MediaMTX
}

func NewPreviewService(repo port.Repository, mtxCli *mediamtx.Client, mtxCfg config.MediaMTX) *PreviewService {
	return &PreviewService{
		repo:   repo,
		mtxCli: mtxCli,
		mtxCfg: mtxCfg,
	}
}

// GetPreviewURLs 获取流的预览 URL
func (s *PreviewService) GetPreviewURLs(ctx context.Context, streamID uuid.UUID) (*PreviewURLs, error) {
	stream, err := s.repo.GetStream(ctx, streamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("stream not found")
		}
		return nil, err
	}

	if !stream.Enabled {
		return nil, errors.New("stream is disabled")
	}

	pathName := s.pathName(stream)

	urls := &PreviewURLs{
		HLS:    fmt.Sprintf("%s/%s/index.m3u8", s.mtxCfg.HLSAddress, pathName),
		RTSP:   fmt.Sprintf("%s/%s", s.mtxCfg.RTSPAddress, pathName),
		RTMP:   fmt.Sprintf("%s/%s", s.mtxCfg.RTMPAddress, pathName),
		WebRTC: fmt.Sprintf("%s/%s/whep", s.mtxCfg.WebRTCAddress, pathName),
	}

	return urls, nil
}

// Start 开始预览（返回所有协议 URL）
func (s *PreviewService) Start(ctx context.Context, streamID uuid.UUID) (*PreviewURLs, error) {
	return s.GetPreviewURLs(ctx, streamID)
}

// IsStreamReady 检查流是否就绪
func (s *PreviewService) IsStreamReady(ctx context.Context, streamID uuid.UUID) (bool, error) {
	stream, err := s.repo.GetStream(ctx, streamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("stream not found")
		}
		return false, err
	}

	if !stream.Enabled {
		return false, nil
	}

	pathName := s.pathName(stream)
	return s.mtxCli.IsPathReady(ctx, pathName)
}

// pathName 生成 MediaMTX 路径名
func (s *PreviewService) pathName(stream *domain.Stream) string {
	name := strings.ReplaceAll(stream.Name, " ", "_")
	name = strings.ToLower(name)
	return name
}
