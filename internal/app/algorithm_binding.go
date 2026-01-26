package app

import (
	"context"
	"errors"

	"goyavision/internal/domain"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AlgorithmBindingService struct {
	repo port.Repository
}

func NewAlgorithmBindingService(repo port.Repository) *AlgorithmBindingService {
	return &AlgorithmBindingService{repo: repo}
}

func (s *AlgorithmBindingService) Create(ctx context.Context, streamID uuid.UUID, binding *domain.AlgorithmBinding) (*domain.AlgorithmBinding, error) {
	if binding.IntervalSec <= 0 {
		return nil, errors.New("interval_sec must be greater than 0")
	}

	_, err := s.repo.GetStream(ctx, streamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("stream not found")
		}
		return nil, err
	}

	_, err = s.repo.GetAlgorithm(ctx, binding.AlgorithmID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("algorithm not found")
		}
		return nil, err
	}

	binding.StreamID = streamID
	if err := s.repo.CreateAlgorithmBinding(ctx, binding); err != nil {
		return nil, err
	}
	return binding, nil
}

func (s *AlgorithmBindingService) Get(ctx context.Context, id uuid.UUID) (*domain.AlgorithmBinding, error) {
	binding, err := s.repo.GetAlgorithmBinding(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("algorithm binding not found")
		}
		return nil, err
	}
	return binding, nil
}

func (s *AlgorithmBindingService) ListByStream(ctx context.Context, streamID uuid.UUID) ([]*domain.AlgorithmBinding, error) {
	_, err := s.repo.GetStream(ctx, streamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("stream not found")
		}
		return nil, err
	}

	return s.repo.ListAlgorithmBindingsByStream(ctx, streamID)
}

func (s *AlgorithmBindingService) Update(ctx context.Context, id uuid.UUID, binding *domain.AlgorithmBinding) (*domain.AlgorithmBinding, error) {
	existing, err := s.repo.GetAlgorithmBinding(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("algorithm binding not found")
		}
		return nil, err
	}

	if binding.AlgorithmID != uuid.Nil {
		_, err = s.repo.GetAlgorithm(ctx, binding.AlgorithmID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("algorithm not found")
			}
			return nil, err
		}
		existing.AlgorithmID = binding.AlgorithmID
	}
	if binding.IntervalSec > 0 {
		existing.IntervalSec = binding.IntervalSec
	}
	if binding.InitialDelaySec >= 0 {
		existing.InitialDelaySec = binding.InitialDelaySec
	}
	if binding.Schedule != nil {
		existing.Schedule = binding.Schedule
	}
	if binding.Config != nil {
		existing.Config = binding.Config
	}
	if binding.Enabled != existing.Enabled {
		existing.Enabled = binding.Enabled
	}

	if err := s.repo.UpdateAlgorithmBinding(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *AlgorithmBindingService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetAlgorithmBinding(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("algorithm binding not found")
		}
		return err
	}
	return s.repo.DeleteAlgorithmBinding(ctx, id)
}
