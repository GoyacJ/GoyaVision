package ai_model

import (
	"fmt"
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
	ID          uuid.UUID
	Name        string
	Description string
	Provider    Provider
	Endpoint    string
	APIKey      string
	ModelName   string
	Config      map[string]interface{}
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m *AIModel) IsActive() bool {
	return m.Status == StatusActive
}

func (m *AIModel) IsDisabled() bool {
	return m.Status == StatusDisabled
}

func (m *AIModel) MaskAPIKey() string {
	if len(m.APIKey) <= 8 {
		return "***"
	}
	return m.APIKey[:3] + "..." + m.APIKey[len(m.APIKey)-3:]
}

func (m *AIModel) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("name is required")
	}
	switch m.Provider {
	case ProviderOpenAI, ProviderAnthropic, ProviderOllama, ProviderLocal, ProviderCustom:
	default:
		return fmt.Errorf("invalid provider: %s", m.Provider)
	}
	return nil
}

type Filter struct {
	Keyword  string
	Provider *Provider
	Status   *Status
	Limit    int
	Offset   int
}
