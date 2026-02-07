package dto

import (
	"time"

	"github.com/google/uuid"
	"goyavision/internal/domain/ai_model"
)

type AIModelCreateReq struct {
	Name      string                 `json:"name"`
	Provider  string                 `json:"provider"`
	Endpoint  string                 `json:"endpoint"`
	APIKey    string                 `json:"api_key"`
	ModelName string                 `json:"model_name"`
	Config    map[string]interface{} `json:"config"`
}

type AIModelUpdateReq struct {
	Name      *string                 `json:"name"`
	Provider  *string                 `json:"provider"`
	Endpoint  *string                 `json:"endpoint"`
	APIKey    *string                 `json:"api_key"`
	ModelName *string                 `json:"model_name"`
	Config    map[string]interface{}  `json:"config"`
	Status    *string                 `json:"status"`
}

type AIModelListQuery struct {
	Keyword string `query:"keyword"`
	Limit   int    `query:"limit"`
	Offset  int    `query:"offset"`
}

type AIModelResponse struct {
	ID        uuid.UUID              `json:"id"`
	Name      string                 `json:"name"`
	Provider  string                 `json:"provider"`
	Endpoint  string                 `json:"endpoint"`
	ModelName string                 `json:"model_name"`
	Config    map[string]interface{} `json:"config"`
	Status    string                 `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type AIModelListResponse struct {
	Items []*AIModelResponse `json:"items"`
	Total int64              `json:"total"`
}

func AIModelToResponse(m *ai_model.AIModel) *AIModelResponse {
	return &AIModelResponse{
		ID:        m.ID,
		Name:      m.Name,
		Provider:  string(m.Provider),
		Endpoint:  m.Endpoint,
		ModelName: m.ModelName,
		Config:    m.Config,
		Status:    string(m.Status),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func AIModelsToResponse(models []*ai_model.AIModel) []*AIModelResponse {
	res := make([]*AIModelResponse, len(models))
	for i, m := range models {
		res[i] = AIModelToResponse(m)
	}
	return res
}
