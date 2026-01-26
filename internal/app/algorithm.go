package app

import (
	"context"
	"errors"

	"goyavision/internal/domain"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AlgorithmService struct {
	repo port.Repository
}

func NewAlgorithmService(repo port.Repository) *AlgorithmService {
	return &AlgorithmService{repo: repo}
}

func (s *AlgorithmService) Create(ctx context.Context, alg *domain.Algorithm) (*domain.Algorithm, error) {
	if alg.Name == "" {
		return nil, errors.New("name is required")
	}
	if alg.Endpoint == "" {
		return nil, errors.New("endpoint is required")
	}

	if err := s.repo.CreateAlgorithm(ctx, alg); err != nil {
		return nil, err
	}
	return alg, nil
}

func (s *AlgorithmService) Get(ctx context.Context, id uuid.UUID) (*domain.Algorithm, error) {
	alg, err := s.repo.GetAlgorithm(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("algorithm not found")
		}
		return nil, err
	}
	return alg, nil
}

func (s *AlgorithmService) List(ctx context.Context) ([]*domain.Algorithm, error) {
	return s.repo.ListAlgorithms(ctx)
}

func (s *AlgorithmService) Update(ctx context.Context, id uuid.UUID, alg *domain.Algorithm) (*domain.Algorithm, error) {
	existing, err := s.repo.GetAlgorithm(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("algorithm not found")
		}
		return nil, err
	}

	if alg.Name != "" {
		existing.Name = alg.Name
	}
	if alg.Endpoint != "" {
		existing.Endpoint = alg.Endpoint
	}
	if alg.InputSpec != nil {
		existing.InputSpec = alg.InputSpec
	}
	if alg.OutputSpec != nil {
		existing.OutputSpec = alg.OutputSpec
	}

	if err := s.repo.UpdateAlgorithm(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *AlgorithmService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetAlgorithm(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("algorithm not found")
		}
		return err
	}
	return s.repo.DeleteAlgorithm(ctx, id)
}
