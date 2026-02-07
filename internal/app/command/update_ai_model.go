package command

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/ai_model"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type UpdateAIModelHandler struct {
	uow port.UnitOfWork
}

func NewUpdateAIModelHandler(uow port.UnitOfWork) *UpdateAIModelHandler {
	return &UpdateAIModelHandler{uow: uow}
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
		if cmd.Provider != nil {
			model.Provider = ai_model.Provider(*cmd.Provider)
		}
		if cmd.Endpoint != nil {
			model.Endpoint = *cmd.Endpoint
		}
		if cmd.APIKey != nil {
			model.APIKey = *cmd.APIKey
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

		if err := repos.AIModels.Update(ctx, model); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to update ai model")
		}

		result = model
		return nil
	})

	return result, err
}
