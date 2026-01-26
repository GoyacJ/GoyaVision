package app

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"goyavision/internal/domain"
	"goyavision/internal/port"
	"goyavision/pkg/ffmpeg"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RecordService struct {
	repo       port.Repository
	manager    *ffmpeg.Manager
	basePath   string
	segmentSec int
	tasks      map[uuid.UUID]*ffmpeg.RecordTask
	tasksMu    sync.RWMutex
}

func NewRecordService(repo port.Repository, manager *ffmpeg.Manager, basePath string, segmentSec int) *RecordService {
	return &RecordService{
		repo:       repo,
		manager:    manager,
		basePath:   basePath,
		segmentSec: segmentSec,
		tasks:      make(map[uuid.UUID]*ffmpeg.RecordTask),
	}
}

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

	basePath := filepath.Join(s.basePath, streamID.String())
	session := &domain.RecordSession{
		ID:        uuid.New(),
		StreamID:  streamID,
		Status:    domain.RecordStatusRunning,
		BasePath:  basePath,
		StartedAt: time.Now(),
	}

	if err := s.repo.CreateRecordSession(ctx, session); err != nil {
		return nil, fmt.Errorf("create record session: %w", err)
	}

	task, err := s.manager.StartRecord(ctx, streamID.String(), stream.URL, s.segmentSec)
	if err != nil {
		session.Status = domain.RecordStatusStopped
		now := time.Now()
		session.StoppedAt = &now
		s.repo.UpdateRecordSession(ctx, session)
		return nil, fmt.Errorf("start recording: %w", err)
	}

	s.tasksMu.Lock()
	s.tasks[session.ID] = task
	s.tasksMu.Unlock()

	go s.monitorTask(ctx, session.ID, task)

	return session, nil
}

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

	s.tasksMu.Lock()
	task, exists := s.tasks[session.ID]
	if exists {
		delete(s.tasks, session.ID)
	}
	s.tasksMu.Unlock()

	if task != nil {
		if err := task.Stop(); err != nil {
			return fmt.Errorf("stop recording task: %w", err)
		}
	}

	now := time.Now()
	session.Status = domain.RecordStatusStopped
	session.StoppedAt = &now

	if err := s.repo.UpdateRecordSession(ctx, session); err != nil {
		return fmt.Errorf("update record session: %w", err)
	}

	return nil
}

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

func (s *RecordService) monitorTask(ctx context.Context, sessionID uuid.UUID, task *ffmpeg.RecordTask) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !task.IsRunning() {
				s.tasksMu.Lock()
				delete(s.tasks, sessionID)
				s.tasksMu.Unlock()

				bgCtx := context.Background()
				session, err := s.repo.GetRecordSession(bgCtx, sessionID)
				if err == nil && session != nil && session.Status == domain.RecordStatusRunning {
					now := time.Now()
					session.Status = domain.RecordStatusStopped
					session.StoppedAt = &now
					s.repo.UpdateRecordSession(bgCtx, session)
				}
				return
			}
		}
	}
}
