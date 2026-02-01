package app

import (
	"context"

	"goyavision/internal/domain"
	"goyavision/internal/port"

	"github.com/google/uuid"
)

type InferenceService struct {
	repo port.Repository
}

func NewInferenceService(repo port.Repository) *InferenceService {
	return &InferenceService{repo: repo}
}

func (s *InferenceService) ListResults(ctx context.Context, streamID, bindingID *uuid.UUID, from, to *int64, limit, offset int) ([]*domain.InferenceResult, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 1000 {
		limit = 1000
	}

	return s.repo.ListInferenceResults(ctx, streamID, bindingID, from, to, limit, offset)
}
