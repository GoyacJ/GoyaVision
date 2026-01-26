package dto

import (
	"encoding/json"
	"time"

	"goyavision/internal/domain"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AlgorithmCreateReq struct {
	Name       string          `json:"name" validate:"required"`
	Endpoint   string          `json:"endpoint" validate:"required"`
	InputSpec  json.RawMessage `json:"input_spec"`
	OutputSpec json.RawMessage `json:"output_spec"`
}

type AlgorithmUpdateReq struct {
	Name       *string         `json:"name"`
	Endpoint   *string         `json:"endpoint"`
	InputSpec  json.RawMessage `json:"input_spec"`
	OutputSpec json.RawMessage `json:"output_spec"`
}

type AlgorithmResponse struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Endpoint   string          `json:"endpoint"`
	InputSpec  json.RawMessage `json:"input_spec"`
	OutputSpec json.RawMessage `json:"output_spec"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

func AlgorithmToResponse(a *domain.Algorithm) *AlgorithmResponse {
	if a == nil {
		return nil
	}
	return &AlgorithmResponse{
		ID:         a.ID,
		Name:       a.Name,
		Endpoint:   a.Endpoint,
		InputSpec:  a.InputSpec,
		OutputSpec: a.OutputSpec,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
}

func AlgorithmsToResponse(algorithms []*domain.Algorithm) []*AlgorithmResponse {
	result := make([]*AlgorithmResponse, len(algorithms))
	for i, a := range algorithms {
		result[i] = AlgorithmToResponse(a)
	}
	return result
}

func (r *AlgorithmCreateReq) ToDomain() (*domain.Algorithm, error) {
	alg := &domain.Algorithm{
		Name:     r.Name,
		Endpoint: r.Endpoint,
	}
	if r.InputSpec != nil {
		alg.InputSpec = datatypes.JSON(r.InputSpec)
	}
	if r.OutputSpec != nil {
		alg.OutputSpec = datatypes.JSON(r.OutputSpec)
	}
	return alg, nil
}
