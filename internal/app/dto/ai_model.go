package dto

import (
	"github.com/google/uuid"
)

type CreateAIModelCommand struct {
	Name      string
	Provider  string
	Endpoint  string
	APIKey    string
	ModelName string
	Config    map[string]interface{}
}

type UpdateAIModelCommand struct {
	ID        uuid.UUID
	Name      *string
	Provider  *string
	Endpoint  *string
	APIKey    *string
	ModelName *string
	Config    map[string]interface{}
	Status    *string
}

type DeleteAIModelCommand struct {
	ID uuid.UUID
}

type GetAIModelQuery struct {
	ID uuid.UUID
}

type ListAIModelsQuery struct {
	Keyword string
	Pagination
}
