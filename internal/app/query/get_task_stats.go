package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/workflow"
	"goyavision/pkg/apperr"
)

type GetTaskStatsHandler struct {
	uow port.UnitOfWork
}

func NewGetTaskStatsHandler(uow port.UnitOfWork) *GetTaskStatsHandler {
	return &GetTaskStatsHandler{uow: uow}
}

func (h *GetTaskStatsHandler) Handle(ctx context.Context, query dto.GetTaskStatsQuery) (*workflow.TaskStats, error) {
	var result *workflow.TaskStats
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		stats, err := repos.Tasks.GetStats(ctx, query.WorkflowID)
		if err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get task stats")
		}
		result = stats
		return nil
	})

	return result, err
}
