package command

import (
	"context"
	"fmt"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/ai_model"
	"goyavision/pkg/apperr"
)

type CreateAIModelHandler struct {
	uow    port.UnitOfWork
	crypto port.CryptoService
}

func NewCreateAIModelHandler(uow port.UnitOfWork, crypto port.CryptoService) *CreateAIModelHandler {
	return &CreateAIModelHandler{uow: uow, crypto: crypto}
}

func (h *CreateAIModelHandler) Handle(ctx context.Context, cmd dto.CreateAIModelCommand) (*ai_model.AIModel, error) {
	if cmd.Name == "" {
		return nil, apperr.InvalidInput("name is required")
	}
	if cmd.Provider == "" {
		return nil, apperr.InvalidInput("provider is required")
	}

	provider := ai_model.Provider(cmd.Provider)
	switch provider {
	case ai_model.ProviderOpenAI, ai_model.ProviderAnthropic, ai_model.ProviderOllama, ai_model.ProviderLocal, ai_model.ProviderCustom:
	default:
		return nil, apperr.InvalidInput(fmt.Sprintf("invalid provider: %s, allowed values: openai|anthropic|ollama|local|custom", cmd.Provider))
	}

	var result *ai_model.AIModel
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		existing, _, err := repos.AIModels.List(ctx, ai_model.Filter{Keyword: cmd.Name, Limit: 100})
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to check name uniqueness")
		}
		for _, m := range existing {
			if m.Name == cmd.Name {
				return apperr.Conflict(fmt.Sprintf("ai model with name %s already exists", cmd.Name))
			}
		}

		apiKey := cmd.APIKey
		if apiKey != "" && h.crypto != nil {
			encrypted, encErr := h.crypto.Encrypt(apiKey)
			if encErr != nil {
				return apperr.Wrap(encErr, apperr.CodeInternal, "failed to encrypt api key")
			}
			apiKey = encrypted
		}

		model := &ai_model.AIModel{
			Name:        cmd.Name,
			Description: cmd.Description,
			Provider:    provider,
			Endpoint:    cmd.Endpoint,
			APIKey:      apiKey,
			ModelName:   cmd.ModelName,
			Config:      cmd.Config,
			Status:      ai_model.StatusActive,
		}

		if err := repos.AIModels.Create(ctx, model); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to create ai model")
		}

		result = model
		return nil
	})

	return result, err
}
