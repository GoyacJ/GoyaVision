package app

import (
	"context"
	"errors"
	"fmt"

	"goyavision/internal/port"
	"goyavision/pkg/preview"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PreviewService struct {
	repo    port.Repository
	manager *preview.Manager
}

func NewPreviewService(repo port.Repository, manager *preview.Manager) *PreviewService {
	return &PreviewService{
		repo:    repo,
		manager: manager,
	}
}

func (s *PreviewService) Start(ctx context.Context, streamID uuid.UUID) (string, error) {
	stream, err := s.repo.GetStream(ctx, streamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("stream not found")
		}
		return "", err
	}

	if !stream.Enabled {
		return "", errors.New("stream is disabled")
	}

	task, exists := s.manager.GetPreview(streamID.String())
	if exists && task.IsRunning() {
		return task.HLSURL, nil
	}

	task, err = s.manager.StartPreview(ctx, streamID.String(), stream.URL)
	if err != nil {
		return "", fmt.Errorf("start preview: %w", err)
	}

	return task.HLSURL, nil
}

func (s *PreviewService) Stop(ctx context.Context, streamID uuid.UUID) error {
	_, err := s.repo.GetStream(ctx, streamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("stream not found")
		}
		return err
	}

	if err := s.manager.StopPreview(streamID.String()); err != nil {
		return fmt.Errorf("stop preview: %w", err)
	}

	return nil
}
