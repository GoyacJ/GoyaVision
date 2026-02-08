package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/workflow"
	"goyavision/pkg/apperr"
)

type ListTaskContextPatchesHandler struct {
	uow port.UnitOfWork
}

func NewListTaskContextPatchesHandler(uow port.UnitOfWork) *ListTaskContextPatchesHandler {
	return &ListTaskContextPatchesHandler{uow: uow}
}

func (h *ListTaskContextPatchesHandler) Handle(ctx context.Context, query dto.ListTaskContextPatchesQuery) (*dto.PagedResult[*workflow.TaskContextPatch], error) {
	query.Pagination.Normalize()

	var items []*workflow.TaskContextPatch
	var total int64
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.Contexts == nil {
			return apperr.Internal("context repository is not configured", nil)
		}
		var err error
		items, total, err = repos.Contexts.ListPatches(ctx, query.TaskID, query.Pagination.Limit, query.Pagination.Offset)
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to list task context patches")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &dto.PagedResult[*workflow.TaskContextPatch]{
		Items:  items,
		Total:  total,
		Limit:  query.Pagination.Limit,
		Offset: query.Pagination.Offset,
	}, nil
}
