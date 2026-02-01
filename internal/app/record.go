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

type RecordService struct {
	repo   port.Repository
	mtxCli *mediamtx.Client
	mtxCfg config.MediaMTX
}

func NewRecordService(repo port.Repository, mtxCli *mediamtx.Client, mtxCfg config.MediaMTX) *RecordService {
	return &RecordService{
		repo:   repo,
		mtxCli: mtxCli,
		mtxCfg: mtxCfg,
	}
}

// Start 开始录制
func (s *RecordService) Start(ctx context.Context, streamID uuid.UUID) (*domain.RecordSession, error) {
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

	existing, err := s.repo.GetRunningRecordSessionByStream(ctx, streamID)
	if err == nil && existing != nil {
		return nil, errors.New("recording already in progress")
	}

	pathName := s.pathName(stream)

	recordPath := fmt.Sprintf("%s/%s/%%Y-%%m-%%d_%%H-%%M-%%S", s.mtxCfg.RecordPath, pathName)
	err = s.mtxCli.EnableRecording(ctx, pathName, recordPath, s.mtxCfg.RecordFormat, s.mtxCfg.SegmentDuration)
	if err != nil {
		return nil, fmt.Errorf("enable mediamtx recording: %w", err)
	}

	session := &domain.RecordSession{
		ID:        uuid.New(),
		StreamID:  streamID,
		Status:    domain.RecordStatusRunning,
		BasePath:  fmt.Sprintf("%s/%s", s.mtxCfg.RecordPath, pathName),
		StartedAt: time.Now(),
	}

	if err := s.repo.CreateRecordSession(ctx, session); err != nil {
		s.mtxCli.DisableRecording(ctx, pathName)
		return nil, fmt.Errorf("create record session: %w", err)
	}

	return session, nil
}

// Stop 停止录制
func (s *RecordService) Stop(ctx context.Context, streamID uuid.UUID) error {
	stream, err := s.repo.GetStream(ctx, streamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("stream not found")
		}
		return err
	}

	session, err := s.repo.GetRunningRecordSessionByStream(ctx, streamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("no active recording session")
		}
		return err
	}

	pathName := s.pathName(stream)
	if err := s.mtxCli.DisableRecording(ctx, pathName); err != nil {
		return fmt.Errorf("disable mediamtx recording: %w", err)
	}

	now := time.Now()
	session.Status = domain.RecordStatusStopped
	session.StoppedAt = &now

	if err := s.repo.UpdateRecordSession(ctx, session); err != nil {
		return fmt.Errorf("update record session: %w", err)
	}

	return nil
}

// ListSessions 列出录制会话
func (s *RecordService) ListSessions(ctx context.Context, streamID uuid.UUID) ([]*domain.RecordSession, error) {
	_, err := s.repo.GetStream(ctx, streamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("stream not found")
		}
		return nil, err
	}

	return s.repo.ListRecordSessionsByStream(ctx, streamID)
}

// GetRecordings 获取流的录制文件列表（从 MediaMTX 获取）
func (s *RecordService) GetRecordings(ctx context.Context, streamID uuid.UUID) (*mediamtx.Recording, error) {
	stream, err := s.repo.GetStream(ctx, streamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("stream not found")
		}
		return nil, err
	}

	pathName := s.pathName(stream)
	return s.mtxCli.GetRecordings(ctx, pathName)
}

// IsRecording 检查流是否正在录制
func (s *RecordService) IsRecording(ctx context.Context, streamID uuid.UUID) (bool, error) {
	session, err := s.repo.GetRunningRecordSessionByStream(ctx, streamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return session != nil && session.Status == domain.RecordStatusRunning, nil
}

// pathName 生成 MediaMTX 路径名
func (s *RecordService) pathName(stream *domain.Stream) string {
	name := strings.ReplaceAll(stream.Name, " ", "_")
	name = strings.ToLower(name)
	return name
}

// SyncRecordingStatus 同步录制状态（从 MediaMTX）
func (s *RecordService) SyncRecordingStatus(ctx context.Context) error {
	streams, err := s.repo.ListStreams(ctx, nil)
	if err != nil {
		return err
	}

	for _, stream := range streams {
		session, err := s.repo.GetRunningRecordSessionByStream(ctx, stream.ID)
		if err != nil || session == nil {
			continue
		}

		pathName := s.pathName(stream)
		cfg, err := s.mtxCli.GetPathConfig(ctx, pathName)
		if err != nil {
			continue
		}

		if cfg.Record == nil || !*cfg.Record {
			now := time.Now()
			session.Status = domain.RecordStatusStopped
			session.StoppedAt = &now
			s.repo.UpdateRecordSession(ctx, session)
		}
	}

	return nil
}
