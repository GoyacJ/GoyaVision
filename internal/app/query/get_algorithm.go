package query

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/algorithm"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type GetAlgorithmHandler struct {
	uow port.UnitOfWork
}

func NewGetAlgorithmHandler(uow port.UnitOfWork) *GetAlgorithmHandler {
	return &GetAlgorithmHandler{uow: uow}
}

func (h *GetAlgorithmHandler) Handle(ctx context.Context, query dto.GetAlgorithmQuery) (*algorithm.Algorithm, error) {
	var result *algorithm.Algorithm
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		a, err := repos.Algorithms.GetWithRelations(ctx, query.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("algorithm", query.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get algorithm")
		}
		result = a
		return nil
	})
	return result, err
}
