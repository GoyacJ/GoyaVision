package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/workflow"
	"goyavision/pkg/apperr"
)

type ListWorkflowRevisionsHandler struct {
	uow port.UnitOfWork
}

func NewListWorkflowRevisionsHandler(uow port.UnitOfWork) *ListWorkflowRevisionsHandler {
	return &ListWorkflowRevisionsHandler{uow: uow}
}

func (h *ListWorkflowRevisionsHandler) Handle(ctx context.Context, query dto.ListWorkflowRevisionsQuery) (*dto.PagedResult[*workflow.WorkflowRevision], error) {
	query.Pagination.Normalize()

	var items []*workflow.WorkflowRevision
	var total int64
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		items, total, err = repos.Workflows.ListRevisions(ctx, workflow.RevisionFilter{
			WorkflowID: query.WorkflowID,
			Limit:      query.Pagination.Limit,
			Offset:     query.Pagination.Offset,
		})
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to list workflow revisions")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &dto.PagedResult[*workflow.WorkflowRevision]{
		Items:  items,
		Total:  total,
		Limit:  query.Pagination.Limit,
		Offset: query.Pagination.Offset,
	}, nil
}
