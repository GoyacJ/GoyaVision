package command

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/pkg/apperr"
)

type DeleteAIModelHandler struct {
	uow port.UnitOfWork
}

func NewDeleteAIModelHandler(uow port.UnitOfWork) *DeleteAIModelHandler {
	return &DeleteAIModelHandler{uow: uow}
}

func (h *DeleteAIModelHandler) Handle(ctx context.Context, cmd dto.DeleteAIModelCommand) error {
	return h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if err := repos.AIModels.Delete(ctx, cmd.ID); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to delete ai model")
		}
		return nil
	})
}
