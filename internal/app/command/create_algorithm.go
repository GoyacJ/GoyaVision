package command

import (
	"context"
	"errors"
	"fmt"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/algorithm"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type CreateAlgorithmHandler struct {
	uow port.UnitOfWork
}

func NewCreateAlgorithmHandler(uow port.UnitOfWork) *CreateAlgorithmHandler {
	return &CreateAlgorithmHandler{uow: uow}
}

func (h *CreateAlgorithmHandler) Handle(ctx context.Context, cmd dto.CreateAlgorithmCommand) (*algorithm.Algorithm, error) {
	if cmd.Code == "" {
		return nil, apperr.InvalidInput("code is required")
	}
	if cmd.Name == "" {
		return nil, apperr.InvalidInput("name is required")
	}

	status := cmd.Status
	if status == "" {
		status = algorithm.StatusDraft
	}

	visibility := cmd.Visibility

	var result *algorithm.Algorithm
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if _, err := repos.Algorithms.GetByCode(ctx, cmd.Code); err == nil {
			return apperr.Conflict(fmt.Sprintf("algorithm with code %s already exists", cmd.Code))
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to check algorithm code")
		}

		a := &algorithm.Algorithm{
			Code:           cmd.Code,
			Name:           cmd.Name,
			Description:    cmd.Description,
			Scenario:       cmd.Scenario,
			Status:         status,
			Tags:           cmd.Tags,
			Visibility:     visibility,
			VisibleRoleIDs: cmd.VisibleRoleIDs,
		}
		if err := repos.Algorithms.Create(ctx, a); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to create algorithm")
		}

		if cmd.InitialVersion != nil {
			version, impls, evals, err := buildAlgorithmVersionFromInput(cmd.InitialVersion, a.ID)
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
		}

		var err error
		result, err = repos.Algorithms.GetWithRelations(ctx, a.ID)
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to load algorithm")
		}
		return nil
	})
	return result, err
}
