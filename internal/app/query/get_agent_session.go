package query

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/agent"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type GetAgentSessionHandler struct {
	uow port.UnitOfWork
}

func NewGetAgentSessionHandler(uow port.UnitOfWork) *GetAgentSessionHandler {
	return &GetAgentSessionHandler{uow: uow}
}

func (h *GetAgentSessionHandler) Handle(ctx context.Context, query dto.GetAgentSessionQuery) (*agent.Session, error) {
	var result *agent.Session
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.AgentSessions == nil {
			return apperr.Internal("agent session repository is not configured", nil)
		}
		session, err := repos.AgentSessions.Get(ctx, query.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("agent_session", query.ID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get agent session")
		}
		result = session
		return nil
	})
	return result, err
}
