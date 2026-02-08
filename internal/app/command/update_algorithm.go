package command

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/algorithm"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type UpdateAlgorithmHandler struct {
	uow port.UnitOfWork
}

func NewUpdateAlgorithmHandler(uow port.UnitOfWork) *UpdateAlgorithmHandler {
	return &UpdateAlgorithmHandler{uow: uow}
}

func (h *UpdateAlgorithmHandler) Handle(ctx context.Context, cmd dto.UpdateAlgorithmCommand) (*algorithm.Algorithm, error) {
	var result *algorithm.Algorithm
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		a, err := repos.Algorithms.Get(ctx, cmd.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("algorithm", cmd.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get algorithm")
		}

		if cmd.Name != nil {
			a.Name = *cmd.Name
		}
		if cmd.Description != nil {
			a.Description = *cmd.Description
		}
		if cmd.Scenario != nil {
			a.Scenario = *cmd.Scenario
		}
		if cmd.Status != nil {
			a.Status = *cmd.Status
		}
		if len(cmd.Tags) > 0 {
			a.Tags = cmd.Tags
		}
		if cmd.Visibility != nil {
			a.Visibility = *cmd.Visibility
		}
		if cmd.VisibleRoleIDs != nil {
			a.VisibleRoleIDs = cmd.VisibleRoleIDs
		}

		if err := repos.Algorithms.Update(ctx, a); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to update algorithm")
		}
		result, err = repos.Algorithms.GetWithRelations(ctx, a.ID)
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to load algorithm")
		}
		return nil
	})
	return result, err
}
