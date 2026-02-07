package ai_model

import (
	"time"

	"github.com/google/uuid"
)

type Provider string

const (
	ProviderOpenAI    Provider = "openai"
	ProviderAnthropic Provider = "anthropic"
	ProviderOllama    Provider = "ollama"
	ProviderLocal     Provider = "local"
	ProviderCustom    Provider = "custom"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusDisabled Status = "disabled"
)

type AIModel struct {
	ID        uuid.UUID
	Name      string
	Provider  Provider
	Endpoint  string
	APIKey    string
	ModelName string
	Config    map[string]interface{}
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Filter struct {
	Keyword string
	Limit   int
	Offset  int
}
