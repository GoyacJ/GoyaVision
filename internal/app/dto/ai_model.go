package dto

import (
	"github.com/google/uuid"
	"goyavision/internal/domain/ai_model"
)

type CreateAIModelCommand struct {
	Name           string
	Description    string
	Provider       string
	Endpoint       string
	APIKey         string
	ModelName      string
	Config         map[string]interface{}
	Visibility     ai_model.Visibility
	VisibleRoleIDs []string
}

type UpdateAIModelCommand struct {
	ID             uuid.UUID
	Name           *string
	Description    *string
	Provider       *string
	Endpoint       *string
	APIKey         *string
	ModelName      *string
	Config         map[string]interface{}
	Status         *string
	Visibility     *ai_model.Visibility
	VisibleRoleIDs []string
}

type DeleteAIModelCommand struct {
	ID uuid.UUID
}

type GetAIModelQuery struct {
	ID uuid.UUID
}

type ListAIModelsQuery struct {
	Keyword  string
	Provider *string
	Status   *string
	Pagination
}

type TestAIModelCommand struct {
	ID uuid.UUID
}

type TestAIModelResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
