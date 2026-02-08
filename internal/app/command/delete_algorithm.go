package command

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type DeleteAlgorithmHandler struct {
	uow port.UnitOfWork
}

func NewDeleteAlgorithmHandler(uow port.UnitOfWork) *DeleteAlgorithmHandler {
	return &DeleteAlgorithmHandler{uow: uow}
}

func (h *DeleteAlgorithmHandler) Handle(ctx context.Context, cmd dto.DeleteAlgorithmCommand) error {
	return h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if _, err := repos.Algorithms.Get(ctx, cmd.ID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("algorithm", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get algorithm")
		}
		if err := repos.Algorithms.Delete(ctx, cmd.ID); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to delete algorithm")
		}
		return nil
	})
}
