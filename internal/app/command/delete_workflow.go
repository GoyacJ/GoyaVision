package command

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type DeleteWorkflowHandler struct {
	uow port.UnitOfWork
}

func NewDeleteWorkflowHandler(uow port.UnitOfWork) *DeleteWorkflowHandler {
	return &DeleteWorkflowHandler{uow: uow}
}

func (h *DeleteWorkflowHandler) Handle(ctx context.Context, cmd dto.DeleteWorkflowCommand) error {
	return h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if _, err := repos.Workflows.Get(ctx, cmd.ID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("workflow", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get workflow")
		}

		if err := repos.Workflows.Delete(ctx, cmd.ID); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to delete workflow")
		}

		return nil
	})
}
