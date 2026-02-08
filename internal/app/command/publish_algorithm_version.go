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

type PublishAlgorithmVersionHandler struct {
	uow port.UnitOfWork
}

func NewPublishAlgorithmVersionHandler(uow port.UnitOfWork) *PublishAlgorithmVersionHandler {
	return &PublishAlgorithmVersionHandler{uow: uow}
}

func (h *PublishAlgorithmVersionHandler) Handle(ctx context.Context, cmd dto.PublishAlgorithmVersionCommand) (*algorithm.Version, error) {
	var result *algorithm.Version
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		a, err := repos.Algorithms.Get(ctx, cmd.AlgorithmID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("algorithm", cmd.AlgorithmID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get algorithm")
		}

		v, err := repos.Algorithms.GetVersion(ctx, cmd.VersionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("algorithm_version", cmd.VersionID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get algorithm version")
		}
		if v.AlgorithmID != cmd.AlgorithmID {
			return apperr.InvalidInput("algorithm version does not belong to algorithm")
		}

		v.Status = algorithm.VersionStatusPublished
		if err := repos.Algorithms.UpdateVersion(ctx, v); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to publish algorithm version")
		}

		a.Status = algorithm.StatusPublished
		if err := repos.Algorithms.Update(ctx, a); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to update algorithm status")
		}

		result = v
		return nil
	})
	return result, err
}
