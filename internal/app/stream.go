package app

import (
	"context"
	"errors"

	"goyavision/internal/domain"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StreamService struct {
	repo port.Repository
}

func NewStreamService(repo port.Repository) *StreamService {
	return &StreamService{repo: repo}
}

func (s *StreamService) Create(ctx context.Context, req *CreateStreamRequest) (*domain.Stream, error) {
	if req.URL == "" {
		return nil, errors.New("url is required")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	stream := &domain.Stream{
		URL:     req.URL,
		Name:    req.Name,
		Enabled: true,
	}
	if req.Enabled != nil {
		stream.Enabled = *req.Enabled
	}

	if err := s.repo.CreateStream(ctx, stream); err != nil {
		return nil, err
	}
	return stream, nil
}

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

func (s *StreamService) List(ctx context.Context, enabled *bool) ([]*domain.Stream, error) {
	return s.repo.ListStreams(ctx, enabled)
}

func (s *StreamService) Update(ctx context.Context, id uuid.UUID, req *UpdateStreamRequest) (*domain.Stream, error) {
	stream, err := s.repo.GetStream(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("stream not found")
		}
		return nil, err
	}

	if req.URL != nil {
		if *req.URL == "" {
			return nil, errors.New("url cannot be empty")
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
	return stream, nil
}

func (s *StreamService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetStream(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("stream not found")
		}
		return err
	}
	return s.repo.DeleteStream(ctx, id)
}

type CreateStreamRequest struct {
	URL     string
	Name    string
	Enabled *bool
}

type UpdateStreamRequest struct {
	URL     *string
	Name    *string
	Enabled *bool
}
