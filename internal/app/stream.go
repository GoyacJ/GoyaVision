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

type StreamService struct {
	repo   port.Repository
	mtxCli *mediamtx.Client
	mtxCfg config.MediaMTX
}

func NewStreamService(repo port.Repository, mtxCli *mediamtx.Client, mtxCfg config.MediaMTX) *StreamService {
	return &StreamService{
		repo:   repo,
		mtxCli: mtxCli,
		mtxCfg: mtxCfg,
	}
}

// Create 创建流（同时在 MediaMTX 中注册路径）
func (s *StreamService) Create(ctx context.Context, req *CreateStreamRequest) (*domain.Stream, error) {
	if req.URL == "" && req.Type != domain.StreamTypePush {
		return nil, errors.New("pull stream requires url")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	streamType := domain.StreamTypePull
	if req.Type != "" {
		streamType = req.Type
	}

	stream := &domain.Stream{
		URL:     req.URL,
		Name:    req.Name,
		Type:    streamType,
		Enabled: true,
	}
	if req.Enabled != nil {
		stream.Enabled = *req.Enabled
	}

	if err := s.repo.CreateStream(ctx, stream); err != nil {
		return nil, err
	}

	if stream.Enabled {
		if err := s.registerPath(ctx, stream); err != nil {
			s.repo.DeleteStream(ctx, stream.ID)
			return nil, fmt.Errorf("register mediamtx path: %w", err)
		}
	}

	return stream, nil
}

// Get 获取流
func (s *StreamService) Get(ctx context.Context, id uuid.UUID) (*domain.Stream, error) {
	stream, err := s.repo.GetStream(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("stream not found")
		}
		return nil, err
	}
	return stream, nil
}

// GetWithStatus 获取流及其实时状态
func (s *StreamService) GetWithStatus(ctx context.Context, id uuid.UUID) (*domain.StreamWithStatus, error) {
	stream, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	status := s.getStreamStatus(ctx, stream)
	return &domain.StreamWithStatus{
		Stream: stream,
		Status: status,
	}, nil
}

// List 列出所有流
func (s *StreamService) List(ctx context.Context, enabled *bool) ([]*domain.Stream, error) {
	return s.repo.ListStreams(ctx, enabled)
}

// ListWithStatus 列出所有流及其实时状态
func (s *StreamService) ListWithStatus(ctx context.Context, enabled *bool) ([]*domain.StreamWithStatus, error) {
	streams, err := s.repo.ListStreams(ctx, enabled)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.StreamWithStatus, len(streams))
	for i, stream := range streams {
		status := s.getStreamStatus(ctx, stream)
		result[i] = &domain.StreamWithStatus{
			Stream: stream,
			Status: status,
		}
	}
	return result, nil
}

// Update 更新流
func (s *StreamService) Update(ctx context.Context, id uuid.UUID, req *UpdateStreamRequest) (*domain.Stream, error) {
	stream, err := s.repo.GetStream(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("stream not found")
		}
		return nil, err
	}

	oldEnabled := stream.Enabled
	oldURL := stream.URL

	if req.URL != nil {
		if *req.URL == "" && stream.Type != domain.StreamTypePush {
			return nil, errors.New("url cannot be empty for pull stream")
		}
		stream.URL = *req.URL
	}
	if req.Name != nil {
		if *req.Name == "" {
			return nil, errors.New("name cannot be empty")
		}
		stream.Name = *req.Name
	}
	if req.Enabled != nil {
		stream.Enabled = *req.Enabled
	}

	if err := s.repo.UpdateStream(ctx, stream); err != nil {
		return nil, err
	}

	needReregister := stream.Enabled && (stream.URL != oldURL || !oldEnabled)
	if needReregister {
		if oldEnabled {
			s.unregisterPath(ctx, stream)
		}
		if err := s.registerPath(ctx, stream); err != nil {
			return stream, fmt.Errorf("register mediamtx path: %w", err)
		}
	} else if !stream.Enabled && oldEnabled {
		s.unregisterPath(ctx, stream)
	}

	return stream, nil
}

// Delete 删除流
func (s *StreamService) Delete(ctx context.Context, id uuid.UUID) error {
	stream, err := s.repo.GetStream(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("stream not found")
		}
		return err
	}

	s.unregisterPath(ctx, stream)

	return s.repo.DeleteStream(ctx, id)
}

// Enable 启用流
func (s *StreamService) Enable(ctx context.Context, id uuid.UUID) error {
	stream, err := s.repo.GetStream(ctx, id)
	if err != nil {
		return err
	}

	if stream.Enabled {
		return nil
	}

	stream.Enabled = true
	if err := s.repo.UpdateStream(ctx, stream); err != nil {
		return err
	}

	return s.registerPath(ctx, stream)
}

// Disable 禁用流
func (s *StreamService) Disable(ctx context.Context, id uuid.UUID) error {
	stream, err := s.repo.GetStream(ctx, id)
	if err != nil {
		return err
	}

	if !stream.Enabled {
		return nil
	}

	stream.Enabled = false
	if err := s.repo.UpdateStream(ctx, stream); err != nil {
		return err
	}

	return s.unregisterPath(ctx, stream)
}

// GetStatus 获取流状态
func (s *StreamService) GetStatus(ctx context.Context, id uuid.UUID) (*domain.StreamStatus, error) {
	stream, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.getStreamStatus(ctx, stream), nil
}

// registerPath 在 MediaMTX 中注册路径
func (s *StreamService) registerPath(ctx context.Context, stream *domain.Stream) error {
	pathName := s.pathName(stream)

	cfg := &mediamtx.PathConfig{
		Name: pathName,
	}

	if stream.Type == domain.StreamTypePull {
		cfg.Source = stream.URL
		cfg.SourceOnDemand = true
		cfg.SourceOnDemandStartTimeout = "10s"
		cfg.SourceOnDemandCloseAfter = "10s"
	}

	return s.mtxCli.AddPath(ctx, pathName, cfg)
}

// unregisterPath 从 MediaMTX 中删除路径
func (s *StreamService) unregisterPath(ctx context.Context, stream *domain.Stream) error {
	pathName := s.pathName(stream)
	return s.mtxCli.DeletePath(ctx, pathName)
}

// getStreamStatus 获取流的实时状态
func (s *StreamService) getStreamStatus(ctx context.Context, stream *domain.Stream) *domain.StreamStatus {
	pathName := s.pathName(stream)

	status := &domain.StreamStatus{
		StreamID:  stream.ID,
		PathName:  pathName,
		RTSPUrl:   fmt.Sprintf("%s/%s", s.mtxCfg.RTSPAddress, pathName),
		RTMPUrl:   fmt.Sprintf("%s/%s", s.mtxCfg.RTMPAddress, pathName),
		HLSUrl:    fmt.Sprintf("%s/%s", s.mtxCfg.HLSAddress, pathName),
		WebRTCUrl: fmt.Sprintf("%s/%s", s.mtxCfg.WebRTCAddress, pathName),
	}

	if !stream.Enabled {
		return status
	}

	path, err := s.mtxCli.GetPath(ctx, pathName)
	if err != nil {
		return status
	}

	status.Ready = path.Ready
	status.Online = path.Online
	status.Tracks = path.Tracks
	status.BytesReceived = path.BytesReceived
	status.BytesSent = path.BytesSent
	status.ReaderCount = len(path.Readers)

	return status
}

// pathName 生成 MediaMTX 路径名
func (s *StreamService) pathName(stream *domain.Stream) string {
	name := strings.ReplaceAll(stream.Name, " ", "_")
	name = strings.ToLower(name)
	return name
}

// SyncPaths 同步数据库流与 MediaMTX 路径
func (s *StreamService) SyncPaths(ctx context.Context) error {
	streams, err := s.repo.ListStreams(ctx, nil)
	if err != nil {
		return err
	}

	for _, stream := range streams {
		if stream.Enabled {
			pathName := s.pathName(stream)
			_, err := s.mtxCli.GetPathConfig(ctx, pathName)
			if err != nil {
				if err := s.registerPath(ctx, stream); err != nil {
					continue
				}
			}
		}
	}

	return nil
}

type CreateStreamRequest struct {
	URL     string
	Name    string
	Type    domain.StreamType
	Enabled *bool
}

type UpdateStreamRequest struct {
	URL     *string
	Name    *string
	Enabled *bool
}
