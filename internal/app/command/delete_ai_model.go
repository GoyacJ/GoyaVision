package command

import (
	"context"
	"fmt"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
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
		aiModelMode := operator.ExecModeAIModel
		operators, _, err := repos.Operators.List(ctx, operator.Filter{
			ExecMode: &aiModelMode,
			Limit:    1000,
		})
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to check ai model references")
		}

		for _, op := range operators {
			if op.ActiveVersion == nil || op.ActiveVersion.ExecConfig == nil || op.ActiveVersion.ExecConfig.AIModel == nil {
				continue
			}
			if op.ActiveVersion.ExecConfig.AIModel.ModelID == cmd.ID {
				return apperr.InvalidInput(fmt.Sprintf("cannot delete: ai model is referenced by operator '%s'", op.Name))
			}
		}

		if err := repos.AIModels.Delete(ctx, cmd.ID); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to delete ai model")
		}
		return nil
	})
}
