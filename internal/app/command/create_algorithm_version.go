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

type CreateAlgorithmVersionHandler struct {
	uow port.UnitOfWork
}

func NewCreateAlgorithmVersionHandler(uow port.UnitOfWork) *CreateAlgorithmVersionHandler {
	return &CreateAlgorithmVersionHandler{uow: uow}
}

func (h *CreateAlgorithmVersionHandler) Handle(ctx context.Context, cmd dto.CreateAlgorithmVersionCommand) (*algorithm.Version, error) {
	var result *algorithm.Version
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if _, err := repos.Algorithms.Get(ctx, cmd.AlgorithmID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("algorithm", cmd.AlgorithmID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get algorithm")
		}

		version, impls, evals, err := buildAlgorithmVersionFromCommand(cmd)
		if err != nil {
			return err
		}

		if err := repos.Algorithms.CreateVersion(ctx, version); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to create algorithm version")
		}
		if err := repos.Algorithms.ReplaceImplementations(ctx, version.ID, impls); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to save algorithm implementations")
		}
		if err := repos.Algorithms.ReplaceEvaluations(ctx, version.ID, evals); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to save algorithm evaluations")
		}

		result, err = repos.Algorithms.GetVersion(ctx, version.ID)
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to load algorithm version")
		}
		return nil
	})
	return result, err
}
