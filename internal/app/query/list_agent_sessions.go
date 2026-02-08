package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/agent"
	"goyavision/pkg/apperr"
)

type ListAgentSessionsHandler struct {
	uow port.UnitOfWork
}

func NewListAgentSessionsHandler(uow port.UnitOfWork) *ListAgentSessionsHandler {
	return &ListAgentSessionsHandler{uow: uow}
}

func (h *ListAgentSessionsHandler) Handle(ctx context.Context, query dto.ListAgentSessionsQuery) (*dto.PagedResult[*agent.Session], error) {
	query.Pagination.Normalize()

	var items []*agent.Session
	var total int64
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.AgentSessions == nil {
			return apperr.Internal("agent session repository is not configured", nil)
		}
		var err error
		items, total, err = repos.AgentSessions.List(ctx, agent.SessionFilter{
			TaskID: query.TaskID,
			Status: query.Status,
			Limit:  query.Pagination.Limit,
			Offset: query.Pagination.Offset,
		})
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to list agent sessions")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &dto.PagedResult[*agent.Session]{
		Items:  items,
		Total:  total,
		Limit:  query.Pagination.Limit,
		Offset: query.Pagination.Offset,
	}, nil
}
