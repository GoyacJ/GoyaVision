package dto

import (
	"github.com/google/uuid"
)

type CreateAIModelCommand struct {
	Name        string
	Description string
	Provider    string
	Endpoint    string
	APIKey      string
	ModelName   string
	Config      map[string]interface{}
}

type UpdateAIModelCommand struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	Provider    *string
	Endpoint    *string
	APIKey      *string
	ModelName   *string
	Config      map[string]interface{}
	Status      *string
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
