package command

import (
	"context"
	"errors"
	"fmt"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/ai_model"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type UpdateAIModelHandler struct {
	uow    port.UnitOfWork
	crypto port.CryptoService
}

func NewUpdateAIModelHandler(uow port.UnitOfWork, crypto port.CryptoService) *UpdateAIModelHandler {
	return &UpdateAIModelHandler{uow: uow, crypto: crypto}
}

func (h *UpdateAIModelHandler) Handle(ctx context.Context, cmd dto.UpdateAIModelCommand) (*ai_model.AIModel, error) {
	var result *ai_model.AIModel
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		model, err := repos.AIModels.Get(ctx, cmd.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("ai model", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get ai model")
		}

		if cmd.Name != nil {
			model.Name = *cmd.Name
		}
		if cmd.Description != nil {
			model.Description = *cmd.Description
		}
		if cmd.Provider != nil {
			p := ai_model.Provider(*cmd.Provider)
			switch p {
			case ai_model.ProviderOpenAI, ai_model.ProviderAnthropic, ai_model.ProviderOllama, ai_model.ProviderLocal, ai_model.ProviderCustom,
				ai_model.ProviderQwen, ai_model.ProviderDoubao, ai_model.ProviderZhipu, ai_model.ProviderVLLM:
			default:
				return apperr.InvalidInput(fmt.Sprintf("invalid provider: %s", *cmd.Provider))
			}
			model.Provider = p
		}
		if cmd.Endpoint != nil {
			model.Endpoint = *cmd.Endpoint
		}
		if cmd.APIKey != nil {
			apiKey := *cmd.APIKey
			if apiKey != "" && h.crypto != nil {
				encrypted, encErr := h.crypto.Encrypt(apiKey)
				if encErr != nil {
					return apperr.Wrap(encErr, apperr.CodeInternal, "failed to encrypt api key")
				}
				apiKey = encrypted
			}
			model.APIKey = apiKey
		}
		if cmd.ModelName != nil {
			model.ModelName = *cmd.ModelName
		}
		if cmd.Config != nil {
			model.Config = cmd.Config
		}
		if cmd.Status != nil {
			model.Status = ai_model.Status(*cmd.Status)
		}
		if cmd.Visibility != nil {
			model.Visibility = *cmd.Visibility
		}

		if err := repos.AIModels.Update(ctx, model); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to update ai model")
		}

		result = model
		return nil
	})

	return result, err
}
