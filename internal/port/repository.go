package port

import (
	"context"

	"goyavision/internal/domain"

	"github.com/google/uuid"
)

type Repository interface {
	// Stream
	CreateStream(ctx context.Context, s *domain.Stream) error
	GetStream(ctx context.Context, id uuid.UUID) (*domain.Stream, error)
	ListStreams(ctx context.Context, enabled *bool) ([]*domain.Stream, error)
	UpdateStream(ctx context.Context, s *domain.Stream) error
	DeleteStream(ctx context.Context, id uuid.UUID) error

	// Algorithm
	CreateAlgorithm(ctx context.Context, a *domain.Algorithm) error
	GetAlgorithm(ctx context.Context, id uuid.UUID) (*domain.Algorithm, error)
	ListAlgorithms(ctx context.Context) ([]*domain.Algorithm, error)
	UpdateAlgorithm(ctx context.Context, a *domain.Algorithm) error
	DeleteAlgorithm(ctx context.Context, id uuid.UUID) error

	// AlgorithmBinding
	CreateAlgorithmBinding(ctx context.Context, b *domain.AlgorithmBinding) error
	GetAlgorithmBinding(ctx context.Context, id uuid.UUID) (*domain.AlgorithmBinding, error)
	ListAlgorithmBindingsByStream(ctx context.Context, streamID uuid.UUID) ([]*domain.AlgorithmBinding, error)
	UpdateAlgorithmBinding(ctx context.Context, b *domain.AlgorithmBinding) error
	DeleteAlgorithmBinding(ctx context.Context, id uuid.UUID) error

	// RecordSession
	CreateRecordSession(ctx context.Context, r *domain.RecordSession) error
	GetRecordSession(ctx context.Context, id uuid.UUID) (*domain.RecordSession, error)
	GetRunningRecordSessionByStream(ctx context.Context, streamID uuid.UUID) (*domain.RecordSession, error)
	ListRecordSessionsByStream(ctx context.Context, streamID uuid.UUID) ([]*domain.RecordSession, error)
	UpdateRecordSession(ctx context.Context, r *domain.RecordSession) error

	// InferenceResult
	CreateInferenceResult(ctx context.Context, ir *domain.InferenceResult) error
	ListInferenceResults(ctx context.Context, streamID, bindingID *uuid.UUID, from, to *int64, limit, offset int) ([]*domain.InferenceResult, int64, error)
}
