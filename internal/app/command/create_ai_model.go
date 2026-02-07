package command

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/ai_model"
	"goyavision/pkg/apperr"
)

type CreateAIModelHandler struct {
	uow port.UnitOfWork
}

func NewCreateAIModelHandler(uow port.UnitOfWork) *CreateAIModelHandler {
	return &CreateAIModelHandler{uow: uow}
}

func (h *CreateAIModelHandler) Handle(ctx context.Context, cmd dto.CreateAIModelCommand) (*ai_model.AIModel, error) {
	if cmd.Name == "" {
		return nil, apperr.InvalidInput("name is required")
	}
	if cmd.Provider == "" {
		return nil, apperr.InvalidInput("provider is required")
	}

	var result *ai_model.AIModel
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		model := &ai_model.AIModel{
			Name:      cmd.Name,
			Provider:  ai_model.Provider(cmd.Provider),
			Endpoint:  cmd.Endpoint,
			APIKey:    cmd.APIKey,
			ModelName: cmd.ModelName,
			Config:    cmd.Config,
			Status:    ai_model.StatusActive,
		}

		if err := repos.AIModels.Create(ctx, model); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to create ai model")
		}

		result = model
		return nil
	})

	return result, err
}
