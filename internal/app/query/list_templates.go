package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
	"goyavision/pkg/apperr"
)

type ListTemplatesHandler struct {
	uow port.UnitOfWork
}

func NewListTemplatesHandler(uow port.UnitOfWork) *ListTemplatesHandler {
	return &ListTemplatesHandler{uow: uow}
}

func (h *ListTemplatesHandler) Handle(ctx context.Context, q dto.ListTemplatesQuery) (*dto.PagedResult[*operator.OperatorTemplate], error) {
	q.Pagination.Normalize()

	filter := operator.TemplateFilter{
		Category: q.Category,
		Type:     q.Type,
		ExecMode: q.ExecMode,
		Keyword:  q.Keyword,
		Tags:     q.Tags,
		Limit:    q.Pagination.Limit,
		Offset:   q.Pagination.Offset,
	}

	var items []*operator.OperatorTemplate
	var total int64
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		items, total, err = repos.OperatorTemplates.List(ctx, filter)
		return err
	})
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeDBError, "failed to list templates")
	}

	return &dto.PagedResult[*operator.OperatorTemplate]{
		Items:  items,
		Total:  total,
		Limit:  q.Pagination.Limit,
		Offset: q.Pagination.Offset,
	}, nil
}
